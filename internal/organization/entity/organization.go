package entity

import "time"

type Organization struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(125);not null" json:"title"`
	Tag         string    `gorm:"type:varchar(15);not null;unique" json:"tag"`
	Description string    `gorm:"type:text;not null" json:"description"`
	OwnerID     uint64    `gorm:"type:bigint;not null" json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OrganizationAccount struct {
	OrganizationID uint64
	AccountID      uint64
}
