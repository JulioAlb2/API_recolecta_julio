package adapters

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresPuntoRecoleccion struct {
	conn *pgxpool.Pool
}

func NewPostgresPuntoRecoleccion() ports.IPuntoRecoleccion {
	conn, _ := core.ConnectPostgres()
	return &PostgresPuntoRecoleccion{conn: conn}
}

//
// SAVE
//
func (pg *PostgresPuntoRecoleccion) Save(p *entities.PuntoRecoleccion) (*entities.PuntoRecoleccion, error) {
	p.CreatedAt = time.Now()
	sql := `
	INSERT INTO punto_recoleccion
	(
		ruta_id,
		direccion,
		created_at
	)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	err := pg.conn.QueryRow(
		context.Background(),
		sql,
		p.RutaID,
		p.CP,
		p.CreatedAt,
	).Scan(&p.PuntoID)

	if err != nil {
		return nil, err
	}

	// Guardar coordenadas geográficas en Redis
	rdb, err := core.ConnectRedis()
	if err == nil {
		ctx := context.Background()
		rdb.HSet(ctx, fmt.Sprintf("point:%d", p.PuntoID), map[string]interface{}{
			"route_id": p.RutaID,
			"lat":      p.Lat,
			"lon":      p.Lon,
			"label":    p.CP,
		})
	}

	p.Eliminado = false
	return p, nil
}


//
// UPDATE
//
func (pg *PostgresPuntoRecoleccion) Update(id int32, p *entities.PuntoRecoleccion) (*entities.PuntoRecoleccion, error) {
	sql := `
	UPDATE punto_recoleccion
	SET
		ruta_id = $1,
		direccion = $2
	WHERE id = $3 AND deleted_at IS NULL
	`

	ct, err := pg.conn.Exec(
		context.Background(),
		sql,
		p.RutaID,
		p.CP,
		id,
	)

	if err != nil {
		return nil, err
	}

	if ct.RowsAffected() == 0 {
		return nil, errors.New("punto de recolección no encontrado")
	}

	// Actualizar coordenadas en Redis
	rdb, err := core.ConnectRedis()
	if err == nil {
		ctx := context.Background()
		rdb.HSet(ctx, fmt.Sprintf("point:%d", id), map[string]interface{}{
			"route_id": p.RutaID,
			"lat":      p.Lat,
			"lon":      p.Lon,
			"label":    p.CP,
		})
	}

	p.PuntoID = id
	return p, nil
}

//
// GET ALL
//
func (pg *PostgresPuntoRecoleccion) ListAll() ([]entities.PuntoRecoleccion, error) {
	sql := `
	SELECT id, ruta_id, direccion, (deleted_at IS NOT NULL) AS eliminado
	FROM punto_recoleccion
	WHERE deleted_at IS NULL
	ORDER BY id DESC
	`

	rows, err := pg.conn.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var puntos []entities.PuntoRecoleccion
	rdb, _ := core.ConnectRedis()

	for rows.Next() {
		var p entities.PuntoRecoleccion
		err := rows.Scan(&p.PuntoID, &p.RutaID, &p.CP, &p.Eliminado)
		if err != nil {
			return nil, err
		}

		// Hydrate coordinates from Redis
		if rdb != nil {
			vals, err := rdb.HGetAll(context.Background(), fmt.Sprintf("point:%d", p.PuntoID)).Result()
			if err == nil && len(vals) > 0 {
				p.Lat, _ = strconv.ParseFloat(vals["lat"], 64)
				p.Lon, _ = strconv.ParseFloat(vals["lon"], 64)
			}
		}

		puntos = append(puntos, p)
	}

	return puntos, nil
}

//
// GET BY ID
//
func (pg *PostgresPuntoRecoleccion) GetById(id int32) (*entities.PuntoRecoleccion, error) {
	var p entities.PuntoRecoleccion

	sql := `
	SELECT id, ruta_id, direccion, (deleted_at IS NOT NULL) AS eliminado
	FROM punto_recoleccion
	WHERE id = $1 AND deleted_at IS NULL
	`

	err := pg.conn.QueryRow(context.Background(), sql, id).Scan(
		&p.PuntoID,
		&p.RutaID,
		&p.CP,
		&p.Eliminado,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("punto de recolección no encontrado")
		}
		return nil, err
	}

	// Hydrate coordinates from Redis
	rdb, err := core.ConnectRedis()
	if err == nil {
		vals, err := rdb.HGetAll(context.Background(), fmt.Sprintf("point:%d", p.PuntoID)).Result()
		if err == nil && len(vals) > 0 {
			p.Lat, _ = strconv.ParseFloat(vals["lat"], 64)
			p.Lon, _ = strconv.ParseFloat(vals["lon"], 64)
		}
	}

	return &p, nil
}

//
// GET BY RUTA
//
func (pg *PostgresPuntoRecoleccion) GetByRuta(rutaId int32) ([]entities.PuntoRecoleccion, error) {
	sql := `
	SELECT id, ruta_id, direccion, (deleted_at IS NOT NULL) AS eliminado
	FROM punto_recoleccion
	WHERE ruta_id = $1 AND deleted_at IS NULL
	ORDER BY id
	`

	rows, err := pg.conn.Query(context.Background(), sql, rutaId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var puntos []entities.PuntoRecoleccion
	rdb, _ := core.ConnectRedis()

	for rows.Next() {
		var p entities.PuntoRecoleccion
		err := rows.Scan(&p.PuntoID, &p.RutaID, &p.CP, &p.Eliminado)
		if err != nil {
			return nil, err
		}

		// Hydrate coordinates from Redis
		if rdb != nil {
			vals, err := rdb.HGetAll(context.Background(), fmt.Sprintf("point:%d", p.PuntoID)).Result()
			if err == nil && len(vals) > 0 {
				p.Lat, _ = strconv.ParseFloat(vals["lat"], 64)
				p.Lon, _ = strconv.ParseFloat(vals["lon"], 64)
			}
		}

		puntos = append(puntos, p)
	}

	return puntos, nil
}

//
// DELETE (Soft delete)
//
func (pg *PostgresPuntoRecoleccion) Delete(id int32) error {
	sql := `
	UPDATE punto_recoleccion
	SET deleted_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL
	`

	ct, err := pg.conn.Exec(context.Background(), sql, id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return errors.New("punto de recolección no encontrado")
	}

	// Eliminar coordenadas de Redis
	rdb, err := core.ConnectRedis()
	if err == nil {
		rdb.Del(context.Background(), fmt.Sprintf("point:%d", id))
	}

	return nil
}