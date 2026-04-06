package usecase

import (
	"context"
	"online-learning-platform-go-api/internal/user/adaptors"
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"
)

type AccountUseCase struct {
	repo *adaptors.Repository
}

func NewAccountUseCase(repo *adaptors.Repository) *AccountUseCase {
	return &AccountUseCase{repo: repo}
}

type Repository interface {
	Create(ctx context.Context, account *entity.Account) error
	Get(ctx context.Context, id int) (dto.AccountResponse, error)
	Update(ctx context.Context, account *entity.Account) error
}
