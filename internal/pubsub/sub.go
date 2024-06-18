package pubsub

// Subscriber
/* Represents a Subscriber that can be notified. */
type Subscriber[T any] interface {
    Listen()
	Notify(msg T)
    Interrupt()
}
