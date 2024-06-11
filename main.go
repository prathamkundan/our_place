package main

import (
	"log"
	"net/http"
	. "space/internal"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var canvas = NewCanvas(512, 512)

var hub = SetupHub()

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not upgrade connection: %v\n", err)
	}

	err = SetupClient(conn, hub, canvas)
	if err != nil {
		log.Printf("Could not set up client for: %s", conn.RemoteAddr())
	}
}

func main() {
    // Maybe replace with a Register function that can call the handler for the client
	hub.Register <- canvas
    go canvas.HandleIncoming()

	http.HandleFunc("/ws", handleConnection)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
