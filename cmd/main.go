package main

import (
	"log"
	"net/http"
	"space/primitive"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var appState = NewCanvas(512, 512)

var hub = Hub{
	Broadcast:   make(chan SMessage, 128),
	Register:    make(chan primitive.Subscriber[SMessage]),
	Deregister:  make(chan primitive.Subscriber[SMessage]),
	subscribers: make(map[primitive.Subscriber[SMessage]]bool),
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Could not upgrade connection: %v\n", err)
	}

	c := Client{conn: conn, Send: make(chan (SMessage), 8)}

	// Register client
	hub.Register <- &c

	go c.HandleIncoming()
	go c.HandleSocketIncoming()
}

func main() {
	hub.Register <- appState
	http.HandleFunc("/ws", handleConnection)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
