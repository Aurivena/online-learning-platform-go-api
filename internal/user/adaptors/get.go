package adaptors

import (
	"context"
	"online-learning-platform-go-api/internal/user/entity"
)

func (r *AccountRepository) GetByID(ctx context.Context, id uint64) (*entity.Account, error) {
	var account entity.Account

	if err := r.db.First(&account, id).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) GetByEmail(ctx context.Context, email string) (*entity.Account, error) {
	var account entity.Account

	if err := r.db.Where("email = ?", email).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) GetByUsername(ctx context.Context, username string) (*entity.Account, error) {
	var account entity.Account

	if err := r.db.Where("username = ?", username).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepository) GetAll(ctx context.Context) ([]entity.Account, error) {
	var accounts []entity.Account
	err := r.db.WithContext(ctx).
		Order("created_at desc").
		Find(&accounts).Error
	return accounts, err
}
