package adaptors

import (
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"
)

func (a *AccountPostgres) Get(id int) (*dto.AccountResponse, error) {
	var account entity.Account

	if err := a.db.First(&account, id).Error; err != nil {
		return nil, err
	}

	response := dto.AccountResponse{
		ID:        int(account.ID),
		Email:     account.Email,
		Username:  account.Username,
		Role:      string(account.Role),
		CreatedAt: account.CreatedAt,
	}

	return &response, nil
}
