package usecase

import (
	"context"
	"online-learning-platform-go-api/internal/user/domain"
	"online-learning-platform-go-api/internal/user/dto"
	"online-learning-platform-go-api/internal/user/entity"
	"time"

	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
)

func (u *AccountUseCase) Registration(ctx context.Context, input dto.RegistrationRequest) *netsp.Response[netsp.ErrorDetail] {
	if input.OrganizationID != nil && *input.OrganizationID != 0 {
		if _, err := u.orgRepo.GetByID(ctx, *input.OrganizationID); err != nil {
			return netsp.BuildError(
				netstatus.CodeBadRequest,
				netsp.ErrorDetail{
					Title:    "Organization Not Found",
					Message:  "The specified organization does not exist",
					Solution: "Please check organization_id and try again",
				},
			)
		}
	}

	password, err := domain.PasswordHash(input.Password)
	if err != nil {
		return netsp.BuildError(
			netstatus.CodeBadRequest,
			netsp.ErrorDetail{
				Title:    "Invalid Password",
				Message:  "The provided password is invalid or has incorrect format",
				Solution: "Please check the password and try again",
			},
		)
	}

	account := entity.Account{
		Email:        input.Email,
		Username:     input.Username,
		PasswordHash: password,
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	if err := u.repo.Create(ctx, &account); err != nil {
		return netsp.BuildError(
			netstatus.CodeBadRequest,
			netsp.ErrorDetail{
				Title:    "Account Creation Error",
				Message:  "The account could not be created, possibly due to a duplicate email or username",
				Solution: "Please check your data and try again",
			},
		)
	}

	if input.OrganizationID != nil && *input.OrganizationID != 0 {
		if err := u.orgRepo.AddAccount(ctx, *input.OrganizationID, uint64(account.ID)); err != nil {
			_ = u.repo.Delete(ctx, uint64(account.ID))
			return netsp.BuildError(
				netstatus.CodeBadRequest,
				netsp.ErrorDetail{
					Title:    "Organization Link Failed",
					Message:  "Could not attach the account to the organization",
					Solution: "The account was not created; try again or contact support",
				},
			)
		}
	}

	return nil
}
