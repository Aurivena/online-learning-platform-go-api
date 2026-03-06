package adaptors

import "gorm.io/gorm"

type AccountPostgres struct {
	db *gorm.DB
}

func NewAccountPostgres(db *gorm.DB) *AccountPostgres {
	return &AccountPostgres{db: db}
}
