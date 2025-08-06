package main

import (
	"net/http"
	"time"

	"golang.org/x/net/websocket"
	"log"
)

func main() {

	mux := http.NewServeMux()

	game := NewGame()
	go game.Start()
	mux.Handle("/", websocket.Handler(game.ServeWS))

	const port = "8080"
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("ListenAndServe on %q", port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}
