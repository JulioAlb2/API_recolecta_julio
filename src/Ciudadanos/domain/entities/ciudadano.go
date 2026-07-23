package entities

import "time"

type Ciudadano struct {
	ID        int       `json:"id"`
	TenantID  int       `json:"tenant_id"`
	Email     string    `json:"email"`
	Alias     string    `json:"alias"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}