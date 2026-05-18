package gateway

import (
	"testing"

	"github.com/gin-gonic/gin"

	userEntity "online-learning-platform-go-api/internal/user/entity"
)

func TestCurrentAuthMissing(t *testing.T) {
	t.Parallel()

	c := &gin.Context{}
	_, _, ok := currentAuth(c)
	if ok {
		t.Fatalf("expected currentAuth to fail on empty context")
	}
}

func TestCurrentAuthSuccess(t *testing.T) {
	t.Parallel()

	c := &gin.Context{}
	c.Set("userId", uint64(42))
	c.Set("role", userEntity.RoleAdmin)

	uid, role, ok := currentAuth(c)
	if !ok {
		t.Fatalf("expected currentAuth success")
	}
	if uid != 42 {
		t.Fatalf("unexpected uid: %d", uid)
	}
	if role != userEntity.RoleAdmin {
		t.Fatalf("unexpected role: %s", role)
	}
}

func TestCurrentAuthWrongUserType(t *testing.T) {
	t.Parallel()

	c := &gin.Context{}
	c.Set("userId", uint(42))
	c.Set("role", userEntity.RoleAdmin)

	_, _, ok := currentAuth(c)
	if ok {
		t.Fatalf("expected currentAuth to fail for non-uint64 userId")
	}
}
