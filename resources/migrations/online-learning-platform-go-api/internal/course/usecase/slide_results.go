package usecase

import (
	"context"
	"net/http"

	"github.com/Aurivena/spond/v4/netsp"
)

func (uc *SlideUseCase) GetTestResultForAccount(ctx context.Context, accountID, slideID uint64) (*TestResultRecord, *netsp.Response[netsp.ErrorDetail]) {
	if uc.resultRepo == nil {
		return nil, nil
	}
	r, err := uc.resultRepo.GetByAccountAndSlide(ctx, accountID, slideID)
	if err != nil {
		return nil, nil
	}
	return r, nil
}

func (uc *SlideUseCase) ListTestResults(ctx context.Context, orgID *uint64) ([]TestResultRecord, *netsp.Response[netsp.ErrorDetail]) {
	if uc.resultRepo == nil {
		return []TestResultRecord{}, nil
	}
	rows, err := uc.resultRepo.List(ctx, orgID)
	if err != nil {
		return nil, netsp.BuildError(http.StatusInternalServerError, netsp.ErrorDetail{
			Title:    "Failed to Fetch Test Results",
			Message:  "Could not load test results from database",
			Solution: "Please try again later",
		})
	}
	if rows == nil {
		rows = []TestResultRecord{}
	}
	return rows, nil
}

func (uc *SlideUseCase) ListCourseProgress(ctx context.Context, orgID *uint64) ([]CourseProgressRecord, *netsp.Response[netsp.ErrorDetail]) {
	if uc.resultRepo == nil {
		return []CourseProgressRecord{}, nil
	}
	rows, err := uc.resultRepo.ListCourseProgress(ctx, orgID)
	if err != nil {
		return nil, netsp.BuildError(http.StatusInternalServerError, netsp.ErrorDetail{
			Title:    "Failed to Fetch Course Progress",
			Message:  "Could not load course progress from database",
			Solution: "Please try again later",
		})
	}
	if rows == nil {
		rows = []CourseProgressRecord{}
	}
	return rows, nil
}
