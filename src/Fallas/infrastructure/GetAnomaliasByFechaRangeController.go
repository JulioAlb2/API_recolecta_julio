// GetAnomaliasByFechaRangeController.go
package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAnomaliasByFechaRangeController struct {
	getByFechaRangeUseCase *application.GetAnomaliasByFechaRangeUseCase
}

func NewGetAnomaliasByFechaRangeController(getByFechaRangeUseCase *application.GetAnomaliasByFechaRangeUseCase) *GetAnomaliasByFechaRangeController {
	return &GetAnomaliasByFechaRangeController{
		getByFechaRangeUseCase: getByFechaRangeUseCase,
	}
}

// @Summary      Anomalías por rango de fecha
// @Tags         Anomalia
// @Produce      json
// @Success      200 {object} entities.AnomaliaResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/por-fecha [get]
func (ctrl *GetAnomaliasByFechaRangeController) Run(c *gin.Context) {
	fechaInicio := c.Query("fecha_inicio")
	fechaFin := c.Query("fecha_fin")

	if fechaInicio == "" || fechaFin == "" {
		core.RespondValidationError(c, "Parámetros de fecha inválidos", map[string]string{
			"fecha_inicio": "requerido",
			"fecha_fin":    "requerido",
		})
		return
	}

	anomalias, err := ctrl.getByFechaRangeUseCase.Run(fechaInicio, fechaFin)
	if err != nil {
		core.RespondInternalServerError(c, "No se pudieron obtener las anomalías para el rango de fechas", err)
		return
	}

	c.JSON(http.StatusOK, anomalias)
}
