package hub

import (
	"fmt"
	"net/http"
	"time"

	"encoding/json"

	"github.com/gorilla/websocket"
)

type message struct {
	client    *client
	timestamp time.Time
	game      *string
	content   []byte
}

type hub struct {
	clients        map[*client]string
	games          map[string]map[*client]bool
	receive        chan message
	broadcast      chan message
	subscription   chan subscription
	unsubscription chan *client
	upgrader       websocket.Upgrader
}

type subscription struct {
	client *client
	game   string
}

var (
	h = hub{
		clients:        make(map[*client]string),
		games:          make(map[string]map[*client]bool),
		receive:        make(chan message),
		broadcast:      make(chan message),
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
		case m := <-h.broadcast:
			if m.game != nil {
				h.broadcastToGame(*m.game, m.content)
			}
		case m := <-h.receive:
			h.handle(&m)
		}
	}
}

func (h *hub) handle(m *message) {
	var f interface{}

	err := json.Unmarshal(m.content, &f)

	if err != nil {
		fmt.Printf("%s: Error decoding json: %s\n", m.timestamp, err)
		return
	}

	msg := f.(map[string]interface{})

	switch msg["action"] {
	case "SUBSCRIBE":
		switch msg["game"].(type) {
		case string:
			g := msg["game"].(string)
			go func() {
				h.subscription <- subscription{client: m.client, game: g}
			}()
		}
	default:
		fmt.Printf("%s: Unknown action requested: %s\n", m.timestamp, msg["action"])
	}
}

func (h *hub) subscribe(c *client, g string) error {
	h.clients[c] = g

	_, ok := h.games[g]
	if !ok {
		h.games[g] = map[*client]bool{}
	}
	h.games[g][c] = true

	fmt.Println("Subscribe client to game: " + g)

	go func() {
		h.broadcast <- message{
			client:    c,
			content:   []byte(fmt.Sprintf("You are now playing %s\n", g)),
			game:      &g,
			timestamp: time.Now(),
		}
	}()

	return nil
}

func (h *hub) unsubscribe(c *client) error {
	g := h.clients[c]
	delete(h.games[g], c)
	if len(h.games[g]) == 0 {
		delete(h.games, g)
	}
	delete(h.clients, c)

	return nil
}

func (h *hub) broadcastToGame(g string, m []byte) {
	for c := range h.games[g] {
		select {
		case c.send <- m:
			break
		// Client unreachable?
		default:
			h.unsubscription <- c
		}
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
	go client.readPump(h.unsubscription, h.receive)
	go client.writePump()
}
