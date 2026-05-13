package dto

import (
	userDTO "online-learning-platform-go-api/internal/user/dto"
	"time"
)

type CreateOrganizationRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=125"`
	Tag         string `json:"tag" binding:"required,min=1,max=15"`
	Description string `json:"description" binding:"required,min=1"`
}

type UpdateOrganizationRequest struct {
	Title       string `json:"title" binding:"min=1,max=125"`
	Description string `json:"description" binding:"min=1"`
}

type OrganizationResponse struct {
	ID          uint64                `json:"id"`
	Title       string                `json:"title"`
	Tag         string                `json:"tag"`
	Description string                `json:"description"`
	Owner       userDTO.AccountResponse `json:"owner"`
	CreatedAt   time.Time             `json:"created_at"`
}

type AddAccountToOrgRequest struct {
	AccountID uint64 `json:"account_id" binding:"required"`
}

type RemoveAccountFromOrgRequest struct {
	AccountID uint64 `json:"account_id" binding:"required"`
}
