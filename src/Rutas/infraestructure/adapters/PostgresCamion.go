package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresCamion struct {
	conn *pgxpool.Pool
}

func NewPostgresCamion() ports.ICamion {
	conn,_ := core.ConnectPostgres()
	return &PostgresCamion{
		conn: conn,
	}
}

func mapDisponibilidadToEstado(id int32) string {
	switch id {
	case 1:
		return "OPERATIVO"
	case 2:
		return "MANTENIMIENTO"
	case 3:
		return "FUERA_SERVICIO"
	case 4:
		return "BAJA"
	default:
		return "OPERATIVO"
	}
}

func mapEstadoToDisponibilidad(estado string) (int32, string, string) {
	switch estado {
	case "OPERATIVO":
		return 1, "OPERATIVO", "green"
	case "MANTENIMIENTO":
		return 2, "MANTENIMIENTO", "orange"
	case "FUERA_SERVICIO":
		return 3, "FUERA_SERVICIO", "red"
	case "BAJA":
		return 4, "BAJA", "grey"
	default:
		return 1, "OPERATIVO", "green"
	}
}

func (pg *PostgresCamion) Save(camion *entities.Camion) (*entities.Camion, error) {
	camion.CreatedAt = time.Now()
	estado := mapDisponibilidadToEstado(camion.DisponibilidadID)

	sql := `
	INSERT INTO camion
	(
		placa,
		modelo,
		tipo_id,
		rentado,
		estado,
		created_at,
		updated_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, NULL)
	RETURNING id
	`

	err := pg.conn.QueryRow(
		context.Background(),
		sql,
		camion.Placa,
		camion.Modelo,
		camion.TipoCamionID,
		camion.EsRentado,
		estado,
		camion.CreatedAt,
	).Scan(&camion.CamionID)

	if err != nil {
		return nil, err
	}

	dispID, nameDisp, colorDisp := mapEstadoToDisponibilidad(estado)
	camion.DisponibilidadID = dispID
	camion.NombreDisponibilidad = nameDisp
	camion.ColorDisponibilidad = colorDisp

	return camion, nil
}


