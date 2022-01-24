package host

import (
	types "Alien/types"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var connectedNodes []*websocket.Conn
var connectedDashboards []*websocket.Conn
var tmp types.Packet

func home(w http.ResponseWriter, r *http.Request) {
	// Handle the home page

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		json.NewEncoder(w).Encode("OK")
		return
	}

	fmt.Fprintf(w, "Status: online (duh)")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Handle incomming websocket connections

	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		json.NewEncoder(w).Encode("OK")
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	go ConnectionHandler(c)
}

// https://gist.github.com/denji/12b3a568f092ab951456
func StartAPI(addr string) {
	// Start the API

	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", home)
	if state.Config.TLS.Active {
		log.Fatal(http.ListenAndServeTLS(addr, state.Config.TLS.Cert, state.Config.TLS.Key, nil))
	} else {
		log.Fatal(http.ListenAndServe(addr, nil))
	}
}

func CheckInitialAuth(p types.Packet) bool {
	// Check the auth packet
	return p.Content.Auth == state.Config.Token
}

func RemoveConnection(c *websocket.Conn) {
	// Remove a connection from the list of connected nodes
	// This is called when a connection is closed

	for i, v := range connectedNodes {
		if v == c {
			connectedNodes = append(connectedNodes[:i], connectedNodes[i+1:]...)
			log.Println("Connected nodes:", len(connectedNodes))
			return
		}
	}

	for i, v := range connectedDashboards {
		if v == c {
			connectedDashboards = append(connectedDashboards[:i], connectedDashboards[i+1:]...)
			log.Println("Connected dashboards:", len(connectedDashboards))
			return
		}
	}
}

func ConnectionHandler(c *websocket.Conn) {
	// Handle a websocket connection
	// First, read the auth packet
	// This returns the config packet
	// Then, handle the normal packets

	defer c.Close()

	// returns msg type, msg, error
	_, authMessage, err := c.ReadMessage()
	if err != nil {
		log.Println("Initial auth read:", err, string(authMessage), c.RemoteAddr().String())
		c.Close()
		return
	}

	var p types.Packet
	err = p.Decode(authMessage)
	if err != nil {
		log.Println("Initial auth decode:", err, c.RemoteAddr().String())
		errp := tmp.MakeError("First packet must be of type auth")
		c.WriteMessage(websocket.TextMessage, errp.Encode())
		c.Close()
		return
	}

	if p.Type != "auth" {
		log.Println("First packet isnt auth", c.RemoteAddr().String())
		errp := tmp.MakeError("First packet must be of type auth")

		c.WriteMessage(websocket.TextMessage, errp.Encode())
		c.Close()
		return
	}

	if !CheckInitialAuth(p) {
		log.Println("auth:", err)
		errp := tmp.MakeError("Invalid auth packet")
		c.WriteMessage(websocket.TextMessage, errp.Encode())
		c.Close()
		return
	}

	// this is needed or it will crash
	var clientType string
	if p.Content.Response == nil {
		clientType = ""
	} else {
		clientType = p.Content.Response.Message
	}

	log.Println("New client connected. Client type:", clientType, " - IP:", c.RemoteAddr().String())
	if clientType == "node" {
		connectedNodes = append(connectedNodes, c)
		log.Println("Connected nodes:", len(connectedNodes))
	} else if clientType == "web" {
		connectedDashboards = append(connectedDashboards, c)
		log.Println("Connected dashboards:", len(connectedDashboards))
	} else {
		log.Println("Unknown client type:", clientType, " - IP:", c.RemoteAddr().String())
	}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			RemoveConnection(c)
			c.Close()
			break
		}

		var p types.Packet
		err = p.Decode(message)
		if err != nil {
			// log.Println("decode:", err, c.RemoteAddr().String())
			// log.Println("message:", string(message))
			errp := tmp.MakeError("Error decoding message")
			c.WriteMessage(websocket.TextMessage, errp.Encode())
			RemoveConnection(c)
			return
		}

		res := HandlePacket(p)
		err = c.WriteMessage(websocket.TextMessage, res.Encode())
		if err != nil {
			log.Println("write:", err, c.RemoteAddr().String())
			return
		}

	}
}
