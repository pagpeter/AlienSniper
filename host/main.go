package host

import types "Alien/types"

var state types.State
var RequestMap map[string]int

func Start() {
	state.LoadState()
	RequestMap = map[string]int{
		"mojang":    state.Config.Requests.Mojang,
		"giftcard":  state.Config.Requests.Giftcard,
		"microsoft": state.Config.Requests.Microsoft,
	}
	StartAPI("localhost:8080")
}
