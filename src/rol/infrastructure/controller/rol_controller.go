package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/rol/application"
)

type RolController struct {
	createRolUC *application.CreateRol
	listRolUC   *application.ListRol
	updateRolUC *application.UpdateRol
	deleteRolUC *application.DeleteRol
}

func NewRolController(
	createRol *application.CreateRol,
	listRol *application.ListRol,
	updateRol *application.UpdateRol,
	deleteRol *application.DeleteRol,
) *RolController {
	return &RolController{
		createRolUC: createRol,
		listRolUC:   listRol,
		updateRolUC: updateRol,
		deleteRolUC: deleteRol,
	}
}

func (c *RolController) RegisterRoutes(group *gin.RouterGroup) {
	group.POST("/roles", core.RequireRole(core.ADMIN), c.Create)
	group.GET("/roles", core.RequireRole(core.ADMIN), c.List)
	group.PUT("/roles/:id", core.RequireRole(core.ADMIN), c.Update)
	group.DELETE("/roles/:id", core.RequireRole(core.ADMIN), c.Delete)
}

// @Summary      Crear rol
// @Tags         Rol
// @Accept       json
// @Produce      json
// @Param        body body map[string]interface{} true "Nombre del rol"
// @Success      201 {object} map[string]interface{}
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/roles [post]
func (c *RolController) Create(ctx *gin.Context) {
	var req struct {
		Nombre string `json:"nombre" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.createRolUC.Execute(req.Nombre)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Rol creado"})
}

// @Summary      Listar roles
// @Tags         Rol
// @Produce      json
// @Success      200 {array} map[string]interface{}
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/roles [get]
func (c *RolController) List(ctx *gin.Context) {
	roles, err := c.listRolUC.Execute()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, roles)
}

// @Summary      Actualizar rol
// @Tags         Rol
// @Accept       json
// @Produce      json
// @Param        id path int true "ID rol"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/roles/{id} [put]
func (c *RolController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var req struct {
		Nombre string `json:"nombre" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = c.updateRolUC.Execute(id, req.Nombre)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Rol actualizado"})
}

// @Summary      Eliminar rol
// @Tags         Rol
// @Produce      json
// @Param        id path int true "ID rol"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/roles/{id} [delete]
func (c *RolController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	err = c.deleteRolUC.Execute(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Rol eliminado"})
}
