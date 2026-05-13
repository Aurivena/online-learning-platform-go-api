package middleware

import (
	"net/http"
	"online-learning-platform-go-api/config"

	"github.com/Aurivena/spond/v4/netoutput"
	"github.com/Aurivena/spond/v4/netsp"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	token *config.TokenConfig
}

func NewMiddleware(token *config.TokenConfig) *Middleware {
	return &Middleware{token: token}
}

func (m *Middleware) AuthRequired(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists || userID == nil {
		errResp := netsp.BuildError(
			http.StatusUnauthorized,
			netsp.ErrorDetail{
				Title:    "Unauthorized",
				Message:  "Authentication required",
				Solution: "Please login first",
			},
		)
		netoutput.WriteHTTP(c.Writer, *errResp)
		c.Abort()
		return
	}
	c.Next()
}
