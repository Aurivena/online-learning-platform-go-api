package adaptors

import (
	"context"
	"online-learning-platform-go-api/internal/user/entity"
)

func (r *Repository) Create(ctx context.Context, account *entity.Account) error {
	err := r.db.Create(&account).Error
	if err != nil {
		return err
	}

	return nil
}
