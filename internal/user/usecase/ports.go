package usecase

import (
	"context"
	orgentity "online-learning-platform-go-api/internal/organization/entity"
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"

	"github.com/Aurivena/spond/v4/netsp"
)

type AccountUseCaseInterface interface {
	Registration(ctx context.Context, dto dto.RegistrationRequest) *netsp.Response[netsp.ErrorDetail]
	Login(ctx context.Context, dto dto.LoginRequest) (*LoginResult, *netsp.Response[netsp.ErrorDetail])
	Get(ctx context.Context, id int) (*entity.Account, *netsp.Response[netsp.ErrorDetail])
	GetProfile(ctx context.Context, id uint64) (*dto.UserProfileResponse, *netsp.Response[netsp.ErrorDetail])
	Update(ctx context.Context, req dto.UpdateRequest, id int) *netsp.Response[netsp.ErrorDetail]
}

type AccountUseCase struct {
	repo    AccountRepository
	orgRepo OrganizationRepository
}

func NewAccountUseCase(repo AccountRepository, orgRepo OrganizationRepository) AccountUseCaseInterface {
	return &AccountUseCase{repo: repo, orgRepo: orgRepo}
}

type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	GetByID(ctx context.Context, id uint64) (*entity.Account, error)
	GetByEmail(ctx context.Context, email string) (*entity.Account, error)
	GetByUsername(ctx context.Context, username string) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id uint64) error
}

type OrganizationRepository interface {
	GetByAccountEntities(ctx context.Context, accountID uint64) ([]orgentity.Organization, error)
	GetByID(ctx context.Context, id uint64) (*orgentity.Organization, error)
	AddAccount(ctx context.Context, orgID, accountID uint64) error
}
