package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/colonia/application"
	"github.com/vicpoo/API_recolecta/src/colonia/domain"
	"github.com/vicpoo/API_recolecta/src/core"
)

type ColoniaController struct {
	create *application.CreateColonia
	get    *application.GetColonia
	list   *application.ListColonias
	update *application.UpdateColonia
	delete *application.DeleteColonia
}

func NewColoniaController(
	create *application.CreateColonia,
	get *application.GetColonia,
	list *application.ListColonias,
	update *application.UpdateColonia,
	delete *application.DeleteColonia,
) *ColoniaController {
	return &ColoniaController{create, get, list, update, delete}
}

// tenantIDFromContext lee el tenant_id que el middleware JWT ya dejo en el
// contexto de Gin. Solo se usa en rutas ya protegidas con JWTAuthMiddleware;
// el chequeo aqui es una garantia extra de que nunca se escribe sin tenant.
func tenantIDFromContext(ctx *gin.Context) (int, bool) {
	val, exists := ctx.Get("tenant_id")
	if !exists {
		return 0, false
	}
	tenantID, ok := val.(int)
	return tenantID, ok
}

func (c *ColoniaController) RegisterRoutes(r *gin.Engine) {

	public := r.Group("/api/colonia")
	{
		public.GET("", c.List)
		public.GET("/:id", c.GetByID)
	}

	// Rutas protegidas SOLO ADMIN
	admin := r.Group(
		"/api/colonia",
		core.JWTAuthMiddleware(),
		core.RequireRole(core.ADMIN),
	)
	{
		admin.POST("", c.Create)
		admin.PUT("/:id", c.Update)
		admin.DELETE("/:id", c.Delete)
	}
}

// @Summary      Crear colonia
// @Description  Crea una nueva colonia. Solo administradores (rol ADMIN)
// @Tags         Colonia
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body domain.CreateColoniaRequest true "Body"
// @Success      201 {object} domain.ColoniaResponse
// @Failure      400 {object} core.ErrorResponse "Datos inválidos"
// @Failure      401 {object} core.ErrorResponse "No autenticado"
// @Failure      403 {object} core.ErrorResponse "No autorizado (requiere rol ADMIN)"
// @Failure      500 {object} core.ErrorResponse "Error interno del servidor"
// @Router       /api/colonia [post]
func (c *ColoniaController) Create(ctx *gin.Context) {
	tenantID, ok := tenantIDFromContext(ctx)
	if !ok {
		core.RespondError(ctx, http.StatusUnauthorized, core.ErrCodeUnauthorized, "tenant no encontrado en token", nil)
		return
	}

	var request domain.CreateColoniaRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		core.RespondValidationError(ctx, "Datos de colonia inválidos", map[string]string{"error": err.Error()})
		return
	}

	body := domain.Colonia{
		Nombre: request.Nombre,
		Zona:   request.Zona,
	}

	if err := c.create.Execute(ctx.Request.Context(), tenantID, &body); err != nil {
		core.RespondInternalServerError(ctx, "Error al crear colonia", err)
		return
	}

	ctx.JSON(http.StatusCreated, domain.ColoniaResponse{
		Success: true,
		Message: "Colonia creada correctamente",
		Data:    body,
		Code:    http.StatusCreated,
	})
}

// @Summary      Obtener colonia por ID
// @Description  Obtiene los detalles de una colonia específica. Endpoint público
// @Tags         Colonia
// @Produce      json
// @Param        id path int true "ID de la colonia"
// @Success      200 {object} domain.ColoniaResponse
// @Failure      400 {object} core.ErrorResponse "ID inválido"
// @Failure      404 {object} core.ErrorResponse "Colonia no encontrada"
// @Failure      500 {object} core.ErrorResponse "Error interno del servidor"
// @Security     BearerAuth
// @Router       /api/colonia/{id} [get]
func (c *ColoniaController) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "ID de colonia inválido")
		return
	}

	colonia, err := c.get.Execute(id)
	if err != nil {
		core.RespondError(ctx, http.StatusNotFound, core.ErrCodeNotFound, "Colonia no encontrada", nil)
		return
	}

	ctx.JSON(http.StatusOK, domain.ColoniaResponse{
		Success: true,
		Message: "Colonia obtenida correctamente",
		Data:    *colonia,
		Code:    http.StatusOK,
	})
}

