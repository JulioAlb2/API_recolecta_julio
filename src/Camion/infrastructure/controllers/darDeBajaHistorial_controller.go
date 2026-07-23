package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
)

type DarDeBajaHistorialAsignacionController struct {
	uc *application.DarDeBajaHistorialAsignacionUseCase
}

func NewDarDeBajaHistorialAsignacionController(uc *application.DarDeBajaHistorialAsignacionUseCase) *DarDeBajaHistorialAsignacionController {
	return &DarDeBajaHistorialAsignacionController{uc: uc}
}

// @Summary      Dar de baja historial
// @Tags         HistorialAsignacion
// @Produce      json
// @Param        id path int true "ID historial"
// @Success      200 {object} entities.HistorialAsignacionCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/historial-asignacion/baja/{id} [put]
func (ctr *DarDeBajaHistorialAsignacionController) Run(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "id inválido"})
		return
	}

	if err := ctr.uc.Run(int32(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "asignación cerrada"})
}
