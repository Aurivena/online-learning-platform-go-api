package gateway

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Aurivena/spond/v4/netoutput"
	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
	"github.com/gin-gonic/gin"

	"online-learning-platform-go-api/internal/course/dto"
)

func parseAdminDepartmentID(c *gin.Context) (*uint64, bool) {
	var orgID *uint64
	rawOrgID := strings.TrimSpace(c.Query("departmentId"))
	if rawOrgID == "" {
		rawOrgID = strings.TrimSpace(c.Query("organizationId"))
	}
	if rawOrgID != "" {
		id, err := strconv.ParseUint(rawOrgID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID подразделения"})
			return nil, false
		}
		orgID = &id
	}
	return orgID, true
}

func (g *CourseGateway) requireAdminForResults(c *gin.Context) bool {
	_, role, ok := currentAuth(c)
	if !ok {
		errResp := authRequiredError()
		netoutput.WriteHTTP(c.Writer, *errResp)
		return false
	}
	if !isAdmin(role) {
		errResp := forbiddenError("Only admins can view employee test results")
		netoutput.WriteHTTP(c.Writer, *errResp)
		return false
	}
	return true
}

func (g *CourseGateway) ListAdminTestResults(c *gin.Context) {
	if !g.requireAdminForResults(c) {
		return
	}
	orgID, ok := parseAdminDepartmentID(c)
	if !ok {
		return
	}

	rows, errResp := g.slideUC.ListTestResults(c, orgID)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	resp := make([]dto.AdminTestResultResponse, 0, len(rows))
	for i := range rows {
		resp = append(resp, dto.AdminTestResultResponse{
			AccountID:        rows[i].AccountID,
			AccountEmail:     rows[i].AccountEmail,
			AccountUsername:  rows[i].AccountUsername,
			CourseID:         rows[i].CourseID,
			CourseTitle:      rows[i].CourseTitle,
			ModuleID:         rows[i].ModuleID,
			ModuleTitle:      rows[i].ModuleTitle,
			SlideID:          rows[i].SlideID,
			SlideTitle:       rows[i].SlideTitle,
			SelectedOptionID: rows[i].SelectedOptionID,
			IsRight:          rows[i].IsRight,
			Attempts:         rows[i].Attempts,
			FirstAttemptAt:   rows[i].FirstAttemptAt,
			LastAttemptAt:    rows[i].LastAttemptAt,
		})
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[[]dto.AdminTestResultResponse]{
		Code: netstatus.CodeSuccess,
		Data: resp,
	})
}

func (g *CourseGateway) ListAdminCourseProgress(c *gin.Context) {
	if !g.requireAdminForResults(c) {
		return
	}
	orgID, ok := parseAdminDepartmentID(c)
	if !ok {
		return
	}

	rows, errResp := g.slideUC.ListCourseProgress(c, orgID)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	resp := make([]dto.AdminCourseProgressResponse, 0, len(rows))
	for i := range rows {
		resp = append(resp, dto.AdminCourseProgressResponse{
			AccountID:       rows[i].AccountID,
			AccountEmail:    rows[i].AccountEmail,
			AccountUsername: rows[i].AccountUsername,
			CourseID:        rows[i].CourseID,
			CourseTitle:     rows[i].CourseTitle,
			TotalTests:      rows[i].TotalTests,
			AttemptedTests:  rows[i].AttemptedTests,
			PassedTests:     rows[i].PassedTests,
			Completed:       rows[i].Completed,
			LastActivityAt:  rows[i].LastActivityAt,
		})
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[[]dto.AdminCourseProgressResponse]{
		Code: netstatus.CodeSuccess,
		Data: resp,
	})
}
