package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
)

type CerrarAsignacionActivaCamionController struct {
	uc *application.CerrarAsignacionActivaCamionUseCase
}

func NewCerrarAsignacionActivaCamionController(uc *application.CerrarAsignacionActivaCamionUseCase) *CerrarAsignacionActivaCamionController {
	return &CerrarAsignacionActivaCamionController{uc: uc}
}

// @Summary      Cerrar asignación activa de camión
// @Tags         HistorialAsignacion
// @Produce      json
// @Param        camionId path int true "ID camión"
// @Success      200 {object} entities.HistorialAsignacionCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/historial-asignacion/cerrar/camion/{camionId} [put]
func (ctr *CerrarAsignacionActivaCamionController) Run(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("camionId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "camionId inválido"})
		return
	}

	if err := ctr.uc.Run(int32(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "asignación activa del camión cerrada"})
}
