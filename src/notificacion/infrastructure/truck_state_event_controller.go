package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type TruckStateEventController struct {
	queryUc   *application.QueryEventTraceUseCase
}

func NewTruckStateEventController(queryUc *application.QueryEventTraceUseCase) *TruckStateEventController {
	return &TruckStateEventController{queryUc: queryUc}
}

func (ctrl *TruckStateEventController) GetByEventID(c *gin.Context) {
	eventID := c.Param("event_id")
	trace, err := ctrl.queryUc.GetByEventID(c.Request.Context(), eventID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, trace)
}

func (ctrl *TruckStateEventController) ListByTruckID(c *gin.Context) {
	truckIDStr := c.Param("truck_id")
	truckID, err := strconv.ParseInt(truckIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid truck_id"})
		return
	}
	traces, err := ctrl.queryUc.ListByTruckID(c.Request.Context(), int32(truckID), 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, traces)
}

