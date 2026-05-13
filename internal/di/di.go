package di

import (
	"online-learning-platform-go-api/internal/gateway"
	userAdaptors "online-learning-platform-go-api/internal/user/adaptors"
	userUsecase "online-learning-platform-go-api/internal/user/usecase"

	courseAdaptors "online-learning-platform-go-api/internal/course/adaptors"
	courseUsecase "online-learning-platform-go-api/internal/course/usecase"
	orgAdaptors "online-learning-platform-go-api/internal/organization/adaptors"
	orgUsecase "online-learning-platform-go-api/internal/organization/usecase"

	"gorm.io/gorm"
)

type Provider struct {
	db *gorm.DB
}

func NewProvider(db *gorm.DB) *Provider {
	return &Provider{db: db}
}

func (p *Provider) User() userUsecase.AccountUseCaseInterface {
	userRepo := userAdaptors.NewAccountRepository(p.db)
	orgRepo := orgAdaptors.NewOrganizationRepository(p.db)
	return userUsecase.NewAccountUseCase(userRepo, orgRepo)
}

func (p *Provider) Organization() orgUsecase.OrganizationUseCaseInterface {
	return orgUsecase.NewOrganizationUseCase(orgAdaptors.NewOrganizationRepository(p.db))
}

func (p *Provider) OrganizationGateway() *gateway.OrganizationGateway {
	orgUC := p.Organization()
	userRepo := userAdaptors.NewAccountRepository(p.db)
	return gateway.NewOrganizationGateway(orgUC, userRepo)
}

func (p *Provider) Course() courseUsecase.CourseUseCaseInterface {
	courseRepo := courseAdaptors.NewCourseRepository(p.db)
	moduleRepo := courseAdaptors.NewModuleRepository(p.db)
	return courseUsecase.NewCourseUseCase(courseRepo, moduleRepo)
}

func (p *Provider) Module() courseUsecase.ModuleUseCaseInterface {
	moduleRepo := courseAdaptors.NewModuleRepository(p.db)
	slideRepo := courseAdaptors.NewSlideRepository(p.db)
	return courseUsecase.NewModuleUseCase(moduleRepo, slideRepo)
}

func (p *Provider) Slide() courseUsecase.SlideUseCaseInterface {
	slideRepo := courseAdaptors.NewSlideRepository(p.db)
	moduleRepo := courseAdaptors.NewModuleRepository(p.db)
	return courseUsecase.NewSlideUseCase(slideRepo, moduleRepo)
}
