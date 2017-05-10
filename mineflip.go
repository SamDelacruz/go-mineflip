package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

type message struct {
	Handle string `json:"handle"`
	Data   string `json:"data"`
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		m := "Unable to upgrade to websocket"
		fmt.Println("Error: " + m)
		return
	}

	for {
		mt, data, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) || err == io.EOF {
				fmt.Println("Websocket closed")
				break
			}

			fmt.Println("Error reading websocket message")
		}
		switch mt {
		case websocket.TextMessage:
			msg, err := validateMessage(data)
			if err != nil {
				fmt.Println("Invalid message")
			}
			fmt.Println(msg.Handle + ":" + msg.Data)
		default:
			fmt.Println("Unknown message type!")
		}
	}
}

func validateMessage(data []byte) (message, error) {
	var msg message

	if err := json.Unmarshal(data, &msg); err != nil {
		return msg, errors.Wrap(err, "Unmarshalling message")
	}

	if msg.Handle == "" && msg.Data == "" {
		return msg, errors.New("Message has no Handle or Text")
	}

	return msg, nil
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT env variable must be set")
	}

	http.HandleFunc("/ws", handleWebsocket)
	log.Println(http.ListenAndServe(":"+port, nil))
}
