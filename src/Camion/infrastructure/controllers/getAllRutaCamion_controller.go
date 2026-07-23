package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
)

type GetAllRutaCamionController struct {
	uc *application.ListAllRutaCamionUseCase
}

func NewGetAllRutaCamionController(
	uc *application.ListAllRutaCamionUseCase,
) *GetAllRutaCamionController {
	return &GetAllRutaCamionController{uc}
}

// @Summary      Listar rutas-camión
// @Tags         RutaCamion
// @Produce      json
// @Success      200 {object} entities.HistorialAsignacionCamionListResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/ruta-camion/ [get]
func (c *GetAllRutaCamionController) Run(ctx *gin.Context) {
	result, err := c.uc.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
