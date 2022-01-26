package node

import (
	types "Alien/types"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Payload struct {
	Payload     []string
	Conns       []*tls.Conn
	AccountType []string
}

type MCbearers struct {
	Bearers     []string
	AccountType []string
	Emails      []string
}

// taken from https://github.com/Liza-Developer/apiGO/blob/main/mcsn.go
// https://github.com/MCGoSnipe/Runtime
// https://github.com/Kqzz/MCsniperGO

func Sleep(dropTime int64, delay float64) {
	dropStamp := time.Unix(dropTime, 0)

	time.Sleep(time.Until(dropStamp.Add(time.Millisecond * time.Duration(0-delay)).Add(time.Duration(-float64(time.Since(time.Now()).Nanoseconds())/1000000.0) * time.Millisecond)))
}

func saveLogs(content string) {

	var logFile *os.File

	if _, err := os.Stat("logs.txt"); errors.Is(err, os.ErrNotExist) {
		logFile, _ = os.Create("logs.txt")

		text := "               _,--=--._\n"
		text += "             ,'    _    `.\n"
		text += "            -    _(_)_o   - \n"
		text += "       ____'    /_  _/]    `____\n"
		text += "-=====::(+):::::::::::::::::(+)::=====-\n"
		text += `         (+).""""""""""""",(+)` + "\n"
		text += "             .           ,\n"
		text += "               `  -=-  '\n\n\n"

		logFile.WriteString(text)
	}

	log.Printf(content + "\n")

	logFile, _ = os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logFile.Close()

	logFile.WriteString(time.Now().Format("2006-01-02 15:04:05") + " " + content + "\n")
}

func PreSleep(dropTime int64) {
	dropStamp := time.Unix(dropTime, 0)
	for {
		time.Sleep(time.Second * 1)
		if time.Until(dropStamp) <= 5*time.Second {
			break
		}
	}
}

// lizas implementation
func pingMojang() float64 {
	var pingTimes float64
	conn, _ := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
	defer conn.Close()
	for i := 0; i < 10; i++ {
		recv := make([]byte, 4096)
		time1 := time.Now()
		conn.Write([]byte("PUT /minecraft/profile/name/test HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer TestToken\r\n\r\n"))
		conn.Read(recv)
		pingTimes += float64(time.Since(time1).Milliseconds())
	}
	return float64(pingTimes/10000) * 5000
}

func (bearers MCbearers) AddAccounts(accounts []types.StoredAccount) MCbearers {
	for _, details := range accounts {
		bearers.Bearers = append(bearers.Bearers, details.Bearer)
		bearers.AccountType = append(bearers.AccountType, details.Type)
		bearers.Emails = append(bearers.Emails, details.Email)
	}

	return bearers
}

func (bearers MCbearers) RemoveAccounts() MCbearers {
	bearers.AccountType = []string{}
	bearers.Bearers = []string{}
	return bearers
}

func (accountBearer MCbearers) CreatePayloads(name string) Payload {
	payload := make([]string, 0)
	var conns []*tls.Conn

	for i, bearer := range accountBearer.Bearers {
		if accountBearer.AccountType[i] == "giftcard" {
			payload = append(payload, fmt.Sprintf("POST /minecraft/profile HTTP/1.1\r\nHost: api.minecraftservices.com\r\nConnection: open\r\nContent-Length:%s\r\nContent-Type: application/json\r\nAccept: application/json\r\nAuthorization: Bearer %s\r\n\r\n"+string([]byte(`{"profileName":"`+name+`"}`))+"\r\n", strconv.Itoa(len(string([]byte(`{"profileName":"`+name+`"}`)))), bearer))
		} else {
			payload = append(payload, "PUT /minecraft/profile/name/"+name+" HTTP/1.1\r\nHost: api.minecraftservices.com\r\nUser-Agent: MCSN/1.0\r\nAuthorization: bearer "+bearer+"\r\n\r\n")
		}
	}

	for range payload {
		conn, _ := tls.Dial("tcp", "api.minecraftservices.com"+":443", nil)
		conns = append(conns, conn)
	}

	return Payload{Payload: payload, Conns: conns, AccountType: accountBearer.AccountType}
}
