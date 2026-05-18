package usecase

import (
	"context"
	"net/http"
	"online-learning-platform-go-api/internal/organization/dto"
	"online-learning-platform-go-api/internal/organization/entity"
	"strings"

	"github.com/Aurivena/spond/v4/netsp"
)

type OrganizationUseCaseInterface interface {
	CreateOrganization(ctx context.Context, ownerID uint64, input dto.CreateOrganizationRequest) (*entity.Organization, *netsp.Response[netsp.ErrorDetail])
	GetOrganization(ctx context.Context, id uint64) (*entity.Organization, *netsp.Response[netsp.ErrorDetail])
	GetOrganizationByTag(ctx context.Context, tag string) (*entity.Organization, *netsp.Response[netsp.ErrorDetail])
	ListMyOrganizations(ctx context.Context, ownerID uint64) ([]entity.Organization, *netsp.Response[netsp.ErrorDetail])
	ListAllOrganizations(ctx context.Context) ([]entity.Organization, *netsp.Response[netsp.ErrorDetail])
	UpdateOrganization(ctx context.Context, id uint64, input dto.UpdateOrganizationRequest) *netsp.Response[netsp.ErrorDetail]
	DeleteOrganization(ctx context.Context, id uint64) *netsp.Response[netsp.ErrorDetail]
	AddAccountToOrganization(ctx context.Context, orgID, accountID uint64) *netsp.Response[netsp.ErrorDetail]
	RemoveAccountFromOrganization(ctx context.Context, orgID, accountID uint64) *netsp.Response[netsp.ErrorDetail]
	ListOrganizationAccounts(ctx context.Context, orgID uint64) ([]dto.OrganizationAccountResponse, *netsp.Response[netsp.ErrorDetail])
}

type OrganizationUseCase struct {
	repo OrganizationRepository
}

func NewOrganizationUseCase(repo OrganizationRepository) *OrganizationUseCase {
	return &OrganizationUseCase{repo: repo}
}

func generateDetailitHeader(seed string) string {
	return `ДетаЛит`
}

func (uc *OrganizationUseCase) CreateOrganization(ctx context.Context, ownerID uint64, input dto.CreateOrganizationRequest) (*entity.Organization, *netsp.Response[netsp.ErrorDetail]) {
	headerTitle := strings.TrimSpace(input.HeaderTitle)
	if headerTitle == "" {
		headerTitle = generateDetailitHeader(input.Tag + ":" + input.Title)
	}
	org := &entity.Organization{
		Title:       input.Title,
		Tag:         input.Tag,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		HeaderTitle: headerTitle,
		OwnerID:     ownerID,
	}

	if err := uc.repo.Create(ctx, org); err != nil {
		return nil, netsp.BuildError(
			http.StatusBadRequest,
			netsp.ErrorDetail{
				Title:    "Не удалось создать подразделение",
				Message:  "Не удалось сохранить подразделение в базе данных",
				Solution: "Проверьте данные: тег может быть уже занят",
			},
		)
	}

	if err := uc.repo.AddAccount(ctx, org.ID, ownerID); err != nil {
		uc.repo.Delete(ctx, org.ID)
		return nil, netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Не удалось добавить владельца",
				Message:  "Не удалось добавить владельца в подразделение",
				Solution: "Повторите попытку позже",
			},
		)
	}

	return org, nil
}

func (uc *OrganizationUseCase) GetOrganization(ctx context.Context, id uint64) (*entity.Organization, *netsp.Response[netsp.ErrorDetail]) {
	org, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Подразделение не найдено",
				Message:  "Запрошенное подразделение не существует",
				Solution: "Проверьте ID подразделения и повторите попытку",
			},
		)
	}

	return org, nil
}

func (uc *OrganizationUseCase) GetOrganizationByTag(ctx context.Context, tag string) (*entity.Organization, *netsp.Response[netsp.ErrorDetail]) {
	org, err := uc.repo.GetByTag(ctx, tag)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Подразделение не найдено",
				Message:  "Подразделение с таким тегом не найдено",
				Solution: "Проверьте тег и повторите попытку",
			},
		)
	}

	return org, nil
}

func (uc *OrganizationUseCase) ListMyOrganizations(ctx context.Context, ownerID uint64) ([]entity.Organization, *netsp.Response[netsp.ErrorDetail]) {
	orgs, err := uc.repo.GetByAccountEntities(ctx, ownerID)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Не удалось загрузить подразделения",
				Message:  "Не удалось получить список подразделений из базы данных",
				Solution: "Повторите попытку позже",
			},
		)
	}

	if orgs == nil {
		orgs = []entity.Organization{}
	}

	return orgs, nil
}

