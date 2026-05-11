package adaptors

import (
	"context"
	"online-learning-platform-go-api/internal/organization/entity"

	"gorm.io/gorm"
)

type OrganizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) Create(ctx context.Context, org *entity.Organization) error {
	return r.db.WithContext(ctx).Create(org).Error
}

func (r *OrganizationRepository) GetByID(ctx context.Context, id uint64) (*entity.Organization, error) {
	var org entity.Organization
	err := r.db.WithContext(ctx).First(&org, id).Error
	return &org, err
}

func (r *OrganizationRepository) GetByTag(ctx context.Context, tag string) (*entity.Organization, error) {
	var org entity.Organization
	err := r.db.WithContext(ctx).Where("tag = ?", tag).First(&org).Error
	return &org, err
}

func (r *OrganizationRepository) GetByOwner(ctx context.Context, ownerID uint64) ([]entity.Organization, error) {
	var orgs []entity.Organization
	err := r.db.WithContext(ctx).
		Where("owner_id = ?", ownerID).
		Order("created_at desc").
		Find(&orgs).Error
	return orgs, err
}

func (r *OrganizationRepository) GetAll(ctx context.Context) ([]entity.Organization, error) {
	var orgs []entity.Organization
	err := r.db.WithContext(ctx).
		Order("created_at desc").
		Find(&orgs).Error
	return orgs, err
}

func (r *OrganizationRepository) Update(ctx context.Context, org *entity.Organization) error {
	return r.db.WithContext(ctx).Model(org).Updates(org).Error
}

func (r *OrganizationRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.Organization{}, id).Error
}

func (r *OrganizationRepository) AddAccount(ctx context.Context, orgID, accountID uint64) error {
	return r.db.WithContext(ctx).Table("organization_accounts").Create(map[string]interface{}{
		"organization_id": orgID,
		"account_id":      accountID,
	}).Error
}

func (r *OrganizationRepository) RemoveAccount(ctx context.Context, orgID, accountID uint64) error {
	return r.db.WithContext(ctx).Table("organization_accounts").
		Where("organization_id = ? AND account_id = ?", orgID, accountID).
		Delete(nil).Error
}

func (r *OrganizationRepository) GetAccounts(ctx context.Context, orgID uint64) ([]uint64, error) {
	var accounts []struct {
		AccountID uint64
	}
	err := r.db.WithContext(ctx).
		Table("organization_accounts").
		Where("organization_id = ?", orgID).
		Scan(&accounts).Error

	result := make([]uint64, len(accounts))
	for i, acc := range accounts {
		result[i] = acc.AccountID
	}

	return result, err
}

func (r *OrganizationRepository) IsMember(ctx context.Context, orgID, accountID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("organization_accounts").
		Where("organization_id = ? AND account_id = ?", orgID, accountID).
		Count(&count).Error
	return count > 0, err
}
