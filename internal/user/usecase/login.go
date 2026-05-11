package usecase

import (
	"context"
	"online-learning-platform-go-api/internal/user/domain"
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"

	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
)

type LoginResult struct {
	Account  *entity.Account
	Response *dto.AuthResponse
}

func (u *AccountUseCase) Login(ctx context.Context, input dto.LoginRequest) (*LoginResult, *netsp.Response[netsp.ErrorDetail]) {
	account, err := u.repo.GetByEmail(ctx, input.Input)
	if err != nil {
		return nil, netsp.BuildError(
			netstatus.CodeNotFound,
			netsp.ErrorDetail{
				Title:    "Account Not Found",
				Message:  "No account found with this email",
				Solution: "Please check your email or register a new account",
			},
		)
	}

	if account == nil {
		account, err = u.repo.GetByUsername(ctx, input.Input)
		if err != nil {
			return nil, netsp.BuildError(
				netstatus.CodeNotFound,
				netsp.ErrorDetail{
					Title:    "Account Not Found",
					Message:  "No account found with this username",
					Solution: "Please check your username or register a new account",
				},
			)
		}
	}

	if err := domain.PasswordVerify(account.PasswordHash, input.Password); err != nil {
		return nil, netsp.BuildError(
			netstatus.CodeUnauthorized,
			netsp.ErrorDetail{
				Title:    "Invalid Password",
				Message:  "The password is incorrect",
				Solution: "Please check your password and try again",
			},
		)
	}

	return &LoginResult{
		Account:  account,
		Response: &dto.AuthResponse{},
	}, nil
}
