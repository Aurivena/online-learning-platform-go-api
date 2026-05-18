package entity

import "time"

type Organization struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(125);not null" json:"title"`
	Tag         string    `gorm:"type:varchar(15);not null;unique" json:"tag"`
	Description string    `gorm:"type:text;not null" json:"description"`
	ImageURL    string    `gorm:"column:image_url;type:text" json:"image_url"`
	HeaderTitle string    `gorm:"column:header_title;type:varchar(80)" json:"header_title"`
	OwnerID     uint64    `gorm:"type:bigint;not null" json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type OrganizationAccount struct {
	OrganizationID uint64
	AccountID      uint64
}
