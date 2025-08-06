package websocket

type ReceiveMessage struct {
	Coord `json:"coord"`
}

type SendMessage struct {
	ID string `json:"id"`
	ReceiveMessage
}

type Coord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
