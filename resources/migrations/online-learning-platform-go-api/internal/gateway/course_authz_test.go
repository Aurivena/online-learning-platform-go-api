package gateway

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Aurivena/spond/v4/netsp"
	"github.com/gin-gonic/gin"

	courseDTO "online-learning-platform-go-api/internal/course/dto"
	courseEntity "online-learning-platform-go-api/internal/course/entity"
	courseUsecase "online-learning-platform-go-api/internal/course/usecase"
	orgDTO "online-learning-platform-go-api/internal/organization/dto"
	orgEntity "online-learning-platform-go-api/internal/organization/entity"
	orgUsecase "online-learning-platform-go-api/internal/organization/usecase"
	userEntity "online-learning-platform-go-api/internal/user/entity"
)

type courseUCMock struct {
	courses map[uint64]*courseEntity.Course
	resp    map[uint64]*netsp.Response[netsp.ErrorDetail]
}

func (m *courseUCMock) CreateCourse(context.Context, uint64, uint64, courseDTO.CreateCourseRequest) (*courseEntity.Course, *netsp.Response[netsp.ErrorDetail]) {
	return nil, nil
}
func (m *courseUCMock) GetCourse(_ context.Context, id uint64) (*courseEntity.Course, *netsp.Response[netsp.ErrorDetail]) {
	if r := m.resp[id]; r != nil {
		return nil, r
	}
	if c := m.courses[id]; c != nil {
		return c, nil
	}
	return nil, netsp.BuildError(http.StatusNotFound, netsp.ErrorDetail{Title: "not found"})
}
func (m *courseUCMock) ListCourses(context.Context, uint64) ([]courseEntity.Course, *netsp.Response[netsp.ErrorDetail]) {
	return nil, nil
}
func (m *courseUCMock) ListAllCourses(context.Context) ([]courseEntity.Course, *netsp.Response[netsp.ErrorDetail]) {
	return nil, nil
}
func (m *courseUCMock) UpdateCourse(context.Context, uint64, courseDTO.UpdateCourseRequest) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *courseUCMock) DeleteCourse(context.Context, uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *courseUCMock) SetCourseOrganizations(context.Context, uint64, []uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *courseUCMock) IsCourseLinkedToOrganization(_ context.Context, courseID, orgID uint64) (bool, *netsp.Response[netsp.ErrorDetail]) {
	course := m.courses[courseID]
	if course == nil {
		return false, nil
	}
	if course.OrganizationID == orgID {
		return true, nil
	}
	for _, id := range course.OrganizationIDs {
		if id == orgID {
			return true, nil
		}
	}
	return false, nil
}
func (m *courseUCMock) AddModuleToCourse(context.Context, uint64, uint64, courseDTO.AddModuleToCourseRequest) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *courseUCMock) AttachModuleToCourse(context.Context, uint64, uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *courseUCMock) ReorderModules(context.Context, uint64, []uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *courseUCMock) RemoveModuleFromCourse(context.Context, uint64, uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}

type moduleUCMock struct {
	courseIDByModule map[uint64]uint64
	modules          map[uint64]*courseEntity.Module
}

func (m *moduleUCMock) CreateModule(context.Context, courseDTO.CreateModuleRequest) (*courseEntity.Module, *netsp.Response[netsp.ErrorDetail]) {
	return nil, nil
}
func (m *moduleUCMock) GetModule(_ context.Context, id uint64) (*courseEntity.Module, *netsp.Response[netsp.ErrorDetail]) {
	if mod, ok := m.modules[id]; ok {
		return mod, nil
	}
	return nil, netsp.BuildError(http.StatusNotFound, netsp.ErrorDetail{Title: "not found"})
}
func (m *moduleUCMock) GetModuleCourseID(_ context.Context, moduleID uint64) (uint64, *netsp.Response[netsp.ErrorDetail]) {
	if cid, ok := m.courseIDByModule[moduleID]; ok {
		return cid, nil
	}
	return 0, netsp.BuildError(http.StatusNotFound, netsp.ErrorDetail{Title: "not found"})
}
func (m *moduleUCMock) UpdateModule(context.Context, uint64, courseDTO.UpdateModuleRequest) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *moduleUCMock) DeleteModule(context.Context, uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *moduleUCMock) AddSlideToModule(context.Context, uint64, uint64, courseDTO.AddSlideTModuleRequest) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *moduleUCMock) AttachSlideToModule(context.Context, uint64, uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *moduleUCMock) ReorderSlides(context.Context, uint64, []uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *moduleUCMock) RemoveSlideFromModule(context.Context, uint64, uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}

