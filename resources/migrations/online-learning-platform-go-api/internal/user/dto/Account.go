package dto

import (
	"online-learning-platform-go-api/internal/user/entity"
	"time"
)

type RegistrationRequest struct {
	Email          string  `json:"email" example:"admin@example.com" binding:"required,email"`
	Username       string  `json:"username" example:"admin" binding:"required,min=3,max=100"`
	Password       string  `json:"password" example:"password" binding:"required,min=8,max=100"`
	OrganizationID *uint64 `json:"organization_id"`
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

type UserProfileResponse struct {
	ID            uint64                     `json:"id"`
	Email         string                     `json:"email"`
	Username      string                     `json:"username"`
	Role          entity.Role                `json:"role"`
	CreatedAt     time.Time                  `json:"created_at"`
	Organizations []UserOrganizationResponse `json:"organizations"`
}

type UserOrganizationResponse struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Tag         string    `json:"tag"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	HeaderTitle string    `json:"header_title"`
	OwnerID     uint64    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}
