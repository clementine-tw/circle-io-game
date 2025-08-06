package main

import (
	"net/http"
	"time"

	"golang.org/x/net/websocket"
	"log"
	"os"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("Env 'PORT' is not set, use default: %q", defaultPort)
		port = defaultPort
	}

	mux := http.NewServeMux()

	game := NewGame()
	go game.Start()
	mux.Handle("/", websocket.Handler(game.ServeWS))

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
