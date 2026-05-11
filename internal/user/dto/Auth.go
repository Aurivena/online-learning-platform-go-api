package dto

type LoginRequest struct {
	Input    string `json:"input" example:"admin@example.com" binding:"required"`
	Password string `json:"password" example:"secret_password_123" binding:"required,min=4,max=100"`
}

type AuthResponse struct {
	AccessToken  string `json:"accessToken" example:"eyJhbGciOiJIUzI1Ni..."`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1Ni..."`
}
