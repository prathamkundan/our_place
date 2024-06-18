package main

import (
	"log"
	"net/http"
	. "space/internal"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var canvas = NewCanvas(512, 512)

var hub = SetupHub()

func handleConnection(w http.ResponseWriter, r *http.Request) {
	log.Println("Connection from :", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not upgrade connection: %v\n", err)
	}

	_, err = SetupClient(conn, hub, canvas)
	if err != nil {
		log.Printf("Could not set up client for: %s", conn.RemoteAddr())
	}
}

func main() {
	// Maybe replace with a Register function that can call the handler for the client
	hub.Register(canvas)

	http.HandleFunc("/ws", handleConnection)
	log.Println("Starting server at:", 8000)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
