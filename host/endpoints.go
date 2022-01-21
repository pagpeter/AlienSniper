package host

import (
	types "Alien/types"
	"strings"
)

// all endpoints for the websocket connection
// functions should have the following arguments:
//
// func f(p types.Packet) types.Packet {
// 	res := types.Packet{}
// 	res.Type = "xxx_response"
//
// }

func add_account_endpoint(p types.Packet) types.Packet {
	res := types.Packet{}
	res.Type = "add_account_response"
	res.Content.Response = &types.Response{}
	if p.Content.Account == (nil) {
		res.Content.Response.Error = "No account provided"
		return res
	}

	// if !utils.IsInMap(p.Content.Account.Type, RequestMap) {
	// 	res.Content.Response.Error = "Invalid account type"
	// 	return res
	// }

	if p.Content.Account == (nil) {
		res.Content.Response.Error = "No account provided"
		return res
	}

	// dont waste time
	go func() {
		acc := types.StoredAccount{
			Email:        p.Content.Account.Email,
			Password:     p.Content.Account.Password,
			Type:         p.Content.Account.Type,
			AuthInterval: 86400,
		}
		state.Accounts = append(state.Accounts, acc)
		// state.Accounts = append(state.Accounts, auth.Auth(p.Content.Account.Email, p.Content.Account.Password, p.Content.Account.Security, p))
		state.SaveState()
	}()
	res.Content.Response.Error = ""
	res.Content.Response.Message = "Account added successfully"
	return res
}

func add_multiple_accounts_endpoint(p types.Packet) types.Packet {
	res := types.Packet{}
	res.Type = "add_multi_response"
	res.Content.Response = &types.Response{}

	if p.Content.Account.Lines == (nil) {
		res.Content.Response.Error = "No accounts provided"
		return res
	}

	// if !utils.IsInMap(p.Content.Account.Type, RequestMap) {
	// 	res.Content.Response.Error = "Invalid account type"
	// 	return res
	// }
	c := 0
	for _, line := range p.Content.Account.Lines {
		data := strings.Split(line, ":")
		if len(data) != 2 {
			res.Content.Response.Error = "Invalid format"
			return res
		}
		c++
		acc := types.StoredAccount{
			Email:        data[0],
			Password:     data[1],
			Type:         p.Content.Account.Type,
			AuthInterval: 86400,
		}
		state.Accounts = append(state.Accounts, acc)
	}
	res.Content.Response.Error = ""
	res.Content.Response.Message = string(c) + " account(s) added successfully"
	return res

}

func remove_account_endpoint(p types.Packet) types.Packet {
	res := types.Packet{}
	var accs []types.StoredAccount
	res.Type = "remove_account_response"
	res.Content.Response = &types.Response{}

	for _, sa := range state.Accounts {
		if sa.Email != p.Content.Remove.Email {
			accs = append(accs, sa)
		}
	}

	state.Accounts = accs
	go func() {
		state.SaveState()
	}()

	res.Content.Response.Error = ""
	res.Content.Response.Message = "Account removed successfully"
	return res
}
