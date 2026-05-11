package usecase

import (
	"context"
	"net/http"
	"online-learning-platform-go-api/internal/organization/dto"
	"online-learning-platform-go-api/internal/organization/entity"

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
}

type OrganizationUseCase struct {
	repo OrganizationRepository
}

func NewOrganizationUseCase(repo OrganizationRepository) *OrganizationUseCase {
	return &OrganizationUseCase{repo: repo}
}

func (uc *OrganizationUseCase) CreateOrganization(ctx context.Context, ownerID uint64, input dto.CreateOrganizationRequest) (*entity.Organization, *netsp.Response[netsp.ErrorDetail]) {
	org := &entity.Organization{
		Title:       input.Title,
		Tag:         input.Tag,
		Description: input.Description,
		OwnerID:     ownerID,
	}

	if err := uc.repo.Create(ctx, org); err != nil {
		return nil, netsp.BuildError(
			http.StatusBadRequest,
			netsp.ErrorDetail{
				Title:    "Failed to Create Organization",
				Message:  "Could not create organization in database",
				Solution: "Tag may already exist or check your input",
			},
		)
	}

	if err := uc.repo.AddAccount(ctx, org.ID, ownerID); err != nil {
		uc.repo.Delete(ctx, org.ID)
		return nil, netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Add Owner to Organization",
				Message:  "Could not add owner account to organization",
				Solution: "Please try again later",
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
				Title:    "Organization Not Found",
				Message:  "The requested organization does not exist",
				Solution: "Please check the organization ID and try again",
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
				Title:    "Organization Not Found",
				Message:  "No organization found with this tag",
				Solution: "Please check the tag and try again",
			},
		)
	}

	return org, nil
}

func (uc *OrganizationUseCase) ListMyOrganizations(ctx context.Context, ownerID uint64) ([]entity.Organization, *netsp.Response[netsp.ErrorDetail]) {
	orgs, err := uc.repo.GetByOwner(ctx, ownerID)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Fetch Organizations",
				Message:  "Could not retrieve organizations from database",
				Solution: "Please try again later",
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
				Title:    "Failed to Fetch Organizations",
				Message:  "Could not retrieve organizations from database",
				Solution: "Please try again later",
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
				Title:    "Organization Not Found",
				Message:  "The requested organization does not exist",
				Solution: "Please check the organization ID and try again",
			},
		)
	}

	if input.Title != "" {
		org.Title = input.Title
	}
	if input.Description != "" {
		org.Description = input.Description
	}

	if err := uc.repo.Update(ctx, org); err != nil {
		return netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Update Organization",
				Message:  "Could not update organization in database",
				Solution: "Please try again later",
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
				Title:    "Failed to Delete Organization",
				Message:  "Could not delete organization from database",
				Solution: "Please try again later",
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
				Title:    "Organization Not Found",
				Message:  "The requested organization does not exist",
				Solution: "Please check the organization ID and try again",
			},
		)
	}

	if err := uc.repo.AddAccount(ctx, orgID, accountID); err != nil {
		return netsp.BuildError(
			http.StatusBadRequest,
			netsp.ErrorDetail{
				Title:    "Failed to Add Account",
				Message:  "Could not add account to organization",
				Solution: "Account may already be in organization",
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
				Title:    "Failed to Remove Account",
				Message:  "Could not remove account from organization",
				Solution: "Please try again later",
			},
		)
	}

	return nil
}
