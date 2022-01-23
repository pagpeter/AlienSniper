package node

import (
	types "Alien/types"
	"fmt"
	"log"
)

var config types.Config

func Start() {
	config.LoadFromFile()
	log.Println("Loaded config")
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	log.Println("Trying to connect to host at", addr)
	c := MakeConnection(addr)
	ListenToEvents(c)
}

