package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application"
)

type GetCamionByIDController struct {
	uc *application.GetCamionByIDUseCase
}

func NewGetCamionByIDController(uc *application.GetCamionByIDUseCase) *GetCamionByIDController {
	return &GetCamionByIDController{
		uc: uc,
	}
}

// @Summary      Obtener camión por ID
// @Tags         Camion
// @Produce      json
// @Param        id path int true "ID camión"
// @Success      200 {object} entities.EstadoCamionResponse
// @Failure      404 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/camion/{id} [get]
func (ctr *GetCamionByIDController) Run(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id inválido",
		})
		return
	}

	camion, err := ctr.uc.Run(int32(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, camion)
}
