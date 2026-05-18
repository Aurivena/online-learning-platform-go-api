package gateway

import (
	"net/http"

	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
	"github.com/gin-gonic/gin"

	userEntity "online-learning-platform-go-api/internal/user/entity"
)

func currentAuth(c *gin.Context) (uint64, userEntity.Role, bool) {
	rawUserID, ok := c.Get("userId")
	if !ok || rawUserID == nil {
		return 0, "", false
	}
	userID, ok := rawUserID.(uint64)
	if !ok {
		return 0, "", false
	}

	rawRole, _ := c.Get("role")
	role, _ := rawRole.(userEntity.Role)
	return userID, role, true
}

func authRequiredError() *netsp.Response[netsp.ErrorDetail] {
	return netsp.BuildError(
		netstatus.CodeUnauthorized,
		netsp.ErrorDetail{
			Title:    "Unauthorized",
			Message:  "Пользователь не авторизован",
			Solution: "Войдите в систему и повторите действие",
		},
	)
}

func forbiddenError(message string) *netsp.Response[netsp.ErrorDetail] {
	return netsp.BuildError(
		http.StatusForbidden,
		netsp.ErrorDetail{
			Title:    "Доступ запрещён",
			Message:  message,
			Solution: "Используйте учетную запись с нужными правами",
		},
	)
}

func isAdmin(role userEntity.Role) bool {
	return role == userEntity.RoleAdmin
}
