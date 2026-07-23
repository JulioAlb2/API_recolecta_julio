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

// FIX 2026-07-22 (detalle y contexto completo en
// docs/09-fix-colonia-id-hardcodeado-ruta.md, en la raíz del monorepo
// 233298_recolecta_web -- fuera de este submódulo gin-backend):
// Antes este método insertaba colonia_id=1 fijo porque en ese momento
// colonia_id no tenía FK y cualquier entero pasaba. El commit 9ed2eb7
// (multitenancy) agregó en db_constraints.sql la restricción
// fk_colonia_ruta: ruta.colonia_id -> colonia.colonia_id, y la tabla
// colonia no trae ningún seed. Resultado: si en el ambiente no existía
// una colonia con id=1, CADA insert a "ruta" violaba la FK y el
// controlador devolvía 500 ("Error creando ruta") al guardar cualquier
// ruta desde el dashboard.
//
// Este cambio deja de asumir un id fijo: resuelve un colonia_id real
// (el primero que exista) y, si la tabla colonia está vacía, crea una
// colonia por defecto para no bloquear la creación de rutas.
//
// TODO(multitenant): esta resolución es global (no filtra por tenant).
// Cuando el módulo de Rutas se conecte al flujo de multitenancy
// (RunInTenantTx, como ya hace colonia_repository.go), reemplazar esto
// por una resolución de colonia por defecto scoped al tenant actual.
func (pg *PostgresRuta) Save(ruta *entities.Ruta) error {
	ruta.CreatedAt = time.Now()
	ctx := context.Background()

	coloniaID, err := pg.resolveDefaultColoniaID(ctx)
	if err != nil {
		return fmt.Errorf("error al resolver colonia por defecto: %w", err)
	}

	sql := `
		INSERT INTO ruta (nombre, descripcion, json_ruta, colonia_id, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err = pg.conn.QueryRow(
		ctx,
		sql,
		ruta.Nombre,
		ruta.Descripcion,
		ruta.JsonRuta,
		coloniaID,
		ruta.CreatedAt,
	).Scan(&ruta.RutaID)

	if err != nil {
		return fmt.Errorf("error al guardar ruta: %w", err)
	}
	ruta.Eliminado = false
	return nil
}

// resolveDefaultColoniaID devuelve un colonia_id válido para asociar a una
// ruta cuando quien llama (hoy: el dashboard web, vía CreateRutaController)
// no especifica una colonia explícita. Ver nota de FIX arriba de Save().
func (pg *PostgresRuta) resolveDefaultColoniaID(ctx context.Context) (int32, error) {
	var coloniaID int32

	err := pg.conn.QueryRow(ctx, `SELECT colonia_id FROM colonia ORDER BY colonia_id LIMIT 1`).Scan(&coloniaID)
	if err == nil {
		return coloniaID, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return 0, fmt.Errorf("error al buscar colonia existente: %w", err)
	}

	// No hay ninguna colonia todavía: creamos una por defecto en vez de
	// fallar, para no bloquear la creación de rutas mientras el módulo de
	// colonias no tenga datos reales cargados.
	err = pg.conn.QueryRow(
		ctx,
		`INSERT INTO colonia (nombre, zona) VALUES ($1, $2) RETURNING colonia_id`,
		"Sin colonia asignada",
		"Sin definir",
	).Scan(&coloniaID)
	if err != nil {
		return 0, fmt.Errorf("error al crear colonia por defecto: %w", err)
	}

	return coloniaID, nil
}

func (pg *PostgresRuta) ListAll() ([]entities.Ruta, error) {
	sql := `
		SELECT id, nombre, descripcion, json_ruta, (deleted_at IS NOT NULL) AS eliminado, created_at
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
		var r entities.Ruta

		err := rows.Scan(
			&r.RutaID,
			&r.Nombre,
			&r.Descripcion,
			&r.JsonRuta,
			&r.Eliminado,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear ruta: %w", err)
		}

		rutas = append(rutas, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración: %w", err)
	}

	return rutas, nil
}

func (pg *PostgresRuta) GetById(id int32) (*entities.Ruta, error) {
	sql := `
		SELECT id, nombre, descripcion, json_ruta, (deleted_at IS NOT NULL) AS eliminado, created_at
		FROM ruta
		WHERE id = $1 AND deleted_at IS NULL
	`

	var r entities.Ruta

	err := pg.conn.QueryRow(context.Background(), sql, id).Scan(
		&r.RutaID,
		&r.Nombre,
		&r.Descripcion,
		&r.JsonRuta,
		&r.Eliminado,
		&r.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, errors.New("ruta no encontrada")
	}
	if err != nil {
		return nil, fmt.Errorf("error al obtener ruta: %w", err)
	}

	return &r, nil
}

func (pg *PostgresRuta) Update(ruta *entities.Ruta) error {
	sql := `
		UPDATE ruta
		SET nombre = $1, descripcion = $2, json_ruta = $3
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
		SET deleted_at = NOW()
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
		SELECT id, nombre, descripcion, json_ruta, (deleted_at IS NOT NULL) AS eliminado, created_at
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
		var r entities.Ruta

		err := rows.Scan(
			&r.RutaID,
			&r.Nombre,
			&r.Descripcion,
			&r.JsonRuta,
			&r.Eliminado,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear ruta: %w", err)
		}

		rutas = append(rutas, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error durante la iteración: %w", err)
	}

	return rutas, nil
}