//GetNotificacionesActivasByUsuarioIDController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type GetNotificacionesActivasByUsuarioIDController struct {
	useCase *application.GetNotificacionesActivasByUsuarioIDUseCase
}

func NewGetNotificacionesActivasByUsuarioIDController(useCase *application.GetNotificacionesActivasByUsuarioIDUseCase) *GetNotificacionesActivasByUsuarioIDController {
	return &GetNotificacionesActivasByUsuarioIDController{useCase: useCase}
}

// @Summary      Notificaciones activas por usuario
// @Tags         Notificacion
// @Produce      json
// @Success      200 {object} entities.NotificacionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/activas/usuario/{usuario_id} [get]
func (ctrl *GetNotificacionesActivasByUsuarioIDController) Run(c *gin.Context) {
	idParam := c.Param("usuario_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido", "error": err.Error()})
		return
	}

	result, err := ctrl.useCase.Run(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudieron obtener las notificaciones activas", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
