package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/clementine-tw/circle-io-game/server/internal/auth"
	"golang.org/x/net/websocket"
)

const (
	mapWidth         = 1280
	mapHeight        = 720
	maxFoods         = 100
	randomFoodBiasis = 20
	initialRadius    = 10
)

type Game struct {
	Players    map[string]*Player
	Foods      map[string]struct{}
	Update     chan UpdateMessage
	Register   chan *Player
	Unregister chan *Player
}

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type UpdateMessage struct {
	Player *Player
	Coord  `json:"coord"`
}

type BroadcastMessage struct {
	ID     string `json:"id"`
	Radius int    `json:"radius"`
	Coord  Coord  `json:"coord"`
}

func NewGame() *Game {
	return &Game{
		Players:    make(map[string]*Player),
		Foods:      make(map[string]struct{}),
		Update:     make(chan UpdateMessage),
		Register:   make(chan *Player),
		Unregister: make(chan *Player),
	}
}

func (g *Game) Start() {
	g.initFoods()

	for {
		select {
		case msg := <-g.Update:
			log.Printf("Received: %+v", msg)
			player := msg.Player
			player.Coord = msg.Coord
			minX := player.Coord.X - player.Radius
			maxX := player.Coord.X + player.Radius
			minY := player.Coord.Y - player.Radius
			maxY := player.Coord.Y + player.Radius
			for x := minX; x <= maxX; x++ {
				for y := minY; y <= maxY; y++ {
					key := getGridKey(x, y)
					if _, ok := g.Foods[key]; ok {
						player.Radius++
						delete(g.Foods, key)
					}
				}
			}
			for _, other := range g.Players {
				err := websocket.JSON.Send(other.Conn, BroadcastMessage{
					ID:     player.ID,
					Radius: player.Radius,
					Coord:  player.Coord,
				})
				if err != nil {
					log.Printf("Send the positions of other player to new user error: %v", err)
				}
			}

		case player := <-g.Register:
			for id, other := range g.Players {
				err := websocket.JSON.Send(player.Conn, BroadcastMessage{
					ID:     id,
					Radius: other.Radius,
					Coord:  other.Coord,
				})
				if err != nil {
					log.Printf("Send the positions of other player to new user error: %v", err)
					break
				}
			}
			log.Printf("%q joined", player.ID)
			g.Players[player.ID] = player

		case player := <-g.Unregister:
			log.Printf("%q left", player.ID)
			delete(g.Players, player.ID)

		}
	}
}

func (g *Game) ServeWS(ws *websocket.Conn) {
	userID, err := auth.GetUrlQueryToken(ws.Request())
	if err != nil {
		log.Printf("Couldn't get user id: %v", err)
		ws.Write([]byte("Couldn't get user id, closing connection..."))
		ws.WriteClose(1000)
		return
	}

	player := &Player{
		ID:     userID,
		Conn:   ws,
		Radius: initialRadius,
		Game:   g,
	}
	g.Register <- player
	player.Read()
}

func (g *Game) initFoods() {
	log.Print("init foods")
	minX := randomFoodBiasis
	maxX := mapWidth - randomFoodBiasis
	minY := randomFoodBiasis
	maxY := mapHeight - randomFoodBiasis
	for len(g.Foods) < maxFoods {
		x := rand.Intn(maxX) + minX
		y := rand.Intn(maxY) + minY
		g.Foods[getGridKey(x, y)] = struct{}{}
	}

}

func getGridKey(x, y int) string {
	return fmt.Sprintf("%d_%d", x, y)
}