func (pg *PostgresCamion) ListAll() ([]entities.Camion, error) {
	rows, err := pg.conn.Query(context.Background(),
		`SELECT id, placa, modelo, tipo_id, rentado,
		        estado,
		        created_at, updated_at
		 FROM camion
		 WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var camiones []entities.Camion

	for rows.Next() {
		var c entities.Camion
		var estado string
		var updatedAtNullable *time.Time
		err := rows.Scan(
			&c.CamionID,
			&c.Placa,
			&c.Modelo,
			&c.TipoCamionID,
			&c.EsRentado,
			&estado,
			&c.CreatedAt,
			&updatedAtNullable,
		)
		if err != nil {
			return nil, err
		}
		if updatedAtNullable != nil {
			c.UpdatedAt = *updatedAtNullable
		}
		dispID, nameDisp, colorDisp := mapEstadoToDisponibilidad(estado)
		c.DisponibilidadID = dispID
		c.NombreDisponibilidad = nameDisp
		c.ColorDisponibilidad = colorDisp
		camiones = append(camiones, c)
	}

	return camiones, nil
}

func (pg *PostgresCamion) Delete(id int32) error {
	cmd, err := pg.conn.Exec(context.Background(),
		`UPDATE camion SET deleted_at = NOW() WHERE id = $1`, id)

	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("camion no encontrado")
	}
	return nil
}

func (pg *PostgresCamion) GetByID(id int32) (*entities.Camion, error) {
	sql := `
	SELECT 
		id,
		placa,
		modelo,
		tipo_id,
		rentado,
		estado,
		created_at,
		updated_at
	FROM camion
	WHERE id = $1 AND deleted_at IS NULL
	`

	var camion entities.Camion
	var estado string
	var updatedAtNullable *time.Time

	err := pg.conn.QueryRow(context.Background(), sql, id).Scan(
		&camion.CamionID,
		&camion.Placa,
		&camion.Modelo,
		&camion.TipoCamionID,
		&camion.EsRentado,
		&estado,
		&camion.CreatedAt,
		&updatedAtNullable,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("camión no encontrado")
		}
		return nil, err
	}
	if updatedAtNullable != nil {
		camion.UpdatedAt = *updatedAtNullable
	}

	dispID, nameDisp, colorDisp := mapEstadoToDisponibilidad(estado)
	camion.DisponibilidadID = dispID
	camion.NombreDisponibilidad = nameDisp
	camion.ColorDisponibilidad = colorDisp

	return &camion, nil
}


func (pg *PostgresCamion) Update(id int32, camion *entities.Camion) (*entities.Camion, error) {
	camion.UpdatedAt = time.Now()
	estado := mapDisponibilidadToEstado(camion.DisponibilidadID)

	sql := `
	UPDATE camion
	SET 
		placa = $1,
		modelo = $2,
		tipo_id = $3,
		rentado = $4,
		estado = $5,
		updated_at = $6
	WHERE id = $7 AND deleted_at IS NULL
	`

	cmdTag, err := pg.conn.Exec(
		context.Background(),
		sql,
		camion.Placa,
		camion.Modelo,
		camion.TipoCamionID,
		camion.EsRentado,
		estado,
		camion.UpdatedAt, 
		id,
	)

	if err != nil {
		return nil, err
	}

	if cmdTag.RowsAffected() == 0 {
		return nil, errors.New("camión no encontrado")
	}

	dispID, nameDisp, colorDisp := mapEstadoToDisponibilidad(estado)
	camion.DisponibilidadID = dispID
	camion.NombreDisponibilidad = nameDisp
	camion.ColorDisponibilidad = colorDisp

	return camion, nil
}


func (pg *PostgresCamion) GetByPlaca(placa string) (*entities.Camion, error) {
	sql := `
	SELECT 
		id,
		placa,
		modelo,
		tipo_id,
		rentado,
		estado,
		created_at,
		updated_at
	FROM camion
	WHERE placa = $1 AND deleted_at IS NULL
	`

	var camion entities.Camion
	var estado string
	var updatedAtNullable *time.Time

	err := pg.conn.QueryRow(
		context.Background(),
		sql,
		placa,
	).Scan(
		&camion.CamionID,
		&camion.Placa,
		&camion.Modelo,
		&camion.TipoCamionID,
		&camion.EsRentado,
		&estado,
		&camion.CreatedAt,
		&updatedAtNullable,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("camión no encontrado")
		}
		return nil, err
	}
	if updatedAtNullable != nil {
		camion.UpdatedAt = *updatedAtNullable
	}

	dispID, nameDisp, colorDisp := mapEstadoToDisponibilidad(estado)
	camion.DisponibilidadID = dispID
	camion.NombreDisponibilidad = nameDisp
	camion.ColorDisponibilidad = colorDisp

	return &camion, nil
}

func (pg *PostgresCamion) GetByModelo(modelo string) ([]entities.Camion, error) {
	sql := `
	SELECT 
		id,
		placa,
		modelo,
		tipo_id,
		rentado,
		estado,
		created_at,
		updated_at
	FROM camion
	WHERE modelo ILIKE '%' || $1 || '%' AND deleted_at IS NULL
	`

	rows, err := pg.conn.Query(context.Background(), sql, modelo)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var camiones []entities.Camion

	for rows.Next() {
		var camion entities.Camion
		var estado string
		var updatedAtNullable *time.Time
		if err := rows.Scan(
			&camion.CamionID,
			&camion.Placa,
			&camion.Modelo,
			&camion.TipoCamionID,
			&camion.EsRentado,
			&estado,
			&camion.CreatedAt,
			&updatedAtNullable,
		); err != nil {
			return nil, err
		}
		if updatedAtNullable != nil {
			camion.UpdatedAt = *updatedAtNullable
		}

		dispID, nameDisp, colorDisp := mapEstadoToDisponibilidad(estado)
		camion.DisponibilidadID = dispID
		camion.NombreDisponibilidad = nameDisp
		camion.ColorDisponibilidad = colorDisp

		camiones = append(camiones, camion)
	}

	return camiones, nil
}