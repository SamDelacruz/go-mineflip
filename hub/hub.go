package hub

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type message struct {
	client  *client
	content string
}

type hub struct {
	clients        map[*client]map[string]bool
	games          map[string]map[*client]bool
	receive        chan message
	broadcast      chan broadcast
	subscription   chan subscription
	unsubscription chan *client
	upgrader       websocket.Upgrader
}

type subscription struct {
	client *client
	game   string
}

type broadcast struct {
	game    string
	message message
}

var (
	h = hub{
		clients:        make(map[*client]map[string]bool),
		games:          make(map[string]map[*client]bool),
		receive:        make(chan message),
		broadcast:      make(chan broadcast),
		subscription:   make(chan subscription),
		unsubscription: make(chan *client),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
)

// Run the hub goroutine
func Run() {
	for {
		select {
		case sub := <-h.subscription:
			h.subscribe(sub.client, sub.game)
		case unsub := <-h.unsubscription:
			h.unsubscribe(unsub)
		case broadcast := <-h.broadcast:
			h.broadcastToGame(broadcast.game, broadcast.message)
		case m := <-h.receive:
			go h.handle(m)
		}
	}
}

func (h *hub) handle(m message) {
	fmt.Printf("Handling m: %s\n", m)
}

func (h *hub) subscribe(c *client, g string) error {
	_, ok := h.clients[c]
	if !ok {
		h.clients[c] = map[string]bool{}
	}
	h.clients[c][g] = true

	_, ok = h.games[g]
	if !ok {
		h.games[g] = map[*client]bool{}
	}
	h.games[g][c] = true

	return nil
}

func (h *hub) unsubscribe(c *client) error {
	for g := range h.clients[c] {
		delete(h.games[g], c)
		if len(h.games[g]) == 0 {
			delete(h.games, g)
		}
	}

	delete(h.clients, c)

	return nil
}

func (h *hub) broadcastToGame(g string, m message) {
	for c := range h.games[g] {
		data, err := json.Marshal(m)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
		select {
		case c.send <- data:
			break
		// Client unreachable?
		default:
			h.unsubscribe(c)
		}
		c.send <- data
	}
}

// HandleWebsocket handles websocket connection
func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ws, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		m := "Unable to upgrade to websocket"
		fmt.Println("Error: " + m)
		return
	}

	client := client{ws: ws, send: make(chan []byte), exit: make(chan struct{})}
	go client.readPump()
	go client.writePump()
}