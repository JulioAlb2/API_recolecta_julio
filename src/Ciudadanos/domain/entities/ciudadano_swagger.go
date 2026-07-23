package entities

// swagger:model CreateCiudadanoRequest
type CreateCiudadanoRequest struct {
	Email    string  `json:"email" binding:"required"`
	Alias    string  `json:"alias" binding:"required"`
	Password string  `json:"password" binding:"required"`
	FCMToken string  `json:"fcm_token" binding:"required"`
}

// swagger:model UpdateCiudadanoRequest
type UpdateCiudadanoRequest struct {
	Email    string `json:"email"`
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

// swagger:model LoginCiudadanoRequest
type LoginCiudadanoRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// swagger:model LoginCiudadanoData
type LoginCiudadanoData struct {
	Ciudadano Ciudadano `json:"ciudadano"`
	Token     string    `json:"token"`
}

// swagger:model LoginCiudadanoResponse
type LoginCiudadanoResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Token   string    `json:"token"`
	Data    Ciudadano `json:"data"`
	Code    int       `json:"code"`
}

// swagger:model CiudadanoResponse
type CiudadanoResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    Ciudadano `json:"data"`
	Code    int       `json:"code"`
}

// swagger:model CiudadanoIDResponse
type CiudadanoIDResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      int    `json:"id"`
	Code    int    `json:"code"`
}

// swagger:model CiudadanoDetailResponse
type CiudadanoDetailResponse struct {
	Success bool      `json:"success"`
	Data    Ciudadano `json:"data"`
}

// swagger:model CiudadanoListResponse
type CiudadanoListResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    []Ciudadano `json:"data"`
	Code    int         `json:"code"`
}

// swagger:model CiudadanoMessageResponse
type CiudadanoMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}