type slideUCMock struct {
	slide          *courseEntity.Slide
	getSlideResp   *netsp.Response[netsp.ErrorDetail]
	testResult     *courseUsecase.TestResultRecord
	testResultResp *netsp.Response[netsp.ErrorDetail]
	listRows       []courseUsecase.TestResultRecord
	listResp       *netsp.Response[netsp.ErrorDetail]
	progressRows   []courseUsecase.CourseProgressRecord
	progressResp   *netsp.Response[netsp.ErrorDetail]
}

func (m *slideUCMock) CreateSlide(context.Context, courseDTO.CreateSlideRequest) (*courseEntity.Slide, *netsp.Response[netsp.ErrorDetail]) {
	return nil, nil
}
func (m *slideUCMock) GetSlide(context.Context, uint64) (*courseEntity.Slide, *netsp.Response[netsp.ErrorDetail]) {
	if m.getSlideResp != nil {
		return nil, m.getSlideResp
	}
	return m.slide, nil
}
func (m *slideUCMock) UpdateSlide(context.Context, uint64, courseDTO.UpdateSlideRequest) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *slideUCMock) DeleteSlide(context.Context, uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *slideUCMock) CheckTestSlideOption(context.Context, uint64, uint64, uint64, uint64) (bool, *netsp.Response[netsp.ErrorDetail]) {
	return false, nil
}
func (m *slideUCMock) GetTestResultForAccount(context.Context, uint64, uint64) (*courseUsecase.TestResultRecord, *netsp.Response[netsp.ErrorDetail]) {
	if m.testResultResp != nil {
		return nil, m.testResultResp
	}
	return m.testResult, nil
}
func (m *slideUCMock) ListTestResults(context.Context, *uint64) ([]courseUsecase.TestResultRecord, *netsp.Response[netsp.ErrorDetail]) {
	if m.listResp != nil {
		return nil, m.listResp
	}
	if m.listRows == nil {
		return []courseUsecase.TestResultRecord{}, nil
	}
	return m.listRows, nil
}
func (m *slideUCMock) ListCourseProgress(context.Context, *uint64) ([]courseUsecase.CourseProgressRecord, *netsp.Response[netsp.ErrorDetail]) {
	if m.progressResp != nil {
		return nil, m.progressResp
	}
	if m.progressRows == nil {
		return []courseUsecase.CourseProgressRecord{}, nil
	}
	return m.progressRows, nil
}

type orgUCMock struct {
	orgByID map[uint64]*orgEntity.Organization
}

func (m *orgUCMock) CreateOrganization(context.Context, uint64, orgDTO.CreateOrganizationRequest) (*orgEntity.Organization, *netsp.Response[netsp.ErrorDetail]) {
	return nil, nil
}
func (m *orgUCMock) GetOrganization(_ context.Context, id uint64) (*orgEntity.Organization, *netsp.Response[netsp.ErrorDetail]) {
	if o := m.orgByID[id]; o != nil {
		return o, nil
	}
	return nil, netsp.BuildError(http.StatusNotFound, netsp.ErrorDetail{Title: "not found"})
}
func (m *orgUCMock) GetOrganizationByTag(context.Context, string) (*orgEntity.Organization, *netsp.Response[netsp.ErrorDetail]) {
	return nil, nil
}
func (m *orgUCMock) ListMyOrganizations(context.Context, uint64) ([]orgEntity.Organization, *netsp.Response[netsp.ErrorDetail]) {
	return nil, nil
}
func (m *orgUCMock) ListAllOrganizations(context.Context) ([]orgEntity.Organization, *netsp.Response[netsp.ErrorDetail]) {
	return nil, nil
}
func (m *orgUCMock) UpdateOrganization(context.Context, uint64, orgDTO.UpdateOrganizationRequest) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *orgUCMock) DeleteOrganization(context.Context, uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *orgUCMock) AddAccountToOrganization(context.Context, uint64, uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *orgUCMock) RemoveAccountFromOrganization(context.Context, uint64, uint64) *netsp.Response[netsp.ErrorDetail] {
	return nil
}
func (m *orgUCMock) ListOrganizationAccounts(context.Context, uint64) ([]orgDTO.OrganizationAccountResponse, *netsp.Response[netsp.ErrorDetail]) {
	return nil, nil
}

func newGatewayWithMocks() *CourseGateway {
	return &CourseGateway{
		courseUC: &courseUCMock{
			courses: map[uint64]*courseEntity.Course{},
			resp:    map[uint64]*netsp.Response[netsp.ErrorDetail]{},
		},
		moduleUC: &moduleUCMock{
			courseIDByModule: map[uint64]uint64{},
			modules:          map[uint64]*courseEntity.Module{},
		},
		slideUC: &slideUCMock{},
		orgUC: &orgUCMock{
			orgByID: map[uint64]*orgEntity.Organization{},
		},
	}
}

