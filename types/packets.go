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
	Remove   *Remove   `json:"remove,omitempty"`
	Logs     []Log     `json:"logs,omitempty"`
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
	Type      string          `json:"type"`
	Name      string          `json:"name"`
	Timestamp int64           `json:"timestamp,omitempty"`
	Searches  string          `json:"searches,omitempty"`
	Group     string          `json:"group,omitempty"`
	Accounts  []StoredAccount `json:"accounts,omitempty"`
}

type Log struct {
	Name     string  `json:"name"`
	Requests float64 `json:"requests"`
	Delay    float64 `json:"delay"`
	Success  bool    `json:"success"`
	Sends    []*Sent `json:"sends"`
}

type RequestLog struct {
	Timestamp  int64 `json:"timestamp"`
	Statuscode int   `json:"statuscode"`
}

type Sent struct {
	Content []RequestLog `json:"content"`
	Email   string       `json:"email"`
	Ip      string       `json:"ip"`
}

type Account struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Type     string   `json:"type"`
	Lines    []string `json:"lines,omitempty"`
	Group    string   `json:"group,omitempty"`
	Status   string   `json:"status,omitempty"`
	Security string   `json:"sq,omitempty"`
	Bearer   string   `json:"bearer,omitempty"`
	Usable   bool     `json:"usable,omitempty"`
}

type Remove struct {
	Email string `json:"email"`
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
