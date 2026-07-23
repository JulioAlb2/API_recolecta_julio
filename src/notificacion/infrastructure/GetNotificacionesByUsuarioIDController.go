//GetNotificacionesByUsuarioIDController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type GetNotificacionesByUsuarioIDController struct {
	useCase *application.GetNotificacionesByUsuarioIDUseCase
}

func NewGetNotificacionesByUsuarioIDController(useCase *application.GetNotificacionesByUsuarioIDUseCase) *GetNotificacionesByUsuarioIDController {
	return &GetNotificacionesByUsuarioIDController{useCase: useCase}
}

// @Summary      Notificaciones por usuario
// @Tags         Notificacion
// @Produce      json
// @Success      200 {object} entities.NotificacionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/usuario/{usuario_id} [get]
func (ctrl *GetNotificacionesByUsuarioIDController) Run(c *gin.Context) {
	idParam := c.Param("usuario_id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID inválido", "error": err.Error()})
		return
	}

	result, err := ctrl.useCase.Run(int32(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudieron obtener las notificaciones", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
