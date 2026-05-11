package adaptors

import (
	"context"
	"online-learning-platform-go-api/internal/user/entity"
)

func (r *Repository) Update(ctx context.Context, account *entity.Account) error {
	return r.db.Save(&account).Error
}
