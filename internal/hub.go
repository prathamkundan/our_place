package internal

import (
	"log"
	. "space/internal/pubsub"
)

type Hub struct {
	// Communication channels only these are used by the clients to ensure safe access.

	// broadcast channel messages sent to this are sent to all subscribers. This needs to be buffered
	broadcast chan SMessage

	// Channel to register subscribers
	register chan Subscriber[SMessage]

	// Channel to deregister subscribers
	deregister chan Subscriber[SMessage]

	// Map to keep track of registered subscribers. Do not access on references.
	subscribers map[Subscriber[SMessage]]bool
}

func SetupHub() *Hub {
	hub := &Hub{
		broadcast:   make(chan SMessage, 128),
		register:    make(chan Subscriber[SMessage]),
		deregister:  make(chan Subscriber[SMessage]),
		subscribers: make(map[Subscriber[SMessage]]bool),
	}
	go hub.HandleMessage()

	return hub
}

func (h *Hub) Register(s Subscriber[SMessage]) {
	h.register <- s
	go s.Listen()
}

func (h *Hub) Deregister(s Subscriber[SMessage]) {
	h.deregister <- s
	s.Interrupt()
}

func (h *Hub) Broadcast(msg SMessage) {
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
