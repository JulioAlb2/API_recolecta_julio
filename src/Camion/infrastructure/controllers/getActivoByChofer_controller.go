package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
)

type GetActivoByChoferController struct {
	uc *application.GetActivoByChoferUseCase
}

func NewGetActivoByChoferController(uc *application.GetActivoByChoferUseCase) *GetActivoByChoferController {
	return &GetActivoByChoferController{uc: uc}
}

// @Summary      Asignación activa por chofer
// @Tags         HistorialAsignacion
// @Produce      json
// @Param        choferId path int true "ID chofer"
// @Success      200 {object} entities.HistorialAsignacionCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      404 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/historial-asignacion/activo/chofer/{choferId} [get]
func (ctr *GetActivoByChoferController) Run(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("choferId"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "choferId inválido"})
		return
	}

	data, err := ctr.uc.Run(int32(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
