package host

import (
	// "Alien/types"

	types "Alien/types"
)

func HandlePacket(p types.Packet) types.Packet {
	// Handle a packet
	// return the response packet

	// auth is handled somwehere else
	res := types.Packet{}

	switch p.Type {
	// Authenticate the client. Returns an error, because it should be the first packet (handled in API.go)
	case "auth":
		res = tmp.MakeError("Already authed")

	// return the config
	case "config":
		res.Type = "config_response"
		res.Content.Config = &types.Config{}
		res.Content.Config.LoadFromFile()

	// Save snipe logs.
	case "save_logs":
		res.Type = "logs_response"
		save_logs(p)

	// Add a single account
	case "add_account":
		return add_account_endpoint(p)

	// Add multiple accounts
	case "add_multiple_accounts":
		return add_multiple_accounts_endpoint(p)

	// Remove an account by email
	case "remove_account":
		return remove_account_endpoint(p)

	// Get full state
	case "get_state":
		res.Type = "state_response"
		res.Content.Response = &types.Response{}
		res.Content.State = &state

	case "add_task":
		return add_task_endpoint(p)

	case "get_tasks":
		res.Type = "tasks_response"
		res.Content.Response = &types.Response{}
		res.Content.State.Tasks = state.Tasks
		return res

	case "get_accounts":
		res.Type = "accounts_response"
		res.Content.Response = &types.Response{}
		res.Content.State.Accounts = state.Accounts
		return res

	case "get_logs":
		res.Type = "logs_response"
		res.Content.Response = &types.Response{}
		res.Content.State.Logs = state.Logs
		return res

	case "get_config":
		res.Type = "config_response"
		res.Content.Response = &types.Response{}
		res.Content.State.Config = state.Config
		return res

	// Invalid packet
	default:
		res.Type = "error_response"
		res.Content.Error = "Unknown packet type: " + p.Type
	}

	return res
}
