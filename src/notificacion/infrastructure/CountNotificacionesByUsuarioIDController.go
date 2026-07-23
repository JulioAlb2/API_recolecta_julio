//CountNotificacionesByUsuarioIDController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type CountNotificacionesByUsuarioIDController struct {
	useCase *application.CountNotificacionesByUsuarioIDUseCase
}

func NewCountNotificacionesByUsuarioIDController(useCase *application.CountNotificacionesByUsuarioIDUseCase) *CountNotificacionesByUsuarioIDController {
	return &CountNotificacionesByUsuarioIDController{useCase: useCase}
}

// @Summary      Contar notificaciones por usuario
// @Tags         Notificacion
// @Produce      json
// @Success      200 {object} entities.NotificacionListResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/count/usuario/{usuario_id} [get]
func (ctrl *CountNotificacionesByUsuarioIDController) Run(c *gin.Context) {
	usuarioIDParam := c.Param("usuario_id")
	usuarioID, err := strconv.Atoi(usuarioIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID de usuario inválido", "error": err.Error()})
		return
	}

	count, err := ctrl.useCase.Run(int32(usuarioID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudo contar las notificaciones", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
