package tracking_ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 8 * 1024
)

// Client es una conexión WebSocket autenticada.
type Client struct {
	hub         *Hub
	conn        *websocket.Conn
	send        chan []byte
	UserID      int
	RoleID      int
	IsPublisher bool
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("tracking_ws: cierre inesperado user_id=%d: %v", c.UserID, err)
			}
			break
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(raw, &msg); err != nil {
			c.sendJSON(map[string]interface{}{
				"type":    "error",
				"message": "JSON inválido",
			})
			continue
		}

		c.handleMessage(msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) sendJSON(payload map[string]interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	select {
	case c.send <- data:
	default:
	}
}

func (c *Client) handleMessage(msg map[string]interface{}) {
	typ, _ := msg["type"].(string)
	switch typ {
	case "ping":
		c.sendJSON(map[string]interface{}{"type": "pong"})

	case "location_update":
		if !c.IsPublisher {
			log.Printf("tracking_ws: location_update rechazada user_id=%d role_id=%d", c.UserID, c.RoleID)
			c.sendJSON(map[string]interface{}{
				"type":    "error",
				"message": "solo conductores pueden publicar ubicación",
			})
			return
		}

		out := map[string]interface{}{
			"type":         "location_update",
			"conductor_id": c.UserID,
			"timestamp":    time.Now().UnixMilli(),
		}
		for _, key := range []string{"lat", "lng", "velocidad", "rumbo", "en_servicio", "ruta_id", "state_code"} {
			if v, ok := msg[key]; ok {
				out[key] = v
			}
		}

		log.Printf("tracking_ws: location_update de conductor_id=%d lat=%v lng=%v",
			c.UserID, out["lat"], out["lng"])
		c.hub.broadcastExcept(c, out)
		c.sendJSON(map[string]interface{}{"type": "location_ack"})

	default:
		// Ignorar tipos desconocidos para no romper clientes futuros.
	}
}
