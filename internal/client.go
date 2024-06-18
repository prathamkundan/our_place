package internal

import (
	"fmt"
	"log"
	. "space/internal/pubsub"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	send   chan (SMessage)
	Hub    Publisher[SMessage]
	Canvas *Canvas

	interrupt          chan bool
	unsuccessful_reads int
}

func (c *Client) String() string {
	return fmt.Sprintf("Client: %s", c.conn.RemoteAddr())
}

func SetupClient(conn *websocket.Conn, hub Publisher[SMessage], canvas *Canvas) (*Client, error) {
	c := Client{
		send:   make(chan (SMessage), 8),
		conn:   conn,
		Hub:    hub,
		Canvas: canvas,

		interrupt:          make(chan bool),
		unsuccessful_reads: 0,
	}

	log.Println("Sent pull response to: ", c.conn.RemoteAddr())
	hub.Register(&c)
	err := c.conn.WriteMessage(websocket.BinaryMessage, c.Canvas.PackCanvas())
	if err != nil {
		log.Printf("Could not send the initial Canvas")
		// Closing the connection
		c.conn.Close()
		return nil, err
	}

	go c.HandleConnection()

	return &c, nil
}

// Implementing the interface Subscriber.
func (c *Client) Notify(msg SMessage) {
	c.send <- msg
}

func (c *Client) Interrupt() {
    c.interrupt <- true
    close(c.interrupt)
}

func (c *Client) Listen() {
	for {
		select {
		case msg := <-c.send:
			log.Printf("Trying to send %s: %s", c.conn.RemoteAddr(), msg)
			bmsg, err := pack(msg)
			if err != nil {
				log.Printf("Got an invalid message from the hub. Please investigate.\n")
			} else if err := c.conn.WriteMessage(websocket.BinaryMessage, bmsg) != nil; err {
				log.Printf("Could not write message: %v\n", err)
			}
		case <-c.interrupt:
			log.Printf("%s Deregistered from hub. Stopping.\n", c.conn.RemoteAddr())
			return
		}
	}
}

func (c *Client) HandleConnection() {
	for {
		msgType, msg, err := c.conn.ReadMessage()
		log.Println("Message from: ", c.conn.RemoteAddr())

		if err != nil {
			if websocket.IsCloseError(
				err, websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseMessage,
			) {
				log.Printf("The client disconnected: %v\n", err)
				c.Hub.Deregister(c)
				close(c.send)
				return
			} else {
				log.Printf("Could not read message from %s: %v\n", c.conn.RemoteAddr(), err)
				if c.unsuccessful_reads >= 5 {
					log.Println("Failed 5 times. Quitting", c.conn.RemoteAddr(), err)
					c.Hub.Deregister(c)
					close(c.send)
					return
				}
				c.unsuccessful_reads++

			}
		}
		if msgType == websocket.BinaryMessage {
			c.unsuccessful_reads = 0
			smsg, err := unpack(msg)
			log.Printf("Got message from %s: %s", c.conn.RemoteAddr(), smsg)
			log.Println(smsg)
			if err != nil {
				log.Printf("Got an invalid message from client %s: %v\n", c.conn.RemoteAddr().String(), err)
			} else {
				c.Hub.Broadcast(smsg)
			}
		}
	}
}
