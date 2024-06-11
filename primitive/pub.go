package primitive

// Publisher (Quite useless currently)
/* Publisher that notifies subscribers */
type Publisher[T any] interface {
    Register(sub *T)
    Deregister(sub *T)
}
