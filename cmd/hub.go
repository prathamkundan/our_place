package main

import "space/primitive"

type Hub struct {
	// Communication channels only these are used by the clients to ensure safe access.

	// Broadcast channel messages sent to this are sent to all subscribers. This needs to be buffered
	Broadcast chan SMessage

	// Channel to register subscribers
	Register chan primitive.Subscriber[SMessage]

	// Channel to deregister subscribers
	Deregister chan primitive.Subscriber[SMessage]

	// Map to keep track of registered subscribers. Do not access on references.
	subscribers map[primitive.Subscriber[SMessage]]bool
}

func (h Hub) HandleMessage() {
	select {
	case msg := <-h.Broadcast:
		for key := range h.subscribers {
			key.Notify(msg)
		}
	case sub := <-h.Register:
		h.subscribers[sub] = true
	case sub := <-h.Deregister:
		delete(h.subscribers, sub)
	}
}
