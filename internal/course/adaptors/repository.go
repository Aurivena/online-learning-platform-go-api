package adaptors

import (
	"context"
	"online-learning-platform-go-api/internal/course/entity"

	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

func (r *CourseRepository) Create(ctx context.Context, course *entity.Course) error {
	return r.db.WithContext(ctx).Create(course).Error
}

func (r *CourseRepository) GetByID(ctx context.Context, id uint64) (*entity.Course, error) {
	var course entity.Course
	err := r.db.WithContext(ctx).
		Preload("Modules", func(db *gorm.DB) *gorm.DB {
			return db.Order("course_modules.index")
		}).
		Preload("Modules.Slides", func(db *gorm.DB) *gorm.DB {
			return db.Order("module_slides.index")
		}).
		First(&course, id).Error
	return &course, err
}

func (r *CourseRepository) GetByOrganization(ctx context.Context, orgID uint64) ([]entity.Course, error) {
	var courses []entity.Course
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Preload("Modules", func(db *gorm.DB) *gorm.DB {
			return db.Order("course_modules.index")
		}).
		Preload("Modules.Slides", func(db *gorm.DB) *gorm.DB {
			return db.Order("module_slides.index")
		}).
		Order("created_at desc").
		Find(&courses).Error
	return courses, err
}

func (r *CourseRepository) Update(ctx context.Context, course *entity.Course) error {
	return r.db.WithContext(ctx).Model(course).Updates(course).Error
}

func (r *CourseRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.Course{}, id).Error
}

func (r *CourseRepository) AddModule(ctx context.Context, courseID, moduleID uint64, index int) error {
	return r.db.WithContext(ctx).Table("course_modules").Create(map[string]interface{}{
		"course_id": courseID,
		"module_id": moduleID,
		"index":     index,
	}).Error
}

func (r *CourseRepository) RemoveModule(ctx context.Context, courseID, moduleID uint64) error {
	return r.db.WithContext(ctx).Table("course_modules").
		Where("course_id = ? AND module_id = ?", courseID, moduleID).
		Delete(nil).Error
}

type ModuleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) *ModuleRepository {
	return &ModuleRepository{db: db}
}

func (r *ModuleRepository) Create(ctx context.Context, module *entity.Module) error {
	return r.db.WithContext(ctx).Create(module).Error
}

func (r *ModuleRepository) GetByID(ctx context.Context, id uint64) (*entity.Module, error) {
	var module entity.Module
	err := r.db.WithContext(ctx).
		Preload("Slides", func(db *gorm.DB) *gorm.DB {
			return db.Order("module_slides.index")
		}).
		First(&module, id).Error
	return &module, err
}

func (r *ModuleRepository) Update(ctx context.Context, module *entity.Module) error {
	return r.db.WithContext(ctx).Model(module).Updates(module).Error
}

func (r *ModuleRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.Module{}, id).Error
}

func (r *ModuleRepository) AddSlide(ctx context.Context, moduleID, slideID uint64, index int) error {
	return r.db.WithContext(ctx).Table("module_slides").Create(map[string]interface{}{
		"module_id": moduleID,
		"slide_id":  slideID,
		"index":     index,
	}).Error
}

func (r *ModuleRepository) RemoveSlide(ctx context.Context, moduleID, slideID uint64) error {
	return r.db.WithContext(ctx).Table("module_slides").
		Where("module_id = ? AND slide_id = ?", moduleID, slideID).
		Delete(nil).Error
}

type SlideRepository struct {
	db *gorm.DB
}

func NewSlideRepository(db *gorm.DB) *SlideRepository {
	return &SlideRepository{db: db}
}

func (r *SlideRepository) Create(ctx context.Context, slide *entity.Slide) error {
	return r.db.WithContext(ctx).Create(slide).Error
}

func (r *SlideRepository) GetByID(ctx context.Context, id uint64) (*entity.Slide, error) {
	var slide entity.Slide
	err := r.db.WithContext(ctx).First(&slide, id).Error
	return &slide, err
}

func (r *SlideRepository) Update(ctx context.Context, slide *entity.Slide) error {
	return r.db.WithContext(ctx).Model(slide).Updates(slide).Error
}

func (r *SlideRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.Slide{}, id).Error
}
