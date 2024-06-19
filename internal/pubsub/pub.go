package pubsub

// Publisher
/* Publisher that notifies subscribers */
type Publisher[T any] interface {
	// Adds the client to the subscriber set
	Register(s Subscriber[T])
	Deregister(s Subscriber[T])
	Broadcast(msg T)
}
