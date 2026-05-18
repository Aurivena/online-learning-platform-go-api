package gateway

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogoutClearsCookies(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	g := &Gateway{}
	g.Logout(c)

	if w.Code != 200 {
		t.Fatalf("unexpected status: %d", w.Code)
	}

	setCookies := w.Result().Header["Set-Cookie"]
	if len(setCookies) < 2 {
		t.Fatalf("expected at least 2 Set-Cookie headers, got %d", len(setCookies))
	}

	var hasAccess, hasRefresh bool
	for _, h := range setCookies {
		if strings.Contains(h, "access_token=") {
			hasAccess = true
		}
		if strings.Contains(h, "refresh_token=") {
			hasRefresh = true
		}
	}
	if !hasAccess || !hasRefresh {
		t.Fatalf("expected both access_token and refresh_token cookies to be cleared; headers: %v", setCookies)
	}
}
