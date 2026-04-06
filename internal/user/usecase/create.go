package usecase

import (
	"context"
	"net/http"
	"online-learning-platform-go-api/internal/pkg"
	"online-learning-platform-go-api/internal/user/domain"
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"
	"time"

	"github.com/Aurivena/spond/v3/netsp"
)

func (u *AccountUseCase) Registration(ctx context.Context, dto dto.RegistrationRequest) (string, string, *netsp.AppError) {
	password, err := domain.PasswordHash(dto.Password)
	if err != nil {
		return "", "", netsp.BuildError(
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
		return "", "", netsp.BuildError(
			http.StatusBadRequest,
			"Token Generation Error",
			"The access token could not be generated",
			"Please check the token generation service and try again",
		)
	}

	accessToken, err := pkg.GenerateToken(account.ID, entity.User, 0, nil)
	if err != nil {
		return "", "", netsp.BuildError(
			http.StatusBadRequest,
			"Token Generation Error",
			"The refresh token could not be generated",
			"Please check the token generation service and try again",
		)
	}

	refreshToken, err := pkg.GenerateToken(account.ID, entity.User, 0, nil)
	if err != nil {
		return "", "", netsp.BuildError(
			http.StatusBadRequest,
			"Invalid Input",
			"The provided data is incorrect",
			"Please check the documentation and try again",
		)
	}

	return accessToken, refreshToken, nil
}
