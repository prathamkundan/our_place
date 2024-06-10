package main

// AppState
/* Represents the state of the canvas in memory */
type AppState struct {
	Updates chan (SMessage)

	height uint32
	width  uint32
	canvas []SColor
}

func NewCanvas(height, width uint32) AppState {
	return AppState{Updates: make(chan SMessage, 128), height: height, width: width, canvas: make([]SColor, height*width)}
}

// Implements the Subscriber interface to allow it to be notified by the hub
func (a AppState) Notify(msg SMessage) {
	a.Updates <- msg
}

func (a AppState) HandleIncoming() {
	for {
		select {
		case msg := <-a.Updates:
			a.canvas[msg.pos] = msg.color
		}
	}
}
