package usecase

import (
	"context"
	"net/http"
	"online-learning-platform-go-api/internal/course/dto"
	"online-learning-platform-go-api/internal/course/entity"

	"github.com/Aurivena/spond/v4/netsp"
)

type CourseUseCaseInterface interface {
	CreateCourse(ctx context.Context, ownerID uint64, input dto.CreateCourseRequest) (*entity.Course, *netsp.Response[netsp.ErrorDetail])
	GetCourse(ctx context.Context, id uint64) (*entity.Course, *netsp.Response[netsp.ErrorDetail])
	ListCourses(ctx context.Context, orgID uint64) ([]entity.Course, *netsp.Response[netsp.ErrorDetail])
	UpdateCourse(ctx context.Context, id uint64, input dto.UpdateCourseRequest) *netsp.Response[netsp.ErrorDetail]
	DeleteCourse(ctx context.Context, id uint64) *netsp.Response[netsp.ErrorDetail]
	AddModuleToCourse(ctx context.Context, courseID uint64, input dto.AddModuleToCourseRequest) *netsp.Response[netsp.ErrorDetail]
	RemoveModuleFromCourse(ctx context.Context, courseID, moduleID uint64) *netsp.Response[netsp.ErrorDetail]
}

type CourseUseCase struct {
	courseRepo CourseRepository
	moduleRepo ModuleRepository
}

func NewCourseUseCase(courseRepo CourseRepository, moduleRepo ModuleRepository) *CourseUseCase {
	return &CourseUseCase{
		courseRepo: courseRepo,
		moduleRepo: moduleRepo,
	}
}

func (uc *CourseUseCase) CreateCourse(ctx context.Context, ownerID uint64, input dto.CreateCourseRequest) (*entity.Course, *netsp.Response[netsp.ErrorDetail]) {
	course := &entity.Course{
		Title:          input.Title,
		Description:    input.Description,
		Owner:          ownerID,
		OrganizationID: input.OrganizationID,
	}

	if err := uc.courseRepo.Create(ctx, course); err != nil {
		return nil, netsp.BuildError(
			http.StatusBadRequest,
			netsp.ErrorDetail{
				Title:    "Failed to Create Course",
				Message:  "Could not create course in database",
				Solution: "Please check your input and try again",
			},
		)
	}

	return course, nil
}

func (uc *CourseUseCase) GetCourse(ctx context.Context, id uint64) (*entity.Course, *netsp.Response[netsp.ErrorDetail]) {
	course, err := uc.courseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Course Not Found",
				Message:  "The requested course does not exist",
				Solution: "Please check the course ID and try again",
			},
		)
	}

	return course, nil
}

func (uc *CourseUseCase) ListCourses(ctx context.Context, orgID uint64) ([]entity.Course, *netsp.Response[netsp.ErrorDetail]) {
	courses, err := uc.courseRepo.GetByOrganization(ctx, orgID)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Fetch Courses",
				Message:  "Could not retrieve courses from database",
				Solution: "Please try again later",
			},
		)
	}

	if courses == nil {
		courses = []entity.Course{}
	}

	return courses, nil
}

func (uc *CourseUseCase) UpdateCourse(ctx context.Context, id uint64, input dto.UpdateCourseRequest) *netsp.Response[netsp.ErrorDetail] {
	course, err := uc.courseRepo.GetByID(ctx, id)
	if err != nil {
		return netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Course Not Found",
				Message:  "The requested course does not exist",
				Solution: "Please check the course ID and try again",
			},
		)
	}

	if input.Title != "" {
		course.Title = input.Title
	}
	if input.Description != "" {
		course.Description = input.Description
	}

	if err := uc.courseRepo.Update(ctx, course); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Update Course",
				Message:  "Could not update course in database",
				Solution: "Please try again later",
			},
		)
	}

	return nil
}

func (uc *CourseUseCase) DeleteCourse(ctx context.Context, id uint64) *netsp.Response[netsp.ErrorDetail] {
	if err := uc.courseRepo.Delete(ctx, id); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Delete Course",
				Message:  "Could not delete course from database",
				Solution: "Please try again later",
			},
		)
	}

	return nil
}

func (uc *CourseUseCase) AddModuleToCourse(ctx context.Context, courseID uint64, input dto.AddModuleToCourseRequest) *netsp.Response[netsp.ErrorDetail] {
	if _, err := uc.moduleRepo.GetByID(ctx, input.ModuleID); err != nil {
		return netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Module Not Found",
				Message:  "The requested module does not exist",
				Solution: "Please check the module ID and try again",
			},
		)
	}

	if err := uc.courseRepo.AddModule(ctx, courseID, input.ModuleID, input.Index); err != nil {
		return netsp.BuildError(
			http.StatusBadRequest,
			netsp.ErrorDetail{
				Title:    "Failed to Add Module",
				Message:  "Could not add module to course",
				Solution: "Module may already be in course or index is invalid",
			},
		)
	}

	return nil
}

func (uc *CourseUseCase) RemoveModuleFromCourse(ctx context.Context, courseID, moduleID uint64) *netsp.Response[netsp.ErrorDetail] {
	if err := uc.courseRepo.RemoveModule(ctx, courseID, moduleID); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Remove Module",
				Message:  "Could not remove module from course",
				Solution: "Please try again later",
			},
		)
	}

	return nil
}
