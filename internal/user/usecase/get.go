package usecase

import (
	"context"
	"online-learning-platform-go-api/internal/user/dto"
)

func (uc *AccountUseCase) Get(ctx context.Context, id int) (*dto.AccountResponse, error) {
	account, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	response := dto.AccountResponse{
		ID:        account.ID,
		Email:     account.Email,
		Username:  account.Username,
		Role:      account.Role,
		CreatedAt: account.CreatedAt,
	}

	return &response, nil
}
