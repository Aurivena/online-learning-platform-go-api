package gateway

import (
	"online-learning-platform-go-api/config"
	"online-learning-platform-go-api/internal/di"
	"online-learning-platform-go-api/internal/middleware"
	"online-learning-platform-go-api/internal/user/usecase"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Gateway struct {
	User usecase.AccountUseCaseInterface
}

func NewGateway(provider *di.Provider) *Gateway {
	return &Gateway{
		User: provider.User(),
	}
}

func SetupRouter(cfg config.Server, mw *middleware.Middleware, gateGateway *Gateway, courseGateway *CourseGateway) *gin.Engine {
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
		auth := api.Group("/auth")
		{
			auth.POST("/register", gateGateway.Registration, mw.SetToken)
			auth.POST("/login", gateGateway.Login, mw.SetToken)
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
				courses.GET("/", courseGateway.ListCourses)
				courses.POST("/", courseGateway.CreateCourse)
				courses.GET("/:courseId", courseGateway.GetCourse)
				courses.PUT("/:courseId", courseGateway.UpdateCourse)
				courses.DELETE("/:courseId", courseGateway.DeleteCourse)
				courses.POST("/:courseId/modules", courseGateway.AddModuleToCourse)
				courses.DELETE("/:courseId/modules/:moduleId", courseGateway.RemoveModuleFromCourse)

				modules := courses.Group("/:courseId/modules")
				{
					modules.POST("/", courseGateway.CreateModule)
					modules.GET("/:moduleId", courseGateway.GetModule)
					modules.PUT("/:moduleId", courseGateway.UpdateModule)
					modules.DELETE("/:moduleId", courseGateway.DeleteModule)
					modules.POST("/:moduleId/slides", courseGateway.AddSlideToModule)
					modules.DELETE("/:moduleId/slides/:slideId", courseGateway.RemoveSlideFromModule)

					slides := modules.Group("/:moduleId/slides")
					{
						slides.POST("/", courseGateway.CreateSlide)
						slides.GET("/:slideId", courseGateway.GetSlide)
						slides.PUT("/:slideId", courseGateway.UpdateSlide)
						slides.DELETE("/:slideId", courseGateway.DeleteSlide)
					}
				}
			}
		}
	}

	return gHttp
}
