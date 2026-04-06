package gateway

import (
	"net/http"
	"online-learning-platform-go-api/internal/user/dto"

	"github.com/Aurivena/spond/v3/netsp"
	"github.com/gin-gonic/gin"
)

func (g *Gateway) Registration(c *gin.Context) {
	var input dto.RegistrationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, errResp := g.User.Registration(c, input)
	if errResp != nil {
		netsp.SendResponseError(c.Writer, errResp)
		return
	}

	c.SetCookie("access_token", accessToken, 0, "", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 0, "", "", false, true)

	netsp.SendResponseSuccess[any](c.Writer, http.StatusCreated, nil)
}
