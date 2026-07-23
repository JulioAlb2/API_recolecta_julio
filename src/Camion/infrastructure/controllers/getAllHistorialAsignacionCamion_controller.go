package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
)

type GetAllHistorialAsignacionCamionController struct {
	uc *application.ListAllHistorialAsignacionCamionUseCase
}

func NewGetAllHistorialAsignacionCamionController(uc *application.ListAllHistorialAsignacionCamionUseCase) *GetAllHistorialAsignacionCamionController {
	return &GetAllHistorialAsignacionCamionController{uc: uc}
}

// @Summary      Listar historiales de asignación
// @Tags         HistorialAsignacion
// @Produce      json
// @Success      200 {object} entities.HistorialAsignacionCamionListResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/historial-asignacion/ [get]
func (ctr *GetAllHistorialAsignacionCamionController) Run(ctx *gin.Context) {
	data, err := ctr.uc.Run()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
