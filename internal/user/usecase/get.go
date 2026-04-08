package usecase

import (
	"context"
	"net/http"
	"online-learning-platform-go-api/internal/user/dto"

	"github.com/Aurivena/spond/v3/netsp"
)

func (uc *AccountUseCase) Get(ctx context.Context, id int) (*dto.AccountResponse, *netsp.AppError) {
	account, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusNotFound,
			"Account Not Found",
			"An account with the specified ID was not found",
			"Please check the account ID and try again",
		)
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
