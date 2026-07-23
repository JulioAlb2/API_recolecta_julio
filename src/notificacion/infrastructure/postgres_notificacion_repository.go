// postgres_notificacion_repository.go
package infrastructure

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/core"
	repositories "github.com/vicpoo/API_recolecta/src/notificacion/domain"
	"github.com/vicpoo/API_recolecta/src/notificacion/domain/entities"
)

type PostgresNotificacionRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresNotificacionRepository() repositories.INotificacion {
	pool := core.GetBD()
	return &PostgresNotificacionRepository{pool: pool}
}

// Proyección de la tabla anomalia al formato de entities.Notificacion
const querySelectNotificacion = `
	SELECT 
		anomalia_id AS notificacion_id,
		conductor_id AS usuario_id,
		tipo_anomalia AS tipo,
		'Notificación de ' || tipo_anomalia AS titulo,
		descripcion AS mensaje,
		COALESCE(estado != 'RESUELTA', true) AS activa,
		camion_id AS id_camion_relacionado,
		NULL::integer AS id_falla_relacionado,
		NULL::integer AS id_mantenimiento_relacionado,
		conductor_id AS creado_por,
		fecha_reporte AS created_at
	FROM anomalia
`

func (pg *PostgresNotificacionRepository) executeListQuery(query string, args ...interface{}) ([]entities.Notificacion, error) {
	ctx := context.Background()
	rows, err := pg.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notificaciones []entities.Notificacion
	for rows.Next() {
		var n entities.Notificacion
		err := rows.Scan(
			&n.NotificacionID,
			&n.UsuarioID,
			&n.Tipo,
			&n.Titulo,
			&n.Mensaje,
			&n.Activa,
			&n.IDCamionRelacionado,
			&n.IDFallaRelacionado,
			&n.IDMantenimientoRelacionado,
			&n.CreadoPor,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		notificaciones = append(notificaciones, n)
	}
	return notificaciones, nil
}

func (pg *PostgresNotificacionRepository) GetByID(id int32) (*entities.Notificacion, error) {
	query := querySelectNotificacion + " WHERE anomalia_id = $1"
	ctx := context.Background()
	row := pg.pool.QueryRow(ctx, query, id)

	var n entities.Notificacion
	err := row.Scan(
		&n.NotificacionID,
		&n.UsuarioID,
		&n.Tipo,
		&n.Titulo,
		&n.Mensaje,
		&n.Activa,
		&n.IDCamionRelacionado,
		&n.IDFallaRelacionado,
		&n.IDMantenimientoRelacionado,
		&n.CreadoPor,
		&n.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("notificación con ID %d no encontrada", id)
		}
		return nil, err
	}
	return &n, nil
}

func (pg *PostgresNotificacionRepository) GetAll() ([]entities.Notificacion, error) {
	query := querySelectNotificacion + " WHERE eliminado = false ORDER BY fecha_reporte DESC"
	return pg.executeListQuery(query)
}

func (pg *PostgresNotificacionRepository) GetActivas() ([]entities.Notificacion, error) {
	query := querySelectNotificacion + " WHERE COALESCE(estado != 'RESUELTA', true) AND eliminado = false ORDER BY fecha_reporte DESC"
	return pg.executeListQuery(query)
}

func (pg *PostgresNotificacionRepository) GetInactivas() ([]entities.Notificacion, error) {
	query := querySelectNotificacion + " WHERE estado = 'RESUELTA' AND eliminado = false ORDER BY fecha_reporte DESC"
	return pg.executeListQuery(query)
}

func (pg *PostgresNotificacionRepository) GetGlobales() ([]entities.Notificacion, error) {
	query := querySelectNotificacion + " WHERE conductor_id IS NULL AND eliminado = false ORDER BY fecha_reporte DESC"
	return pg.executeListQuery(query)
}

func (pg *PostgresNotificacionRepository) GetByUsuarioID(usuarioID int32) ([]entities.Notificacion, error) {
	query := querySelectNotificacion + " WHERE conductor_id = $1 AND eliminado = false ORDER BY fecha_reporte DESC"
	return pg.executeListQuery(query, usuarioID)
}

func (pg *PostgresNotificacionRepository) GetActivasByUsuarioID(usuarioID int32) ([]entities.Notificacion, error) {
	query := querySelectNotificacion + " WHERE conductor_id = $1 AND COALESCE(estado != 'RESUELTA', true) AND eliminado = false ORDER BY fecha_reporte DESC"
	return pg.executeListQuery(query, usuarioID)
}

