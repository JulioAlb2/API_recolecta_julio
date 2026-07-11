package infrastructure

type RecorridoResponse struct {
	Success bool       `json:"success" example:"true"`
	Message string     `json:"message,omitempty" example:"no hay recorrido activo"`
	Data    *Recorrido `json:"data,omitempty"`
}

type IniciarRecorridoRequest struct {
	RutaID   int `json:"ruta_id" example:"1"`
	ChoferID int `json:"chofer_id" example:"7"`
}
