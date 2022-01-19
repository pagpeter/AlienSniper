package host

import (
	utils "Alien/shared"
	types "Alien/types"
)

func HandlePacket(p types.Packet) types.Packet {
	// Handle a packet
	// return the response packet

	// auth is handled somwehere else
	res := types.Packet{}

	switch p.Type {
	case "auth":
		res = tmp.MakeError("Already authed")
	case "config":
		res.Type = "config_response"
		res.Content.Config = &types.Config{}
		res.Content.Config.LoadFromFile()
	case "add_account":
		res.Type = "add_account_response"
		res.Content.Response = &types.Response{}
		if p.Content.Account == (nil) {
			res.Content.Response.Error = "No account provided"
			break
		}

		if !utils.IsInMap(p.Content.Account.Type, RequestMap) {
			res.Content.Response.Error = "Invalid account type"
			break
		}

		acc := types.StoredAccount{
			Email:    p.Content.Account.Email,
			Password: p.Content.Account.Password,
			Type:     p.Content.Account.Type,
			Requests: RequestMap[p.Content.Account.Type],
		}
		state.Accounts = append(state.Accounts, acc)
		res.Content.Response.Message = "Account added"
		state.SaveState()
	case "remove_account":
		res.Type = "remove_account_response"
		res.Content.Response = &types.Response{}
		res.Content.Response.Message = "Account removed"
	case "get_state":
		res.Type = "state_response"
		res.Content.Response = &types.Response{}
		res.Content.State = &state
	default:
		res.Type = "error_response"
		res.Content.Error = "Unknown packet type: " + p.Type
	}

	return res
}
