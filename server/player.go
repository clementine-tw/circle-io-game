package main

import (
	"log"

	"golang.org/x/net/websocket"
)

type Player struct {
	ID     string
	Radius int
	Coord
	Conn *websocket.Conn
	Game *Game
}

func (p *Player) Read() {
	log.Printf("Reading message from player %q", p.ID)
	defer func() {
		p.Game.Unregister <- p
		p.Conn.Close()
	}()

	for {
		var msg UpdateMessage
		err := websocket.JSON.Receive(p.Conn, &msg)
		if err != nil {
			log.Printf("Couldn't receive message: %v", err)
			break
		}
		msg.Player = p
		p.Game.Update <- msg
	}
}