func TestEnsureCourseBelongsToOrganization(t *testing.T) {
	t.Parallel()

	course := &courseEntity.Course{ID: 10, OrganizationID: 5}
	if errResp := ensureCourseBelongsToOrganization(course, 5); errResp != nil {
		t.Fatalf("expected nil error, got %+v", errResp)
	}
	if errResp := ensureCourseBelongsToOrganization(course, 7); errResp == nil {
		t.Fatalf("expected error for mismatch organization")
	}
	if errResp := ensureCourseBelongsToOrganization(nil, 1); errResp == nil {
		t.Fatalf("expected error for nil course")
	}
}

func TestEnsureCourseWritePermissionByIDOwner(t *testing.T) {
	t.Parallel()

	g := newGatewayWithMocks()
	mock := g.courseUC.(*courseUCMock)
	mock.courses[100] = &courseEntity.Course{ID: 100, Owner: 42, OrganizationID: 1, CreatedAt: time.Now()}

	c := &gin.Context{}
	c.Set("userId", uint64(42))
	c.Set("role", userEntity.RoleUser)

	course, errResp := g.ensureCourseWritePermissionByID(c, 100)
	if errResp != nil {
		t.Fatalf("expected permission success, got error %+v", errResp)
	}
	if course == nil || course.ID != 100 {
		t.Fatalf("unexpected course result: %+v", course)
	}
}

func TestEnsureCourseWritePermissionByIDAdmin(t *testing.T) {
	t.Parallel()

	g := newGatewayWithMocks()
	mock := g.courseUC.(*courseUCMock)
	mock.courses[101] = &courseEntity.Course{ID: 101, Owner: 999}

	c := &gin.Context{}
	c.Set("userId", uint64(1))
	c.Set("role", userEntity.RoleAdmin)

	_, errResp := g.ensureCourseWritePermissionByID(c, 101)
	if errResp != nil {
		t.Fatalf("admin must be allowed, got error %+v", errResp)
	}
}

func TestEnsureCourseWritePermissionByIDForbidden(t *testing.T) {
	t.Parallel()

	g := newGatewayWithMocks()
	mock := g.courseUC.(*courseUCMock)
	mock.courses[102] = &courseEntity.Course{ID: 102, Owner: 777}

	c := &gin.Context{}
	c.Set("userId", uint64(42))
	c.Set("role", userEntity.RoleUser)

	_, errResp := g.ensureCourseWritePermissionByID(c, 102)
	if errResp == nil {
		t.Fatalf("expected forbidden error")
	}
	if errResp.Code != http.StatusForbidden {
		t.Fatalf("unexpected code: %d", errResp.Code)
	}
}

func TestEnsureCourseWritePermissionFromRequestByModule(t *testing.T) {
	t.Parallel()

	g := newGatewayWithMocks()
	courseMock := g.courseUC.(*courseUCMock)
	moduleMock := g.moduleUC.(*moduleUCMock)

	moduleMock.courseIDByModule[9] = 200
	courseMock.courses[200] = &courseEntity.Course{ID: 200, Owner: 42}

	c := &gin.Context{}
	c.Params = gin.Params{{Key: "moduleId", Value: "9"}}
	c.Set("userId", uint64(42))
	c.Set("role", userEntity.RoleUser)

	course, errResp := g.ensureCourseWritePermissionFromRequest(c)
	if errResp != nil {
		t.Fatalf("expected success, got error %+v", errResp)
	}
	if course == nil || course.ID != 200 {
		t.Fatalf("unexpected course: %+v", course)
	}
}

func TestEnsureCourseWritePermissionFromRequestBadCourseID(t *testing.T) {
	t.Parallel()

	g := newGatewayWithMocks()
	c := &gin.Context{}
	c.Params = gin.Params{{Key: "courseId", Value: "abc"}}
	c.Set("userId", uint64(1))
	c.Set("role", userEntity.RoleAdmin)

	_, errResp := g.ensureCourseWritePermissionFromRequest(c)
	if errResp == nil {
		t.Fatalf("expected bad request error")
	}
	if errResp.Code != http.StatusBadRequest {
		t.Fatalf("unexpected error code: %d", errResp.Code)
	}
}

var (
	_ courseUsecase.CourseUseCaseInterface    = (*courseUCMock)(nil)
	_ courseUsecase.ModuleUseCaseInterface    = (*moduleUCMock)(nil)
	_ courseUsecase.SlideUseCaseInterface     = (*slideUCMock)(nil)
	_ orgUsecase.OrganizationUseCaseInterface = (*orgUCMock)(nil)
)
