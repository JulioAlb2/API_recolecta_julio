package entities

import "time"

type Empleado struct {
	ID          int        `json:"id"`
	TenantID    int        `json:"tenant_id"`
	Nombre      string     `json:"nombre"`
	Apellidos   string     `json:"apellidos"`
	Mail        string     `json:"mail"`
	Username    string     `json:"username"`
	Password    string     `json:"-"`
	Desactivado bool       `json:"desactivado"`
	RolID       int        `json:"rol_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
