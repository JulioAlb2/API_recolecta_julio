// postgres_anomalia_repository.go
package infrastructure

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresAnomaliaRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresAnomaliaRepository() repositories.IAnomalia {
	pool := core.GetBD()
	return &PostgresAnomaliaRepository{pool: pool}
}

const anomaliaColumnas = `
	anomalia_id,
	tipo_anomalia,
	punto_id,
	conductor_id,
	camion_id,
	ruta_id,
	anomalia_referencia_id,
	descripcion,
	json_ruta,
	estado,
	eliminado,
	fecha_reporte,
	fecha_resolucion,
	created_at,
	updated_at,
	estado_pipeline,
	nivel_riesgo,
	inferencia_id,
	categoria_clasificada,
	subtipo_clasificado,
	accion_sugerida,
	pipeline_error
`

// scanner es satisfecho tanto por pgx.Row (QueryRow) como por pgx.Rows (Query + Next).
type scanner interface {
	Scan(dest ...interface{}) error
}

func scanAnomalia(s scanner, a *entities.Anomalia) error {
	return s.Scan(
		&a.AnomaliaID,
		&a.TipoAnomalia,
		&a.PuntoID,
		&a.ConductorID,
		&a.CamionID,
		&a.RutaID,
		&a.AnomaliaReferenciaID,
		&a.Descripcion,
		&a.JsonRuta,
		&a.Estado,
		&a.Eliminado,
		&a.FechaReporte,
		&a.FechaResolucion,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.EstadoPipeline,
		&a.NivelRiesgo,
		&a.InferenciaID,
		&a.CategoriaClasificada,
		&a.SubtipoClasificado,
		&a.AccionSugerida,
		&a.PipelineError,
	)
}

func (pg *PostgresAnomaliaRepository) collect(ctx context.Context, query string, args ...interface{}) ([]entities.Anomalia, error) {
	rows, err := pg.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var anomalias []entities.Anomalia
	for rows.Next() {
		var anomalia entities.Anomalia
		if err := scanAnomalia(rows, &anomalia); err != nil {
			log.Println("Error al escanear la anomalía:", err)
			return nil, err
		}
		anomalias = append(anomalias, anomalia)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error después de iterar filas:", err)
		return nil, err
	}

	return anomalias, nil
}

func (pg *PostgresAnomaliaRepository) Save(anomalia *entities.Anomalia) error {
	query := `
		INSERT INTO anomalia (
			tipo_anomalia,
			punto_id,
			conductor_id,
			camion_id,
			ruta_id,
			anomalia_referencia_id,
			descripcion,
			json_ruta,
			estado,
			eliminado,
			fecha_reporte,
			fecha_resolucion,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING anomalia_id
	`

	ctx := context.Background()

	var id int32
	err := pg.pool.QueryRow(
		ctx,
		query,
		anomalia.TipoAnomalia,
		anomalia.PuntoID,
		anomalia.ConductorID,
		anomalia.CamionID,
		anomalia.RutaID,
		anomalia.AnomaliaReferenciaID,
		anomalia.Descripcion,
		anomalia.JsonRuta,
		anomalia.Estado,
		anomalia.Eliminado,
		anomalia.FechaReporte,
		anomalia.FechaResolucion,
		anomalia.CreatedAt,
		anomalia.UpdatedAt,
	).Scan(&id)

	if err != nil {
		log.Println("Error al guardar la anomalía:", err)
		return err
	}

	anomalia.AnomaliaID = id
	return nil
}

func (pg *PostgresAnomaliaRepository) Update(anomalia *entities.Anomalia) error {
	query := `
		UPDATE anomalia
		SET
			tipo_anomalia = $1,
			punto_id = $2,
			conductor_id = $3,
			camion_id = $4,
			ruta_id = $5,
			anomalia_referencia_id = $6,
			descripcion = $7,
			json_ruta = $8,
			estado = $9,
			eliminado = $10,
			fecha_reporte = $11,
			fecha_resolucion = $12,
			updated_at = $13
		WHERE anomalia_id = $14
	`

	ctx := context.Background()
	cmdTag, err := pg.pool.Exec(
		ctx,
		query,
		anomalia.TipoAnomalia,
		anomalia.PuntoID,
		anomalia.ConductorID,
		anomalia.CamionID,
		anomalia.RutaID,
		anomalia.AnomaliaReferenciaID,
		anomalia.Descripcion,
		anomalia.JsonRuta,
		anomalia.Estado,
		anomalia.Eliminado,
		anomalia.FechaReporte,
		anomalia.FechaResolucion,
		time.Now(),
		anomalia.AnomaliaID,
	)

	if err != nil {
		log.Println("Error al actualizar la anomalía:", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("anomalía con ID %d no encontrada", anomalia.AnomaliaID)
	}

	return nil
}

// Delete realiza un borrado lógico (eliminado = true). Se conserva el
// registro porque otras filas (p. ej. SEGUIMIENTO_FALLA_CRITICA) pueden
// referenciarlo a través de anomalia_referencia_id.
func (pg *PostgresAnomaliaRepository) Delete(id int32) error {
	query := `
		UPDATE anomalia
		SET eliminado = true, updated_at = $1
		WHERE anomalia_id = $2
	`

	ctx := context.Background()
	cmdTag, err := pg.pool.Exec(ctx, query, time.Now(), id)
	if err != nil {
		log.Println("Error al eliminar la anomalía:", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("anomalía con ID %d no encontrada", id)
	}

	return nil
}

func (pg *PostgresAnomaliaRepository) GetByID(id int32) (*entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE anomalia_id = $1`

	ctx := context.Background()
	row := pg.pool.QueryRow(ctx, query, id)

	var anomalia entities.Anomalia
	err := scanAnomalia(row, &anomalia)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("anomalía con ID %d no encontrada", id)
		}
		log.Println("Error al buscar la anomalía por ID:", err)
		return nil, err
	}

	return &anomalia, nil
}

func (pg *PostgresAnomaliaRepository) GetAll() ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE eliminado = false ORDER BY fecha_reporte DESC`
	return pg.collect(context.Background(), query)
}

func (pg *PostgresAnomaliaRepository) GetByTipoAnomalia(tipoAnomalia entities.TipoAnomalia) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE tipo_anomalia = $1 AND eliminado = false ORDER BY fecha_reporte DESC`
	return pg.collect(context.Background(), query, tipoAnomalia)
}

func (pg *PostgresAnomaliaRepository) GetByEstado(estado string) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE estado = $1 AND eliminado = false ORDER BY fecha_reporte DESC`
	return pg.collect(context.Background(), query, estado)
}

func (pg *PostgresAnomaliaRepository) GetByPuntoID(puntoID int32) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE punto_id = $1 AND eliminado = false ORDER BY fecha_reporte DESC`
	return pg.collect(context.Background(), query, puntoID)
}

func (pg *PostgresAnomaliaRepository) GetByConductorID(conductorID int32) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE conductor_id = $1 AND eliminado = false ORDER BY fecha_reporte DESC`
	return pg.collect(context.Background(), query, conductorID)
}

