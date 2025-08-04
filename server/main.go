package main

import (
	"net/http"
	"time"

	"log"

	"github.com/clementine-tw/circle-io-game/server/internal/websocket"
)

func main() {

	mux := http.NewServeMux()

	pool := websocket.NewPool()
	go pool.Start()
	mux.Handle("/", websocket.Handler(pool.ServeWS))

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
