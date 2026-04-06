package usecase

import (
	"context"
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"
	"time"
)

func (uc *AccountUseCase) Update(ctx context.Context, req dto.UpdateRequest, id int) error {
	account := &entity.Account{
		ID:        uint(id),
		Username:  req.Username,
		Email:     req.Email,
		Role:      req.Role,
		UpdatedAt: time.Now().UTC(),
	}
	return uc.repo.Update(ctx, account)
}
