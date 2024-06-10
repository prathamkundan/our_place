package primitive

// Subscriber
/* Represents a Subscriber that can be notified. */
type Subscriber[T any] interface {
    Notify(msg T)
}
