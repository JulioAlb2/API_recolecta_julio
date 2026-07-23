package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
)

type GetRutaCamionByRutaIDController struct {
	uc *application.GetRutaCamionByRutaIDUseCase
}

func NewGetRutaCamionByRutaIDController(
	uc *application.GetRutaCamionByRutaIDUseCase,
) *GetRutaCamionByRutaIDController {
	return &GetRutaCamionByRutaIDController{uc}
}

// @Summary      Rutas-camión por ID de ruta
// @Tags         RutaCamion
// @Produce      json
// @Param        ruta_id path int true "ID ruta"
// @Success      200 {object} entities.HistorialAsignacionCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/ruta-camion/ruta/{ruta_id} [get]
func (c *GetRutaCamionByRutaIDController) Run(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("ruta_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Ruta ID inválido"})
		return
	}

	result, err := c.uc.Execute(int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
