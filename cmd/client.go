package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	Send chan (SMessage)
	hub  *Hub
}

// Implementing the interface Subscriber.
func (c *Client) Notify(msg SMessage) {
	c.Send <- msg
}

func (c Client) HandleIncoming() {
	for {
		select {
		case msg := <-c.hub.Broadcast:
			bmsg, err := pack(msg, &appState)
			if err != nil {
				log.Printf("Got an invalid message from the hub. Please investigate.\n")
			} else if err := c.conn.WriteMessage(websocket.BinaryMessage, bmsg) != nil; err {
				log.Printf("Could not write message: %v\n", err)
			}
		}
	}
}

func (c Client) HandleSocketIncoming() {
	for {
		msgType, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("The client disconnected: %v\n", err)
				c.hub.Deregister <- &c
				return
			} else {
				log.Printf("Could not read message: %v\n", err)
			}
		}
		if msgType == websocket.BinaryMessage {
			smsg, err := unpack(msg, &appState)
			if err != nil {
				log.Printf("Got an invalid message from client %s: %v\n", c.conn.RemoteAddr().String(), err)
			} else {
				c.hub.Broadcast <- smsg
			}
		}
	}
}
