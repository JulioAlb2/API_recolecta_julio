package tracking_ws

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/vicpoo/API_recolecta/src/core"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Móvil / ngrok / web local: el origen no es crítico; auth es por JWT.
		return true
	},
}

// Handler expone el upgrade WebSocket de tracking.
type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// ServeWS GET /ws?token=<JWT>&role=conductor (role opcional; el JWT manda).
func (h *Handler) ServeWS(c *gin.Context) {
	token := strings.TrimSpace(c.Query("token"))
	if token == "" {
		// Fallback por si algún cliente manda Authorization en el handshake.
		auth := c.GetHeader("Authorization")
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			token = strings.TrimSpace(auth[7:])
		}
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token requerido"})
		return
	}

	claims, err := core.ParseToken(token)
	if err != nil || claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return
	}

	roleHint := strings.ToLower(strings.TrimSpace(c.Query("role")))
	// Ciudadano JWT trae role_id=0 → solo oyente.
	// Empleado conductor: rol 4 (y 2 por seeds antiguos), o query role=conductor.
	isPublisher := claims.RoleID != 0 && (
		claims.RoleID == core.CONDUCTOR ||
			claims.RoleID == 2 ||
			roleHint == "conductor")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		hub:         h.hub,
		conn:        conn,
		send:        make(chan []byte, 64),
		UserID:      claims.UserID,
		RoleID:      claims.RoleID,
		IsPublisher: isPublisher,
	}

	h.hub.register <- client

	client.sendJSON(map[string]interface{}{
		"type":         "connected",
		"user_id":      claims.UserID,
		"role_id":      claims.RoleID,
		"is_publisher": isPublisher,
	})

	go client.writePump()
	go client.readPump()
}
