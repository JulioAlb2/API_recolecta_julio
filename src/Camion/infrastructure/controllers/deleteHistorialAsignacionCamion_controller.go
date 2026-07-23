package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
)

type DeleteHistorialAsignacionCamionController struct {
	uc *application.DeleteHistorialAsignacionCamionUseCase
}

func NewDeleteHistorialAsignacionCamionController(uc *application.DeleteHistorialAsignacionCamionUseCase) *DeleteHistorialAsignacionCamionController {
	return &DeleteHistorialAsignacionCamionController{uc: uc}
}

// @Summary      Eliminar historial de asignación
// @Tags         HistorialAsignacion
// @Produce      json
// @Param        id path int true "ID historial"
// @Success      200 {object} entities.HistorialAsignacionCamionMessageResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/historial-asignacion/{id} [delete]
func (ctr *DeleteHistorialAsignacionCamionController) Run(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id inválido"})
		return
	}

	if err := ctr.uc.Run(int32(id)); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "registro eliminado"})
}
