package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
	"github.com/vicpoo/API_recolecta/src/core"
)

type GetAllPuntoRecoleccionController struct {
	uc *application.ListAllPuntoRecoleccionUseCase
}

func NewGetAllPuntoRecoleccionController(uc *application.ListAllPuntoRecoleccionUseCase) *GetAllPuntoRecoleccionController {
	return &GetAllPuntoRecoleccionController{uc: uc}
}

// @Summary      Listar puntos de recolección
// @Tags         PuntoRecoleccion
// @Produce      json
// @Success      200 {object} entities.PuntoRecoleccionListResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/puntos-recoleccion/ [get]
func (c *GetAllPuntoRecoleccionController) Run(ctx *gin.Context) {
	result, err := c.uc.Execute()
	if err != nil {
		core.RespondInternalServerError(ctx, "No se pudieron obtener los puntos de recolección", err)
		return
	}

	core.RespondOK(ctx, gin.H{
		"data": result,
	})
}