func (uc *OrganizationUseCase) ListAllOrganizations(ctx context.Context) ([]entity.Organization, *netsp.Response[netsp.ErrorDetail]) {
	orgs, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Не удалось загрузить подразделения",
				Message:  "Не удалось получить список подразделений из базы данных",
				Solution: "Повторите попытку позже",
			},
		)
	}

	if orgs == nil {
		orgs = []entity.Organization{}
	}

	return orgs, nil
}

func (uc *OrganizationUseCase) UpdateOrganization(ctx context.Context, id uint64, input dto.UpdateOrganizationRequest) *netsp.Response[netsp.ErrorDetail] {
	org, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Подразделение не найдено",
				Message:  "Запрошенное подразделение не существует",
				Solution: "Проверьте ID подразделения и повторите попытку",
			},
		)
	}

	if input.Title != "" {
		org.Title = input.Title
	}
	if input.Tag != "" {
		org.Tag = input.Tag
	}
	if input.Description != "" {
		org.Description = input.Description
	}
	if input.ImageURL != nil {
		org.ImageURL = *input.ImageURL
	}
	if input.HeaderTitle != nil {
		headerTitle := strings.TrimSpace(*input.HeaderTitle)
		if headerTitle == "" {
			headerTitle = generateDetailitHeader(org.Tag + ":" + org.Title)
		}
		org.HeaderTitle = headerTitle
	}

	if err := uc.repo.Update(ctx, org); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Не удалось обновить подразделение",
				Message:  "Не удалось обновить подразделение в базе данных",
				Solution: "Повторите попытку позже",
			},
		)
	}

	return nil
}

func (uc *OrganizationUseCase) DeleteOrganization(ctx context.Context, id uint64) *netsp.Response[netsp.ErrorDetail] {
	if err := uc.repo.Delete(ctx, id); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Не удалось удалить подразделение",
				Message:  "Не удалось удалить подразделение из базы данных",
				Solution: "Повторите попытку позже",
			},
		)
	}

	return nil
}

func (uc *OrganizationUseCase) AddAccountToOrganization(ctx context.Context, orgID, accountID uint64) *netsp.Response[netsp.ErrorDetail] {
	if _, err := uc.repo.GetByID(ctx, orgID); err != nil {
		return netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Подразделение не найдено",
				Message:  "Запрошенное подразделение не существует",
				Solution: "Проверьте ID подразделения и повторите попытку",
			},
		)
	}

	if err := uc.repo.AddAccount(ctx, orgID, accountID); err != nil {
		return netsp.BuildError(
			http.StatusBadRequest,
			netsp.ErrorDetail{
				Title:    "Не удалось добавить сотрудника",
				Message:  "Не удалось добавить учетную запись в подразделение",
				Solution: "Возможно, сотрудник уже добавлен в это подразделение",
			},
		)
	}

	return nil
}

func (uc *OrganizationUseCase) RemoveAccountFromOrganization(ctx context.Context, orgID, accountID uint64) *netsp.Response[netsp.ErrorDetail] {
	if err := uc.repo.RemoveAccount(ctx, orgID, accountID); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Не удалось удалить сотрудника",
				Message:  "Не удалось удалить учетную запись из подразделения",
				Solution: "Повторите попытку позже",
			},
		)
	}

	return nil
}

func (uc *OrganizationUseCase) ListOrganizationAccounts(ctx context.Context, orgID uint64) ([]dto.OrganizationAccountResponse, *netsp.Response[netsp.ErrorDetail]) {
	if _, err := uc.repo.GetByID(ctx, orgID); err != nil {
		return nil, netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Подразделение не найдено",
				Message:  "Запрошенное подразделение не существует",
				Solution: "Проверьте ID подразделения и повторите попытку",
			},
		)
	}

	accounts, err := uc.repo.GetAccountEntities(ctx, orgID)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Не удалось загрузить сотрудников",
				Message:  "Не удалось получить список участников подразделения",
				Solution: "Повторите попытку позже",
			},
		)
	}

	result := make([]dto.OrganizationAccountResponse, len(accounts))
	for i, account := range accounts {
		result[i] = dto.OrganizationAccountResponse{
			ID:        uint64(account.ID),
			Email:     account.Email,
			Username:  account.Username,
			Role:      account.Role,
			CreatedAt: account.CreatedAt,
		}
	}
	return result, nil
}
