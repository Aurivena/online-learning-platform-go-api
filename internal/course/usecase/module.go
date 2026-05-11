package usecase

import (
	"context"
	"net/http"
	"online-learning-platform-go-api/internal/course/dto"
	"online-learning-platform-go-api/internal/course/entity"

	"github.com/Aurivena/spond/v4/netsp"
)

type ModuleUseCaseInterface interface {
	CreateModule(ctx context.Context, input dto.CreateModuleRequest) (*entity.Module, *netsp.Response[netsp.ErrorDetail])
	GetModule(ctx context.Context, id uint64) (*entity.Module, *netsp.Response[netsp.ErrorDetail])
	UpdateModule(ctx context.Context, id uint64, input dto.UpdateModuleRequest) *netsp.Response[netsp.ErrorDetail]
	DeleteModule(ctx context.Context, id uint64) *netsp.Response[netsp.ErrorDetail]
	AddSlideToModule(ctx context.Context, moduleID uint64, input dto.AddSlideTModuleRequest) *netsp.Response[netsp.ErrorDetail]
	RemoveSlideFromModule(ctx context.Context, moduleID, slideID uint64) *netsp.Response[netsp.ErrorDetail]
}

type ModuleUseCase struct {
	moduleRepo ModuleRepository
	slideRepo  SlideRepository
}

func NewModuleUseCase(moduleRepo ModuleRepository, slideRepo SlideRepository) *ModuleUseCase {
	return &ModuleUseCase{
		moduleRepo: moduleRepo,
		slideRepo:  slideRepo,
	}
}

func (uc *ModuleUseCase) CreateModule(ctx context.Context, input dto.CreateModuleRequest) (*entity.Module, *netsp.Response[netsp.ErrorDetail]) {
	module := &entity.Module{
		Title: input.Title,
	}

	if err := uc.moduleRepo.Create(ctx, module); err != nil {
		return nil, netsp.BuildError(
			http.StatusBadRequest,
			netsp.ErrorDetail{
				Title:    "Failed to Create Module",
				Message:  "Could not create module in database",
				Solution: "Please check your input and try again",
			},
		)
	}

	return module, nil
}

func (uc *ModuleUseCase) GetModule(ctx context.Context, id uint64) (*entity.Module, *netsp.Response[netsp.ErrorDetail]) {
	module, err := uc.moduleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Module Not Found",
				Message:  "The requested module does not exist",
				Solution: "Please check the module ID and try again",
			},
		)
	}

	return module, nil
}

func (uc *ModuleUseCase) UpdateModule(ctx context.Context, id uint64, input dto.UpdateModuleRequest) *netsp.Response[netsp.ErrorDetail] {
	module, err := uc.moduleRepo.GetByID(ctx, id)
	if err != nil {
		return netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Module Not Found",
				Message:  "The requested module does not exist",
				Solution: "Please check the module ID and try again",
			},
		)
	}

	if input.Title != "" {
		module.Title = input.Title
	}

	if err := uc.moduleRepo.Update(ctx, module); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Update Module",
				Message:  "Could not update module in database",
				Solution: "Please try again later",
			},
		)
	}

	return nil
}

func (uc *ModuleUseCase) DeleteModule(ctx context.Context, id uint64) *netsp.Response[netsp.ErrorDetail] {
	if err := uc.moduleRepo.Delete(ctx, id); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Delete Module",
				Message:  "Could not delete module from database",
				Solution: "Please try again later",
			},
		)
	}

	return nil
}

func (uc *ModuleUseCase) AddSlideToModule(ctx context.Context, moduleID uint64, input dto.AddSlideTModuleRequest) *netsp.Response[netsp.ErrorDetail] {
	if _, err := uc.slideRepo.GetByID(ctx, input.SlideID); err != nil {
		return netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Slide Not Found",
				Message:  "The requested slide does not exist",
				Solution: "Please check the slide ID and try again",
			},
		)
	}

	if err := uc.moduleRepo.AddSlide(ctx, moduleID, input.SlideID, input.Index); err != nil {
		return netsp.BuildError(
			http.StatusBadRequest,
			netsp.ErrorDetail{
				Title:    "Failed to Add Slide",
				Message:  "Could not add slide to module",
				Solution: "Slide may already be in module or index is invalid",
			},
		)
	}

	return nil
}

