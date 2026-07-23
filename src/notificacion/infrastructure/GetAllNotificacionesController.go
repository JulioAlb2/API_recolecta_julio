//GetAllNotificacionesController.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type GetAllNotificacionesController struct {
	getAllUseCase *application.GetAllNotificacionesUseCase
}

func NewGetAllNotificacionesController(getAllUseCase *application.GetAllNotificacionesUseCase) *GetAllNotificacionesController {
	return &GetAllNotificacionesController{
		getAllUseCase: getAllUseCase,
	}
}

// @Summary      Listar notificaciones
// @Tags         Notificacion
// @Produce      json
// @Success      200 {object} entities.NotificacionListResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/notificaciones/ [get]
func (ctrl *GetAllNotificacionesController) Run(c *gin.Context) {
	notificaciones, err := ctrl.getAllUseCase.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "No se pudieron obtener las notificaciones",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, notificaciones)
}
