package host

import types "Alien/types"

var state types.State
var RequestMap map[string]int

func init() {
	state.LoadState()
	RequestMap = map[string]int{
		"mojang":    state.Config.Requests.Mojang,
		"giftcard":  state.Config.Requests.Giftcard,
		"microsoft": state.Config.Requests.Microsoft,
	}
}

func Start() {
	state.LoadState()
	StartAPI("localhost:8080")
}
