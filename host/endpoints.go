package host

import (
	types "Alien/types"
	"bytes"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
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

	if p.Content.Account == (nil) || p.Content.Account.Email == "" || p.Content.Account.Password == "" {
		res.Type = "error"
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

	if len(p.Content.Account.Lines) == 0 {
		res.Type = "error"
		res.Content.Response.Error = "Invalid Format"
		return res
	}

	for _, line := range p.Content.Account.Lines {
		data := strings.Split(line, ":")

		if len(data) != 2 || data[0] == "" || data[1] == "" {
			res.Type = "error"
			res.Content.Response.Error = "Make sure your acccount(s) have a email and password"
			return res
		} else {
			acc := types.StoredAccount{
				Email:        data[0],
				Password:     data[1],
				Type:         p.Content.Account.Type,
				AuthInterval: 86400,
			}
			state.Accounts = append(state.Accounts, acc)
		}
	}

	state.SaveState()
	res.Content.Response.Message = "account(s) added successfully"

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

	if p.Content.Task == (nil) {
		res.Type = "error"
		res.Content.Response.Error = "No task provided"
		return res
	}

	if p.Content.Task.Type == "" {
		res.Type = "error"
		res.Content.Response.Error = "No task type provided"
		return res
	}

	if p.Content.Task.Type == "snipe" {
		name := p.Content.Task.Name
		group := p.Content.Task.Group
		if name == "" {
			res.Type = "error"
			res.Content.Response.Error = "No name provided"
			return res
		}

		for _, t := range state.Tasks {
			if t.Name == name {
				res.Type = "error"
				res.Content.Response.Error = "Task with that name already exists"
				return res
			}
		}

		drop, err := getDroptime(name, "droptime.site")
		if err != nil {
			res.Type = "error"
			res.Content.Response.Error = err.Error()
			return res
		}

		searches, err := droptimeSiteSearches(name)
		if err != nil {
			res.Type = "error"
			res.Content.Response.Error = err.Error()
		}

		res.Content.Task.Searches = searches
		res.Content.Task.Timestamp = drop.Unix()
		res.Content.Task.Name = name
		res.Content.Task.Type = "snipe"
		res.Content.Task.Group = group

		t := types.QueuedTask{
			Type:      "snipe",
			Name:      name,
			Timestamp: drop.Unix(),
			Group:     group,
			Searches:  searches,
		}

		state.Tasks = append(state.Tasks, t)
		state.SaveState()
	}

	return res
}

func add_session(p types.Packet) types.Packet {
	res := types.Packet{}
	res.Type = "add_session_response"
	res.Content.Response = &types.Response{}

	if len(p.Content.Vps) == 0 {
		res.Type = "error"
		res.Content.Response.Error = "No details provided."
		return res
	}

	for _, details := range p.Content.Vps {
		if strings.ToLower(details.Group) != "giftcard" && strings.ToLower(details.Group) != "microsoft" {
			res.Type = "error"
			res.Content.Response.Error = "Vps must have `Giftcard` or `Microsoft` selected."
			return res
		} else {
			if AddVps(details.Ip, details.Port, details.Password, details.Host) {
				state.Vps = append(state.Vps, details)
			} else {
				res.Type = "error"
				res.Content.Response.Error = "Failed to connect to the VPS"
				return res
			}
		}
	}

	res.Content.Vps = p.Content.Vps
	state.SaveState()

	return res
}

func save_logs(p types.Packet) types.Packet {
	res := types.Packet{}
	res.Type = "save_logs_response"
	res.Content.Response = &types.Response{Message: "Saved accounts"}
	var isin bool = false

	// for every log in the DB
	for i, l := range state.Logs {

		// for every new log
		for _, nl := range p.Content.Logs {

			// if the new log is already in the DB
			if l.Name == nl.Name {
				isin = true
				state.Logs[i].Sends = append(state.Logs[i].Sends, nl.Sends...)
				state.Logs[i].Requests += nl.Requests
				if nl.Success {
					state.Logs[i].Success = true
				}
			}
		}
	}

	if !isin {
		state.Logs = append(state.Logs, p.Content.Logs...)
	}

	state.SaveState()

	return res
}

func start_vps(ip, port, password, user string) bool {
	conn, err := ssh.Dial("tcp", ip+":"+port, &ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	})

	if err != nil {
		log.Println("Error: " + err.Error())
	} else {
		sesh, _ := conn.NewSession()
		defer sesh.Close()

		var stdoutBuf bytes.Buffer
		sesh.Stdout = &stdoutBuf
		err := sesh.Run(fmt.Sprintf("nohup ./Alien node -ip %v:%v", state.Config.Host, state.Config.Port))
		if err == nil {
			return true
		}
	}

	return false
}

func AddVps(ip, port, password, user string) bool {
	conn, err := ssh.Dial("tcp", ip+":"+port, &ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
	})
	if err != nil {
		log.Println("Error: " + err.Error())
	} else {
		sesh, _ := conn.NewSession()
		defer sesh.Close()
		var stdoutBuf bytes.Buffer
		sesh.Stdout = &stdoutBuf
		err := sesh.Run(fmt.Sprintf("wget linkhere\nchmod +x ./Alien\nnohup ./Alien node -ip %v:%v", state.Config.Host, state.Config.Port))
		if err == nil {
			return true
		}
	}

	return false
}
