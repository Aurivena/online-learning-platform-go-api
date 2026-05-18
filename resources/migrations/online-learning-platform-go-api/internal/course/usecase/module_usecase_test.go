package usecase

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Aurivena/spond/v4/netsp"
	"gorm.io/gorm"

	"online-learning-platform-go-api/internal/course/dto"
	"online-learning-platform-go-api/internal/course/entity"
)

type moduleRepoMock struct {
	courseIDByModule map[uint64]uint64
	getCourseErr     error
}

func (m *moduleRepoMock) Create(context.Context, *entity.Module) error { return nil }
func (m *moduleRepoMock) GetByID(context.Context, uint64) (*entity.Module, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *moduleRepoMock) Update(context.Context, *entity.Module) error          { return nil }
func (m *moduleRepoMock) Delete(context.Context, uint64) error                  { return nil }
func (m *moduleRepoMock) AddSlide(context.Context, uint64, uint64, int) error   { return nil }
func (m *moduleRepoMock) RemoveSlide(context.Context, uint64, uint64) error     { return nil }
func (m *moduleRepoMock) NextSlideIndex(context.Context, uint64) (int, error)   { return 1, nil }
func (m *moduleRepoMock) ReorderSlides(context.Context, uint64, []uint64) error { return nil }
func (m *moduleRepoMock) GetCourseIDByModuleID(_ context.Context, moduleID uint64) (uint64, error) {
	if m.getCourseErr != nil {
		return 0, m.getCourseErr
	}
	if id, ok := m.courseIDByModule[moduleID]; ok {
		return id, nil
	}
	return 0, gorm.ErrRecordNotFound
}

type slideRepoMock struct {
	created *entity.Slide
}

func (m *slideRepoMock) Create(_ context.Context, slide *entity.Slide) error {
	m.created = slide
	return nil
}
func (m *slideRepoMock) GetByID(context.Context, uint64) (*entity.Slide, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *slideRepoMock) Update(context.Context, *entity.Slide) error { return nil }
func (m *slideRepoMock) Delete(context.Context, uint64) error        { return nil }

type testResultRepoMock struct{}

func (m *testResultRepoMock) Upsert(context.Context, uint64, uint64, uint64, uint64, bool) error {
	return nil
}
func (m *testResultRepoMock) GetByAccountAndSlide(context.Context, uint64, uint64) (*TestResultRecord, error) {
	return nil, gorm.ErrRecordNotFound
}
func (m *testResultRepoMock) List(context.Context, *uint64) ([]TestResultRecord, error) {
	return nil, nil
}
func (m *testResultRepoMock) ListCourseProgress(context.Context, *uint64) ([]CourseProgressRecord, error) {
	return nil, nil
}

func TestGetModuleCourseIDSuccess(t *testing.T) {
	t.Parallel()

	mr := &moduleRepoMock{
		courseIDByModule: map[uint64]uint64{10: 77},
	}
	uc := NewModuleUseCase(mr, &slideRepoMock{})

	courseID, errResp := uc.GetModuleCourseID(context.Background(), 10)
	if errResp != nil {
		t.Fatalf("unexpected error: %+v", errResp)
	}
	if courseID != 77 {
		t.Fatalf("unexpected course id: %d", courseID)
	}
}

func TestGetModuleCourseIDNotFound(t *testing.T) {
	t.Parallel()

	mr := &moduleRepoMock{
		getCourseErr: errors.New("not found"),
	}
	uc := NewModuleUseCase(mr, &slideRepoMock{})

	_, errResp := uc.GetModuleCourseID(context.Background(), 10)
	if errResp == nil {
		t.Fatalf("expected error")
	}
	if errResp.Code != http.StatusNotFound {
		t.Fatalf("unexpected code: %d", errResp.Code)
	}
}

func TestSlideCreateInitializesEmptyPayload(t *testing.T) {
	t.Parallel()

	sr := &slideRepoMock{}
	uc := NewSlideUseCase(sr, &moduleRepoMock{}, &testResultRepoMock{})

	slide, errResp := uc.CreateSlide(context.Background(), dto.CreateSlideRequest{
		Title:       "Intro",
		Description: "desc",
		SlideType:   entity.SlideTypeText,
		Payload:     nil,
	})
	if errResp != nil {
		t.Fatalf("unexpected error: %+v", errResp)
	}
	if slide == nil {
		t.Fatalf("expected slide")
	}
	if slide.Payload == nil {
		t.Fatalf("expected payload to be initialized")
	}
	if sr.created == nil || sr.created.Payload == nil {
		t.Fatalf("expected created slide payload to be initialized")
	}
}

var (
	_ ModuleRepository     = (*moduleRepoMock)(nil)
	_ SlideRepository      = (*slideRepoMock)(nil)
	_ TestResultRepository = (*testResultRepoMock)(nil)
	_                      = netsp.Response[netsp.ErrorDetail]{}
)
