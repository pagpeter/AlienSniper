package node

import (
	types "Alien/types"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

var err error

// https://github.com/gorilla/websocket/blob/master/examples/echo/client.go
func MakeConnection(addr string) *websocket.Conn {
	// connect to the host
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}

	c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
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
	case "send_logs":
		res = send_logs(p.Content.Logs)
	default:
		res = res.MakeError("Cant handle packet")
	}
	return res
}

func send_logs(Logs []types.Log) types.Packet {
	res := types.Packet{
		Type: "save_logs",
		Content: types.Content{
			Logs: Logs,
			Response: &types.Response{
				Message: "node",
			},
		},
	}

	c.WriteMessage(websocket.TextMessage, res.Encode())

	return res
}

func ListenToEvents() {
	tmp := types.Packet{}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
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
