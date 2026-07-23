//GetNotificacionesByFechaRangeController.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type GetNotificacionesByFechaRangeController struct {
	useCase *application.GetNotificacionesByFechaRangeUseCase
}

func NewGetNotificacionesByFechaRangeController(useCase *application.GetNotificacionesByFechaRangeUseCase) *GetNotificacionesByFechaRangeController {
	return &GetNotificacionesByFechaRangeController{useCase: useCase}
}

// @Summary      Notificaciones por rango de fecha
// @Tags         Notificacion
// @Produce      json
// @Success      200 {object} entities.NotificacionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/rango-fecha [get]
func (ctrl *GetNotificacionesByFechaRangeController) Run(c *gin.Context) {
	fechaInicio := c.Query("fecha_inicio")
	fechaFin := c.Query("fecha_fin")
	
	if fechaInicio == "" || fechaFin == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Se requieren fecha_inicio y fecha_fin"})
		return
	}

	result, err := ctrl.useCase.Run(fechaInicio, fechaFin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudieron obtener las notificaciones", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
