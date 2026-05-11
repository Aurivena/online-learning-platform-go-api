package usecase

import (
	"context"
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"

	"github.com/Aurivena/spond/v4/netsp"
)

type AccountUseCaseInterface interface {
	Registration(ctx context.Context, dto dto.RegistrationRequest) (*dto.RegistrationResponse, *netsp.Response[netsp.ErrorDetail])
	Login(ctx context.Context, dto dto.LoginRequest) (*LoginResult, *netsp.Response[netsp.ErrorDetail])
	Get(ctx context.Context, id int) (*entity.Account, *netsp.Response[netsp.ErrorDetail])
	Update(ctx context.Context, req dto.UpdateRequest, id int) *netsp.Response[netsp.ErrorDetail]
}

type AccountUseCase struct {
	repo AccountRepository
}

func NewAccountUseCase(repo AccountRepository) *AccountUseCase {
	return &AccountUseCase{repo: repo}
}

type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	Get(ctx context.Context, id int) (*entity.Account, error)
	GetByEmail(ctx context.Context, email string) (*entity.Account, error)
	GetByUsername(ctx context.Context, username string) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
}
