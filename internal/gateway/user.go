package gateway

import (
	"online-learning-platform-go-api/internal/user/dto"

	"github.com/Aurivena/spond/v4/netoutput"
	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
	"github.com/gin-gonic/gin"
)

func (g *Gateway) Registration(c *gin.Context) {
	var input dto.RegistrationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if errResp := g.User.Registration(c, input); errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *Gateway) Login(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	loginResult, errResp := g.User.Login(c, input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	c.Set("userId", loginResult.Account.ID)
	c.Set("role", loginResult.Account.Role)
	c.Next()

	accessToken, _ := c.Get("accessToken")
	refreshToken, _ := c.Get("refreshToken")

	loginResult.Response.AccessToken = accessToken.(string)
	loginResult.Response.RefreshToken = refreshToken.(string)

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.AuthResponse]{
		Code: netstatus.CodeSuccess,
		Data: *loginResult.Response,
	})
}

func (g *Gateway) GetProfile(c *gin.Context) {
	userID, ok := c.Get("userId")
	if !ok {
		errResp := netsp.BuildError(
			401,
			netsp.ErrorDetail{
				Title:    "Unauthorized",
				Message:  "User ID not found in context",
				Solution: "Please authenticate first",
			},
		)
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	profile, errResp := g.User.GetProfile(c, userID.(uint64))
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.UserProfileResponse]{
		Code: netstatus.CodeSuccess,
		Data: *profile,
	})
}
