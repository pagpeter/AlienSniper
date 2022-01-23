package host

import (
	types "Alien/types"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var lastAuthedGlobal time.Time
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

	if state.Config.Host != "localhost" && state.Config.Host != "127.0.0.1" && state.Config.Host != "0.0.0.0" {
		log.Println("host can only be localhost or 0.0.0.0. hosting on 0.0.0.0")
		log.Println("You can change the host and port in the config.")
		state.Config.Host = "0.0.0.0"
	}

	addr := fmt.Sprintf("%s:%d", state.Config.Host, state.Config.Port)
	StartAPI(addr)
}

// Check if any account has to be authenticated
func AuthThread() {
	for {
		time.Sleep(time.Second * 10)
		// check if the last auth was more than a minute ago
		for i, acc := range state.Accounts {
			if time.Now().Unix() > acc.LastAuthed+acc.AuthInterval {
				log.Println("[Auth]", acc.Email, "is due for auth")

				// by default, the account isnt usable
				acc.Usable = false

				// authenticating account
				acc.Bearer, acc.Type = Auth(acc.Email, acc.Password, acc.Type, types.Packet{})
				log.Println("[Auth]", acc.Email, "is usable:", acc.Usable)
				lastAuthedGlobal = time.Now()

				// if the account is usable, update the last authed time
				if acc.Bearer != "" {
					acc.LastAuthed = time.Now().Unix()
					acc.Usable = true
					state.Accounts[i] = acc
					state.SaveState()
					break // break the loop to update the state.Accounts info.
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

				break // break the loop to update the state.Accounts info.
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
				// go func(task *types.QueuedTask) {
				// TODO

				// get account that should be used
				// get accounts

				// assign each VPS a account

				// sending to all VPSs
				for _, vps := range connectedNodes {
					log.Println("Sending to VPS")
					p := types.Packet{}
					p.Type = "task"
					p.Content.Task = &types.Task{
						Type:      task.Type,
						Name:      task.Name,
						Timestamp: task.Timestamp,
					}
					vps.WriteMessage(websocket.TextMessage, p.Encode())
				}
				// remove task from queue
				var ts []types.QueuedTask
				for _, i := range state.Tasks {
					if i.Name != task.Name {
						ts = append(ts, i)
					}
				}
				state.Tasks = ts
				state.SaveState()
				// }(&task)
			}
		}
	}
}
