package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type CitizenInboxController struct {
	uc *application.GetCitizenInboxUseCase
}

func NewCitizenInboxController(uc *application.GetCitizenInboxUseCase) *CitizenInboxController {
	return &CitizenInboxController{uc: uc}
}

// GetInbox
// @Summary      Obtener bandeja de entrada del ciudadano
// @Description  Retorna las últimas 50 notificaciones recibidas o fallidas enviadas al ciudadano, leídas de la lista en Redis.
// @Tags         PushNotification
// @Produce      json
// @Param        citizen_id path string true "ID del Ciudadano"
// @Success      200 {array} application.InboxRecord
// @Failure      500 {object} map[string]string "error"
// @Router       /api/notificaciones-push/ciudadanos/{citizen_id}/historial [get]
func (ctrl *CitizenInboxController) GetInbox(c *gin.Context) {
	citizenID := c.Param("citizen_id")
	inbox, err := ctrl.uc.Execute(c.Request.Context(), citizenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, inbox)
}
