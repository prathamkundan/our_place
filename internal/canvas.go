package internal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sync"
)

// Canvas
/* Represents the state of the canvas in memory */
type Canvas struct {
	Updates chan (SMessage)

	height uint32
	width  uint32
	canvas []SColor
	mu     sync.Mutex
}

func NewCanvas(height, width uint32) *Canvas {
	return &Canvas{Updates: make(chan SMessage, 128), height: height, width: width, canvas: make([]SColor, height*width)}
}

// Implements the Subscriber interface to allow it to be notified by the hub
func (a *Canvas) Notify(msg SMessage) {
	a.Updates <- msg
}

func (c *Canvas) HandleIncoming() {
	for {
		select {
		case msg := <-c.Updates:
			c.mu.Lock()
			c.canvas[msg.pos] = msg.color
			c.mu.Unlock()
		}
	}
}

func (c *Canvas) String() string {
	return fmt.Sprintf("Canvas: %d x %d", c.height, c.width)
}

func (c *Canvas) at(pos int) SColor {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.canvas[pos]
}

func (c *Canvas) PackCanvas() []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, []byte("PULL"))

	binary.Write(&buf, binary.LittleEndian, c.height)
	binary.Write(&buf, binary.LittleEndian, c.width)

	// Pack canvas array
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, color := range c.canvas {
		binary.Write(&buf, binary.LittleEndian, color)
	}

	return buf.Bytes()
}
