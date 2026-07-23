package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Camion/domain/ports"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresHistorialAsignacionCamion struct {
	conn *pgxpool.Pool
}

func NewPostgresHistorialAsignacionCamion() ports.IHistorialAsignacionCamion {
	conn, _ := core.ConnectPostgres()
	return &PostgresHistorialAsignacionCamion{conn: conn}
}

//
// CREATE
//
func (pg *PostgresHistorialAsignacionCamion) Save(h *entities.HistorialAsignacionCamion) (*entities.HistorialAsignacionCamion, error) {
	h.CreatedAt = time.Now()
	sql := `
	INSERT INTO historial_asignacion_camion
	(
		id_chofer,
		id_camion,
		fecha_asignacion,
		fecha_baja,
		eliminado,
		created_at,
		updated_at
	)
	VALUES ($1, $2, $3, $4, false, $5, $5)
	RETURNING id_historial
	`

	err := pg.conn.QueryRow(
		context.Background(),
		sql,
		h.IDChofer,
		h.IDCamion,
		h.FechaAsignacion,
		h.FechaBaja,
		h.CreatedAt,
	).Scan(&h.IDHistorial)

	if err != nil {
		return nil, err
	}

	return h, nil
}



//
// GET BY ID
//
func (pg *PostgresHistorialAsignacionCamion) GetById(id int32) (*entities.HistorialAsignacionCamion, error) {
	var h entities.HistorialAsignacionCamion

	sql := `
	SELECT id_historial, id_chofer, id_camion, fecha_asignacion, fecha_baja, eliminado, created_at, updated_at
	FROM historial_asignacion_camion
	WHERE id_historial = $1
	`

	var updatedAtNullable *time.Time
	err := pg.conn.QueryRow(context.Background(), sql, id).Scan(
		&h.IDHistorial,
		&h.IDChofer,
		&h.IDCamion,
		&h.FechaAsignacion,
		&h.FechaBaja,
		&h.Eliminado,
		&h.CreatedAt,
		&updatedAtNullable,
	)
	if err == nil && updatedAtNullable != nil {
		h.UpdatedAt = *updatedAtNullable
	}

	if err == pgx.ErrNoRows {
		return nil, errors.New("historial no encontrado")
	}
	if err != nil {
		return nil, err
	}

	return &h, nil
}

