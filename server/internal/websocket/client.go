package websocket

import (
	"log"

	"golang.org/x/net/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
	Coord
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg ReceiveMessage
		err := websocket.JSON.Receive(c.Conn, &msg)
		if err != nil {
			log.Printf("Receive message error: %v", err)
			break
		}
		log.Printf("Received message: %+v\n", msg)
		c.Coord = msg.Coord
		c.Pool.Broadcast <- SendMessage{
			ID:             c.ID,
			ReceiveMessage: msg,
		}
	}
}
