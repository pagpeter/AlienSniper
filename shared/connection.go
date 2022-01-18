package shared

import "encoding/json"

type Content struct {
	Error  string  `json:"error,omitempty"`
	Auth   string  `json:"auth,omitempty"`
	Config *Config `json:"config,omitempty"`
}

type Packet struct {
	Type    string  `json:"type"`
	Content Content `json:"content"`
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
