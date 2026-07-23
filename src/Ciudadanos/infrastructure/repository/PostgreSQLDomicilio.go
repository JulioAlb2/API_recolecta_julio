package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
)

type DomicilioPostgresRepository struct {
	db *pgxpool.Pool
}

func NewDomicilioPostgresRepository(db *pgxpool.Pool) *DomicilioPostgresRepository {
	return &DomicilioPostgresRepository{db: db}
}

func (r *DomicilioPostgresRepository) Create(ctx context.Context, d *entities.Domicilio) (int, error) {
	const q = `
		INSERT INTO domicilio (ciudadano_id, colonia_id, alias, calle, numero, referencia, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(
		ctx,
		q,
		d.CiudadanoID,
		d.ColoniaID,
		d.Alias,
		d.Calle,
		d.Numero,
		d.Referencia,
		d.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *DomicilioPostgresRepository) GetByID(ctx context.Context, id int) (*entities.Domicilio, error) {
	const q = `
		SELECT id, ciudadano_id, colonia_id, alias, calle, numero, referencia, created_at
		FROM domicilio
		WHERE id = $1
	`

	var d entities.Domicilio
	err := r.db.QueryRow(ctx, q, id).Scan(
		&d.ID,
		&d.CiudadanoID,
		&d.ColoniaID,
		&d.Alias,
		&d.Calle,
		&d.Numero,
		&d.Referencia,
		&d.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}

func (r *DomicilioPostgresRepository) List(ctx context.Context) ([]entities.Domicilio, error) {
	const q = `
		SELECT id, ciudadano_id, colonia_id, alias, calle, numero, referencia, created_at
		FROM domicilio
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domicilios []entities.Domicilio

	for rows.Next() {
		var d entities.Domicilio
		if err := rows.Scan(
			&d.ID,
			&d.CiudadanoID,
			&d.ColoniaID,
			&d.Alias,
			&d.Calle,
			&d.Numero,
			&d.Referencia,
			&d.CreatedAt,
		); err != nil {
			return nil, err
		}
		domicilios = append(domicilios, d)
	}

	return domicilios, rows.Err()
}

func (r *DomicilioPostgresRepository) ListByCiudadanoID(ctx context.Context, ciudadanoID int) ([]entities.Domicilio, error) {
	const q = `
		SELECT id, ciudadano_id, colonia_id, alias, calle, numero, referencia, created_at
		FROM domicilio
		WHERE ciudadano_id = $1
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, q, ciudadanoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domicilios []entities.Domicilio

	for rows.Next() {
		var d entities.Domicilio
		if err := rows.Scan(
			&d.ID,
			&d.CiudadanoID,
			&d.ColoniaID,
			&d.Alias,
			&d.Calle,
			&d.Numero,
			&d.Referencia,
			&d.CreatedAt,
		); err != nil {
			return nil, err
		}
		domicilios = append(domicilios, d)
	}

	return domicilios, rows.Err()
}

func (r *DomicilioPostgresRepository) Update(ctx context.Context, d *entities.Domicilio) error {
	const q = `
		UPDATE domicilio
		SET colonia_id = $1,
		    alias = $2,
		    calle = $3,
		    numero = $4,
		    referencia = $5
		WHERE id = $6
	`

	cmd, err := r.db.Exec(ctx, q, d.ColoniaID, d.Alias, d.Calle, d.Numero, d.Referencia, d.ID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("domicilio no encontrado")
	}

	return nil
}

func (r *DomicilioPostgresRepository) DeleteByCiudadano(ctx context.Context, id int, ciudadanoID int) error {
	const q = `
		DELETE FROM domicilio
		WHERE id = $1 AND ciudadano_id = $2
	`

	cmd, err := r.db.Exec(ctx, q, id, ciudadanoID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("domicilio no encontrado o no pertenece al ciudadano")
	}

	return nil
}

func (r *DomicilioPostgresRepository) FindByAlias(ctx context.Context, alias string, ciudadanoID int) (*entities.Domicilio, error) {
	const q = `
		SELECT id, ciudadano_id, colonia_id, alias, calle, numero, referencia, created_at
		FROM domicilio
		WHERE alias = $1 AND ciudadano_id = $2
	`

	var d entities.Domicilio
	err := r.db.QueryRow(ctx, q, alias, ciudadanoID).Scan(
		&d.ID,
		&d.CiudadanoID,
		&d.ColoniaID,
		&d.Alias,
		&d.Calle,
		&d.Numero,
		&d.Referencia,
		&d.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &d, nil
}