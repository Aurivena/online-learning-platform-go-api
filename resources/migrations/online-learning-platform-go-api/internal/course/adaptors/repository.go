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
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(course).Error; err != nil {
			return err
		}
		if course.OrganizationID == 0 {
			return nil
		}
		return tx.Table("course_organizations").Create(map[string]interface{}{
			"course_id":       course.ID,
			"organization_id": course.OrganizationID,
		}).Error
	})
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
	ids, err := r.GetOrganizationIDs(ctx, course.ID)
	if err != nil {
		return nil, err
	}
	course.OrganizationIDs = ids
	return &course, nil
}

func (r *CourseRepository) GetByOrganization(ctx context.Context, orgID uint64) ([]entity.Course, error) {
	var courses []entity.Course
	if err := r.db.WithContext(ctx).
		Table("courses").
		Select("courses.*").
		Joins("inner join course_organizations co on co.course_id = courses.id").
		Where("co.organization_id = ?", orgID).
		Order("courses.created_at desc").
		Find(&courses).Error; err != nil {
		return nil, err
	}
	if err := r.loadCourses(ctx, courses); err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *CourseRepository) GetAll(ctx context.Context) ([]entity.Course, error) {
	var courses []entity.Course
	if err := r.db.WithContext(ctx).
		Order("created_at desc").
		Find(&courses).Error; err != nil {
		return nil, err
	}
	if err := r.loadCourses(ctx, courses); err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *CourseRepository) loadCourses(ctx context.Context, courses []entity.Course) error {
	for i := range courses {
		if err := r.db.WithContext(ctx).Model(&courses[i]).Order("course_modules.index").Association("Modules").Find(&courses[i].Modules); err != nil {
			return err
		}
		for j := range courses[i].Modules {
			if err := r.db.WithContext(ctx).Model(&courses[i].Modules[j]).Order("module_slides.index").Association("Slides").Find(&courses[i].Modules[j].Slides); err != nil {
				return err
			}
		}
		ids, err := r.GetOrganizationIDs(ctx, courses[i].ID)
		if err != nil {
			return err
		}
		courses[i].OrganizationIDs = ids
	}
	return nil
}

func (r *CourseRepository) Update(ctx context.Context, course *entity.Course) error {
	return r.db.WithContext(ctx).Model(course).Updates(course).Error
}

func (r *CourseRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.Course{}, id).Error
}

func (r *CourseRepository) SetOrganizations(ctx context.Context, courseID uint64, organizationIDs []uint64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("course_organizations").Where("course_id = ?", courseID).Delete(&struct{}{}).Error; err != nil {
			return err
		}
		for _, orgID := range organizationIDs {
			if orgID == 0 {
				continue
			}
			if err := tx.Table("course_organizations").Create(map[string]interface{}{
				"course_id":       courseID,
				"organization_id": orgID,
			}).Error; err != nil {
				return err
			}
		}
		if len(organizationIDs) > 0 {
			return tx.Model(&entity.Course{}).Where("id = ?", courseID).Update("organization_id", organizationIDs[0]).Error
		}
		return nil
	})
}

func (r *CourseRepository) GetOrganizationIDs(ctx context.Context, courseID uint64) ([]uint64, error) {
	var rows []struct {
		OrganizationID uint64 `gorm:"column:organization_id"`
	}
	if err := r.db.WithContext(ctx).
		Table("course_organizations").
		Select("organization_id").
		Where("course_id = ?", courseID).
		Order("organization_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	ids := make([]uint64, len(rows))
	for i, row := range rows {
		ids[i] = row.OrganizationID
	}
	return ids, nil
}

func (r *CourseRepository) IsLinkedToOrganization(ctx context.Context, courseID, orgID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("course_organizations").
		Where("course_id = ? AND organization_id = ?", courseID, orgID).
		Count(&count).Error
	return count > 0, err
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
		Delete(&struct{}{}).Error
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

func (r *ModuleRepository) GetCourseIDByModuleID(ctx context.Context, moduleID uint64) (uint64, error) {
	var row struct {
		CourseID uint64 `gorm:"column:course_id"`
	}
	if err := r.db.WithContext(ctx).
		Table("course_modules").
		Select("course_id").
		Where("module_id = ?", moduleID).
		Limit(1).
		Scan(&row).Error; err != nil {
		return 0, err
	}
	if row.CourseID == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return row.CourseID, nil
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
		Delete(&struct{}{}).Error
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
