package entities

// swagger:model CreateRutaRequest
type CreateRutaRequest struct {
	Nombre      string      `json:"nombre" binding:"required"`
	Descripcion string      `json:"descripcion"`
	JsonRuta    interface{} `json:"json_ruta" binding:"required"`
}

// swagger:model RutaResponse
type RutaResponse struct {
	Success bool   `json:"success"`
	Data    Ruta   `json:"data"`
}
