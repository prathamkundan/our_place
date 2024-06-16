package internal

import (
	"log"
)


type Hub struct {
	// Communication channels only these are used by the clients to ensure safe access.

	// Broadcast channel messages sent to this are sent to all subscribers. This needs to be buffered
	Broadcast chan SMessage

	// Channel to register subscribers
	Register chan Subscriber[SMessage]

	// Channel to deregister subscribers
	Deregister chan Subscriber[SMessage]

	// Map to keep track of registered Subscribers. Do not access on references.
	Subscribers map[Subscriber[SMessage]]bool
}

func SetupHub() *Hub {
	hub := &Hub{
		Broadcast:   make(chan SMessage, 128),
		Register:    make(chan Subscriber[SMessage]),
		Deregister:  make(chan Subscriber[SMessage]),
		Subscribers: make(map[Subscriber[SMessage]]bool),
	}
	go hub.HandleMessage()

	return hub
}

func (h *Hub) HandleMessage() {
	for {
		select {
		case msg := <-h.Broadcast:
			log.Println("Broadcasting: ", msg)
			for key, val := range h.Subscribers {
				log.Println("Notifying:", key)
				if val {
					key.Notify(msg)
				}
			}
		case sub := <-h.Register:
			log.Println("Registering: ", sub)
			h.Subscribers[sub] = true
		case sub := <-h.Deregister:
			log.Println("Deregistering: ", sub)
			delete(h.Subscribers, sub)
		}
	}
}
