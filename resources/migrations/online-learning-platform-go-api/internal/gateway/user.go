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
	if c.IsAborted() {
		return
	}

	accessToken, okAccess := c.Get("accessToken")
	refreshToken, okRefresh := c.Get("refreshToken")
	accessStr, okAccessType := accessToken.(string)
	refreshStr, okRefreshType := refreshToken.(string)
	if !okAccess || !okRefresh || !okAccessType || !okRefreshType || accessStr == "" || refreshStr == "" {
		errResp := netsp.BuildError(
			500,
			netsp.ErrorDetail{
				Title:    "Token Generation Failed",
				Message:  "Login succeeded but token generation failed",
				Solution: "Please retry login later",
			},
		)
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	loginResult.Response.AccessToken = accessStr
	loginResult.Response.RefreshToken = refreshStr

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

func (g *Gateway) Logout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}
