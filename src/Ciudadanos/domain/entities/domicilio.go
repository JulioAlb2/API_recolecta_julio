package entities

import "time"

// Domicilio representa la dirección registrada de un ciudadano, con coordenadas opcionales.
type Domicilio struct {
	ID          int       `json:"id" example:"1"`
	CiudadanoID int       `json:"ciudadano_id" example:"42"`
	ColoniaID   int       `json:"colonia_id" example:"1"`
	Alias       string    `json:"alias" example:"Casa"`
	Calle       string    `json:"calle" example:"Calle Olmo"`
	Numero      string    `json:"numero" example:"123"`
	Referencia  *string   `json:"referencia,omitempty" example:"Frente al parque"`
	Latitud     *float64  `json:"latitud,omitempty" example:"16.6278"`
	Longitud    *float64  `json:"longitud,omitempty" example:"-93.1045"`
	CreatedAt   time.Time `json:"created_at" example:"2026-07-04T16:00:00Z"`
}
