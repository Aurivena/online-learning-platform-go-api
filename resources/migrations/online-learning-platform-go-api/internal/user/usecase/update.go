package usecase

import (
	"context"
	"net/http"
	"online-learning-platform-go-api/internal/user/domain"
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"
	"time"

	"github.com/Aurivena/spond/v4/netsp"
)

func (uc *AccountUseCase) Update(ctx context.Context, req dto.UpdateRequest, id int) *netsp.Response[netsp.ErrorDetail] {
	hashPassword, err := domain.PasswordHash(req.Password)
	if err != nil {
		return netsp.BuildError(
			http.StatusBadRequest,
			netsp.ErrorDetail{
				Title:    "Invalid Password",
				Message:  "The provided password is invalid or has incorrect format",
				Solution: "Please check the password and try again",
			},
		)
	}
	account := &entity.Account{
		ID:           uint(id),
		Username:     req.Username,
		Email:        req.Email,
		Role:         req.Role,
		PasswordHash: hashPassword,
		UpdatedAt:    time.Now().UTC(),
	}
	if err := uc.repo.Update(ctx, account); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Account Update Error",
				Message:  "Failed to update account in the database",
				Solution: "Please try again later",
			},
		)
	}
	return nil
}
