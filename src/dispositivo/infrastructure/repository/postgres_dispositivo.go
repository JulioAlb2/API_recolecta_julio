package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/dispositivo/domain/entities"
)

type PostgresDispositivoRepository struct {
	db *pgxpool.Pool
}

func NewPostgresDispositivoRepository(db *pgxpool.Pool) *PostgresDispositivoRepository {
	return &PostgresDispositivoRepository{db: db}
}

func (r *PostgresDispositivoRepository) Solicitar(ctx context.Context, d *entities.Dispositivo) error {
	query := `
		INSERT INTO dispositivos (conductor_id, mac_address, serial_number, api_key, nombre_dispositivo, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, FALSE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (conductor_id) 
		DO UPDATE SET 
			mac_address = EXCLUDED.mac_address,
			serial_number = EXCLUDED.serial_number,
			api_key = EXCLUDED.api_key,
			nombre_dispositivo = EXCLUDED.nombre_dispositivo,
			active = FALSE,
			updated_at = CURRENT_TIMESTAMP,
			deleted_at = NULL
	`
	_, err := r.db.Exec(ctx, query, d.ConductorID, d.MacAddress, d.SerialNumber, d.ApiKey, d.NombreDispositivo)
	return err
}

func (r *PostgresDispositivoRepository) FindByConductorID(ctx context.Context, conductorID int) (*entities.Dispositivo, error) {
	query := `
		SELECT id, conductor_id, mac_address, serial_number, api_key, nombre_dispositivo, sistema_operativo, active, created_at, updated_at, deleted_at
		FROM dispositivos
		WHERE conductor_id = $1 AND deleted_at IS NULL
	`
	var d entities.Dispositivo
	err := r.db.QueryRow(ctx, query, conductorID).Scan(
		&d.ID, &d.ConductorID, &d.MacAddress, &d.SerialNumber, &d.ApiKey, &d.NombreDispositivo, &d.SistemaOperativo, &d.Active, &d.CreatedAt, &d.UpdatedAt, &d.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &d, nil
}

func (r *PostgresDispositivoRepository) Aprobar(ctx context.Context, conductorID int) error {
	query := `
		UPDATE dispositivos
		SET active = TRUE, updated_at = CURRENT_TIMESTAMP
		WHERE conductor_id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, conductorID)
	return err
}

func (r *PostgresDispositivoRepository) Desvincular(ctx context.Context, conductorID int) error {
	query := `
		UPDATE dispositivos
		SET deleted_at = CURRENT_TIMESTAMP, active = FALSE
		WHERE conductor_id = $1 AND deleted_at IS NULL
	`
	_, err := r.db.Exec(ctx, query, conductorID)
	return err
}

func (r *PostgresDispositivoRepository) ListarPendientes(ctx context.Context) ([]*entities.DispositivoConductorResponse, error) {
	query := `
		SELECT d.id, d.conductor_id, e.nombre, e.apellidos, e.mail, d.mac_address, d.serial_number, d.api_key, d.nombre_dispositivo, d.active, d.created_at
		FROM dispositivos d
		JOIN empleado e ON d.conductor_id = e.id
		WHERE d.active = FALSE AND d.deleted_at IS NULL
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*entities.DispositivoConductorResponse
	for rows.Next() {
		var res entities.DispositivoConductorResponse
		err := rows.Scan(
			&res.ID, &res.ConductorID, &res.ConductorNombre, &res.ConductorApellido, &res.ConductorMail,
			&res.MacAddress, &res.SerialNumber, &res.ApiKey, &res.NombreDispositivo, &res.Active, &res.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, &res)
	}
	return result, nil
}
