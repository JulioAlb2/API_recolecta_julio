// CreateAnomaliaController.go
package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CreateAnomaliaController struct {
	createUseCase *application.CreateAnomaliaUseCase
}

func NewCreateAnomaliaController(createUseCase *application.CreateAnomaliaUseCase) *CreateAnomaliaController {
	return &CreateAnomaliaController{
		createUseCase: createUseCase,
	}
}

// @Summary      Crear anomalía
// @Description  Abierto a cualquier usuario autenticado (ciudadano, conductor
// @Description  o staff) -- solo el resto del CRUD de anomalias queda
// @Description  restringido a ADMIN/SUPERVISOR/COORDINADOR.
// @Tags         Anomalia
// @Produce      json
// @Param        body body entities.CreateAnomaliaRequest true "Body"
// @Success      200 {object} entities.AnomaliaResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/anomalias/ [post]
func (ctrl *CreateAnomaliaController) Run(c *gin.Context) {
	var request struct {
		TipoAnomalia         string `json:"tipo_anomalia" binding:"required"`
		PuntoID              *int32 `json:"punto_id"`
		ConductorID          *int32 `json:"conductor_id"`
		CamionID             *int32 `json:"camion_id"`
		RutaID               *int32 `json:"ruta_id"`
		AnomaliaReferenciaID *int32 `json:"anomalia_referencia_id"`
		Descripcion          string `json:"descripcion" binding:"required"`
		JsonRuta             string `json:"json_ruta"`
		Estado               string `json:"estado" binding:"required"`
		FechaReporte         string `json:"fecha_reporte" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		core.RespondValidationError(c, "Datos de entrada inválidos", map[string]string{
			"error": err.Error(),
		})
		return
	}

	tipoAnomalia, err := entities.ParseTipoAnomalia(request.TipoAnomalia)
	if err != nil {
		core.RespondValidationError(c, "Tipo de anomalía inválido", map[string]string{
			"tipo_anomalia": err.Error(),
		})
		return
	}

	// Parsear fecha
	fechaReporte, err := parseFecha(request.FechaReporte)
	if err != nil {
		core.RespondBadRequest(c, "Formato de fecha inválido. Use ISO 8601: YYYY-MM-DDTHH:MM:SSZ", map[string]string{
			"field": "fecha_reporte",
			"error": err.Error(),
		})
		return
	}

	anomalia := entities.NewAnomalia(
		tipoAnomalia,
		request.PuntoID,
		request.ConductorID,
		request.CamionID,
		request.RutaID,
		request.AnomaliaReferenciaID,
		request.Descripcion,
		request.JsonRuta,
		request.Estado,
		fechaReporte,
	)

	createdAnomalia, err := ctrl.createUseCase.Run(anomalia)
	if err != nil {
		core.RespondInternalServerError(c, "No se pudo crear la anomalía", err)
		return
	}

	core.RespondCreated(c, createdAnomalia)
}
