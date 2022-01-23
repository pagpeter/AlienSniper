package node

import (
	types "Alien/types"
	"log"
	"sync"
)

const delay = 0

func StartGCSnipe(timestamp int64, delay int, accounts []*types.Account, name string) types.Log {
	var l types.Log
	l.Name = name
	return l
}

func StartNormalSnipe(timestamp int64, delay int, accounts []*types.Account, name string) types.Log {
	var l types.Log
	l.Name = name
	return l
}

func StartSnipe(task types.Task) {
	accounts := task.Accounts
	droptime := task.Timestamp
	// chans := make([]chan types.Logs, len(accounts))
	var logs []types.Log
	var wg sync.WaitGroup

	for _, account := range accounts {
		if account.Type == "gc" || account.Type == "giftcard" {
			wg.Add(1)
			go func() {
				tmp := StartGCSnipe(droptime, delay, accounts, task.Name)
				logs = append(logs, tmp)
				// chans = append(chans, tmp)
				wg.Done()
			}()
		} else {
			wg.Add(1)
			go func() {
				tmp := StartNormalSnipe(droptime, delay, accounts, task.Name)
				logs = append(logs, tmp)
				// chans = append(chans, tmp)
				wg.Done()
			}()
		}
	}
	wg.Wait()


	log.Println("got logs", len(logs), logs)

}