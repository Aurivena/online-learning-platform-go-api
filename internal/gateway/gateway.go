package gateway

import (
	"online-learning-platform-go-api/config"
	"online-learning-platform-go-api/internal/middleware"
	"online-learning-platform-go-api/internal/user/adaptors"
	"online-learning-platform-go-api/internal/user/usecase"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Gateway struct {
	User *usecase.AccountUseCase
}

func NewGateway(gorm *gorm.DB) *Gateway {
	return &Gateway{
		User: usecase.NewAccountUseCase(adaptors.NewRepository(gorm)),
	}
}

func SetupRouter(cfg config.Server, mw *middleware.Middleware, gateGateway *Gateway) *gin.Engine {
	gHttp := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	domain := cfg.Addr + ":" + cfg.Port
	allowOrigins := strings.Split(domain, ",")

	gHttp.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := gHttp.Group("/api")
	{
		auth := api.Group("/auth", mw.SetToken)
		{
			auth.POST("/login", nil)
			auth.POST("/register", gateGateway.Registration)
			auth.POST("/logout", nil)
		}

		account := api.Group("/account")
		{
			account.GET("/me", nil)
		}

		organizations := api.Group("/organizations")
		{
			organizations.GET("/", nil)
			organizations.POST("/", nil)
			organizations.GET("/:tag", nil)
			organizations.PUT("/:tag", nil)
			organizations.DELETE("/:tag", nil)

			courses := organizations.Group("/:id/courses")
			{
				courses.GET("/", nil)
				courses.POST("/", nil)
				courses.GET("/:courseId", nil)
				courses.PUT("/:courseId", nil)
				courses.DELETE("/:courseId", nil)

				modules := courses.Group("/:courseId/modules")
				{
					modules.POST("/", nil)
					modules.GET("/:moduleId", nil)
					modules.PUT("/:moduleId", nil)
					modules.DELETE("/:moduleId", nil)
					modules.PUT("/recorder", nil)

					slides := modules.Group("/:moduleId/slides")
					{
						slides.POST("/", nil)
						slides.GET("/:slideId", nil)
						slides.PUT("/:slideId", nil)
						slides.DELETE("/:slideId", nil)
						slides.PUT("/recorder", nil)
					}
				}
			}
		}
	}

	return gHttp
}
