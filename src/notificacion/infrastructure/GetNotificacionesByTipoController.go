//GetNotificacionesByTipoController.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type GetNotificacionesByTipoController struct {
	useCase *application.GetNotificacionesByTipoUseCase
}

func NewGetNotificacionesByTipoController(useCase *application.GetNotificacionesByTipoUseCase) *GetNotificacionesByTipoController {
	return &GetNotificacionesByTipoController{useCase: useCase}
}

// @Summary      Notificaciones por tipo
// @Tags         Notificacion
// @Produce      json
// @Success      200 {object} entities.NotificacionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/tipo/{tipo} [get]
func (ctrl *GetNotificacionesByTipoController) Run(c *gin.Context) {
	tipo := c.Param("tipo")
	result, err := ctrl.useCase.Run(tipo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudieron obtener las notificaciones", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
