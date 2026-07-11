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

const domicilioSelectColumns = `
	id, ciudadano_id, colonia_id, alias, calle, numero, referencia, latitud, longitud, created_at
`

func (r *DomicilioPostgresRepository) scanDomicilio(row pgx.Row) (*entities.Domicilio, error) {
	var d entities.Domicilio
	err := row.Scan(
		&d.ID,
		&d.CiudadanoID,
		&d.ColoniaID,
		&d.Alias,
		&d.Calle,
		&d.Numero,
		&d.Referencia,
		&d.Latitud,
		&d.Longitud,
		&d.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DomicilioPostgresRepository) Create(ctx context.Context, d *entities.Domicilio) (int, error) {
	const q = `
		INSERT INTO domicilio (ciudadano_id, colonia_id, alias, calle, numero, referencia, latitud, longitud, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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
		d.Latitud,
		d.Longitud,
		d.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *DomicilioPostgresRepository) GetByID(ctx context.Context, id int) (*entities.Domicilio, error) {
	const q = `
		SELECT ` + domicilioSelectColumns + `
		FROM domicilio
		WHERE id = $1
	`

	d, err := r.scanDomicilio(r.db.QueryRow(ctx, q, id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return d, nil
}

func (r *DomicilioPostgresRepository) List(ctx context.Context) ([]entities.Domicilio, error) {
	const q = `
		SELECT ` + domicilioSelectColumns + `
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
		d, err := r.scanDomicilio(rows)
		if err != nil {
			return nil, err
		}
		domicilios = append(domicilios, *d)
	}

	return domicilios, rows.Err()
}

func (r *DomicilioPostgresRepository) ListByCiudadanoID(ctx context.Context, ciudadanoID int) ([]entities.Domicilio, error) {
	const q = `
		SELECT ` + domicilioSelectColumns + `
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
		d, err := r.scanDomicilio(rows)
		if err != nil {
			return nil, err
		}
		domicilios = append(domicilios, *d)
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
		    referencia = $5,
		    latitud = $6,
		    longitud = $7
		WHERE id = $8
	`

	cmd, err := r.db.Exec(ctx, q, d.ColoniaID, d.Alias, d.Calle, d.Numero, d.Referencia, d.Latitud, d.Longitud, d.ID)
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

func (r *DomicilioPostgresRepository) FindByAlias(ctx context.Context, alias string) (*entities.Domicilio, error) {
	const q = `
		SELECT ` + domicilioSelectColumns + `
		FROM domicilio
		WHERE alias = $1
	`

	d, err := r.scanDomicilio(r.db.QueryRow(ctx, q, alias))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return d, nil
}
