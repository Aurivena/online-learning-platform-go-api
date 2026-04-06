package dto

import (
	"time"
)

type Role string

type RegistrationRequest struct {
	Email    string `json:"email" example:"admin@example.com"`
	Username string `json:"username" example:"admin" min:"3" max:"100"`
	Password string `json:"password" example:"password" min:"8" max:"100"`
	Role     Role   `json:"role" example:"admin" default:"USER"`
}

type AccountResponse struct {
	ID        uint      `json:"id" example:"1"`
	Email     string    `json:"email" example:"admin@example.com"`
	Username  string    `json:"username" example:"admin"`
	Role      Role      `json:"role" example:"admin"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

type UpdateRequest struct {
	Username string `json:"username" example:"admin" min:"3" max:"100"`
	Email    string `json:"email" example:"admin@example.com"`
	Password string `json:"password" example:"password" min:"8" max:"100"`
	Role     Role   `json:"role" example:"admin"`
}
