package websocket

type ReceiveMessage struct {
	Coord struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"coord"`
}

type SendMessage struct {
	ID string `json:"id"`
	ReceiveMessage
}
