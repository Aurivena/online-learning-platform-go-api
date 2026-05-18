package gateway

import (
	"strconv"
	"strings"

	"github.com/Aurivena/spond/v4/netoutput"
	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
	"github.com/gin-gonic/gin"

	"online-learning-platform-go-api/internal/course/dto"
	"online-learning-platform-go-api/internal/course/entity"
	"online-learning-platform-go-api/internal/course/usecase"
	orgUsecase "online-learning-platform-go-api/internal/organization/usecase"
	"online-learning-platform-go-api/internal/storage"
)

type CourseGateway struct {
	courseUC        usecase.CourseUseCaseInterface
	moduleUC        usecase.ModuleUseCaseInterface
	slideUC         usecase.SlideUseCaseInterface
	orgUC           orgUsecase.OrganizationUseCaseInterface
	files           *storage.Bucket
	filesPublicBase string
}

func NewCourseGateway(
	courseUC usecase.CourseUseCaseInterface,
	moduleUC usecase.ModuleUseCaseInterface,
	slideUC usecase.SlideUseCaseInterface,
	orgUC orgUsecase.OrganizationUseCaseInterface,
	files *storage.Bucket,
	filesPublicBase string,
) *CourseGateway {
	return &CourseGateway{
		courseUC:        courseUC,
		moduleUC:        moduleUC,
		slideUC:         slideUC,
		orgUC:           orgUC,
		files:           files,
		filesPublicBase: strings.TrimSpace(filesPublicBase),
	}
}

func (g *CourseGateway) CreateCourse(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID подразделения"})
		return
	}

	var input dto.CreateCourseRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if errResp := g.ensureOrganizationManagePermission(c, orgID); errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	userID, _, ok := currentAuth(c)
	if !ok {
		errResp := authRequiredError()
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	course, errResp := g.courseUC.CreateCourse(c, userID, orgID, input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.CourseResponse]{
		Code: netstatus.CodeSuccess,
		Data: convertToCourseResponse(course),
	})
}

func (g *CourseGateway) GetCourse(c *gin.Context) {
	orgID, errOrg := strconv.ParseUint(c.Param("id"), 10, 64)
	id, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID курса"})
		return
	}

	course, errResp := g.courseUC.GetCourse(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}
	if errOrg == nil {
		if errResp := g.ensureCourseLinkedToOrganization(c, course.ID, orgID); errResp != nil {
			netoutput.WriteHTTP(c.Writer, *errResp)
			return
		}
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.CourseResponse]{
		Code: netstatus.CodeSuccess,
		Data: convertToCourseResponse(course),
	})
}

func (g *CourseGateway) ListCourses(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID подразделения"})
		return
	}

	courses, errResp := g.courseUC.ListCourses(c, orgID)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	response := make([]dto.CourseResponse, 0, len(courses))
	for i := range courses {
		response = append(response, convertToCourseResponse(&courses[i]))
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[[]dto.CourseResponse]{
		Code: netstatus.CodeSuccess,
		Data: response,
	})
}

func (g *CourseGateway) UpdateCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID курса"})
		return
	}
	course, errResp := g.ensureCourseWritePermissionByID(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}
	if rawOrgID := c.Param("id"); rawOrgID != "" {
		orgID, parseErr := strconv.ParseUint(rawOrgID, 10, 64)
		if parseErr != nil {
			c.JSON(400, gin.H{"error": "Некорректный ID подразделения"})
			return
		}
		if errResp = g.ensureCourseLinkedToOrganization(c, course.ID, orgID); errResp != nil {
			netoutput.WriteHTTP(c.Writer, *errResp)
			return
		}
	}

	var input dto.UpdateCourseRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errResp = g.courseUC.UpdateCourse(c, id, input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *CourseGateway) DeleteCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID курса"})
		return
	}
	course, errResp := g.ensureCourseWritePermissionByID(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}
	if rawOrgID := c.Param("id"); rawOrgID != "" {
		orgID, parseErr := strconv.ParseUint(rawOrgID, 10, 64)
		if parseErr != nil {
			c.JSON(400, gin.H{"error": "Некорректный ID подразделения"})
			return
		}
		if errResp = g.ensureCourseLinkedToOrganization(c, course.ID, orgID); errResp != nil {
			netoutput.WriteHTTP(c.Writer, *errResp)
			return
		}
	}

	errResp = g.courseUC.DeleteCourse(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *CourseGateway) AddModuleToCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID курса"})
		return
	}
	if _, errResp := g.ensureCourseWritePermissionByID(c, courseID); errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	moduleID, err := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID модуля"})
		return
	}

	var input dto.AddModuleToCourseRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errResp := g.courseUC.AddModuleToCourse(c, courseID, moduleID, input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *CourseGateway) RemoveModuleFromCourse(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID курса"})
		return
	}
	if _, errResp := g.ensureCourseWritePermissionByID(c, courseID); errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	moduleID, err := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID модуля"})
		return
	}

	errResp := g.courseUC.RemoveModuleFromCourse(c, courseID, moduleID)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *CourseGateway) ListCoursePool(c *gin.Context) {
	_, role, ok := currentAuth(c)
	if !ok {
		errResp := authRequiredError()
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}
	if !isAdmin(role) {
		errResp := forbiddenError("Пул курсов доступен только администратору")
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	courses, errResp := g.courseUC.ListAllCourses(c)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	response := make([]dto.CourseResponse, len(courses))
	for i := range courses {
		response[i] = convertToCourseResponse(&courses[i])
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[[]dto.CourseResponse]{
		Code: netstatus.CodeSuccess,
		Data: response,
	})
}

func (g *CourseGateway) UpdateCourseOrganizations(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID курса"})
		return
	}
	if _, errResp := g.ensureCourseWritePermissionByID(c, courseID); errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	var input dto.UpdateCourseOrganizationsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errResp := g.courseUC.SetCourseOrganizations(c, courseID, input.OrganizationIDs)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func convertToCourseResponse(course *entity.Course) dto.CourseResponse {
	modules := make([]dto.ModuleResponse, len(course.Modules))
	for i, m := range course.Modules {
		slides := make([]dto.SlideResponse, len(m.Slides))
		for j, s := range m.Slides {
			slides[j] = slideEntityToResponse(s)
		}
		modules[i] = dto.ModuleResponse{
			ID:        m.ID,
			Title:     m.Title,
			CreatedAt: m.CreatedAt,
			Slides:    slides,
		}
	}

	return dto.CourseResponse{
		ID:              course.ID,
		Title:           course.Title,
		Description:     course.Description,
		Owner:           course.Owner,
		OrganizationID:  course.OrganizationID,
		OrganizationIDs: course.OrganizationIDs,
		CreatedAt:       course.CreatedAt,
		Modules:         modules,
	}
}
