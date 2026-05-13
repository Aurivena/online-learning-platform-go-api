package gateway

import (
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/Aurivena/spond/v4/netoutput"
	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
	"github.com/gin-gonic/gin"

	"online-learning-platform-go-api/internal/course/dto"
	"online-learning-platform-go-api/internal/course/entity"
)

func (g *CourseGateway) bindCreateSlideRequest(c *gin.Context) (dto.CreateSlideRequest, *multipart.FileHeader, error) {
	var input dto.CreateSlideRequest
	var fileHdr *multipart.FileHeader
	ct := c.GetHeader("Content-Type")
	if strings.HasPrefix(ct, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			return input, nil, err
		}
		fh, err := c.FormFile("request")
		if err != nil {
			return input, nil, err
		}
		src, err := fh.Open()
		if err != nil {
			return input, nil, err
		}
		body, readErr := io.ReadAll(src)
		_ = src.Close()
		if readErr != nil {
			return input, nil, readErr
		}
		if err := json.Unmarshal(body, &input); err != nil {
			return input, nil, err
		}
		if file, err := c.FormFile("file"); err == nil {
			fileHdr = file
			if input.Payload == nil {
				input.Payload = entity.PayloadJSON{}
			}
			input.Payload["filename"] = file.Filename
		}
		return input, fileHdr, nil
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		return input, nil, err
	}
	return input, nil, nil
}

func (g *CourseGateway) bindUpdateSlideRequest(c *gin.Context) (dto.UpdateSlideRequest, *multipart.FileHeader, error) {
	var input dto.UpdateSlideRequest
	var fileHdr *multipart.FileHeader
	ct := c.GetHeader("Content-Type")
	if strings.HasPrefix(ct, "multipart/form-data") {
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			return input, nil, err
		}
		fh, err := c.FormFile("request")
		if err != nil {
			return input, nil, err
		}
		src, err := fh.Open()
		if err != nil {
			return input, nil, err
		}
		body, readErr := io.ReadAll(src)
		_ = src.Close()
		if readErr != nil {
			return input, nil, readErr
		}
		if err := json.Unmarshal(body, &input); err != nil {
			return input, nil, err
		}
		if file, err := c.FormFile("file"); err == nil {
			fileHdr = file
			if input.Payload == nil {
				input.Payload = entity.PayloadJSON{}
			}
			input.Payload["filename"] = file.Filename
		}
		return input, fileHdr, nil
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		return input, nil, err
	}
	return input, nil, nil
}

func (g *CourseGateway) CreateModule(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var input dto.CreateModuleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	module, errResp := g.moduleUC.CreateModule(c, input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	if errResp := g.courseUC.AttachModuleToCourse(c, courseID, module.ID); errResp != nil {
		_ = g.moduleUC.DeleteModule(c, module.ID)
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.ModuleResponse]{
		Code: netstatus.CodeSuccess,
		Data: dto.ModuleResponse{
			ID:        module.ID,
			Title:     module.Title,
			CreatedAt: module.CreatedAt,
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
		}
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.ModuleResponse]{
		Code: netstatus.CodeSuccess,
		Data: dto.ModuleResponse{
			ID:        module.ID,
			Title:     module.Title,
			CreatedAt: module.CreatedAt,
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

	mod, errResp := g.moduleUC.GetModule(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	for i := range mod.Slides {
		sid := mod.Slides[i].ID
		if errResp := g.slideUC.DeleteSlide(c, sid); errResp != nil {
			netoutput.WriteHTTP(c.Writer, *errResp)
			return
		}
	}

	if errResp := g.moduleUC.DeleteModule(c, id); errResp != nil {
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

	slideID, err := strconv.ParseUint(c.Param("slideId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid slide ID"})
		return
	}

	var input dto.AddSlideTModuleRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errResp := g.moduleUC.AddSlideToModule(c, moduleID, slideID, input)
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
	moduleID, err := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	input, fileHdr, err := g.bindCreateSlideRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.SlideType == entity.SlideTypeFile && fileHdr != nil {
		if errResp := g.uploadSlideFile(c.Request.Context(), moduleID, &input.Payload, fileHdr); errResp != nil {
			netoutput.WriteHTTP(c.Writer, *errResp)
			return
		}
	}

	slide, errResp := g.slideUC.CreateSlide(c, input)
	if errResp != nil {
		g.removeSlideObject(c.Request.Context(), input.Payload)
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	if errResp := g.moduleUC.AttachSlideToModule(c, moduleID, slide.ID); errResp != nil {
		_ = g.slideUC.DeleteSlide(c, slide.ID)
		g.removeSlideObject(c.Request.Context(), slide.Payload)
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
		},
	})
}

func (g *CourseGateway) GetSlide(c *gin.Context) {
	moduleID, errMod := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	id, err := strconv.ParseUint(c.Param("slideId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid slide ID"})
		return
	}

	if errMod == nil {
		mod, errResp := g.moduleUC.GetModule(c, moduleID)
		if errResp != nil {
			netoutput.WriteHTTP(c.Writer, *errResp)
			return
		}
		found := false
		for i := range mod.Slides {
			if mod.Slides[i].ID == id {
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "Slide not found in this module"})
			return
		}
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
		},
	})
}

