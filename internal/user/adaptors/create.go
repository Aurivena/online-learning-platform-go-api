package adaptors

import (
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"
	"time"
)

func (a *AccountPostgres) Create(req dto.RegistrationRequest) error {
	account := entity.Account{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		Role:      req.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return a.db.Create(&account).Error
}
