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
				Title:    "Учетная запись не найдена",
				Message:  "Пользователь с таким email не найден",
				Solution: "Проверьте email и пароль или обратитесь к администратору подразделения",
			},
		)
	}

	if err := domain.PasswordVerify(account.PasswordHash, input.Password); err != nil {
		return nil, netsp.BuildError(
			netstatus.CodeUnauthorized,
			netsp.ErrorDetail{
				Title:    "Неверный пароль",
				Message:  "Пароль указан неверно",
				Solution: "Проверьте пароль и повторите попытку",
			},
		)
	}

	return &LoginResult{
		Account:  account,
		Response: &dto.AuthResponse{},
	}, nil
}
