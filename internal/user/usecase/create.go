package usecase

import (
	"context"
	"net/http"
	"online-learning-platform-go-api/internal/user/domain"
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"
	"time"

	"github.com/Aurivena/spond/v3/netsp"
)

func (u *AccountUseCase) Registration(ctx context.Context, dto dto.RegistrationRequest) (*entity.Account, *netsp.AppError) {
	password, err := domain.PasswordHash(dto.Password)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusBadRequest,
			"Invalid Password",
			"The provided password is invalid or has incorrect format",
			"Please check the password and try again",
		)
	}

	account := entity.Account{
		Email:     dto.Email,
		Username:  dto.Username,
		Password:  password,
		Role:      dto.Role,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := u.repo.Create(ctx, &account); err != nil {
		return nil, netsp.BuildError(
			http.StatusBadRequest,
			"Token Generation Error",
			"The access token could not be generated",
			"Please check the token generation service and try again",
		)
	}
	return &account, nil
}
