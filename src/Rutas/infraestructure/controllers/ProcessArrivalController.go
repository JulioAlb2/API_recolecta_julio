package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
)

type ProcessArrivalController struct {
	uc *application.ProcessTruckArrivalUseCase
}

func NewProcessArrivalController(uc *application.ProcessTruckArrivalUseCase) *ProcessArrivalController {
	return &ProcessArrivalController{uc: uc}
}

type ProcessArrivalRequest struct {
	TruckID int32  `json:"truck_id" binding:"required"`
	PointID string `json:"point_id" binding:"required"`
}

// Run
// @Summary      Procesar arribo de camión a parada
// @Description  Actualiza el estado a ARRIVAL y registra el último punto visitado en Redis. Además, realiza geofencing en el radio de la parada y dispara notificaciones push en tiempo real a los ciudadanos cercanos.
// @Tags         Ruta
// @Accept       json
// @Produce      json
// @Param        X-Device-MAC header string true "Dirección MAC del dispositivo Android"
// @Param        X-Device-Serial header string true "Número de serie del dispositivo Android"
// @Param        X-Device-API-Key header string true "API Key única del dispositivo"
// @Param        body body ProcessArrivalRequest true "Datos de arribo a parada"
// @Success      200 {object} map[string]string "Status OK"
// @Failure      400 {object} map[string]string "Bad Request"
// @Failure      500 {object} map[string]string "Internal Server Error"
// @Security     BearerAuth
// @Router       /api/rutas/arrival [post]
func (ctrl *ProcessArrivalController) Run(c *gin.Context) {
	tenantIDVal, exists := c.Get("tenant_id")
	tenantID, ok := tenantIDVal.(int)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant no encontrado en token"})
		return
	}

	var req ProcessArrivalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ctrl.uc.Execute(c.Request.Context(), tenantID, req.TruckID, req.PointID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "arrival_processed_and_citizens_notified"})
}
