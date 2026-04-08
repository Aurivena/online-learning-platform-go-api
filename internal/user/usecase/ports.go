package usecase

import (
	"context"
	"online-learning-platform-go-api/internal/user/entity"
)

type AccountUseCase struct {
	repo AccountRepository
}

func NewAccountUseCase(repo AccountRepository) *AccountUseCase {
	return &AccountUseCase{repo: repo}
}

type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	Get(ctx context.Context, id int) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
}
