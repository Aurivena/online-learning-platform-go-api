package gateway

import (
	"strconv"

	"github.com/Aurivena/spond/v4/netoutput"
	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
	"github.com/gin-gonic/gin"

	"online-learning-platform-go-api/internal/course/dto"
)

func (g *CourseGateway) CreateModule(c *gin.Context) {
	var input dto.CreateModuleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	module, errResp := g.moduleUC.CreateModule(c, input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.ModuleResponse]{
		Code: netstatus.CodeSuccess,
		Data: dto.ModuleResponse{
			ID:        module.ID,
			Title:     module.Title,
			CreatedAt: module.CreatedAt,
			UpdatedAt: module.UpdatedAt,
		},
	})
}

func (g *CourseGateway) GetModule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid module ID"})
		return
	}

	module, errResp := g.moduleUC.GetModule(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	slides := make([]dto.SlideResponse, len(module.Slides))
	for i, s := range module.Slides {
		slides[i] = dto.SlideResponse{
			ID:          s.ID,
			Title:       s.Title,
			Description: s.Description,
			SlideType:   s.SlideType,
			Payload:     s.Payload,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		}
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.ModuleResponse]{
		Code: netstatus.CodeSuccess,
		Data: dto.ModuleResponse{
			ID:        module.ID,
			Title:     module.Title,
			CreatedAt: module.CreatedAt,
			UpdatedAt: module.UpdatedAt,
			Slides:    slides,
		},
	})
}

func (g *CourseGateway) UpdateModule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid module ID"})
		return
	}

	var input dto.UpdateModuleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errResp := g.moduleUC.UpdateModule(c, id, input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *CourseGateway) DeleteModule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid module ID"})
		return
	}

	errResp := g.moduleUC.DeleteModule(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *CourseGateway) AddSlideToModule(c *gin.Context) {
	moduleID, err := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid module ID"})
		return
	}

	var input dto.AddSlideTModuleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errResp := g.moduleUC.AddSlideToModule(c, moduleID, input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *CourseGateway) RemoveSlideFromModule(c *gin.Context) {
	moduleID, err := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid module ID"})
		return
	}

	slideID, err := strconv.ParseUint(c.Param("slideId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid slide ID"})
		return
	}

	errResp := g.moduleUC.RemoveSlideFromModule(c, moduleID, slideID)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *CourseGateway) CreateSlide(c *gin.Context) {
	var input dto.CreateSlideRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	slide, errResp := g.slideUC.CreateSlide(c, input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.SlideResponse]{
		Code: netstatus.CodeSuccess,
		Data: dto.SlideResponse{
			ID:          slide.ID,
			Title:       slide.Title,
			Description: slide.Description,
			SlideType:   slide.SlideType,
			Payload:     slide.Payload,
			CreatedAt:   slide.CreatedAt,
			UpdatedAt:   slide.UpdatedAt,
		},
	})
}

func (g *CourseGateway) GetSlide(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("slideId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid slide ID"})
		return
	}

	slide, errResp := g.slideUC.GetSlide(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.SlideResponse]{
		Code: netstatus.CodeSuccess,
		Data: dto.SlideResponse{
			ID:          slide.ID,
			Title:       slide.Title,
			Description: slide.Description,
			SlideType:   slide.SlideType,
			Payload:     slide.Payload,
			CreatedAt:   slide.CreatedAt,
			UpdatedAt:   slide.UpdatedAt,
		},
	})
}

func (g *CourseGateway) UpdateSlide(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("slideId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid slide ID"})
		return
	}

	var input dto.UpdateSlideRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errResp := g.slideUC.UpdateSlide(c, id, input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *CourseGateway) DeleteSlide(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("slideId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid slide ID"})
		return
	}

	errResp := g.slideUC.DeleteSlide(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}
