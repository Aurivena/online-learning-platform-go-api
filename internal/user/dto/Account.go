package dto

import (
	"online-learning-platform-go-api/internal/user/entity"
	"time"
)

type RegistrationRequest struct {
	Email    string      `json:"email" example:"admin@example.com" binding:"required,email"`
	Username string      `json:"username" example:"admin" binding:"required,min=3,max=100"`
	Password string      `json:"password" example:"password" binding:"required,min=8,max=100"`
	Role     entity.Role `json:"role" example:"USER" default:"USER" binding:"required,oneof=USER ADMIN"`
}

type RegistrationResponse struct {
	ID           uint        `json:"id" example:"1"`
	Email        string      `json:"email" example:"admin@example.com"`
	Username     string      `json:"username" example:"admin"`
	Role         entity.Role `json:"role" example:"USER"`
	AccessToken  string      `json:"access_token" example:"eyJhbGciOiJIUzI1Ni..."`
	RefreshToken string      `json:"refresh_token" example:"eyJhbGciOiJIUzI1Ni..."`
	CreatedAt    time.Time   `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

type AccountResponse struct {
	ID        uint        `json:"id" example:"1"`
	Email     string      `json:"email" example:"admin@example.com"`
	Username  string      `json:"username" example:"admin"`
	Role      entity.Role `json:"role" example:"USER"`
	CreatedAt time.Time   `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

type UpdateRequest struct {
	Username string      `json:"username" example:"admin" min:"3" max:"100"`
	Email    string      `json:"email" example:"admin@example.com"`
	Password string      `json:"password" example:"password" min:"8" max:"100"`
	Role     entity.Role `json:"role" example:"USER"`
}
