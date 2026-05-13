package usecase

import (
	"context"
	"net/http"
	"online-learning-platform-go-api/internal/user/dto"
	userEntity "online-learning-platform-go-api/internal/user/entity"

	"github.com/Aurivena/spond/v4/netsp"
)

func (uc *AccountUseCase) Get(ctx context.Context, id int) (*userEntity.Account, *netsp.Response[netsp.ErrorDetail]) {
	account, err := uc.repo.GetByID(ctx, uint64(id))
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Account Not Found",
				Message:  "An account with the specified ID was not found",
				Solution: "Please check the account ID and try again",
			},
		)
	}
	return account, nil
}

func (uc *AccountUseCase) GetProfile(ctx context.Context, id uint64) (*dto.UserProfileResponse, *netsp.Response[netsp.ErrorDetail]) {
	account, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusNotFound,
			netsp.ErrorDetail{
				Title:    "Account Not Found",
				Message:  "User account not found",
				Solution: "Please check the account ID and try again",
			},
		)
	}

	orgsData, err := uc.orgRepo.GetByAccountEntities(ctx, id)
	if err != nil {
		return nil, netsp.BuildError(
			http.StatusInternalServerError,
			netsp.ErrorDetail{
				Title:    "Failed to Fetch Organizations",
				Message:  "Could not retrieve user organizations",
				Solution: "Please try again later",
			},
		)
	}

	orgs := make([]dto.UserOrganizationResponse, 0, len(orgsData))
	for i := range orgsData {
		org := orgsData[i]
		orgs = append(orgs, dto.UserOrganizationResponse{
			ID:          org.ID,
			Title:       org.Title,
			Tag:         org.Tag,
			Description: org.Description,
			OwnerID:     org.OwnerID,
			CreatedAt:   org.CreatedAt,
		})
	}

	return &dto.UserProfileResponse{
		ID:            uint64(account.ID),
		Email:         account.Email,
		Username:      account.Username,
		Role:          account.Role,
		CreatedAt:     account.CreatedAt,
		Organizations: orgs,
	}, nil
}
