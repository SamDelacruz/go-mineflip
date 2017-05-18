package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/samdelacruz/go-mineflip/api"
	"github.com/samdelacruz/go-mineflip/hub"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT env variable must be set")
	}

	go hub.Run()

	r := mux.NewRouter()
	r.StrictSlash(true) // Redirect trailing slashes
	r.HandleFunc("/ws", hub.HandleWebsocket)
	r.HandleFunc("/games", api.CreateGameHandler).Methods("POST")
	r.HandleFunc("/games/{id:[a-zA-Z0-9]+}", api.GetGameHandler).Methods("GET")
	r.HandleFunc("/games/{id:[a-zA-Z0-9]+}/tiles/{x:[0-4]}/{y:[0-4]}", api.MoveHandler).Methods("GET")
	http.Handle("/", r)
	log.Println(http.ListenAndServe(":"+port, nil))
}
