package tracking_ws

import "github.com/gin-gonic/gin"

func RegisterRoutes(engine *gin.Engine, handler *Handler) {
	engine.GET("/ws", handler.ServeWS)
}
