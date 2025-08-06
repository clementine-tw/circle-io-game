package websocket

import (
	"fmt"
	"log"

	"github.com/clementine-tw/circle-io-game/server/internal/auth"
	"golang.org/x/net/websocket"
)

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan SendMessage
	Clients    map[*Client]struct{}
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan SendMessage),
		Clients:    make(map[*Client]struct{}),
	}
}

func (p *Pool) Start() {
	for {
		select {
		case client := <-p.Register:
			log.Printf("New client %q joined", client.ID)
			for other := range p.Clients {
				err := websocket.JSON.Send(client.Conn, SendMessage{
					ID:             other.ID,
					ReceiveMessage: ReceiveMessage{Coord: other.Coord},
				})
				if err != nil {
					log.Printf("Send the positions of other player to new user error: %v", err)
					break
				}
			}
			p.Clients[client] = struct{}{}
		case client := <-p.Unregister:
			log.Printf("Client %q left", client.ID)
			delete(p.Clients, client)
		case msg := <-p.Broadcast:
			for client := range p.Clients {
				err := websocket.JSON.Send(client.Conn, msg)
				if err != nil {
					log.Printf("Broadcast to %q failed: %v", client.ID, err)
				}
			}
		}
	}
}

func (p *Pool) ServeWS(ws *websocket.Conn) {
	fmt.Println("Serve new client")
	userID, err := auth.GetUrlQueryToken(ws.Request())
	if err != nil {
		log.Printf("Couldn't get user id: %v", err)
		ws.Write([]byte("Couldn't get user id, closing connection..."))
		ws.WriteClose(1000)
		return
	}
	client := &Client{
		ID:   userID,
		Conn: ws,
		Pool: p,
	}
	p.Register <- client
	client.Read()
}
