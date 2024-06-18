package pubsub

// Publisher (Quite useless currently)
/* Publisher that notifies subscribers */
type Publisher[T any] interface {
	// Adds the client to the subscriber set (and calls the listen method on it? would decouple CLient and Hub allowing multiple topics)
	Register(s Subscriber[T])
	Deregister(s Subscriber[T])
    Broadcast(msg T)
}
