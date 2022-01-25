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

func handleTask(p types.Packet) {
	switch p.Content.Task.Type {
	case "snipe":
		StartSnipe(*p.Content.Task)
	}
}

func handleMessage(p types.Packet) {
	switch p.Type {
	case "task":

		log.Printf("Starting task %v\n", p.Content.Task.Type)

		switch p.Content.Task.Type {
		case "snipe":
			StartSnipe(*p.Content.Task)
		}

	case "send_logs":
		log.Println("Sending logs to the host.")
		send_logs(p.Content.Logs)
	}
}

func send_logs(Logs []types.Log) {
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
			errp := tmp.MakeError("Error decoding message")
			c.WriteMessage(websocket.TextMessage, errp.Encode())
			continue
		}

		handleMessage(p)

	}
}
