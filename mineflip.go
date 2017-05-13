package main

import (
	"log"
	"net/http"
	"os"

	"github.com/samdelacruz/go-mineflip/hub"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT env variable must be set")
	}

	// Each request/connection handled in a new goroutine
	go hub.Run()
	http.HandleFunc("/ws", hub.HandleWebsocket)
	log.Println(http.ListenAndServe(":"+port, nil))
}
