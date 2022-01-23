package types

import (
	utils "Alien/shared"
	"encoding/json"
	"log"
)

type StoredAccount struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
	// Requests     int    `json:"requests,omitempty"`
	Bearer       string `json:"bearer,omitempty"`
	LastAuthed   int64  `json:"last_authed,omitempty"`
	AuthInterval int64  `json:"auth_interval,omitempty"`
	Group        string `json:"group,omitempty"`
	Usable       bool   `json:"usable,omitempty"`
	Security     string `json:"sq,omitempty"`
}

type QueuedTask struct {
	Type      string         `json:"type"`
	Name      string         `json:"name,omitempty"`
	Account   *StoredAccount `json:"account,omitempty"`
	Timestamp int64          `json:"timestamp,omitempty"`
	Searches  string         `json:"searches,omitempty"`
	Group     string         `json:"group,omitempty"`
}

type State struct {
	Config   Config          `json:"config,omitempty"`
	Accounts []StoredAccount `json:"accounts,omitempty"`
	Tasks    []QueuedTask    `json:"tasks,omitempty"`
}

func (s *State) ToJson() []byte {
	b, _ := json.MarshalIndent(s, "", "  ")
	return b
}

func (s *State) SaveState() {
	utils.WriteFile("host_state.json", string(s.ToJson()))
}

func (s *State) LoadState() {
	data, err := utils.ReadFile("host_state.json")
	if err != nil {
		log.Println("No state file found, creating new one.")
		s.Config = Config{}
		s.Config.LoadFromFile()
		s.SaveState()
		return
	}
	json.Unmarshal([]byte(data), s)

	s.Config = Config{}
	s.Config.LoadFromFile()
}
