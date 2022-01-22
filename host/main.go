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
	StartAPI("localhost:8080")
}

func AuthThread() {
	time.Sleep(time.Second * 10)
	for _, acc := range state.Accounts {
		if acc.AuthInterval > 0 {
			if time.Now().Unix() > acc.LastAuthed+acc.AuthInterval {
				log.Println("[Auth]", acc.Username, "is due for auth")
				go func(acc *types.StoredAccount) {
					acc.Usable = false
					acc.Bearer, acc.Type = Auth(acc.Email, acc.Password, acc.Type, types.Packet{})
					if acc.Bearer != "" {
						acc.Usable = true
					}
					acc.LastAuthed = time.Now().Unix()
					state.SaveState()
				}(&acc)
			}
		}
	}
}
