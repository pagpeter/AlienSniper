package node

import (
	types "Alien/types"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

// https://github.com/gorilla/websocket/blob/master/examples/echo/client.go
func MakeConnection(addr string) *websocket.Conn {
	// connect to the host
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	log.Println("Connected to host. Sending auth packet")

	// send auth msg
	authmsg := types.Packet{
		Type: "auth",
		Content: types.Content{
			Auth: config.Token,
			Response: &types.Response{
				Message: "node",
			},
		},
	}
	c.WriteMessage(websocket.TextMessage, authmsg.Encode())
	return c
}

func handleTask(p types.Packet) types.Packet {
	tmp := types.Packet{}
	res := types.Packet{}
	switch p.Content.Task.Type {
	case "snipe":
		StartSnipe(*p.Content.Task)
	default:
		res = tmp.MakeError("Unknown task type")
	}
	return res
}

func handleMessage(p types.Packet) types.Packet {
	res := types.Packet{}
	switch p.Type {
	case "task":
		res = handleTask(p)
	default:
		res = res.MakeError("Cant handle packet")
	}
	return res
}

func ListenToEvents(c *websocket.Conn) {
	tmp := types.Packet{}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			continue
		}
		var p types.Packet
		err = p.Decode(message)
		if err != nil {
			// log.Println("decode:", err, c.RemoteAddr().String())
			// log.Println("message:", string(message))
			errp := tmp.MakeError("Error decoding message")
			c.WriteMessage(websocket.TextMessage, errp.Encode())
			continue
		}

		m := handleMessage(p)
		c.WriteMessage(websocket.TextMessage, m.Encode())

		log.Printf("recv: %s", message)
	}
}
