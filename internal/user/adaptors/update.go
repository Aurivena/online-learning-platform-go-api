package adaptors

import (
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"
	"time"
)

func (a *AccountPostgres) Update(req dto.UpdateRequest, id int) error {
	var account entity.Account
	if err := a.db.First(&account, id).Error; err != nil {
		return err
	}
	account.Username = req.Username
	account.Email = req.Email
	account.Password = req.Password
	account.Role = req.Role
	account.UpdatedAt = time.Now()
	return a.db.Save(&account).Error
}
