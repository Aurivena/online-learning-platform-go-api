package adaptors

import (
	"context"
	"online-learning-platform-go-api/internal/user/entity"
)

func (r *AccountRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.Account{}, id).Error
}
