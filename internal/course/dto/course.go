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
	ID              uint64           `json:"id"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	Owner           uint64           `json:"owner"`
	OrganizationID  uint64           `json:"organization_id"`
	OrganizationIDs []uint64         `json:"organization_ids,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	Modules         []ModuleResponse `json:"modules,omitempty"`
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

type UpdateCourseOrganizationsRequest struct {
	OrganizationIDs []uint64 `json:"organization_ids"`
}

type AddSlideTModuleRequest struct {
	Index int `json:"index" binding:"required,min=0"`
}

type AdminTestResultResponse struct {
	AccountID        uint64    `json:"account_id"`
	AccountEmail     string    `json:"account_email"`
	AccountUsername  string    `json:"account_username"`
	CourseID         uint64    `json:"course_id"`
	CourseTitle      string    `json:"course_title"`
	ModuleID         uint64    `json:"module_id"`
	ModuleTitle      string    `json:"module_title"`
	SlideID          uint64    `json:"slide_id"`
	SlideTitle       string    `json:"slide_title"`
	SelectedOptionID uint64    `json:"selected_option_id"`
	IsRight          bool      `json:"is_right"`
	Attempts         int       `json:"attempts"`
	FirstAttemptAt   time.Time `json:"first_attempt_at"`
	LastAttemptAt    time.Time `json:"last_attempt_at"`
}

type AdminCourseProgressResponse struct {
	AccountID       uint64    `json:"account_id"`
	AccountEmail    string    `json:"account_email"`
	AccountUsername string    `json:"account_username"`
	CourseID        uint64    `json:"course_id"`
	CourseTitle     string    `json:"course_title"`
	TotalTests      int       `json:"total_tests"`
	AttemptedTests  int       `json:"attempted_tests"`
	PassedTests     int       `json:"passed_tests"`
	Completed       bool      `json:"completed"`
	LastActivityAt  time.Time `json:"last_activity_at"`
}
