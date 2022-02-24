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

	// checking for TLS
	prefix := "ws://"
	if state.Config.TLS.Active {
		if state.Config.TLS.Cert == "" || state.Config.TLS.Key == "" {
			log.Fatalln("TLS is active but no cert or key is set.")
		}
		fmt.Println("TLS is active")
		prefix = "wss://"
	}
	addr := fmt.Sprintf("%s:%d", state.Config.Host, state.Config.Port)
	log.Println("Listening on", prefix+addr+"/ws")
	StartAPI(addr)

	if len(state.Vps) != 0 {
		for _, details := range state.Vps {
			start_vps(details.Ip, details.Port, details.Password, details.Host)
		}
	}
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

				if acc.Bearer != "" {
					acc.Usable = true
				}

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
			if len(connectedNodes) != 0 {
				// if less than minute is left
				if task.Timestamp-time.Now().Unix() < 60 {
					log.Println("Task", task.Type, "is due for execution. Name:", task.Name)

					var data [][]types.StoredAccount
					var incr int
					var use int

					for _, inp := range state.Accounts {
						if len(data) == 0 {
							data = append(data, []types.StoredAccount{
								inp,
							})
						} else {
							if incr == 4 {
								incr = 0
								use++
								data = append(data, []types.StoredAccount{})
							}

							data[use] = append(data[use], inp)
						}

						incr++
					}

					p := types.Packet{}
					p.Type = "task"

					log.Println("Sending to VPS(s)")

					var cancel bool = false
					for _, info := range data {
						if !cancel {
							p.Content.Task = &types.Task{
								Type:      task.Type,
								Name:      task.Name,
								Timestamp: task.Timestamp,
								Group:     task.Group,
								Accounts:  info,
							}
							if i != len(connectedNodes) {
								connectedNodes[i].WriteMessage(websocket.TextMessage, p.Encode())
							} else {
								cancel = true
							}
						}
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
				}
			}
		}
	}
}
