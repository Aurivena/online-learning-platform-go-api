package dto

import (
	"online-learning-platform-go-api/internal/course/entity"
	"time"
)

type CreateCourseRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"required,min=1"`
}

type UpdateCourseRequest struct {
	Title       string `json:"title" binding:"min=1,max=255"`
	Description string `json:"description" binding:"min=1"`
}

type CourseResponse struct {
	ID             uint64           `json:"id"`
	Title          string           `json:"title"`
	Description    string           `json:"description"`
	Owner          uint64           `json:"owner"`
	OrganizationID uint64           `json:"organization_id"`
	CreatedAt      time.Time        `json:"created_at"`
	Modules        []ModuleResponse `json:"modules,omitempty"`
}

type ModuleResponse struct {
	ID        uint64          `json:"id"`
	Title     string          `json:"title"`
	CreatedAt time.Time       `json:"created_at"`
	Slides    []SlideResponse `json:"slides,omitempty"`
}

type SlideResponse struct {
	ID          uint64             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	SlideType   entity.SlideType   `json:"slide_type"`
	Payload     entity.PayloadJSON `json:"payload"`
	CreatedAt   time.Time          `json:"created_at"`
}

type CreateModuleRequest struct {
	Title string `json:"title" binding:"required,min=1,max=125"`
}

type UpdateModuleRequest struct {
	Title string `json:"title" binding:"min=1,max=125"`
}

type CreateSlideRequest struct {
	Title       string             `json:"title" binding:"required,min=1,max=255"`
	Description string             `json:"description"`
	SlideType   entity.SlideType   `json:"slide_type" binding:"required,oneof=TEXT VIDEO_URL TEST FILE"`
	Payload     entity.PayloadJSON `json:"payload"`
}

type UpdateSlideRequest struct {
	Title       string             `json:"title" binding:"min=1,max=255"`
	Description string             `json:"description"`
	SlideType   entity.SlideType   `json:"slide_type" binding:"oneof=TEXT VIDEO_URL TEST FILE"`
	Payload     entity.PayloadJSON `json:"payload"`
}

type AddModuleToCourseRequest struct {
	Index int `json:"index" binding:"required,min=0"`
}

type AddSlideTModuleRequest struct {
	Index int `json:"index" binding:"required,min=0"`
}
