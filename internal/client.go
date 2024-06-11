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

	interrupt chan bool
}

// Implementing the interface Subscriber.
func (c *Client) Notify(msg SMessage) {
	c.Send <- msg
}

func SetupClient(conn *websocket.Conn, hub *Hub, canvas *Canvas) error {
	c := Client{
		Send:   make(chan (SMessage), 8),
		Conn:   conn,
		Hub:    hub,
		Canvas: canvas,
	}

	c.Hub.Register <- &c
    err := c.Conn.WriteMessage(websocket.BinaryMessage, c.Canvas.PackCanvas())
	if err != nil {
		log.Printf("Could not send the initial Canvas")
		// Closing the connection
		c.Conn.Close()
		return err
	}

	go c.HandleIncoming()
	go c.HandleSocketIncoming()

    return nil
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
		case <-c.interrupt:
			log.Printf("%s Deregistered from hub. Stopping.", c.Conn.RemoteAddr())
			return
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
				c.interrupt <- true
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
