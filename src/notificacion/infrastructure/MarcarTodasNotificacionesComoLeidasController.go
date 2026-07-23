//MarcarTodasNotificacionesComoLeidasController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type MarcarTodasNotificacionesComoLeidasController struct {
	useCase *application.MarcarTodasNotificacionesComoLeidasUseCase
}

func NewMarcarTodasNotificacionesComoLeidasController(useCase *application.MarcarTodasNotificacionesComoLeidasUseCase) *MarcarTodasNotificacionesComoLeidasController {
	return &MarcarTodasNotificacionesComoLeidasController{useCase: useCase}
}

// @Summary      Marcar todas como leídas
// @Tags         Notificacion
// @Produce      json
// @Success      200 {object} entities.NotificacionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/usuario/{usuario_id}/marcar-todas-leidas [patch]
func (ctrl *MarcarTodasNotificacionesComoLeidasController) Run(c *gin.Context) {
	usuarioIDParam := c.Param("usuario_id")
	usuarioID, err := strconv.Atoi(usuarioIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID de usuario inválido", "error": err.Error()})
		return
	}

	err = ctrl.useCase.Run(int32(usuarioID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudieron marcar todas las notificaciones como leídas", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Todas las notificaciones fueron marcadas como leídas"})
}
