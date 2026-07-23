package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Camion/application"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
)


type CreateTipoCamionController struct {
	uc *application.SaveTipoCamionUseCase
}

func NewCreateTipoCamionController(uc *application.SaveTipoCamionUseCase) *CreateTipoCamionController {
	return &CreateTipoCamionController{
		uc: uc,
	}
}

// @Summary      Crear tipo de camión
// @Tags         TipoCamion
// @Accept       json
// @Produce      json
// @Param        body body entities.CreateHistorialAsignacionCamionRequest true "Body"
// @Success      201 {object} entities.HistorialAsignacionCamionResponse
// @Failure      400 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/tipo-camion/ [post]
func (ctr *CreateTipoCamionController) Run(ctx *gin.Context) {
	var tipoCamion entities.TipoCamion

	if err := ctx.ShouldBindBodyWithJSON(&tipoCamion); err != nil {
		fmt.Printf("error to create product  %s", err)
		ctx.JSON(201, gin.H{
			"message": "error al crear el tipo de camion",
			"error": err.Error(), 
			"success": false, 
		})
		return 
	}

	tipoCamionCreated, errC := ctr.uc.Run(&tipoCamion)

	if errC != nil {
		fmt.Printf("error to create tipo camion %s", errC)
		ctx.JSON(500, gin.H{
			"message": "error to create tipo camion", 
			"error": errC.Error(),
		})
	}

	ctx.JSON(201, tipoCamionCreated)
}
