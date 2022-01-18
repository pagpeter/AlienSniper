package shared

import (
	"strings"
)

type Packet struct {
	Type string
	Data string
}

func (p *Packet) Encode() []byte {
	// Encode the packet to a byte array
	return []byte(p.Type + "**" + p.Data)
}

func (p *Packet) Decode(data []byte) (err string) {
	// Decode the packet from a byte array
	pckt := string(data)
	pieces := strings.Split(pckt, "**")

	if len(pieces) != 2 {
		return // Invalid packet
	}

	p.Type = pieces[0]
	p.Data = pieces[1]
	return ""
}

// type Connection struct {
// 	Conn websocket.Conn // The websocket connection
// }
