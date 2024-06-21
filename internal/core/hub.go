package core

import (
	"log"
	. "space/internal/core/pubsub"
)

type Hub struct {
	// Communication channels only these are used by the clients to ensure safe access.

	// broadcast channel messages sent to this are sent to all subscribers. This needs to be buffered
	broadcast chan Message

	// Channel to register subscribers
	register chan Subscriber[Message]

	// Channel to deregister subscribers
	deregister chan Subscriber[Message]

	// Map to keep track of registered subscribers. Do not access on references.
	subscribers map[Subscriber[Message]]bool
}

func SetupHub() *Hub {
	hub := &Hub{
		broadcast:   make(chan Message, 128),
		register:    make(chan Subscriber[Message]),
		deregister:  make(chan Subscriber[Message]),
		subscribers: make(map[Subscriber[Message]]bool),
	}
	go hub.HandleMessage()

	return hub
}

func (h *Hub) Register(s Subscriber[Message]) {
	h.register <- s
	go s.Listen()
}

func (h *Hub) Deregister(s Subscriber[Message]) {
	h.deregister <- s
	s.Interrupt()
}

func (h *Hub) Broadcast(msg Message) {
	h.broadcast <- msg
}

func (h *Hub) HandleMessage() {
	for {
		select {
		case msg := <-h.broadcast:
			log.Println("Broadcasting: ", msg)
			for key, val := range h.subscribers {
				log.Println("Notifying:", key)
				if val {
					key.Notify(msg)
				}
			}
		case sub := <-h.register:
			log.Println("Registering: ", sub)
			h.subscribers[sub] = true
		case sub := <-h.deregister:
			log.Println("Deregistering: ", sub)
			delete(h.subscribers, sub)
		}
	}
}
