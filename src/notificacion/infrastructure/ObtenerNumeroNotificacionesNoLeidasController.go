//ObtenerNumeroNotificacionesNoLeidasController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type ObtenerNumeroNotificacionesNoLeidasController struct {
	useCase *application.ObtenerNumeroNotificacionesNoLeidasUseCase
}

func NewObtenerNumeroNotificacionesNoLeidasController(useCase *application.ObtenerNumeroNotificacionesNoLeidasUseCase) *ObtenerNumeroNotificacionesNoLeidasController {
	return &ObtenerNumeroNotificacionesNoLeidasController{useCase: useCase}
}

// @Summary      Número de notificaciones no leídas
// @Tags         Notificacion
// @Produce      json
// @Success      200 {object} entities.NotificacionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/no-leidas/usuario/{usuario_id} [get]
func (ctrl *ObtenerNumeroNotificacionesNoLeidasController) Run(c *gin.Context) {
	usuarioIDParam := c.Param("usuario_id")
	usuarioID, err := strconv.Atoi(usuarioIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID de usuario inválido", "error": err.Error()})
		return
	}

	count, err := ctrl.useCase.Run(int32(usuarioID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudo obtener el número de notificaciones no leídas", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notificaciones_no_leidas": count})
}
