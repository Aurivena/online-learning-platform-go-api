package di

import (
	"online-learning-platform-go-api/internal/user/adaptors"
	"online-learning-platform-go-api/internal/user/usecase"

	"gorm.io/gorm"
)

type Provider struct {
	db *gorm.DB
}

func NewProvider(db *gorm.DB) *Provider {
	return &Provider{db: db}
}

func (p *Provider) User() *usecase.AccountUseCase {
	return usecase.NewAccountUseCase(adaptors.NewRepository(p.db))
}
