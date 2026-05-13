package adaptors

import (
	"context"
	"online-learning-platform-go-api/internal/user/entity"
)

func (r *AccountRepository) Create(ctx context.Context, account *entity.Account) error {
	account.Role = entity.RoleUser
	err := r.db.Create(&account).Error
	if err != nil {
		return err
	}

	return nil
}
