package gateway

import (
	"online-learning-platform-go-api/config"
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

func NewGateway(user usecase.AccountUseCaseInterface) *Gateway {
	return &Gateway{
		User: user,
	}
}

func SetupRouter(cfg config.Server, mw *middleware.Middleware, userGateway *Gateway, orgGateway *OrganizationGateway, courseGateway *CourseGateway) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	gHttp := gin.New() // Используем New вместо Default, чтобы контролировать всё
	gHttp.Use(gin.Logger(), gin.Recovery())

	// КРИТИЧНО: Отключаем автоматические редиректы, которые ломают CORS
	gHttp.RedirectTrailingSlash = false
	gHttp.RedirectFixedPath = false

	domain := cfg.ServerDomain
	allowOrigins := strings.Split(domain, ",")
	for i, origin := range allowOrigins {
		allowOrigins[i] = strings.TrimSpace(origin)
	}
	gHttp.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := gHttp.Group("/api")
	{
		api.GET("/files/*filepath", mw.DecodeToken, mw.AuthRequired, courseGateway.ServeUploadedObject)

		auth := api.Group("/auth")
		{
			auth.POST("/register", userGateway.Registration)
			auth.POST("/login", userGateway.Login, mw.SetToken)
			auth.POST("/logout", userGateway.Logout)
		}

		account := api.Group("/account")
		{
			account.GET("/me", mw.DecodeToken, mw.AuthRequired, userGateway.GetProfile)
		}

		registerDepartmentRoutes := func(departments *gin.RouterGroup) {
			departments.GET("", mw.DecodeToken, mw.AuthRequired, orgGateway.ListAllOrganizations)
			departments.POST("", mw.DecodeToken, mw.AuthRequired, orgGateway.CreateOrganization)
			departments.GET("/my", mw.DecodeToken, mw.AuthRequired, orgGateway.ListMyOrganizations)
			departments.GET("/:id", mw.DecodeToken, mw.AuthRequired, orgGateway.GetOrganizationByID)
			departments.GET("/tag/:tag", mw.DecodeToken, mw.AuthRequired, orgGateway.GetOrganizationByTag)
			departments.PUT("/:id", mw.DecodeToken, mw.AuthRequired, orgGateway.UpdateOrganization)
			departments.DELETE("/:id", mw.DecodeToken, mw.AuthRequired, orgGateway.DeleteOrganization)
			departments.GET("/:id/accounts", mw.DecodeToken, mw.AuthRequired, orgGateway.ListOrganizationAccounts)
			departments.POST("/:id/accounts", mw.DecodeToken, mw.AuthRequired, orgGateway.AddAccountToOrganization)
			departments.DELETE("/:id/accounts", mw.DecodeToken, mw.AuthRequired, orgGateway.RemoveAccountFromOrganization)

			courses := departments.Group("/:id/courses")
			{
				courses.GET("", mw.DecodeToken, mw.AuthRequired, courseGateway.ListCourses)
				courses.POST("", mw.DecodeToken, mw.AuthRequired, courseGateway.CreateCourse)
				courses.GET("/:courseId", mw.DecodeToken, mw.AuthRequired, courseGateway.GetCourse)
				courses.PUT("/:courseId", mw.DecodeToken, mw.AuthRequired, courseGateway.UpdateCourse)
				courses.DELETE("/:courseId", mw.DecodeToken, mw.AuthRequired, courseGateway.DeleteCourse)

				modules := courses.Group("/:courseId/modules")
				{
					modules.PUT("/reorder", mw.DecodeToken, mw.AuthRequired, courseGateway.ReorderCourseModules)
					modules.POST("", mw.DecodeToken, mw.AuthRequired, courseGateway.CreateModule)
					modules.GET("/:moduleId", mw.DecodeToken, mw.AuthRequired, courseGateway.GetModule)
					modules.PUT("/:moduleId", mw.DecodeToken, mw.AuthRequired, courseGateway.UpdateModule)
					modules.DELETE("/:moduleId", mw.DecodeToken, mw.AuthRequired, courseGateway.DeleteModule)

					slides := modules.Group("/:moduleId/slides")
					{
						slides.PUT("/reorder", mw.DecodeToken, mw.AuthRequired, courseGateway.ReorderModuleSlides)
						slides.POST("", mw.DecodeToken, mw.AuthRequired, courseGateway.CreateSlide)
						slides.GET("/:slideId/file", mw.DecodeToken, mw.AuthRequired, courseGateway.GetSlideFile)
						slides.GET("/:slideId/:optionId", mw.DecodeToken, mw.AuthRequired, courseGateway.CheckSlideOption)
						slides.GET("/:slideId", mw.DecodeToken, mw.AuthRequired, courseGateway.GetSlide)
						slides.PUT("/:slideId", mw.DecodeToken, mw.AuthRequired, courseGateway.UpdateSlide)
						slides.DELETE("/:slideId", mw.DecodeToken, mw.AuthRequired, courseGateway.DeleteSlide)
					}
				}
			}
		}

		registerDepartmentRoutes(api.Group("/departments"))
		registerDepartmentRoutes(api.Group("/organizations"))

		courseShortcuts := api.Group("/courses")
		{
			courseShortcuts.GET("", mw.DecodeToken, mw.AuthRequired, courseGateway.ListCoursePool)
			courseShortcuts.PUT("/:courseId/departments", mw.DecodeToken, mw.AuthRequired, courseGateway.UpdateCourseOrganizations)
			courseShortcuts.POST("/:courseId/modules", mw.DecodeToken, mw.AuthRequired, courseGateway.CreateModule)
			courseShortcuts.PUT("/:courseId/modules/reorder", mw.DecodeToken, mw.AuthRequired, courseGateway.ReorderCourseModules)
			courseShortcuts.GET("/:courseId/modules/:moduleId", mw.DecodeToken, mw.AuthRequired, courseGateway.GetModule)
			courseShortcuts.PUT("/:courseId/modules/:moduleId", mw.DecodeToken, mw.AuthRequired, courseGateway.UpdateModule)
			courseShortcuts.DELETE("/:courseId/modules/:moduleId", mw.DecodeToken, mw.AuthRequired, courseGateway.DeleteModule)
		}

		moduleShortcuts := api.Group("/modules")
		{
			moduleShortcuts.POST("/:moduleId/slides", mw.DecodeToken, mw.AuthRequired, courseGateway.CreateSlide)
			moduleShortcuts.PUT("/:moduleId/slides/reorder", mw.DecodeToken, mw.AuthRequired, courseGateway.ReorderModuleSlides)
			moduleShortcuts.GET("/:moduleId/slides/:slideId/file", mw.DecodeToken, mw.AuthRequired, courseGateway.GetSlideFile)
			moduleShortcuts.GET("/:moduleId/slides/:slideId/:optionId", mw.DecodeToken, mw.AuthRequired, courseGateway.CheckSlideOption)
			moduleShortcuts.GET("/:moduleId/slides/:slideId", mw.DecodeToken, mw.AuthRequired, courseGateway.GetSlide)
			moduleShortcuts.PUT("/:moduleId/slides/:slideId", mw.DecodeToken, mw.AuthRequired, courseGateway.UpdateSlide)
			moduleShortcuts.DELETE("/:moduleId/slides/:slideId", mw.DecodeToken, mw.AuthRequired, courseGateway.DeleteSlide)
		}

		admin := api.Group("/admin")
		{
			admin.GET("/accounts", mw.DecodeToken, mw.AuthRequired, orgGateway.ListAccounts)
			admin.GET("/test-results", mw.DecodeToken, mw.AuthRequired, courseGateway.ListAdminTestResults)
			admin.GET("/course-progress", mw.DecodeToken, mw.AuthRequired, courseGateway.ListAdminCourseProgress)
		}
	}

	return gHttp
}
