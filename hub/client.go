package hub

import "github.com/gorilla/websocket"

type client struct {
	ws   *websocket.Conn
	send chan []byte
	exit chan struct{}
}

func (c *client) writePump() {
	defer func() {
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-c.exit:
			return
		}
	}
}

func (c *client) write(mt int, message []byte) error {
	return c.ws.WriteMessage(mt, message)
}

func (c *client) readPump() {
	defer func() {
		var s struct{}
		c.exit <- s
		h.unsubscription <- c
		c.ws.Close()
	}()

	for {
		_, m, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		h.receive <- message{c, string(m)}
	}
}
