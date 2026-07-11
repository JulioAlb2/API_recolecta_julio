package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/entities"
	"github.com/vicpoo/API_recolecta/src/Rutas/domain/ports"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresRuta struct {
	conn *pgxpool.Pool
}

func NewPostgresRuta() ports.IRuta {
	conn, _ := core.ConnectPostgres()
	return &PostgresRuta{conn: conn}
}

func (pg *PostgresRuta) Save(ruta *entities.Ruta) error {
	ruta.CreatedAt = time.Now()
	sql := `
		INSERT INTO ruta (nombre, descripcion, colonia_id, json_ruta, created_at)
		VALUES ($1, $2, 1, $3, $4)
		RETURNING id
	`

	err := pg.conn.QueryRow(
		context.Background(),
		sql,
		ruta.Nombre,
		ruta.Descripcion,
		ruta.JsonRuta,
		ruta.CreatedAt,
	).Scan(&ruta.RutaID)

	if err != nil {
		return fmt.Errorf("error al guardar ruta: %w", err)
	}
	return nil
}

func (pg *PostgresRuta) scanRuta(row pgx.Row) (*entities.Ruta, error) {
	var r entities.Ruta
	var deletedAt *time.Time

	err := row.Scan(
		&r.RutaID,
		&r.Nombre,
		&r.Descripcion,
		&r.JsonRuta,
		&deletedAt,
		&r.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	r.Eliminado = deletedAt != nil
	return &r, nil
}

func (pg *PostgresRuta) ListAll() ([]entities.Ruta, error) {
	sql := `
		SELECT id, nombre, descripcion, json_ruta::text, deleted_at, created_at
		FROM ruta
		WHERE deleted_at IS NULL
		ORDER BY id DESC
	`

	rows, err := pg.conn.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf("error al consultar rutas: %w", err)
	}
	defer rows.Close()

	var rutas []entities.Ruta
	for rows.Next() {
		r, err := pg.scanRuta(rows)
		if err != nil {
			return nil, fmt.Errorf("error al escanear ruta: %w", err)
		}
		rutas = append(rutas, *r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración: %w", err)
	}

	return rutas, nil
}

func (pg *PostgresRuta) GetById(id int32) (*entities.Ruta, error) {
	sql := `
		SELECT id, nombre, descripcion, json_ruta::text, deleted_at, created_at
		FROM ruta
		WHERE id = $1
	`

	r, err := pg.scanRuta(pg.conn.QueryRow(context.Background(), sql, id))
	if err == pgx.ErrNoRows {
		return nil, errors.New("ruta no encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener ruta: %w", err)
	}

	return r, nil
}

func (pg *PostgresRuta) Update(ruta *entities.Ruta) error {
	sql := `
		UPDATE ruta
		SET nombre = $1, descripcion = $2, json_ruta = $3, updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
	`

	cmd, err := pg.conn.Exec(
		context.Background(),
		sql,
		ruta.Nombre,
		ruta.Descripcion,
		ruta.JsonRuta,
		ruta.RutaID,
	)
	if err != nil {
		return fmt.Errorf("error al actualizar ruta: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("ruta no encontrada o ya eliminada")
	}

	return nil
}

func (pg *PostgresRuta) Delete(id int32) error {
	sql := `
		UPDATE ruta
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	cmd, err := pg.conn.Exec(context.Background(), sql, id)
	if err != nil {
		return fmt.Errorf("error al eliminar ruta: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("ruta no encontrada o ya eliminada")
	}

	return nil
}

func (pg *PostgresRuta) GetActivas() ([]entities.Ruta, error) {
	sql := `
		SELECT id, nombre, descripcion, json_ruta::text, deleted_at, created_at
		FROM ruta
		WHERE deleted_at IS NULL
		ORDER BY id DESC
	`

	rows, err := pg.conn.Query(context.Background(), sql)
	if err != nil {
		return nil, fmt.Errorf("error al consultar rutas activas: %w", err)
	}
	defer rows.Close()

	var rutas []entities.Ruta
	for rows.Next() {
		r, err := pg.scanRuta(rows)
		if err != nil {
			return nil, fmt.Errorf("error al escanear ruta: %w", err)
		}
		rutas = append(rutas, *r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración: %w", err)
	}

	return rutas, nil
}
