package gateway

import (
	"online-learning-platform-go-api/config"
	"online-learning-platform-go-api/internal/user/usecase"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Gateway struct {
	User *usecase.AccountUseCase
}

func NewGateway(cfg config.Server, gateGateway *Gateway) *gin.Engine {
	gHttp := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	domain := cfg.Addr + ":" + cfg.Port
	allowOrigins := strings.Split(domain, ",")

	gHttp.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"X-Session-ID", "X-Password", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := gHttp.Group("/api")
	{
		authorization := api.Group("/authorization")
		{
			authorization.POST("/login", nil)
			authorization.POST("/register", gateGateway.Registration)
		}
	}

	return gHttp
}
