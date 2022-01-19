package types

import (
	"encoding/json"
)

type Content struct {
	Error string `json:"error,omitempty"`
	Auth  string `json:"auth,omitempty"`

	State    *State    `json:"state,omitempty"`
	Config   *Config   `json:"config,omitempty"`
	Response *Response `json:"response,omitempty"`
	Task     *Task     `json:"task,omitempty"`
	Account  *Account  `json:"account,omitempty"`
}

type Packet struct {
	Type    string  `json:"type"`
	Content Content `json:"content"`
}

type Response struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

type Task struct {
	Type string `json:"type"`
}

type Account struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

func (p *Packet) Encode() []byte {
	b, _ := json.Marshal(p)
	return b
}

func (p *Packet) Decode(data []byte) (err error) {
	return json.Unmarshal(data, p)
}

func (p *Packet) MakeError(err string) Packet {
	p.Type = "error"
	p.Content = Content{Error: err}
	return *p
}
