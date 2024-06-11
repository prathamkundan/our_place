package internal

import (
	"bytes"
	"encoding/binary"
)

// Canvas
/* Represents the state of the canvas in memory */
type Canvas struct {
	Updates chan (SMessage)

	height uint32
	width  uint32
	canvas []SColor
}

func NewCanvas(height, width uint32) *Canvas {
	return &Canvas{Updates: make(chan SMessage, 128), height: height, width: width, canvas: make([]SColor, height*width)}
}

// Implements the Subscriber interface to allow it to be notified by the hub
func (a Canvas) Notify(msg SMessage) {
	a.Updates <- msg
}

func (a Canvas) HandleIncoming() {
	for {
		select {
		case msg := <-a.Updates:
			a.canvas[msg.pos] = msg.color
		}
	}
}

func (canvas Canvas) PackCanvas() []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, []byte("PULL"))

	binary.Write(&buf, binary.LittleEndian, canvas.height)
	binary.Write(&buf, binary.LittleEndian, canvas.width)

	// Pack canvas array
	for _, color := range canvas.canvas {
		binary.Write(&buf, binary.LittleEndian, color)
	}

	return buf.Bytes()
}