func (pg *PostgresNotificacionRepository) GetByTipo(tipo string) ([]entities.Notificacion, error) {
	query := querySelectNotificacion + " WHERE tipo_anomalia = $1 AND eliminado = false ORDER BY fecha_reporte DESC"
	return pg.executeListQuery(query, tipo)
}

func (pg *PostgresNotificacionRepository) GetByUsuarioYTipo(usuarioID int32, tipo string) ([]entities.Notificacion, error) {
	query := querySelectNotificacion + " WHERE conductor_id = $1 AND tipo_anomalia = $2 AND eliminado = false ORDER BY fecha_reporte DESC"
	return pg.executeListQuery(query, usuarioID, tipo)
}

func (pg *PostgresNotificacionRepository) GetByCamionID(camionID int32) ([]entities.Notificacion, error) {
	query := querySelectNotificacion + " WHERE camion_id = $1 AND eliminado = false ORDER BY fecha_reporte DESC"
	return pg.executeListQuery(query, camionID)
}

func (pg *PostgresNotificacionRepository) GetByCamionYTipo(camionID int32, tipo string) ([]entities.Notificacion, error) {
	query := querySelectNotificacion + " WHERE camion_id = $1 AND tipo_anomalia = $2 AND eliminado = false ORDER BY fecha_reporte DESC"
	return pg.executeListQuery(query, camionID, tipo)
}

func (pg *PostgresNotificacionRepository) GetByFechaRange(fechaInicio, fechaFin string) ([]entities.Notificacion, error) {
	startTime, errStart := time.Parse("2006-01-02", fechaInicio)
	endTime, errEnd := time.Parse("2006-01-02", fechaFin)
	if errStart != nil || errEnd != nil {
		return nil, fmt.Errorf("formato de fecha inválido. Usar YYYY-MM-DD")
	}
	endTime = endTime.Add(24 * time.Hour)

	query := querySelectNotificacion + " WHERE fecha_reporte BETWEEN $1 AND $2 AND eliminado = false ORDER BY fecha_reporte DESC"
	return pg.executeListQuery(query, startTime, endTime)
}

func (pg *PostgresNotificacionRepository) MarcarComoLeida(id int32) error {
	query := "UPDATE anomalia SET estado = 'RESUELTA', fecha_resolucion = NOW() WHERE anomalia_id = $1"
	_, err := pg.pool.Exec(context.Background(), query, id)
	return err
}

func (pg *PostgresNotificacionRepository) MarcarComoActiva(id int32) error {
	query := "UPDATE anomalia SET estado = 'PENDIENTE', fecha_resolucion = NULL WHERE anomalia_id = $1"
	_, err := pg.pool.Exec(context.Background(), query, id)
	return err
}

func (pg *PostgresNotificacionRepository) MarcarTodasComoLeidas(usuarioID int32) error {
	query := "UPDATE anomalia SET estado = 'RESUELTA', fecha_resolucion = NOW() WHERE conductor_id = $1 AND COALESCE(estado != 'RESUELTA', true)"
	_, err := pg.pool.Exec(context.Background(), query, usuarioID)
	return err
}

func (pg *PostgresNotificacionRepository) CountActivasByUsuarioID(usuarioID int32) (int64, error) {
	query := `SELECT COUNT(*) FROM anomalia WHERE conductor_id = $1 AND COALESCE(estado != 'RESUELTA', true) AND eliminado = false`
	var count int64
	err := pg.pool.QueryRow(context.Background(), query, usuarioID).Scan(&count)
	return count, err
}

func (pg *PostgresNotificacionRepository) CountByUsuarioID(usuarioID int32) (int64, error) {
	query := `SELECT COUNT(*) FROM anomalia WHERE conductor_id = $1 AND eliminado = false`
	var count int64
	err := pg.pool.QueryRow(context.Background(), query, usuarioID).Scan(&count)
	return count, err
}

func (pg *PostgresNotificacionRepository) CountByTipo(tipo string) (int64, error) {
	query := `SELECT COUNT(*) FROM anomalia WHERE tipo_anomalia = $1 AND eliminado = false`
	var count int64
	err := pg.pool.QueryRow(context.Background(), query, tipo).Scan(&count)
	return count, err
}

func (pg *PostgresNotificacionRepository) CountByCamionID(camionID int32) (int64, error) {
	query := `SELECT COUNT(*) FROM anomalia WHERE camion_id = $1 AND eliminado = false`
	var count int64
	err := pg.pool.QueryRow(context.Background(), query, camionID).Scan(&count)
	return count, err
}