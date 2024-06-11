package internal

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	Send   chan (SMessage)
	Hub    *Hub
	Canvas *Canvas
}

// Implementing the interface Subscriber.
func (c *Client) Notify(msg SMessage) {
	c.Send <- msg
}

func (c Client) HandleIncoming() {
	for {
		select {
		case msg := <-c.Hub.Broadcast:
			bmsg, err := pack(msg, c.Canvas)
			if err != nil {
				log.Printf("Got an invalid message from the hub. Please investigate.\n")
			} else if err := c.Conn.WriteMessage(websocket.BinaryMessage, bmsg) != nil; err {
				log.Printf("Could not write message: %v\n", err)
			}
		}
	}
}

func (c Client) HandleSocketIncoming() {
	for {
		msgType, msg, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("The client disconnected: %v\n", err)
				c.Hub.Deregister <- &c
				return
			} else {
				log.Printf("Could not read message: %v\n", err)
			}
		}
		if msgType == websocket.BinaryMessage {
			smsg, err := unpack(msg, c.Canvas)
			if err != nil {
				log.Printf("Got an invalid message from client %s: %v\n", c.Conn.RemoteAddr().String(), err)
			} else {
				c.Hub.Broadcast <- smsg
			}
		}
	}
}
