package usecase

import (
	"context"
	"net/http"
	"online-learning-platform-go-api/internal/user/entity"

	"github.com/Aurivena/spond/v4/netsp"
)

func (uc *AccountUseCase) Get(ctx context.Context, id int) (*entity.Account, *netsp.Response[netsp.ErrorDetail]) {
	account, err := uc.repo.Get(ctx, id)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Account Not Found",
				Message:  "An account with the specified ID was not found",
				Solution: "Please check the account ID and try again",
			},
		)
	}
	return account, nil
}
