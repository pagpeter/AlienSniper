package host

import (
	utils "Alien/shared"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var tmp utils.Packet

func home(w http.ResponseWriter, r *http.Request) {
	// Handle the home page
	fmt.Fprintf(w, "Status: online (duh)")
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Handle incomming websocket connections

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	go ConnectionHandler(c)
}

func StartAPI(addr string) {
	// Start the API

	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func CheckInitialAuth(p utils.Packet) bool {
	// TODO: Check if the auth packet is valid
	return true
}

func HandlePacket(p utils.Packet) utils.Packet {
	// Handle a packet
	// return the response packet
	return utils.Packet{Content: utils.Content{}, Type: "response"}
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

	var p utils.Packet
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

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("message read:", err, c.RemoteAddr().String())
			errp := tmp.MakeError("Error reading message")
			c.WriteMessage(websocket.TextMessage, errp.Encode())
			break
		}
		// log.Printf("recv: %s", message)
		var p utils.Packet
		err = p.Decode(message)
		if err != nil {
			log.Println("decode:", err, c.RemoteAddr().String())
			errp := tmp.MakeError("Error decoding message")
			c.WriteMessage(websocket.TextMessage, errp.Encode())
			break
		}
		res := HandlePacket(p)
		err = c.WriteMessage(websocket.TextMessage, res.Encode())
	}
}
