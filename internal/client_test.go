package internal

import (
	"bytes"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var canvas = NewCanvas(10, 10)

var hub = SetupHub()

var client *Client = nil

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	c, err := SetupClient(conn, hub, canvas)
	if err != nil {
		log.Println("Error starting client")
	}
	client = c
}

func run_server() {
	go canvas.HandleIncoming()
	hub.Register <- canvas
	http.HandleFunc("/ws", handleConnection)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func TestClient(t *testing.T) {
	log.Println("Running client Test")
	canvas.canvas[1] = VIOLET

	go run_server()
	dialer := websocket.DefaultDialer

	dummyUser, _, err := dialer.Dial("ws://localhost:8000/ws", nil)
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}

	_, msg, err := dummyUser.ReadMessage()
	if !bytes.Equal(msg, canvas.PackCanvas()) {
		t.Fatalf("The canvas was not sent properly")
	}

	colorVal, posVal := BLACK, 12
	action, _ := pack(SMessage{
		pos:       uint32(posVal),
		color:     SColor(colorVal),
		timestamp: time.Now(),
	})

	dummyUser.WriteMessage(websocket.BinaryMessage, action)
	time.Sleep(1 * time.Millisecond)

	if canvas.at(posVal) != colorVal {
		t.Fatalf("Update not propagated expected: %d got: %d", colorVal, canvas.canvas[12])
	}

	dummyUser.Close()
}
