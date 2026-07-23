//CountNotificacionesByTipoController.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type CountNotificacionesByTipoController struct {
	useCase *application.CountNotificacionesByTipoUseCase
}

func NewCountNotificacionesByTipoController(useCase *application.CountNotificacionesByTipoUseCase) *CountNotificacionesByTipoController {
	return &CountNotificacionesByTipoController{useCase: useCase}
}

// @Summary      Contar notificaciones por tipo
// @Tags         Notificacion
// @Produce      json
// @Success      200 {object} entities.NotificacionListResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/count/tipo/{tipo} [get]
func (ctrl *CountNotificacionesByTipoController) Run(c *gin.Context) {
	tipo := c.Param("tipo")
	
	count, err := ctrl.useCase.Run(tipo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudo contar las notificaciones", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
