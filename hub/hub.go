package hub

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Message maps to/from json object sent/received via websocket
type Message struct {
	Handle string `json:"handle"`
	Data   string `json:"data"`
}

var (
	sockets = map[*websocket.Conn]map[string]bool{}
	games   = map[string]map[*websocket.Conn]bool{}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	l sync.RWMutex
)

// Subscribe a websocket client connection to a game
func Subscribe(ws *websocket.Conn, g string) error {
	l.Lock()
	defer l.Unlock()

	_, ok := sockets[ws]
	if !ok {
		sockets[ws] = map[string]bool{}
	}
	sockets[ws][g] = true

	_, ok = games[g]
	if !ok {
		games[g] = map[*websocket.Conn]bool{}
	}
	games[g][ws] = true

	return nil
}

// UnsubscribeAll removes a client from all games
func UnsubscribeAll(ws *websocket.Conn) error {
	l.Lock()
	defer l.Unlock()

	for g := range sockets[ws] {
		delete(games[g], ws)
		if len(games[g]) == 0 {
			delete(games, g)
		}
	}

	delete(sockets, ws)

	return nil
}

// EmitToSockets sends a message to all sockets interested in game g
func EmitToSockets(g string, m Message) {
	l.RLock()
	defer l.RUnlock()

	for s := range games[g] {
		s.WriteJSON(m)
	}
}

// HandleWebsocket handles websocket connection
func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
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

	// TODO: Create read/write pumps to prevent simultaneous r/w on same socket
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) || err == io.EOF {
				fmt.Println("Websocket closed")
				e := UnsubscribeAll(ws)
				if e != nil {
					fmt.Println(e)
				}
				break
			}

			fmt.Println("Error reading websocket message")
		}

		fmt.Println(msg.Handle + ":" + msg.Data)

		switch msg.Handle {
		case "subscribe":
			// TODO: Validate that msg.Data is a valid game
			Subscribe(ws, msg.Data)
		default:
			fmt.Println("Unknown handle " + msg.Handle)
		}
	}
}
