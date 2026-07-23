//MarcarNotificacionComoActivaController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type MarcarNotificacionComoActivaController struct {
	useCase *application.MarcarNotificacionComoActivaUseCase
}

func NewMarcarNotificacionComoActivaController(useCase *application.MarcarNotificacionComoActivaUseCase) *MarcarNotificacionComoActivaController {
	return &MarcarNotificacionComoActivaController{useCase: useCase}
}

// @Summary      Reactivar notificación
// @Tags         Notificacion
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.NotificacionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/{id}/reactivar [patch]
func (ctrl *MarcarNotificacionComoActivaController) Run(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido", "error": err.Error()})
		return
	}

	err = ctrl.useCase.Run(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudo reactivar la notificación", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Notificación reactivada"})
}
