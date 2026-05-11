package usecase

import (
	"context"
	"online-learning-platform-go-api/internal/user/domain"
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"
	"time"

	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
)

func (u *AccountUseCase) Registration(ctx context.Context, input dto.RegistrationRequest) (*dto.RegistrationResponse, *netsp.Response[netsp.ErrorDetail]) {
	if input.Role != entity.RoleUser && input.Role != entity.RoleAdmin {
		return nil, netsp.BuildError(
			netstatus.CodeBadRequest,
			netsp.ErrorDetail{
				Title:    "Invalid Role",
				Message:  "Role must be either USER or ADMIN",
				Solution: "Please provide a valid role",
			},
		)
	}

	password, err := domain.PasswordHash(input.Password)
	if err != nil {
		return nil, netsp.BuildError(
			netstatus.CodeBadRequest,
			netsp.ErrorDetail{
				Title:    "Invalid Password",
				Message:  "The provided password is invalid or has incorrect format",
				Solution: "Please check the password and try again",
			},
		)
	}

	account := entity.Account{
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: password,
		Role:         input.Role,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := u.repo.Create(ctx, &account); err != nil {
		return nil, netsp.BuildError(
			netstatus.CodeBadRequest,
			netsp.ErrorDetail{
				Title:    "Account Creation Error",
				Message:  "The account could not be created, possibly due to a duplicate email or username",
				Solution: "Please check your data and try again",
			},
		)
	}

	return &dto.RegistrationResponse{
		ID:        account.ID,
		Email:     account.Email,
		Username:  account.Username,
		Role:      account.Role,
		CreatedAt: account.CreatedAt,
	}, nil
}
