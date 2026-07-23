package entities

import "time"

type Dispositivo struct {
	ID                int        `json:"id"`
	ConductorID       int        `json:"conductor_id"`
	MacAddress        string     `json:"mac_address"`
	SerialNumber      string     `json:"serial_number"`
	ApiKey            string     `json:"api_key"`
	NombreDispositivo string     `json:"nombre_dispositivo"`
	SistemaOperativo  string     `json:"sistema_operativo"`
	Active            bool       `json:"active"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

type SolicitarDispositivoRequest struct {
	MacAddress        string `json:"mac_address" binding:"required"`
	SerialNumber      string `json:"serial_number" binding:"required"`
	NombreDispositivo string `json:"nombre_dispositivo"`
}

type DispositivoConductorResponse struct {
	ID                int       `json:"id"`
	ConductorID       int       `json:"conductor_id"`
	ConductorNombre   string    `json:"conductor_nombre"`
	ConductorApellido string    `json:"conductor_apellido"`
	ConductorMail     string    `json:"conductor_mail"`
	MacAddress        string    `json:"mac_address"`
	SerialNumber      string    `json:"serial_number"`
	ApiKey            string    `json:"api_key"`
	NombreDispositivo string    `json:"nombre_dispositivo"`
	Active            bool      `json:"active"`
	CreatedAt         time.Time `json:"created_at"`
}
