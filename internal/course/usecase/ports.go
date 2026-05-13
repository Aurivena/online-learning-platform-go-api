package usecase

import (
	"context"
	"online-learning-platform-go-api/internal/course/entity"
)

type CourseRepository interface {
	Create(ctx context.Context, course *entity.Course) error
	GetByID(ctx context.Context, id uint64) (*entity.Course, error)
	GetByOrganization(ctx context.Context, orgID uint64) ([]entity.Course, error)
	Update(ctx context.Context, course *entity.Course) error
	Delete(ctx context.Context, id uint64) error
	AddModule(ctx context.Context, courseID, moduleID uint64, index int) error
	RemoveModule(ctx context.Context, courseID, moduleID uint64) error
	NextModuleIndex(ctx context.Context, courseID uint64) (int, error)
	ReorderModules(ctx context.Context, courseID uint64, moduleIDsInOrder []uint64) error
}

type ModuleRepository interface {
	Create(ctx context.Context, module *entity.Module) error
	GetByID(ctx context.Context, id uint64) (*entity.Module, error)
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
