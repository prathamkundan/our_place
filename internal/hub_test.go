package internal

import (
	"testing"
	"time"
)

type DummySub struct {
	invoked int
	h       *Hub
}

func (d *DummySub) Notify(msg SMessage) {
	d.invoked++
}

func (d *DummySub) publish(n int) {
	for i := 0; i < n; i++ {
		d.h.Broadcast <- SMessage{pos: uint32(i), color: WHITE}
	}
	d.h.Deregister <- d
}

func TestHub(t *testing.T) {
	h := &Hub{
		Broadcast:   make(chan SMessage, 128),
		Register:    make(chan Subscriber[SMessage]),
		Deregister:  make(chan Subscriber[SMessage]),
		Subscribers: make(map[Subscriber[SMessage]]bool),
	}
	go h.HandleMessage()

	d1, d2 := DummySub{0, h}, DummySub{0, h}
	n1, n2 := 10, 30

	h.Register <- &d1
	h.Register <- &d2

	// Wait a bit to make hub process the registration
	time.Sleep(1 * time.Millisecond)
	if !h.Subscribers[&d1] {
		t.Fatalf("Could not register the subscriber: %v", d1)
	}
	if !h.Subscribers[&d2] {
		t.Fatalf("Could not register the subscriber: %v", d2)
	}

	d1.publish(n1)
	d2.publish(n2)

	if d1.invoked == n1 {
		t.Fatalf("Invocations expected at least: %d, got: %d", n1*2, d1.invoked)
	}
	if d2.invoked == n2 {
		t.Fatalf("Invocations expected at least: %d, got: %d", n2*2, d2.invoked)
	}

	time.Sleep(1 * time.Millisecond)
	if h.Subscribers[&d1] {
		t.Fatalf("Could not deregister the subscriber: %v", d1)
	}
	if h.Subscribers[&d2] {
		t.Fatalf("Could not deregister the subscriber: %v", d2)
	}
}