func (pg *PostgresAnomaliaRepository) GetByCamionID(camionID int32) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE camion_id = $1 AND eliminado = false ORDER BY fecha_reporte DESC`
	return pg.collect(context.Background(), query, camionID)
}

func (pg *PostgresAnomaliaRepository) GetByRutaID(rutaID int32) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE ruta_id = $1 AND eliminado = false ORDER BY fecha_reporte DESC`
	return pg.collect(context.Background(), query, rutaID)
}

// GetByReferenciaID obtiene los registros que dan seguimiento a otra
// anomalía (uso típico: seguimientos de un REPORTE_FALLA_CRITICA).
func (pg *PostgresAnomaliaRepository) GetByReferenciaID(referenciaID int32) ([]entities.Anomalia, error) {
	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE anomalia_referencia_id = $1 AND eliminado = false ORDER BY fecha_reporte DESC`
	return pg.collect(context.Background(), query, referenciaID)
}

// ActualizarPipeline hace un UPDATE dirigido solo a las columnas del pipeline
// modelo_reportes -> clasificador_reportes (ver migrations/2026-07-22_pipeline_reportes_anomalia.sql).
// No pasa por scanAnomalia/entities.Anomalia a proposito: es un efecto
// secundario del pipeline en background, separado del CRUD normal de
// anomalias, para no arriesgar romper Save/Update/GetAll existentes.
// Los punteros nil dejan la columna correspondiente sin tocar.
func (pg *PostgresAnomaliaRepository) ActualizarPipeline(anomaliaID int32, estadoPipeline string, nivelRiesgo *string, inferenciaID *int32, categoria *string, subtipo *string, accion *string, pipelineError *string) error {
	setClauses := []string{"estado_pipeline = $1", "updated_at = $2"}
	args := []interface{}{estadoPipeline, time.Now()}
	nextParam := 3

	addField := func(column string, value interface{}) {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", column, nextParam))
		args = append(args, value)
		nextParam++
	}

	if nivelRiesgo != nil {
		addField("nivel_riesgo", *nivelRiesgo)
	}
	if inferenciaID != nil {
		addField("inferencia_id", *inferenciaID)
	}
	if categoria != nil {
		addField("categoria_clasificada", *categoria)
	}
	if subtipo != nil {
		addField("subtipo_clasificado", *subtipo)
	}
	if accion != nil {
		addField("accion_sugerida", *accion)
	}
	if pipelineError != nil {
		addField("pipeline_error", *pipelineError)
	}

	query := fmt.Sprintf("UPDATE anomalia SET %s WHERE anomalia_id = $%d", strings.Join(setClauses, ", "), nextParam)
	args = append(args, anomaliaID)

	ctx := context.Background()
	cmdTag, err := pg.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Println("Error al actualizar el pipeline de la anomalía:", err)
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("anomalía con ID %d no encontrada", anomaliaID)
	}

	return nil
}

func (pg *PostgresAnomaliaRepository) GetByFechaRange(fechaInicio, fechaFin string) ([]entities.Anomalia, error) {
	startTime, err := parseFecha(fechaInicio)
	if err != nil {
		return nil, fmt.Errorf("formato de fecha_inicio inválido: %v", err)
	}

	endTime, err := parseFecha(fechaFin)
	if err != nil {
		return nil, fmt.Errorf("formato de fecha_fin inválido: %v", err)
	}

	if len(fechaFin) <= 10 { // Solo fecha (YYYY-MM-DD): incluir todo el día final
		endTime = endTime.Add(24 * time.Hour)
	}

	query := `SELECT ` + anomaliaColumnas + ` FROM anomalia WHERE fecha_reporte BETWEEN $1 AND $2 AND eliminado = false ORDER BY fecha_reporte DESC`
	return pg.collect(context.Background(), query, startTime, endTime)
}
