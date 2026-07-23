package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
)

type GetAllTipoCamionController struct {
	uc *application.ListAllTipoCamionUseCase
}

func NewGetAllTipoCamionController(
	uc *application.ListAllTipoCamionUseCase,
) *GetAllTipoCamionController {
	return &GetAllTipoCamionController{
		uc: uc,
	}
}

// @Summary      Listar tipos de camión
// @Tags         TipoCamion
// @Produce      json
// @Success      200 {object} entities.HistorialAsignacionCamionListResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/tipo-camion/ [get]
func (ctr *GetAllTipoCamionController) Run(ctx *gin.Context) {
	tiposCamion, err := ctr.uc.Run()

	if err != nil {
		fmt.Printf("error to get tipos camion: %s\n", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "error al obtener los tipos de camion",
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tiposCamion,
	})
}
