package dto

type LoginRequest struct {
	Email    string `json:"email" example:"admin@example.com" binding:"required,email"`
	Password string `json:"password" example:"secret_password_123" binding:"required,min=8,max=100"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1Ni..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1Ni..."`
}
