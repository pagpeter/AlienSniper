package host

import (
	"Alien/types"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type bearerMs struct {
	Bearer string `json:"access_token"`
}

type ServerInfo struct {
	Webhook string
	SkinUrl string
}

type securityRes struct {
	Answer answerRes `json:"answer"`
}

type answerRes struct {
	ID int `json:"id"`
}

type accessTokenReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type accessTokenResp struct {
	AccessToken *string `json:"accessToken"`
}

var (
	// RequestMap map[string]int = map[string]int{
	// 	"mojang":    state.Config.Requests.Mojang,
	// 	"giftcard":  state.Config.Requests.Giftcard,
	// 	"microsoft": state.Config.Requests.Microsoft,
	// }
	redirect string
	i        int
	g        int
)

func Auth(email, password, info string, p types.Packet) (string, string) {
	// returns account type, bearer
	var use bool
	var acctype string
	var bearer string

	if i == 3 {
		time.Sleep(time.Minute)
		i = 0
	}

	time.Sleep(time.Second)

	jar, _ := cookiejar.New(nil)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirect = req.URL.String()
			return nil
		},
		Jar: jar,
	}

	resp, _ := http.NewRequest("GET", "https://login.live.com/oauth20_authorize.srf?client_id=000000004C12AE6F&redirect_uri=https://login.live.com/oauth20_desktop.srf&scope=service::user.auth.xboxlive.com::MBI_SSL&display=touch&response_type=token&locale=en", nil)

	resp.Header.Set("User-Agent", "Mozilla/5.0 (XboxReplay; XboxLiveAuth/3.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")

	response, _ := client.Do(resp)

	jar.Cookies(resp.URL)
	bodyByte, _ := ioutil.ReadAll(response.Body)
	myString := string(bodyByte[:])

	search1 := regexp.MustCompile(`value="(.*?)"`)
	search3 := regexp.MustCompile(`urlPost:'(.+?)'`)

	value := search1.FindAllStringSubmatch(myString, -1)[0][1]
	urlPost := search3.FindAllStringSubmatch(myString, -1)[0][1]

	emailEncode := url.QueryEscape(email)
	passwordEncode := url.QueryEscape(password)

	body := []byte(fmt.Sprintf("login=%v&loginfmt=%v&passwd=%v&PPFT=%v", emailEncode, emailEncode, passwordEncode, value))

	req, _ := http.NewRequest("POST", urlPost, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (XboxReplay; XboxLiveAuth/3.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
	client.Do(req)

	respBytes, _ := ioutil.ReadAll(response.Body)

	if strings.Contains(string(respBytes), "Sign in to") {
		bearer = "Invalid"
	}

	if strings.Contains(string(respBytes), "Help us protect your account") {
		bearer = "Invalid"
	}

	if !strings.Contains(redirect, "access_token") || redirect == urlPost {
		bearer = "Invalid"
	}

	if bearer != "Invalid" {
		client := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Renegotiation: tls.RenegotiateFreelyAsClient,
				},
			},
		}

		splitBear := strings.Split(redirect, "#")[1]
		splitValues := strings.Split(splitBear, "&")

		//refresh_token := strings.Split(splitValues[4], "=")[1]
		access_token := strings.Split(splitValues[0], "=")[1]
		//expires_in := strings.Split(splitValues[2], "=")[1]

		body := []byte(`{"Properties": {"AuthMethod": "RPS", "SiteName": "user.auth.xboxlive.com", "RpsTicket": "` + access_token + `"}, "RelyingParty": "http://auth.xboxlive.com", "TokenType": "JWT"}`)
		post, _ := http.NewRequest("POST", "https://user.auth.xboxlive.com/user/authenticate", bytes.NewBuffer(body))

		post.Header.Set("Content-Type", "application/json")
		post.Header.Set("Accept", "application/json")

		bodyRP, _ := client.Do(post)

		rpBody, _ := ioutil.ReadAll(bodyRP.Body)

		Token := func(body string, key string) string {
			keystr := "\"" + key + "\":[^,;\\]}]*"
			r, _ := regexp.Compile(keystr)
			match := r.FindString(body)
			keyValMatch := strings.Split(match, ":")
			return strings.ReplaceAll(keyValMatch[1], "\"", "")
		}(string(rpBody), "Token")
		uhs := func(body string, key string) string {
			keystr := "\"" + key + "\":[^,;\\]}]*"
			r, _ := regexp.Compile(keystr)
			match := r.FindString(body)
			keyValMatch := strings.Split(match, ":")
			return strings.ReplaceAll(keyValMatch[1], "\"", "")
		}(string(rpBody), "uhs")

		payload := []byte(`{"Properties": {"SandboxId": "RETAIL", "UserTokens": ["` + Token + `"]}, "RelyingParty": "rp://api.minecraftservices.com/", "TokenType": "JWT"}`)
		xstsPost, _ := http.NewRequest("POST", "https://xsts.auth.xboxlive.com/xsts/authorize", bytes.NewBuffer(payload))
		xstsPost.Header.Set("Content-Type", "application/json")
		xstsPost.Header.Set("Accept", "application/json")

		bodyXS, _ := client.Do(xstsPost)

		xsBody, _ := ioutil.ReadAll(bodyXS.Body)

		var cont bool = true
		switch bodyXS.StatusCode {
		case 401:
			switch !strings.Contains(string(xsBody), "XErr") {
			case !strings.Contains(string(xsBody), "2148916238"):
				cont = false
			case !strings.Contains(string(xsBody), "2148916233"):
				cont = false
			}
		}

		if cont {
			xsToken := func(body string, key string) string {
				keystr := "\"" + key + "\":[^,;\\]}]*"
				r, _ := regexp.Compile(keystr)
				match := r.FindString(body)
				keyValMatch := strings.Split(match, ":")
				return strings.ReplaceAll(keyValMatch[1], "\"", "")
			}(string(xsBody), "Token")

			mcBearer := []byte(`{"identityToken" : "XBL3.0 x=` + uhs + `;` + xsToken + `", "ensureLegacyEnabled" : true}`)
			mcBPOST, _ := http.NewRequest("POST", "https://api.minecraftservices.com/authentication/login_with_xbox", bytes.NewBuffer(mcBearer))

			mcBPOST.Header.Set("Content-Type", "application/json")

			bodyBearer, _ := client.Do(mcBPOST)

			bearerValue, _ := ioutil.ReadAll(bodyBearer.Body)

			var bearerMS bearerMs
			json.Unmarshal(bearerValue, &bearerMS)

			accountType := func() string {
				var accountT string
				conn, _ := tls.Dial("tcp", "api.minecraftservices.com"+":443", nil)

				fmt.Fprintln(conn, "GET /minecraft/profile/namechange HTTP/1.1\r\nHost: api.minecraftservices.com\r\nUser-Agent: Alien/1.0\r\nAuthorization: Bearer "+bearerMS.Bearer+"\r\n\r\n")

				e := make([]byte, 12)
				conn.Read(e)

				switch string(e[9:12]) {
				case `404`:
					accountT = "giftcard"
				default:
					accountT = "microsoft"
				}
				return accountT
			}()

			bearer = bearerMS.Bearer
			acctype = accountType
			use = true
			i++
		} else {
			use = false
		}
	} else {
		if g == 10 {
			time.Sleep(30 * time.Second)
			g = 0
		}

		var access accessTokenResp

		splitLogin := strings.Split(info, ":")

		data := accessTokenReq{
			Username: email,
			Password: password,
		}

		bytesToSend, _ := json.Marshal(data)

		req, _ := http.NewRequest("POST", "https://authserver.mojang.com/authenticate", bytes.NewBuffer(bytesToSend))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Alien/1.0")

		res, _ := http.DefaultClient.Do(req)

		if res.Status == "200 OK" {
			respData, _ := ioutil.ReadAll(res.Body)

			json.Unmarshal(respData, &access)

			if len(strings.Split(info, ":")) != 5 {
				bearer = *access.AccessToken
				acctype = "microsoft"
				use = true
			} else {
				req, _ = http.NewRequest("GET", "https://api.mojang.com/user/security/challenges", nil)

				req.Header.Set("Authorization", "Bearer "+*access.AccessToken)
				res, _ = http.DefaultClient.Do(req)

				respData, _ = ioutil.ReadAll(res.Body)

				var security []securityRes
				json.Unmarshal(respData, &security)

				if len(security) == 3 {
					dataBytes := []byte(`[{"id": ` + strconv.Itoa(security[0].Answer.ID) + `, "answer": "` + splitLogin[2] + `"}, {"id": ` + strconv.Itoa(security[1].Answer.ID) + `, "answer": "` + splitLogin[3] + `"}, {"id": ` + strconv.Itoa(security[2].Answer.ID) + `, "answer": "` + splitLogin[4] + `"}]`)
					req, _ = http.NewRequest("POST", "https://api.mojang.com/user/security/location", bytes.NewReader(dataBytes))

					req.Header.Set("Authorization", "Bearer "+*access.AccessToken)
					resp, _ := http.DefaultClient.Do(req)
					if resp.StatusCode == 204 {
						bearer = *access.AccessToken
						acctype = "mojang"
						use = true
					}
				}
			}
		} else {
			use = false
		}
		g++
	}

	if use {
		return bearer, acctype
	}
	return "", ""
	// acc := types.StoredAccount{
	// 	Email:        email,
	// 	Password:     password,
	// 	Type:         acctype,
	// 	Group:        p.Content.Account.Group,
	// 	Usable:       use,
	// 	Bearer:       bearer,
	// 	LastAuthed:   time.Now().Unix(),
	// 	AuthInterval: 86400,
	// }

}