//
// LIST ALL
//
func (pg *PostgresHistorialAsignacionCamion) ListAll() ([]entities.HistorialAsignacionCamion, error) {
	sql := `
	SELECT id_historial, id_chofer, id_camion, fecha_asignacion, fecha_baja, eliminado, created_at, updated_at
	FROM historial_asignacion_camion
	ORDER BY id_historial DESC
	`

	rows, err := pg.conn.Query(context.Background(), sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entities.HistorialAsignacionCamion

	for rows.Next() {
		var h entities.HistorialAsignacionCamion
		var updatedAtNullable *time.Time
		err := rows.Scan(
			&h.IDHistorial,
			&h.IDChofer,
			&h.IDCamion,
			&h.FechaAsignacion,
			&h.FechaBaja,
			&h.Eliminado,
			&h.CreatedAt,
			&updatedAtNullable,
		)
		if err != nil {
			return nil, err
		}
		if updatedAtNullable != nil {
			h.UpdatedAt = *updatedAtNullable
		}
		list = append(list, h)
	}

	return list, nil
}

//
// UPDATE
//
func (pg *PostgresHistorialAsignacionCamion) Update(id int32, h *entities.HistorialAsignacionCamion) (*entities.HistorialAsignacionCamion, error) {
	h.UpdatedAt = time.Now()
	sql := `
	UPDATE historial_asignacion_camion
	SET
		id_chofer = $1,
		id_camion = $2,
		fecha_baja = $3,
		eliminado = $4,
		updated_at = $5
	WHERE id_historial = $6
	`

	cmdTag, err := pg.conn.Exec(
		context.Background(),
		sql,
		h.IDChofer,
		h.IDCamion,
		h.FechaBaja,
		h.Eliminado,
		h.UpdatedAt, // 👈 tú la defines
		id,
	)

	if err != nil {
		return nil, err
	}

	if cmdTag.RowsAffected() == 0 {
		return nil, errors.New("historial no encontrado")
	}

	h.IDHistorial = int(id)
	return h, nil
}


//
// DELETE (Soft)
//
func (pg *PostgresHistorialAsignacionCamion) Delete(id int32) error {
	_, err := pg.conn.Exec(context.Background(),
		`UPDATE historial_asignacion_camion SET eliminado=true, updated_at=now() WHERE id_historial=$1`, id)
	return err
}

//
// ============ MÉTODOS AVANZADOS ============
//

func (pg *PostgresHistorialAsignacionCamion) GetByCamionId(camionId int32) ([]entities.HistorialAsignacionCamion, error) {
	sql := `
	SELECT id_historial, id_chofer, id_camion, fecha_asignacion, fecha_baja, eliminado, created_at, updated_at
	FROM historial_asignacion_camion
	WHERE id_camion=$1
	ORDER BY fecha_asignacion DESC
	`
	return pg.fetchMany(sql, camionId)
}

func (pg *PostgresHistorialAsignacionCamion) GetByChoferId(choferId int32) ([]entities.HistorialAsignacionCamion, error) {
	sql := `
	SELECT id_historial, id_chofer, id_camion, fecha_asignacion, fecha_baja, eliminado, created_at, updated_at
	FROM historial_asignacion_camion
	WHERE id_chofer=$1
	ORDER BY fecha_asignacion DESC
	`
	return pg.fetchMany(sql, choferId)
}

func (pg *PostgresHistorialAsignacionCamion) GetActivoByCamionId(camionId int32) (*entities.HistorialAsignacionCamion, error) {
	sql := `
	SELECT id_historial, id_chofer, id_camion, fecha_asignacion, fecha_baja, eliminado, created_at, updated_at
	FROM historial_asignacion_camion
	WHERE id_camion=$1 AND fecha_baja IS NULL AND eliminado=false
	LIMIT 1
	`
	return pg.fetchOne(sql, camionId)
}

func (pg *PostgresHistorialAsignacionCamion) GetActivoByChoferId(choferId int32) (*entities.HistorialAsignacionCamion, error) {
	sql := `
	SELECT id_historial, id_chofer, id_camion, fecha_asignacion, fecha_baja, eliminado, created_at, updated_at
	FROM historial_asignacion_camion
	WHERE id_chofer=$1 AND fecha_baja IS NULL AND eliminado=false
	LIMIT 1
	`
	return pg.fetchOne(sql, choferId)
}

func (pg *PostgresHistorialAsignacionCamion) DarDeBaja(id int32) error {
	_, err := pg.conn.Exec(context.Background(),
		`UPDATE historial_asignacion_camion SET fecha_baja=now(), updated_at=now() WHERE id_historial=$1`, id)
	return err
}

func (pg *PostgresHistorialAsignacionCamion) CerrarAsignacionActivaCamion(camionId int32) error {
	_, err := pg.conn.Exec(context.Background(),
		`UPDATE historial_asignacion_camion SET fecha_baja=now(), updated_at=now() WHERE id_camion=$1 AND fecha_baja IS NULL`, camionId)
	return err
}

func (pg *PostgresHistorialAsignacionCamion) CerrarAsignacionActivaChofer(choferId int32) error {
	_, err := pg.conn.Exec(context.Background(),
		`UPDATE historial_asignacion_camion SET fecha_baja=now(), updated_at=now() WHERE id_chofer=$1 AND fecha_baja IS NULL`, choferId)
	return err
}

//
// Helpers
//
func (pg *PostgresHistorialAsignacionCamion) fetchOne(sql string, param any) (*entities.HistorialAsignacionCamion, error) {
	var h entities.HistorialAsignacionCamion
	var updatedAtNullable *time.Time
	err := pg.conn.QueryRow(context.Background(), sql, param).Scan(
		&h.IDHistorial, &h.IDChofer, &h.IDCamion, &h.FechaAsignacion,
		&h.FechaBaja, &h.Eliminado, &h.CreatedAt, &updatedAtNullable,
	)
	if err == nil && updatedAtNullable != nil {
		h.UpdatedAt = *updatedAtNullable
	}
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("historial_asignacion_camion no encontrado")
		}

		return nil, err
	}
	return &h, nil
}

func (pg *PostgresHistorialAsignacionCamion) fetchMany(sql string, param any) ([]entities.HistorialAsignacionCamion, error) {
	rows, err := pg.conn.Query(context.Background(), sql, param)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []entities.HistorialAsignacionCamion
	for rows.Next() {
		var h entities.HistorialAsignacionCamion
		var updatedAtNullable *time.Time
		err := rows.Scan(
			&h.IDHistorial, &h.IDChofer, &h.IDCamion, &h.FechaAsignacion,
			&h.FechaBaja, &h.Eliminado, &h.CreatedAt, &updatedAtNullable,
		)
		if err != nil {
			return nil, err
		}
		if updatedAtNullable != nil {
			h.UpdatedAt = *updatedAtNullable
		}
		list = append(list, h)
	}
	return list, nil
}