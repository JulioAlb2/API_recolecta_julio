package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type BroadcastCitizenNotificationController struct {
	uc *application.BroadcastCitizenNotificationUseCase
}

func NewBroadcastCitizenNotificationController(uc *application.BroadcastCitizenNotificationUseCase) *BroadcastCitizenNotificationController {
	return &BroadcastCitizenNotificationController{uc: uc}
}

type broadcastMessageRequest struct {
	Title string `json:"title" binding:"required"`
	Body  string `json:"body" binding:"required"`
}

// Broadcast
// @Summary      Enviar notificación general a todos los ciudadanos (Broadcast)
// @Description  Envía una notificación push con título y cuerpo a todos los ciudadanos registrados en el sistema.
// @Tags         PushNotification
// @Accept       json
// @Produce      json
// @Param        body body broadcastMessageRequest true "Datos del mensaje general"
// @Success      200 {object} map[string]domain.SendResult "Detalles de envío para cada ciudadano"
// @Failure      400 {object} map[string]string "error"
// @Failure      500 {object} map[string]string "error"
// @Security     BearerAuth
// @Router       /api/notificaciones-push/ciudadanos/difusion [post]
func (ctrl *BroadcastCitizenNotificationController) Broadcast(c *gin.Context) {
	var req broadcastMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := ctrl.uc.ExecuteBroadcast(c.Request.Context(), req.Title, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// BroadcastRoute
// @Summary      Enviar notificación general a ciudadanos de una ruta
// @Description  Envía una notificación push con título y cuerpo a todos los ciudadanos que tengan un domicilio registrado en la colonia de la ruta especificada.
// @Tags         PushNotification
// @Accept       json
// @Produce      json
// @Param        route_id path int true "ID de la Ruta"
// @Param        body body broadcastMessageRequest true "Datos del mensaje general"
// @Success      200 {object} map[string]domain.SendResult "Detalles de envío para cada ciudadano"
// @Failure      400 {object} map[string]string "error"
// @Failure      500 {object} map[string]string "error"
// @Security     BearerAuth
// @Router       /api/notificaciones-push/ciudadanos/difusion/ruta/{route_id} [post]
func (ctrl *BroadcastCitizenNotificationController) BroadcastRoute(c *gin.Context) {
	routeIDStr := c.Param("route_id")
	routeID, err := strconv.Atoi(routeIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "route_id inválido"})
		return
	}

	var req broadcastMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := ctrl.uc.ExecuteRouteBroadcast(c.Request.Context(), int32(routeID), req.Title, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// BroadcastPoint
// @Summary      Enviar notificación general a ciudadanos cerca de un punto
// @Description  Envía una notificación push con título y cuerpo a todos los ciudadanos con domicilios en un radio de 200m del punto de parada especificado.
// @Tags         PushNotification
// @Accept       json
// @Produce      json
// @Param        point_id path string true "ID del Punto de Parada (ej. 15)"
// @Param        body body broadcastMessageRequest true "Datos del mensaje general"
// @Success      200 {object} map[string]domain.SendResult "Detalles de envío para cada ciudadano"
// @Failure      400 {object} map[string]string "error"
// @Failure      500 {object} map[string]string "error"
// @Security     BearerAuth
// @Router       /api/notificaciones-push/ciudadanos/difusion/punto/{point_id} [post]
func (ctrl *BroadcastCitizenNotificationController) BroadcastPoint(c *gin.Context) {
	pointID := c.Param("point_id")

	var req broadcastMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	results, err := ctrl.uc.ExecutePointBroadcast(c.Request.Context(), pointID, req.Title, req.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
