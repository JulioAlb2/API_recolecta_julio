package infrastructure

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/vicpoo/API_recolecta/src/core"
)

const recorridoKey = "recorrido:activo"

const recorridoTTL = 6 * time.Hour

type Recorrido struct {
	RutaID           int    `json:"ruta_id"`
	ChoferID         int    `json:"chofer_id"`
	PuntoActualIndex int    `json:"punto_actual_index"`
	Iniciada         bool   `json:"iniciada"`
	UpdatedAt        string `json:"updated_at"`
}

type RecorridoRouter struct {
	engine *gin.Engine
}

func NewRecorridoRouter(engine *gin.Engine) *RecorridoRouter {
	return &RecorridoRouter{engine: engine}
}

func redisSafe() (*redis.Client, error) {
	return core.ConnectRedis()
}

func (r *RecorridoRouter) leer(ctx context.Context, rdb *redis.Client) (*Recorrido, error) {
	val, err := rdb.Get(ctx, recorridoKey).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var rec Recorrido
	if err := json.Unmarshal([]byte(val), &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *RecorridoRouter) guardar(ctx context.Context, rdb *redis.Client, rec *Recorrido) error {
	rec.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	b, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return rdb.Set(ctx, recorridoKey, b, recorridoTTL).Err()
}

func (r *RecorridoRouter) Run() {
	g := r.engine.Group("/api/recorrido")
	{
		g.GET("/activo", r.handleActivo)
		g.POST("/iniciar", r.handleIniciar)
		g.PUT("/avanzar", r.handleAvanzar)
		g.PUT("/finalizar", r.handleFinalizar)
	}
}

func (r *RecorridoRouter) handleActivo(c *gin.Context) {
	rdb, err := redisSafe()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": "redis no disponible"})
		return
	}
	rec, err := r.leer(c.Request.Context(), rdb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rec})
}

func (r *RecorridoRouter) handleIniciar(c *gin.Context) {
	var body struct {
		RutaID   int `json:"ruta_id"`
		ChoferID int `json:"chofer_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	rdb, err := redisSafe()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": "redis no disponible"})
		return
	}

	rec := &Recorrido{
		RutaID:           body.RutaID,
		ChoferID:         body.ChoferID,
		PuntoActualIndex: 0,
		Iniciada:         true,
	}
	if err := r.guardar(c.Request.Context(), rdb, rec); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rec})
}

func (r *RecorridoRouter) handleAvanzar(c *gin.Context) {
	rdb, err := redisSafe()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": "redis no disponible"})
		return
	}

	ctx := c.Request.Context()
	rec, err := r.leer(ctx, rdb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if rec == nil || !rec.Iniciada {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "no hay recorrido activo"})
		return
	}

	rec.PuntoActualIndex++
	if err := r.guardar(ctx, rdb, rec); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rec})
}

func (r *RecorridoRouter) handleFinalizar(c *gin.Context) {
	rdb, err := redisSafe()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"success": false, "message": "redis no disponible"})
		return
	}

	ctx := c.Request.Context()
	rec, err := r.leer(ctx, rdb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	if rec == nil {
		rec = &Recorrido{}
	}
	rec.Iniciada = false
	rec.PuntoActualIndex = 0
	if err := r.guardar(ctx, rdb, rec); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": rec})
}
