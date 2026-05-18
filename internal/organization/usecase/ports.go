package usecase

import (
	"context"
	"online-learning-platform-go-api/internal/organization/entity"
	userEntity "online-learning-platform-go-api/internal/user/entity"
)

type OrganizationRepository interface {
	Create(ctx context.Context, org *entity.Organization) error
	GetByID(ctx context.Context, id uint64) (*entity.Organization, error)
	GetByTag(ctx context.Context, tag string) (*entity.Organization, error)
	GetByOwner(ctx context.Context, ownerID uint64) ([]entity.Organization, error)
	GetByAccountEntities(ctx context.Context, accountID uint64) ([]entity.Organization, error)
	GetAll(ctx context.Context) ([]entity.Organization, error)
	Update(ctx context.Context, org *entity.Organization) error
	Delete(ctx context.Context, id uint64) error
	AddAccount(ctx context.Context, orgID, accountID uint64) error
	RemoveAccount(ctx context.Context, orgID, accountID uint64) error
	GetAccounts(ctx context.Context, orgID uint64) ([]uint64, error)
	GetAccountEntities(ctx context.Context, orgID uint64) ([]userEntity.Account, error)
	IsMember(ctx context.Context, orgID, accountID uint64) (bool, error)
}
