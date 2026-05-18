package gateway

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	courseEntity "online-learning-platform-go-api/internal/course/entity"
	courseUsecase "online-learning-platform-go-api/internal/course/usecase"
	userEntity "online-learning-platform-go-api/internal/user/entity"
)

func unwrapDataField(t *testing.T, body []byte) interface{} {
	t.Helper()
	var root map[string]interface{}
	if err := json.Unmarshal(body, &root); err != nil {
		t.Fatalf("failed to unmarshal body: %v; body=%s", err, string(body))
	}
	if v, ok := root["content"]; ok {
		return v
	}
	if v, ok := root["data"]; ok {
		return v
	}
	return root
}

func TestGetSlideSetsIsRightNullWithoutPersonalResult(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "moduleId", Value: "5"},
		{Key: "slideId", Value: "10"},
	}
	c.Set("userId", uint64(7))
	c.Set("role", userEntity.RoleUser)

	g := newGatewayWithMocks()
	moduleMock := g.moduleUC.(*moduleUCMock)
	slideMock := g.slideUC.(*slideUCMock)

	moduleMock.modules[5] = &courseEntity.Module{
		ID: 5,
		Slides: []courseEntity.Slide{
			{ID: 10},
		},
	}
	slideMock.slide = &courseEntity.Slide{
		ID:        10,
		SlideType: courseEntity.SlideTypeTest,
		Payload: courseEntity.PayloadJSON{
			"question": "Q?",
			"is_right": true,
		},
	}

	g.GetSlide(c)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", w.Code, w.Body.String())
	}

	data, ok := unwrapDataField(t, w.Body.Bytes()).(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected data format: %s", w.Body.String())
	}
	payload, ok := data["payload"].(map[string]interface{})
	if !ok {
		t.Fatalf("payload missing: %v", data)
	}
	if v, exists := payload["is_right"]; !exists || v != nil {
		t.Fatalf("expected payload.is_right to be null when user has no result, got %v", payload["is_right"])
	}
}

func TestGetSlideSetsPersonalIsRight(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{
		{Key: "moduleId", Value: "5"},
		{Key: "slideId", Value: "10"},
	}
	c.Set("userId", uint64(7))
	c.Set("role", userEntity.RoleUser)

	g := newGatewayWithMocks()
	moduleMock := g.moduleUC.(*moduleUCMock)
	slideMock := g.slideUC.(*slideUCMock)

	moduleMock.modules[5] = &courseEntity.Module{
		ID: 5,
		Slides: []courseEntity.Slide{
			{ID: 10},
		},
	}
	slideMock.slide = &courseEntity.Slide{
		ID:        10,
		SlideType: courseEntity.SlideTypeTest,
		Payload: courseEntity.PayloadJSON{
			"question": "Q?",
		},
	}
	slideMock.testResult = &courseUsecase.TestResultRecord{
		AccountID: 7,
		SlideID:   10,
		IsRight:   true,
	}

	g.GetSlide(c)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", w.Code, w.Body.String())
	}

	data, ok := unwrapDataField(t, w.Body.Bytes()).(map[string]interface{})
	if !ok {
		t.Fatalf("unexpected data format: %s", w.Body.String())
	}
	payload, ok := data["payload"].(map[string]interface{})
	if !ok {
		t.Fatalf("payload missing: %v", data)
	}
	if payload["is_right"] != true {
		t.Fatalf("expected payload.is_right=true for personal result, got %v", payload["is_right"])
	}
}

func TestListAdminTestResultsForbiddenForNonAdmin(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/admin/test-results", nil)
	c.Set("userId", uint64(42))
	c.Set("role", userEntity.RoleUser)

	g := newGatewayWithMocks()
	g.ListAdminTestResults(c)

	if w.Code == http.StatusOK {
		t.Fatalf("expected forbidden response, got 200 body=%s", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Only admins can view employee test results") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}
