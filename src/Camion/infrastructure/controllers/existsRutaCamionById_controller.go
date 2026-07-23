package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
)

type ExistsRutaCamionByIDController struct {
	uc *application.ExistsRutaCamionByIDUseCase
}

func NewExistsRutaCamionByIDController(
	uc *application.ExistsRutaCamionByIDUseCase,
) *ExistsRutaCamionByIDController {
	return &ExistsRutaCamionByIDController{uc}
}

// @Summary      Verificar existencia ruta-camión
// @Tags         RutaCamion
// @Produce      json
// @Param        id path int true "ID"
// @Success      200 {object} entities.HistorialAsignacionCamionListResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/ruta-camion/exists/{id} [get]
func (c *ExistsRutaCamionByIDController) Run(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	exists, err := c.uc.Execute(int32(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":     id,
		"exists": exists,
	})
}
