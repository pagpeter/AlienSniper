package node

import (
	types "Alien/types"
	"fmt"
	"strconv"
	"sync"
	"time"
)

const delay = 0

var (
	bearers MCbearers
)

func StartSniper(timestamp int64, delay int, name string, i int, payload Payload, email string) types.Log {
	var l types.Log
	var send []time.Time
	var statuscode []int
	var requests []types.RequestLog

	for g := 0; g < 2; {
		recvd := make([]byte, 4069)
		fmt.Fprintln(payload.Conns[i], payload.Payload[i])
		send = append(send, time.Now())
		payload.Conns[i].Read(recvd)
		status, _ := strconv.Atoi(string(recvd[9:12]))
		statuscode = append(statuscode, status)

		g++
	}

	l.Name = name
	l.Delay = float64(delay)
	l.Success = false

	for g, status := range statuscode {
		requests = append(requests, types.RequestLog{
			Timestamp:  send[g].UnixNano(),
			Statuscode: status,
		})

		if status == 200 {
			l.Success = true
		}
	}

	sent := types.Sent{
		Content: requests,
		Email:   email,
		Ip:      "Not Available Atm",
	}

	l.Sends = append(l.Sends, &sent)
	l.Requests = float64(len(statuscode))

	return l
}

func StartSnipe(task types.Task) {
	accounts := task.Accounts
	droptime := task.Timestamp

	// chans := make([]chan types.Logs, len(accounts))
	var logs []types.Log
	var wg sync.WaitGroup

	bearers = bearers.AddAccounts(accounts)

	PreSleep(droptime)

	payload := bearers.CreatePayloads(task.Name)

	Sleep(droptime, delay)

	for i, _ := range payload.AccountType {
		wg.Add(1)
		go func(i int) {
			tmp := StartSniper(droptime, delay, task.Name, i, payload, accounts[i].Email)
			logs = append(logs, tmp)
			// chans = append(chans, tmp)
			wg.Done()
		}(i)
	}

	wg.Wait()

	bearers = bearers.RemoveAccounts()

	var p types.Packet

	p.Type = "send_logs"
	p.Content.Logs = logs

	handleMessage(p)

}
