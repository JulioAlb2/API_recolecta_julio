//CountNotificacionesByCamionIDController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type CountNotificacionesByCamionIDController struct {
	useCase *application.CountNotificacionesByCamionIDUseCase
}

func NewCountNotificacionesByCamionIDController(useCase *application.CountNotificacionesByCamionIDUseCase) *CountNotificacionesByCamionIDController {
	return &CountNotificacionesByCamionIDController{useCase: useCase}
}

// @Summary      Contar notificaciones por camión
// @Tags         Notificacion
// @Produce      json
// @Success      200 {object} entities.NotificacionListResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/count/camion/{camion_id} [get]
func (ctrl *CountNotificacionesByCamionIDController) Run(c *gin.Context) {
	camionIDParam := c.Param("camion_id")
	camionID, err := strconv.Atoi(camionIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID de camión inválido", "error": err.Error()})
		return
	}

	count, err := ctrl.useCase.Run(int32(camionID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudo contar las notificaciones", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
