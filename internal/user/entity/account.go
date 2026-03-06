package entity

import (
	"online-learning-platform-go-api/internal/user/dto"
	"time"
)

type Account struct {
	ID        uint     `gorm:"primaryKey"`
	Email     string   `gorm:"unique;not null"`
	Username  string   `gorm:"unique;not null"`
	Password  string   `gorm:"not null"`
	Role      dto.Role `gorm:"type:varchar(20);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
