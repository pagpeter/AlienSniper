package node

import (
	types "Alien/types"
	"log"

	"github.com/gorilla/websocket"
)

var config types.Config
var c *websocket.Conn

func Start(ip, key string) {
	log.Println("Trying to connect to host at", ip)
	c = MakeConnection(ip, key)
	ListenToEvents()
}
