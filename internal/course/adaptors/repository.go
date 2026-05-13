package adaptors

import (
	"context"
	"database/sql"
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
	if err := r.db.WithContext(ctx).First(&course, id).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&course).Order("course_modules.index").Association("Modules").Find(&course.Modules); err != nil {
		return nil, err
	}
	for i := range course.Modules {
		if err := r.db.WithContext(ctx).Model(&course.Modules[i]).Order("module_slides.index").Association("Slides").Find(&course.Modules[i].Slides); err != nil {
			return nil, err
		}
	}
	return &course, nil
}

func (r *CourseRepository) GetByOrganization(ctx context.Context, orgID uint64) ([]entity.Course, error) {
	var courses []entity.Course
	if err := r.db.WithContext(ctx).
		Where("organization_id = ?", orgID).
		Order("created_at desc").
		Find(&courses).Error; err != nil {
		return nil, err
	}
	for i := range courses {
		if err := r.db.WithContext(ctx).Model(&courses[i]).Order("course_modules.index").Association("Modules").Find(&courses[i].Modules); err != nil {
			return nil, err
		}
		for j := range courses[i].Modules {
			if err := r.db.WithContext(ctx).Model(&courses[i].Modules[j]).Order("module_slides.index").Association("Slides").Find(&courses[i].Modules[j].Slides); err != nil {
				return nil, err
			}
		}
	}
	return courses, nil
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

func (r *CourseRepository) NextModuleIndex(ctx context.Context, courseID uint64) (int, error) {
	var max sql.NullInt64
	err := r.db.WithContext(ctx).Raw(
		`SELECT MAX("index") FROM course_modules WHERE course_id = ?`,
		courseID,
	).Scan(&max).Error
	if err != nil {
		return 0, err
	}
	if !max.Valid {
		return 1, nil
	}
	return int(max.Int64) + 1, nil
}

func (r *CourseRepository) ReorderModules(ctx context.Context, courseID uint64, moduleIDsInOrder []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, mid := range moduleIDsInOrder {
			if res := tx.Exec(
				`UPDATE course_modules SET "index" = ? WHERE course_id = ? AND module_id = ?`,
				-(i + 1), courseID, mid,
			); res.Error != nil {
				return res.Error
			} else if res.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}
		for i, mid := range moduleIDsInOrder {
			if res := tx.Exec(
				`UPDATE course_modules SET "index" = ? WHERE course_id = ? AND module_id = ?`,
				i+1, courseID, mid,
			); res.Error != nil {
				return res.Error
			} else if res.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}
		return nil
	})
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
	if err := r.db.WithContext(ctx).First(&module, id).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&module).Order("module_slides.index").Association("Slides").Find(&module.Slides); err != nil {
		return nil, err
	}
	return &module, nil
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

func (r *ModuleRepository) NextSlideIndex(ctx context.Context, moduleID uint64) (int, error) {
	var max sql.NullInt64
	err := r.db.WithContext(ctx).Raw(
		`SELECT MAX("index") FROM module_slides WHERE module_id = ?`,
		moduleID,
	).Scan(&max).Error
	if err != nil {
		return 0, err
	}
	if !max.Valid {
		return 1, nil
	}
	return int(max.Int64) + 1, nil
}

func (r *ModuleRepository) ReorderSlides(ctx context.Context, moduleID uint64, slideIDsInOrder []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for i, sid := range slideIDsInOrder {
			if res := tx.Exec(
				`UPDATE module_slides SET "index" = ? WHERE module_id = ? AND slide_id = ?`,
				-(i + 1), moduleID, sid,
			); res.Error != nil {
				return res.Error
			} else if res.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}
		for i, sid := range slideIDsInOrder {
			if res := tx.Exec(
				`UPDATE module_slides SET "index" = ? WHERE module_id = ? AND slide_id = ?`,
				i+1, moduleID, sid,
			); res.Error != nil {
				return res.Error
			} else if res.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		}
		return nil
	})
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
