package node

import (
	"crypto/tls"
	"time"
)

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

func pingMojang() (float64, error) {
	payload := "PUT /minecraft/profile/name/test HTTP/1.1\r\nHost: api.minecraftservices.com\r\nAuthorization: Bearer BEARER" + "\r\n"
	conn, err := tls.Dial("tcp", "api.minecraftservices.com:443", nil)
	if err != nil {
	}
	var sumNanos int64
	for i := 0; i < 3; i++ {
		junk := make([]byte, 4096)
		conn.Write([]byte(payload))
		time1 := time.Now()
		conn.Write([]byte("\r\n"))
		conn.Read(junk)
		duration := time.Now().Sub(time1)
		sumNanos += duration.Nanoseconds()
	}
	conn.Close()
	sumNanos /= 3
	avgMillis := float64(sumNanos)/float64(1000000)
	return avgMillis, nil
}