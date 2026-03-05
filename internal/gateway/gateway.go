package gateway

import (
	"online-learning-platform-go-api/internal/pkg"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewGateway(cfg pkg.Server) *gin.Engine {
	gHttp := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	domain := cfg.Addr + ":" + cfg.Port
	allowOrigins := strings.Split(domain, ",")

	gHttp.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"X-Session-ID", "X-Password", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return gHttp
}
