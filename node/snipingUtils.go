package node

import (
	types "Alien/types"
	"crypto/tls"
	"fmt"
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
}

// taken from https://github.com/Liza-Developer/apiGO/blob/main/mcsn.go
// https://github.com/MCGoSnipe/Runtime
// https://github.com/Kqzz/MCsniperGO

func Sleep(dropTime int64, delay float64) {
	dropStamp := time.Unix(dropTime, 0)

	time.Sleep(time.Until(dropStamp.Add(time.Millisecond * time.Duration(0-delay)).Add(time.Duration(-float64(time.Since(time.Now()).Nanoseconds())/1000000.0) * time.Millisecond)))
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

	for i := 0; i < 10; i++ {
		junk := make([]byte, 4096)
		time1 := time.Now()
		conn.Write([]byte("PUT /minecraft/profile/name/test HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer TestToken\r\n\r\n"))
		conn.Read(junk)
		time2 := time.Since(time1)
		pingTimes += float64(time2.Milliseconds())
	}

	conn.Close()

	return float64(pingTimes/10000) * 5000
}

// GoSnipe/cowbos/kqzz's implementation
// func pingMojang() (float64) {
// 	payload := "PUT /minecraft/profile/name/test HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer BEARER" + "\r\n"
// 	conn, _ := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
// 	var sumNanos int64
// 	for i := 0; i < 3; i++ {
// 		junk := make([]byte, 4096)
// 		conn.Write([]byte(payload))
// 		time1 := time.Now()
// 		conn.Write([]byte("\r\n"))
// 		conn.Read(junk)
// 		duration := time.Now().Sub(time1)
// 		sumNanos += duration.Nanoseconds()
// 	}
// 	conn.Close()
// 	sumNanos /= 3
// 	avgMillis := float64(sumNanos)/float64(1000000)
// 	return avgMillis, nil
// }

func (bearers MCbearers) AddAccounts(accounts []types.StoredAccount) MCbearers {
	for _, acc := range accounts {
		bearers.Bearers = append(bearers.Bearers, acc.Bearer)
		bearers.AccountType = append(bearers.AccountType, acc.Type)
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
		if accountBearer.AccountType[i] == "Giftcard" {
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
