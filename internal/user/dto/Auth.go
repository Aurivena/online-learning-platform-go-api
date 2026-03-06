package dto

type AuthRequest struct {
	Input    string `json:"input" example:"admin" min:"3" max:"100"`
	Password string `json:"password" example:"secret_password_123" min:"8" max:"100"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1Ni..."`
}
