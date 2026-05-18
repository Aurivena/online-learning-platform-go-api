package dto

import (
	userDTO "online-learning-platform-go-api/internal/user/dto"
	userEntity "online-learning-platform-go-api/internal/user/entity"
	"time"
)

type CreateOrganizationRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=125"`
	Tag         string `json:"tag" binding:"required,min=1,max=15"`
	Description string `json:"description" binding:"required,min=1"`
	ImageURL    string `json:"image_url"`
	HeaderTitle string `json:"header_title" binding:"omitempty,max=80"`
}

type UpdateOrganizationRequest struct {
	Title       string  `json:"title" binding:"min=1,max=125"`
	Tag         string  `json:"tag" binding:"omitempty,min=1,max=15"`
	Description string  `json:"description" binding:"min=1"`
	ImageURL    *string `json:"image_url"`
	HeaderTitle *string `json:"header_title" binding:"omitempty,max=80"`
}

type OrganizationResponse struct {
	ID          uint64                  `json:"id"`
	Title       string                  `json:"title"`
	Tag         string                  `json:"tag"`
	Description string                  `json:"description"`
	ImageURL    string                  `json:"image_url"`
	HeaderTitle string                  `json:"header_title"`
	Owner       userDTO.AccountResponse `json:"owner"`
	CreatedAt   time.Time               `json:"created_at"`
}

type AddAccountToOrgRequest struct {
	AccountID uint64 `json:"account_id" binding:"required"`
}

type RemoveAccountFromOrgRequest struct {
	AccountID uint64 `json:"account_id" binding:"required"`
}

type OrganizationAccountResponse struct {
	ID        uint64          `json:"id"`
	Email     string          `json:"email"`
	Username  string          `json:"username"`
	Role      userEntity.Role `json:"role"`
	CreatedAt time.Time       `json:"created_at"`
}