// @Summary      Listar colonias
// @Description  Obtiene el listado de todas las colonias disponibles. Endpoint público
// @Tags         Colonia
// @Produce      json
// @Success      200 {object} domain.ColoniaResponse
// @Failure      500 {object} core.ErrorResponse "Error interno del servidor"
// @Security     BearerAuth
// @Router       /api/colonia [get]
func (c *ColoniaController) List(ctx *gin.Context) {
	colonias, err := c.list.Execute()
	if err != nil {
		core.RespondInternalServerError(ctx, "Error al listar colonias", err)
		return
	}

	ctx.JSON(http.StatusOK, domain.ColoniaListResponse{
		Success: true,
		Message: "Colonias listadas correctamente",
		Data:    colonias,
		Code:    http.StatusOK,
	})
}

// @Summary      Actualizar colonia
// @Description  Actualiza los datos de una colonia existente. Solo administradores (rol ADMIN)
// @Tags         Colonia
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID de la colonia"
// @Param        body body domain.CreateColoniaRequest true "Body"
// @Success      200 {object} domain.ColoniaResponse
// @Failure      400 {object} core.ErrorResponse "Datos inválidos"
// @Failure      401 {object} core.ErrorResponse "No autenticado"
// @Failure      403 {object} core.ErrorResponse "No autorizado (requiere rol ADMIN)"
// @Failure      500 {object} core.ErrorResponse "Error interno del servidor"
// @Router       /api/colonia/{id} [put]
func (c *ColoniaController) Update(ctx *gin.Context) {
	tenantID, ok := tenantIDFromContext(ctx)
	if !ok {
		core.RespondError(ctx, http.StatusUnauthorized, core.ErrCodeUnauthorized, "tenant no encontrado en token", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "ID de colonia inválido")
		return
	}

	var request domain.UpdateColoniaRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		core.RespondValidationError(ctx, "Datos de colonia inválidos", map[string]string{"error": err.Error()})
		return
	}

	body := domain.Colonia{
		ColoniaID: id,
		Nombre:    request.Nombre,
		Zona:      request.Zona,
	}

	if err := c.update.Execute(ctx.Request.Context(), tenantID, &body); err != nil {
		core.RespondInternalServerError(ctx, "Error al actualizar colonia", err)
		return
	}

	ctx.JSON(http.StatusOK, domain.ColoniaMessageResponse{
		Success: true,
		Message: "Colonia actualizada correctamente",
		Code:    http.StatusOK,
	})
}

// @Summary      Eliminar colonia
// @Description  Elimina una colonia de la base de datos. Solo administradores (rol ADMIN)
// @Tags         Colonia
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID de la colonia"
// @Success      204 "Colonia eliminada correctamente"
// @Failure      400 {object} core.ErrorResponse "ID inválido"
// @Failure      401 {object} core.ErrorResponse "No autenticado"
// @Failure      403 {object} core.ErrorResponse "No autorizado (requiere rol ADMIN)"
// @Failure      500 {object} core.ErrorResponse "Error interno del servidor"
// @Router       /api/colonia/{id} [delete]
func (c *ColoniaController) Delete(ctx *gin.Context) {
	tenantID, ok := tenantIDFromContext(ctx)
	if !ok {
		core.RespondError(ctx, http.StatusUnauthorized, core.ErrCodeUnauthorized, "tenant no encontrado en token", nil)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		core.RespondInvalidInput(ctx, "ID de colonia inválido")
		return
	}

	if err := c.delete.Execute(ctx.Request.Context(), tenantID, id); err != nil {
		core.RespondInternalServerError(ctx, "Error al eliminar colonia", err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
