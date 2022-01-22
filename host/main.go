package host

import (
	types "Alien/types"
	"log"
	"time"
)

var state types.State
var RequestMap map[string]int

func Start() {
	state.LoadState()
	RequestMap = map[string]int{
		"mojang":    state.Config.Requests.Mojang,
		"giftcard":  state.Config.Requests.Giftcard,
		"microsoft": state.Config.Requests.Microsoft,
	}
	go AuthThread()
	go TaskThread()
	StartAPI("localhost:8080")
}

// Check if any account has to be authenticated
func AuthThread() {
	for {
		time.Sleep(time.Second * 10)
		for _, acc := range state.Accounts {
			if acc.AuthInterval > 0 {
				if time.Now().Unix() > acc.LastAuthed+acc.AuthInterval {
					log.Println("[Auth]", acc.Email, "is due for auth")
					// go func(tmpacc *types.StoredAccount) {

					// by default, the account isnt usable
					acc.Usable = false
					// authenticated account
					acc.Bearer, acc.Type = Auth(acc.Email, acc.Password, acc.Type, types.Packet{})
					log.Println("[Auth]", acc.Email, "is usable:", acc.Usable)

					// if the account is usable, update the last authed time
					if acc.Bearer != "" {
						acc.LastAuthed = time.Now().Unix()
						acc.Usable = true
						state.SaveState()
						continue // continue to next account
					}

					// if the account isnt usable, remove it from the list
					var ts []types.StoredAccount
					for _, i := range state.Accounts {
						if i.Email != acc.Email {
							ts = append(ts, i)
						}
					}
					state.Accounts = ts
					state.SaveState()

					// if tmpacc.Bearer != "" {
					// 	tmpacc.Usable = true
					// } else {
					// 	var tmpaccs []types.StoredAccount
					// 	for _, i := range state.Accounts {
					// 		if i.Email != tmpacc.Email {
					// 			tmpaccs = append(tmpaccs, i)
					// 		}
					// 		state.Accounts = tmpaccs
					// 		state.Accounts = append(state.Accounts, *tmpacc)
					// 		state.SaveState()
					// 		return
					// 	}
					// }
					// tmpacc.LastAuthed = time.Now().Unix()
					// // log.Println(tmpacc, tmpacc.Bearer, tmpacc.LastAuthed)

					// var tmpaccs []types.StoredAccount
					// for _, i := range state.Accounts {
					// 	if i.Email != tmpacc.Email {
					// 		tmpaccs = append(tmpaccs, i)
					// 	}
					// }
					// state.Accounts = tmpaccs
					// state.Accounts = append(state.Accounts, *tmpacc)

					// state.SaveState()
					// }(&acc)
				}
			}
		}
	}

}

// Check if any tasks are due in the next 60 secs
func TaskThread() {
	for {
		time.Sleep(time.Second * 10)
		for _, task := range state.Tasks {
			// if less than minute is left
			if task.Timestamp-time.Now().Unix() < 60 {
				log.Println("Task", task.Type, "is due for execution. Name:", task.Name)
				go func(task *types.QueuedTask) {
					// TODO
					// get account that should be used
					// get accounts
					// assign each VPS a account
					// execute task
					// sending to all VPSs

					// remove task from queue
					var ts []types.QueuedTask
					for _, i := range state.Tasks {
						if i.Name != task.Name {
							ts = append(ts, i)
						}
					}
					state.Tasks = ts
					state.SaveState()
				}(&task)
			}
		}
	}
}
