package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type SlideType string

const (
	SlideTypeText     SlideType = "TEXT"
	SlideTypeVideoURL SlideType = "VIDEO_URL"
	SlideTypeTest     SlideType = "TEST"
	SlideTypeFile     SlideType = "FILE"
)

type Course struct {
	ID              uint64    `gorm:"primaryKey" json:"id"`
	Title           string    `gorm:"type:varchar(255);not null" json:"title"`
	Description     string    `gorm:"type:text;not null" json:"description"`
	Owner           uint64    `gorm:"type:bigint" json:"owner"`
	OrganizationID  uint64    `gorm:"type:bigint" json:"organization_id"`
	CreatedAt       time.Time `json:"created_at"`
	Modules         []Module  `gorm:"many2many:course_modules;joinForeignKey:course_id;joinReferences:module_id" json:"modules,omitempty"`
	OrganizationIDs []uint64  `gorm:"-" json:"organization_ids,omitempty"`
}

type Module struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(125);not null" json:"title"`
	CreatedAt time.Time `json:"created_at"`
	Slides    []Slide   `gorm:"many2many:module_slides;joinForeignKey:module_id;joinReferences:slide_id" json:"slides,omitempty"`
}

type Slide struct {
	ID          uint64      `gorm:"primaryKey" json:"id"`
	Title       string      `gorm:"type:varchar(255)" json:"title"`
	Description string      `gorm:"type:text" json:"description"`
	SlideType   SlideType   `gorm:"type:slide_variation;not null" json:"slide_type"`
	Payload     PayloadJSON `gorm:"type:jsonb" json:"payload"`
	CreatedAt   time.Time   `json:"created_at"`
}

type PayloadJSON map[string]interface{}

func (p PayloadJSON) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *PayloadJSON) Scan(value interface{}) error {
	bytes, _ := value.([]byte)
	return json.Unmarshal(bytes, &p)
}

type CourseModule struct {
	CourseID uint64
	ModuleID uint64
	Index    int
}

type ModuleSlide struct {
	ModuleID uint64
	SlideID  uint64
	Index    int
}
