package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
)

type DeleteTipoCamionController struct {
	uc *application.DeleteTipoCamionUseCase
}

func NewDeleteTipoCamionController(
	uc *application.DeleteTipoCamionUseCase,
) *DeleteTipoCamionController {
	return &DeleteTipoCamionController{
		uc: uc,
	}
}

// @Summary      Eliminar tipo de camión
// @Tags         TipoCamion
// @Produce      json
// @Param        id path int true "ID tipo camión"
// @Success      200 {object} entities.HistorialAsignacionCamionMessageResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/tipo-camion/{id} [delete]
func (ctr *DeleteTipoCamionController) Run(ctx *gin.Context) {
	idParam := ctx.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id inválido",
		})
		return
	}

	err = ctr.uc.Run(int32(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "tipo camion eliminado correctamente",
	})
}
