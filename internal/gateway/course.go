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
	"online-learning-platform-go-api/internal/storage"
)

type CourseGateway struct {
	courseUC        usecase.CourseUseCaseInterface
	moduleUC        usecase.ModuleUseCaseInterface
	slideUC         usecase.SlideUseCaseInterface
	files           *storage.Bucket
	filesPublicBase string
}

func NewCourseGateway(
	courseUC usecase.CourseUseCaseInterface,
	moduleUC usecase.ModuleUseCaseInterface,
	slideUC usecase.SlideUseCaseInterface,
	files *storage.Bucket,
	filesPublicBase string,
) *CourseGateway {
	return &CourseGateway{
		courseUC:        courseUC,
		moduleUC:        moduleUC,
		slideUC:         slideUC,
		files:           files,
		filesPublicBase: strings.TrimSpace(filesPublicBase),
	}
}

func (g *CourseGateway) CreateCourse(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid organization ID"})
		return
	}

	var input dto.CreateCourseRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("userId")
	if !exists {
		errResp := netsp.BuildError(
			netstatus.CodeUnauthorized,
			netsp.ErrorDetail{
				Title:    "Unauthorized",
				Message:  "User not authenticated",
				Solution: "Please login first",
			},
		)
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	course, errResp := g.courseUC.CreateCourse(c, userID.(uint64), orgID, input)
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
	id, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid course ID"})
		return
	}

	course, errResp := g.courseUC.GetCourse(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.CourseResponse]{
		Code: netstatus.CodeSuccess,
		Data: convertToCourseResponse(course),
	})
}

func (g *CourseGateway) ListCourses(c *gin.Context) {
	orgID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid organization ID"})
		return
	}

	courses, errResp := g.courseUC.ListCourses(c, orgID)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	response := make([]dto.CourseResponse, len(courses))
	for i, course := range courses {
		response[i] = convertToCourseResponse(&course)
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[[]dto.CourseResponse]{
		Code: netstatus.CodeSuccess,
		Data: response,
	})
}

func (g *CourseGateway) UpdateCourse(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid course ID"})
		return
	}

	var input dto.UpdateCourseRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errResp := g.courseUC.UpdateCourse(c, id, input)
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
		c.JSON(400, gin.H{"error": "Invalid course ID"})
		return
	}

	errResp := g.courseUC.DeleteCourse(c, id)
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
		c.JSON(400, gin.H{"error": "Invalid course ID"})
		return
	}

	moduleID, err := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid module ID"})
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
		c.JSON(400, gin.H{"error": "Invalid course ID"})
		return
	}

	moduleID, err := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid module ID"})
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

func convertToCourseResponse(course *entity.Course) dto.CourseResponse {
	modules := make([]dto.ModuleResponse, len(course.Modules))
	for i, m := range course.Modules {
		slides := make([]dto.SlideResponse, len(m.Slides))
		for j, s := range m.Slides {
			slides[j] = dto.SlideResponse{
				ID:          s.ID,
				Title:       s.Title,
				Description: s.Description,
				SlideType:   s.SlideType,
				Payload:     s.Payload,
				CreatedAt:   s.CreatedAt,
			}
		}
		modules[i] = dto.ModuleResponse{
			ID:        m.ID,
			Title:     m.Title,
			CreatedAt: m.CreatedAt,
			Slides:    slides,
		}
	}

	return dto.CourseResponse{
		ID:             course.ID,
		Title:          course.Title,
		Description:    course.Description,
		Owner:          course.Owner,
		OrganizationID: course.OrganizationID,
		CreatedAt:      course.CreatedAt,
		Modules:        modules,
	}
}
