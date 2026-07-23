package entities

import "time"

type PuntoRecoleccion struct {
	PuntoID   int32     `json:"punto_id"`
	RutaID    int32     `json:"ruta_id"`
	CP        string    `json:"cp"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	Eliminado bool      `json:"eliminado"`
	CreatedAt time.Time `json:"created_at"`
}
