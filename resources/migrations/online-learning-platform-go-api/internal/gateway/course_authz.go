package gateway

import (
	"net/http"
	"strconv"

	"github.com/Aurivena/spond/v4/netsp"
	"github.com/gin-gonic/gin"

	courseEntity "online-learning-platform-go-api/internal/course/entity"
)

func (g *CourseGateway) ensureOrganizationManagePermission(c *gin.Context, orgID uint64) *netsp.Response[netsp.ErrorDetail] {
	userID, role, ok := currentAuth(c)
	if !ok {
		return authRequiredError()
	}
	if isAdmin(role) {
		return nil
	}
	if g.orgUC == nil {
		return forbiddenError("Проверка прав подразделения не настроена")
	}
	org, errResp := g.orgUC.GetOrganization(c, orgID)
	if errResp != nil {
		return errResp
	}
	if org.OwnerID != userID {
		return forbiddenError("Изменять материалы подразделения может только владелец подразделения или администратор")
	}
	return nil
}

func (g *CourseGateway) ensureCourseWritePermissionByID(c *gin.Context, courseID uint64) (*courseEntity.Course, *netsp.Response[netsp.ErrorDetail]) {
	userID, role, ok := currentAuth(c)
	if !ok {
		return nil, authRequiredError()
	}
	course, errResp := g.courseUC.GetCourse(c, courseID)
	if errResp != nil {
		return nil, errResp
	}
	if isAdmin(role) || course.Owner == userID {
		return course, nil
	}
	return nil, forbiddenError("Изменять курс может только владелец курса или администратор")
}

func (g *CourseGateway) ensureCourseWritePermissionFromRequest(c *gin.Context) (*courseEntity.Course, *netsp.Response[netsp.ErrorDetail]) {
	if raw := c.Param("courseId"); raw != "" {
		courseID, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return nil, netsp.BuildError(http.StatusBadRequest, netsp.ErrorDetail{
				Title:    "Некорректный ID курса",
				Message:  "Не удалось прочитать courseId из адреса запроса",
				Solution: "Передайте корректный числовой ID курса",
			})
		}
		return g.ensureCourseWritePermissionByID(c, courseID)
	}

	if raw := c.Param("moduleId"); raw != "" {
		moduleID, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return nil, netsp.BuildError(http.StatusBadRequest, netsp.ErrorDetail{
				Title:    "Некорректный ID модуля",
				Message:  "Не удалось прочитать moduleId из адреса запроса",
				Solution: "Передайте корректный числовой ID модуля",
			})
		}
		courseID, errResp := g.moduleUC.GetModuleCourseID(c, moduleID)
		if errResp != nil {
			return nil, errResp
		}
		return g.ensureCourseWritePermissionByID(c, courseID)
	}

	return nil, netsp.BuildError(http.StatusBadRequest, netsp.ErrorDetail{
		Title:    "Некорректный запрос",
		Message:  "В адресе запроса нет courseId или moduleId",
		Solution: "Используйте адрес курса или модуля для этого действия",
	})
}

func ensureCourseBelongsToOrganization(course *courseEntity.Course, orgID uint64) *netsp.Response[netsp.ErrorDetail] {
	if course == nil {
		return netsp.BuildError(http.StatusNotFound, netsp.ErrorDetail{
			Title:    "Курс не найден",
			Message:  "Данные курса отсутствуют",
			Solution: "Обновите страницу и повторите действие",
		})
	}
	if course.OrganizationID != orgID {
		return netsp.BuildError(http.StatusNotFound, netsp.ErrorDetail{
			Title:    "Курс не найден",
			Message:  "Курс не относится к этому подразделению",
			Solution: "Проверьте ID подразделения и курса",
		})
	}
	return nil
}

func (g *CourseGateway) ensureCourseLinkedToOrganization(c *gin.Context, courseID, orgID uint64) *netsp.Response[netsp.ErrorDetail] {
	linked, errResp := g.courseUC.IsCourseLinkedToOrganization(c, courseID, orgID)
	if errResp != nil {
		return errResp
	}
	if linked {
		return nil
	}
	return netsp.BuildError(http.StatusNotFound, netsp.ErrorDetail{
		Title:    "Курс не найден",
		Message:  "Курс не привязан к этому подразделению",
		Solution: "Проверьте ID подразделения и курса",
	})
}
