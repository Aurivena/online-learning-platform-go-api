package usecase

import (
	"context"
	"online-learning-platform-go-api/internal/course/entity"
	"time"
)

type CourseRepository interface {
	Create(ctx context.Context, course *entity.Course) error
	GetByID(ctx context.Context, id uint64) (*entity.Course, error)
	GetByOrganization(ctx context.Context, orgID uint64) ([]entity.Course, error)
	GetAll(ctx context.Context) ([]entity.Course, error)
	Update(ctx context.Context, course *entity.Course) error
	Delete(ctx context.Context, id uint64) error
	SetOrganizations(ctx context.Context, courseID uint64, organizationIDs []uint64) error
	GetOrganizationIDs(ctx context.Context, courseID uint64) ([]uint64, error)
	IsLinkedToOrganization(ctx context.Context, courseID, orgID uint64) (bool, error)
	AddModule(ctx context.Context, courseID, moduleID uint64, index int) error
	RemoveModule(ctx context.Context, courseID, moduleID uint64) error
	NextModuleIndex(ctx context.Context, courseID uint64) (int, error)
	ReorderModules(ctx context.Context, courseID uint64, moduleIDsInOrder []uint64) error
}

type ModuleRepository interface {
	Create(ctx context.Context, module *entity.Module) error
	GetByID(ctx context.Context, id uint64) (*entity.Module, error)
	GetCourseIDByModuleID(ctx context.Context, moduleID uint64) (uint64, error)
	Update(ctx context.Context, module *entity.Module) error
	Delete(ctx context.Context, id uint64) error
	AddSlide(ctx context.Context, moduleID, slideID uint64, index int) error
	RemoveSlide(ctx context.Context, moduleID, slideID uint64) error
	NextSlideIndex(ctx context.Context, moduleID uint64) (int, error)
	ReorderSlides(ctx context.Context, moduleID uint64, slideIDsInOrder []uint64) error
}

type SlideRepository interface {
	Create(ctx context.Context, slide *entity.Slide) error
	GetByID(ctx context.Context, id uint64) (*entity.Slide, error)
	Update(ctx context.Context, slide *entity.Slide) error
	Delete(ctx context.Context, id uint64) error
}

type TestResultRecord struct {
	AccountID        uint64
	ModuleID         uint64
	SlideID          uint64
	SelectedOptionID uint64
	IsRight          bool
	Attempts         int
	FirstAttemptAt   time.Time
	LastAttemptAt    time.Time
	CourseID         uint64
	CourseTitle      string
	ModuleTitle      string
	SlideTitle       string
	AccountEmail     string
	AccountUsername  string
}

type CourseProgressRecord struct {
	AccountID       uint64
	AccountEmail    string
	AccountUsername string
	CourseID        uint64
	CourseTitle     string
	TotalTests      int
	AttemptedTests  int
	PassedTests     int
	Completed       bool
	LastActivityAt  time.Time
}

type TestResultRepository interface {
	Upsert(ctx context.Context, accountID, moduleID, slideID, selectedOptionID uint64, isRight bool) error
	GetByAccountAndSlide(ctx context.Context, accountID, slideID uint64) (*TestResultRecord, error)
	List(ctx context.Context, orgID *uint64) ([]TestResultRecord, error)
	ListCourseProgress(ctx context.Context, orgID *uint64) ([]CourseProgressRecord, error)
}
