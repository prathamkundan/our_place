package main

import (
	"log"
	"net/http"
	. "space/internal"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var canvas = NewCanvas(512, 512)

var hub = InitHub()

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not upgrade connection: %v\n", err)
	}

	c := Client{
		Send:   make(chan (SMessage), 8),
		Conn:   conn,
		Hub:    hub,
		Canvas: &canvas,
	}

	// Register client
	hub.Register <- &c

	go c.HandleIncoming()
	go c.HandleSocketIncoming()
}

func main() {
	hub.Register <- canvas
	http.HandleFunc("/ws", handleConnection)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
