package entity

import (
	"time"
)

type Role string

const (
	RoleUser  Role = "USER"
	RoleAdmin Role = "ADMIN"
)

type Account struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"column:email;type:varchar(255);uniqueIndex;not null"`
	Username     string `gorm:"column:username;type:varchar(125)"`
	PasswordHash string `gorm:"column:password_hash;type:varchar(255);not null"`
	Role         Role   `gorm:"column:role;type:roles;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
