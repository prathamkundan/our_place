package internal

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	Send   chan (SMessage)
	Hub    *Hub
	Canvas *Canvas

	interrupt          chan bool
	unsuccessful_reads int
}

func (c *Client) String() string {
	return fmt.Sprintf("Client: %s", c.Conn.RemoteAddr())
}

// Implementing the interface Subscriber.
func (c *Client) Notify(msg SMessage) {
	c.Send <- msg
}

func SetupClient(conn *websocket.Conn, hub *Hub, canvas *Canvas) (*Client, error) {
	c := Client{
		Send:   make(chan (SMessage), 8),
		Conn:   conn,
		Hub:    hub,
		Canvas: canvas,

		interrupt:          make(chan bool),
		unsuccessful_reads: 0,
	}

	log.Println("Sent pull response to: ", c.Conn.RemoteAddr())
	hub.Register <- &c
	err := c.Conn.WriteMessage(websocket.BinaryMessage, c.Canvas.PackCanvas())
	if err != nil {
		log.Printf("Could not send the initial Canvas")
		// Closing the connection
		c.Conn.Close()
		return nil, err
	}

	go c.HandleIncoming()
	go c.HandleSocketIncoming()

	return &c, nil
}

func (c *Client) HandleIncoming() {
	for {
		select {
		case msg := <-c.Send:
			log.Printf("Trying to send %s: %s", c.Conn.RemoteAddr(), msg)
			bmsg, err := pack(msg)
			if err != nil {
				log.Printf("Got an invalid message from the hub. Please investigate.\n")
			} else if err := c.Conn.WriteMessage(websocket.BinaryMessage, bmsg) != nil; err {
				log.Printf("Could not write message: %v\n", err)
			}
		case <-c.interrupt:
			log.Printf("%s Deregistered from hub. Stopping.\n", c.Conn.RemoteAddr())
			return
		}
	}
}

func (c *Client) HandleSocketIncoming() {
	for {
		msgType, msg, err := c.Conn.ReadMessage()
		log.Println("Message from: ", c.Conn.RemoteAddr())

		if err != nil {
			if websocket.IsCloseError(
				err, websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseMessage,
			) {
				log.Printf("The client disconnected: %v\n", err)
				c.Hub.Deregister <- c
				c.interrupt <- true
				close(c.Send)
				close(c.interrupt)
				return
			} else {
				log.Printf("Could not read message from %s: %v\n", c.Conn.RemoteAddr(), err)
				if c.unsuccessful_reads >= 5 {
					log.Println("Failed 5 times. Quitting", c.Conn.RemoteAddr(), err)
					c.Hub.Deregister <- c
					c.interrupt <- true
					close(c.Send)
					close(c.interrupt)
					return
				}
				c.unsuccessful_reads++

			}
		}
		if msgType == websocket.BinaryMessage {
			c.unsuccessful_reads = 0
			smsg, err := unpack(msg)
			log.Printf("Got message from %s: %s", c.Conn.RemoteAddr(), smsg)
			log.Println(smsg)
			if err != nil {
				log.Printf("Got an invalid message from client %s: %v\n", c.Conn.RemoteAddr().String(), err)
			} else {
				c.Hub.Broadcast <- smsg
			}
		}
	}
}
