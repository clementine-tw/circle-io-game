package websocket

import (
	"golang.org/x/net/websocket"
)

func Handler(handler func(*websocket.Conn)) websocket.Handler {
	return websocket.Handler(handler)
}