func (uc *ModuleUseCase) RemoveSlideFromModule(ctx context.Context, moduleID, slideID uint64) *netsp.Response[netsp.ErrorDetail] {
	if err := uc.moduleRepo.RemoveSlide(ctx, moduleID, slideID); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Remove Slide",
				Message:  "Could not remove slide from module",
				Solution: "Please try again later",
			},
		)
	}

	return nil
}

type SlideUseCaseInterface interface {
	CreateSlide(ctx context.Context, input dto.CreateSlideRequest) (*entity.Slide, *netsp.Response[netsp.ErrorDetail])
	GetSlide(ctx context.Context, id uint64) (*entity.Slide, *netsp.Response[netsp.ErrorDetail])
	UpdateSlide(ctx context.Context, id uint64, input dto.UpdateSlideRequest) *netsp.Response[netsp.ErrorDetail]
	DeleteSlide(ctx context.Context, id uint64) *netsp.Response[netsp.ErrorDetail]
}

type SlideUseCase struct {
	slideRepo  SlideRepository
	moduleRepo ModuleRepository
}

func NewSlideUseCase(slideRepo SlideRepository, moduleRepo ModuleRepository) *SlideUseCase {
	return &SlideUseCase{
		slideRepo:  slideRepo,
		moduleRepo: moduleRepo,
	}
}

func (uc *SlideUseCase) CreateSlide(ctx context.Context, input dto.CreateSlideRequest) (*entity.Slide, *netsp.Response[netsp.ErrorDetail]) {
	slide := &entity.Slide{
		Title:       input.Title,
		Description: input.Description,
		SlideType:   input.SlideType,
		Payload:     input.Payload,
	}

	if err := uc.slideRepo.Create(ctx, slide); err != nil {
		return nil, netsp.BuildError(
			http.StatusBadRequest,
			netsp.ErrorDetail{
				Title:    "Failed to Create Slide",
				Message:  "Could not create slide in database",
				Solution: "Please check your input and try again",
			},
		)
	}

	if input.ModuleID > 0 {
		if err := uc.moduleRepo.AddSlide(ctx, input.ModuleID, slide.ID, 0); err != nil {
			return nil, netsp.BuildError(
				http.StatusBadRequest,
				netsp.ErrorDetail{
					Title:    "Failed to Add Slide to Module",
					Message:  "Slide created but could not add to module",
					Solution: "Module ID may be invalid",
				},
			)
		}
	}

	return slide, nil
}

func (uc *SlideUseCase) GetSlide(ctx context.Context, id uint64) (*entity.Slide, *netsp.Response[netsp.ErrorDetail]) {
	slide, err := uc.slideRepo.GetByID(ctx, id)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Slide Not Found",
				Message:  "The requested slide does not exist",
				Solution: "Please check the slide ID and try again",
			},
		)
	}

	return slide, nil
}

func (uc *SlideUseCase) UpdateSlide(ctx context.Context, id uint64, input dto.UpdateSlideRequest) *netsp.Response[netsp.ErrorDetail] {
	slide, err := uc.slideRepo.GetByID(ctx, id)
	if err != nil {
		return netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Slide Not Found",
				Message:  "The requested slide does not exist",
				Solution: "Please check the slide ID and try again",
			},
		)
	}

	if input.Title != "" {
		slide.Title = input.Title
	}
	if input.Description != "" {
		slide.Description = input.Description
	}
	if input.SlideType != "" {
		slide.SlideType = input.SlideType
	}
	if len(input.Payload) > 0 {
		slide.Payload = input.Payload
	}

	if err := uc.slideRepo.Update(ctx, slide); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Update Slide",
				Message:  "Could not update slide in database",
				Solution: "Please try again later",
			},
		)
	}

	return nil
}

func (uc *SlideUseCase) DeleteSlide(ctx context.Context, id uint64) *netsp.Response[netsp.ErrorDetail] {
	if err := uc.slideRepo.Delete(ctx, id); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Delete Slide",
				Message:  "Could not delete slide from database",
				Solution: "Please try again later",
			},
		)
	}

	return nil
}
