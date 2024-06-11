package internal

import "space/primitive"

type Hub struct {
	// Communication channels only these are used by the clients to ensure safe access.

	// Broadcast channel messages sent to this are sent to all subscribers. This needs to be buffered
	Broadcast chan SMessage

	// Channel to register subscribers
	Register chan primitive.Subscriber[SMessage]

	// Channel to deregister subscribers
	Deregister chan primitive.Subscriber[SMessage]

	// Map to keep track of registered Subscribers. Do not access on references.
	Subscribers map[primitive.Subscriber[SMessage]]bool
}

func InitHub() *Hub {
	return &Hub{
		Broadcast:   make(chan SMessage, 128),
		Register:    make(chan primitive.Subscriber[SMessage]),
		Deregister:  make(chan primitive.Subscriber[SMessage]),
		Subscribers: make(map[primitive.Subscriber[SMessage]]bool),
	}

}

func (h Hub) HandleMessage() {
	for {
		select {
		case msg := <-h.Broadcast:
			for key := range h.Subscribers {
				key.Notify(msg)
			}
		case sub := <-h.Register:
			h.Subscribers[sub] = true
		case sub := <-h.Deregister:
			delete(h.Subscribers, sub)
		}
	}
}
