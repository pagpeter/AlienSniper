package host

import (
	types "Alien/types"
	"log"
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

func add_task_endpoint(p types.Packet) types.Packet {
	res := types.Packet{}
	res.Type = "add_task_response"
	res.Content.Response = &types.Response{}
	res.Content.Task = &types.Task{}

	log.Println(p)

	if p.Content.Task == (nil) {
		res.Content.Response.Error = "No task provided"
		return res
	}

	if p.Content.Task.Type == "" {
		res.Content.Response.Error = "No task type provided"
		return res
	}

	if p.Content.Task.Type == "snipe" {
		name := p.Content.Task.Name

		if name == "" {
			res.Content.Response.Error = "No name provided"
			return res
		}

		drop, err := getDroptime(name, "ckm")
		if err != nil {
			res.Content.Response.Error = err.Error()
			return res
		}
		res.Content.Task.Unix = drop.Unix()
		res.Content.Task.Name = name
		res.Content.Task.Type = "snipe"

		t := types.QueuedTask{
			Type:      "snipe",
			Name:      name,
			Timestamp: drop.Unix(),
		}
		state.Tasks = append(state.Tasks, t)
		go func() {
			state.SaveState()
		}()
	}

	res.Content.Response.Error = ""
	return res
}
