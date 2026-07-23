//GetNotificacionesByUsuarioYTipoController.go
package infrastructure

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type GetNotificacionesByUsuarioYTipoController struct {
	useCase *application.GetNotificacionesByUsuarioYTipoUseCase
}

func NewGetNotificacionesByUsuarioYTipoController(useCase *application.GetNotificacionesByUsuarioYTipoUseCase) *GetNotificacionesByUsuarioYTipoController {
	return &GetNotificacionesByUsuarioYTipoController{useCase: useCase}
}

// @Summary      Notificaciones por usuario y tipo
// @Tags         Notificacion
// @Produce      json
// @Success      200 {object} entities.NotificacionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/usuario/{usuario_id}/tipo/{tipo} [get]
func (ctrl *GetNotificacionesByUsuarioYTipoController) Run(c *gin.Context) {
	usuarioIDParam := c.Param("usuario_id")
	usuarioID, err := strconv.Atoi(usuarioIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID de usuario inválido", "error": err.Error()})
		return
	}

	tipo := c.Param("tipo")
	
	result, err := ctrl.useCase.Run(int32(usuarioID), tipo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudieron obtener las notificaciones", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
