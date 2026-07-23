package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/application/recorrido"
	"github.com/vicpoo/API_recolecta/src/core"
)

type RecorridoController struct {
	store *recorrido.RedisStore
}

func NewRecorridoController(store *recorrido.RedisStore) *RecorridoController {
	return &RecorridoController{store: store}
}

type iniciarRecorridoRequest struct {
	RutaID   int32 `json:"ruta_id" binding:"required"`
	ChoferID int32 `json:"chofer_id" binding:"required"`
	CamionID int32 `json:"camion_id" binding:"required"`
}

// Iniciar POST /api/recorrido/iniciar
func (c *RecorridoController) Iniciar(ctx *gin.Context) {
	var req iniciarRecorridoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		core.RespondBadRequest(ctx, "json inválido", err.Error())
		return
	}

	if err := c.store.Iniciar(ctx.Request.Context(), req.ChoferID, req.CamionID, req.RutaID); err != nil {
		core.RespondInternalServerError(ctx, "no se pudo iniciar el recorrido", err)
		return
	}

	core.RespondCreated(ctx, gin.H{
		"iniciada":           true,
		"ruta_id":            req.RutaID,
		"chofer_id":          req.ChoferID,
		"camion_id":          req.CamionID,
		"punto_actual_index": 0,
		"pausado":            false,
	})
}

// Finalizar PUT /api/recorrido/finalizar
func (c *RecorridoController) Finalizar(ctx *gin.Context) {
	choferID := int32(ctx.GetInt("user_id"))
	if choferID == 0 {
		core.RespondBadRequest(ctx, "chofer no identificado", nil)
		return
	}

	if err := c.store.FinalizarByChofer(ctx.Request.Context(), choferID); err != nil {
		core.RespondInternalServerError(ctx, "no se pudo finalizar el recorrido", err)
		return
	}

	core.RespondOK(ctx, gin.H{"message": "recorrido finalizado"})
}

// Avanzar PUT /api/recorrido/avanzar
func (c *RecorridoController) Avanzar(ctx *gin.Context) {
	choferID := int32(ctx.GetInt("user_id"))
	if choferID == 0 {
		core.RespondBadRequest(ctx, "chofer no identificado", nil)
		return
	}

	if err := c.store.AvanzarByChofer(ctx.Request.Context(), choferID); err != nil {
		core.RespondBadRequest(ctx, err.Error(), nil)
		return
	}

	snap, _ := c.store.GetActivoByChofer(ctx.Request.Context(), choferID)
	core.RespondOK(ctx, gin.H{"data": snap})
}

// GetActivo GET /api/recorrido/activo
func (c *RecorridoController) GetActivo(ctx *gin.Context) {
	choferParam := ctx.Query("chofer_id")
	var choferID int32

	if choferParam != "" {
		id, err := strconv.Atoi(choferParam)
		if err != nil {
			core.RespondBadRequest(ctx, "chofer_id inválido", err.Error())
			return
		}
		choferID = int32(id)
	} else {
		choferID = int32(ctx.GetInt("user_id"))
	}

	if choferID == 0 {
		core.RespondOK(ctx, gin.H{"data": nil})
		return
	}

	snap, err := c.store.GetActivoByChofer(ctx.Request.Context(), choferID)
	if err != nil {
		core.RespondInternalServerError(ctx, "error consultando recorrido", err)
		return
	}

	if snap == nil {
		core.RespondOK(ctx, gin.H{"data": nil})
		return
	}

	core.RespondOK(ctx, gin.H{"data": snap})
}

// GetActivoPublic GET /api/recorrido/activo/publico?camion_id|chofer_id|ruta_id
func (c *RecorridoController) GetActivoPublic(ctx *gin.Context) {
	if rutaParam := ctx.Query("ruta_id"); rutaParam != "" {
		id, err := strconv.Atoi(rutaParam)
		if err != nil {
			core.RespondBadRequest(ctx, "ruta_id inválido", err.Error())
			return
		}
		snap, err := c.store.GetByRuta(ctx.Request.Context(), int32(id))
		if err != nil {
			core.RespondInternalServerError(ctx, "error consultando recorrido", err)
			return
		}
		core.RespondOK(ctx, gin.H{"data": snap})
		return
	}

	if camionParam := ctx.Query("camion_id"); camionParam != "" {
		id, err := strconv.Atoi(camionParam)
		if err != nil {
			core.RespondBadRequest(ctx, "camion_id inválido", err.Error())
			return
		}
		snap, err := c.store.GetByCamion(ctx.Request.Context(), int32(id))
		if err != nil {
			core.RespondInternalServerError(ctx, "error consultando recorrido", err)
			return
		}
		core.RespondOK(ctx, gin.H{"data": snap})
		return
	}

	c.GetActivo(ctx)
}

var _ = http.StatusOK