func (g *CourseGateway) CheckSlideOption(c *gin.Context) {
	moduleID, errMod := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	slideID, errSlide := strconv.ParseUint(c.Param("slideId"), 10, 64)
	optionID, errOpt := strconv.ParseUint(c.Param("optionId"), 10, 64)
	if errSlide != nil || errOpt != nil {
		slog.Warn("CheckSlideOption: invalid slide or option id", "slide_err", errSlide, "option_err", errOpt)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid slide or option ID"})
		return
	}
	if errMod != nil {
		slog.Warn("CheckSlideOption: invalid module id", "err", errMod)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	slog.Info("CheckSlideOption", "module_id", moduleID, "slide_id", slideID, "option_id", optionID)

	ok, errResp := g.slideUC.CheckTestSlideOption(c, moduleID, slideID, optionID)
	if errResp != nil {
		slog.Warn("CheckSlideOption: usecase error", "module_id", moduleID, "slide_id", slideID, "option_id", optionID, "net_code", errResp.Code)
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	slog.Info("CheckSlideOption: success", "module_id", moduleID, "slide_id", slideID, "option_id", optionID, "is_correct", ok)
	netoutput.WriteHTTP(c.Writer, netsp.Response[bool]{
		Code: netstatus.CodeSuccess,
		Data: ok,
	})
}

func (g *CourseGateway) UpdateSlide(c *gin.Context) {
	moduleID, errMod := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	id, err := strconv.ParseUint(c.Param("slideId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid slide ID"})
		return
	}

	if errMod == nil {
		mod, errResp := g.moduleUC.GetModule(c, moduleID)
		if errResp != nil {
			netoutput.WriteHTTP(c.Writer, *errResp)
			return
		}
		found := false
		for i := range mod.Slides {
			if mod.Slides[i].ID == id {
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "Slide not found in this module"})
			return
		}
	}

	existing, errResp := g.slideUC.GetSlide(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	input, fileHdr, err := g.bindUpdateSlideRequest(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	st := input.SlideType
	if st == "" {
		st = existing.SlideType
	}

	oldKey := objectKeyFromPayload(existing.Payload)
	merged := clonePayload(existing.Payload)
	overlayPayload(merged, input.Payload)

	if g.files != nil && fileHdr != nil && st == entity.SlideTypeFile {
		if errResp := g.uploadSlideFile(c.Request.Context(), moduleID, &merged, fileHdr); errResp != nil {
			netoutput.WriteHTTP(c.Writer, *errResp)
			return
		}
	}
	input.Payload = merged

	newKey := objectKeyFromPayload(merged)
	errResp = g.slideUC.UpdateSlide(c, id, input)
	if errResp != nil {
		if g.files != nil && newKey != "" && newKey != oldKey {
			_ = g.files.Remove(c.Request.Context(), newKey)
		}
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	if g.files != nil && oldKey != "" && oldKey != objectKeyFromPayload(merged) {
		_ = g.files.Remove(c.Request.Context(), oldKey)
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *CourseGateway) DeleteSlide(c *gin.Context) {
	moduleID, errMod := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	id, err := strconv.ParseUint(c.Param("slideId"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid slide ID"})
		return
	}

	var slidePayload entity.PayloadJSON
	if errMod == nil {
		mod, errResp := g.moduleUC.GetModule(c, moduleID)
		if errResp != nil {
			netoutput.WriteHTTP(c.Writer, *errResp)
			return
		}
		found := false
		for i := range mod.Slides {
			if mod.Slides[i].ID == id {
				slidePayload = mod.Slides[i].Payload
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "Slide not found in this module"})
			return
		}
	} else {
		slide, errResp := g.slideUC.GetSlide(c, id)
		if errResp != nil {
			netoutput.WriteHTTP(c.Writer, *errResp)
			return
		}
		slidePayload = slide.Payload
	}

	errResp := g.slideUC.DeleteSlide(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	g.removeSlideObject(c.Request.Context(), slidePayload)

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *CourseGateway) ReorderCourseModules(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var ids []uint64
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errResp := g.courseUC.ReorderModules(c, courseID, ids); errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *CourseGateway) ReorderModuleSlides(c *gin.Context) {
	moduleID, err := strconv.ParseUint(c.Param("moduleId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	var ids []uint64
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if errResp := g.moduleUC.ReorderSlides(c, moduleID, ids); errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}
