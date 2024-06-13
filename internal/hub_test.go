package internal

import (
	"log"
	"sync"
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

func (d DummySub) publish(wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		d.h.Broadcast <- SMessage{pos: uint32(i), color: WHITE}
	}
}

func TestHub(t *testing.T) {
	log.Println("Running Hub Test")
	h := SetupHub()

	d1 := DummySub{0, h}
	d2 := DummySub{0, h}

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

	var wg = sync.WaitGroup{}
	wg.Add(2)

	go d1.publish(&wg)
	go d2.publish(&wg)

	wg.Wait()

	// Wait for the Boradcast queue to clear up
	time.Sleep(10 * time.Millisecond)

	if d1.invoked != 2000 {
		t.Fatalf("Invocations espected: %d, got: %d", 2000, d1.invoked)
	}
	if d2.invoked != 2000 {
		t.Fatalf("Invocations espected: %d, got: %d", 2000, d2.invoked)
	}

	h.Deregister <- &d1
	h.Deregister <- &d2

	time.Sleep(1 * time.Millisecond)
	if h.Subscribers[&d1] {
		t.Fatalf("Could not deregister the subscriber: %v", d1)
	}
	if h.Subscribers[&d2] {
		t.Fatalf("Could not deregister the subscriber: %v", d2)
	}
}
