package adaptors

import (
	"context"
	"online-learning-platform-go-api/internal/user/entity"
)

func (r *Repository) Get(ctx context.Context, id int) (*entity.Account, error) {
	var account entity.Account

	if err := r.db.First(&account, id).Error; err != nil {
		return nil, err
	}

	return &account, nil
}
