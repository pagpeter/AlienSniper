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

func ListenToEvents(c *websocket.Conn) {
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)
	}
}