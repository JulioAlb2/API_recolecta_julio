//CountNotificacionesActivasByUsuarioIDController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type CountNotificacionesActivasByUsuarioIDController struct {
	useCase *application.CountNotificacionesActivasByUsuarioIDUseCase
}

func NewCountNotificacionesActivasByUsuarioIDController(useCase *application.CountNotificacionesActivasByUsuarioIDUseCase) *CountNotificacionesActivasByUsuarioIDController {
	return &CountNotificacionesActivasByUsuarioIDController{useCase: useCase}
}

// @Summary      Contar notificaciones activas por usuario
// @Tags         Notificacion
// @Produce      json
// @Success      200 {object} entities.NotificacionListResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/count/activas/usuario/{usuario_id} [get]
func (ctrl *CountNotificacionesActivasByUsuarioIDController) Run(c *gin.Context) {
	usuarioIDParam := c.Param("usuario_id")
	usuarioID, err := strconv.Atoi(usuarioIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID de usuario inválido", "error": err.Error()})
		return
	}

	count, err := ctrl.useCase.Run(int32(usuarioID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudo contar las notificaciones activas", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
