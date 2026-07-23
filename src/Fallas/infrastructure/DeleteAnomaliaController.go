// DeleteAnomaliaController.go
package infrastructure

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type DeleteAnomaliaController struct {
	deleteUseCase *application.DeleteAnomaliaUseCase
}

func NewDeleteAnomaliaController(deleteUseCase *application.DeleteAnomaliaUseCase) *DeleteAnomaliaController {
	return &DeleteAnomaliaController{
		deleteUseCase: deleteUseCase,
	}
}

// @Summary      Eliminar anomalía
// @Tags         Anomalia
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.AnomaliaMessageResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/{id} [delete]
func (ctrl *DeleteAnomaliaController) Run(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		core.RespondInvalidInput(c, "El parámetro 'id' debe ser un número entero válido")
		return
	}

	errDelete := ctrl.deleteUseCase.Run(int32(id))
	if errDelete != nil {
		if errDelete.Error() == "anomalia not found" {
			core.RespondNotFound(c, "Anomalía", idParam)
		} else {
			core.RespondInternalServerError(c, "Error al eliminar la anomalía", errDelete)
		}
		return
	}

	core.RespondOK(c, map[string]string{
		"message": "Anomalía eliminada exitosamente",
	})
}
