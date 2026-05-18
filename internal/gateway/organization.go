package gateway

import (
	"strconv"

	"github.com/Aurivena/spond/v4/netoutput"
	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"
	"github.com/gin-gonic/gin"

	"online-learning-platform-go-api/internal/organization/dto"
	"online-learning-platform-go-api/internal/organization/entity"
	"online-learning-platform-go-api/internal/organization/usecase"
	userAdaptors "online-learning-platform-go-api/internal/user/adaptors"
	userDTO "online-learning-platform-go-api/internal/user/dto"
)

type OrganizationGateway struct {
	orgUC    usecase.OrganizationUseCaseInterface
	userRepo *userAdaptors.AccountRepository
}

func NewOrganizationGateway(orgUC usecase.OrganizationUseCaseInterface, userRepo *userAdaptors.AccountRepository) *OrganizationGateway {
	return &OrganizationGateway{
		orgUC:    orgUC,
		userRepo: userRepo,
	}
}

func (g *OrganizationGateway) CreateOrganization(c *gin.Context) {
	var input dto.CreateOrganizationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID, _, ok := currentAuth(c)
	if !ok {
		errResp := authRequiredError()
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	org, errResp := g.orgUC.CreateOrganization(c, userID, input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.OrganizationResponse]{
		Code: netstatus.CodeSuccess,
		Data: g.convertToOrgResponse(c, org),
	})
}

func (g *OrganizationGateway) ListAllOrganizations(c *gin.Context) {
	userID, role, ok := currentAuth(c)
	if !ok {
		errResp := authRequiredError()
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	accountID := c.Query("accountId")

	var orgs []entity.Organization
	var errResp *netsp.Response[netsp.ErrorDetail]

	if accountID != "" {
		id, err := strconv.ParseUint(accountID, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid accountId"})
			return
		}
		if !isAdmin(role) && id != userID {
			errResp := forbiddenError("Вы можете запрашивать только свои подразделения")
			netoutput.WriteHTTP(c.Writer, *errResp)
			return
		}
		orgs, errResp = g.orgUC.ListMyOrganizations(c, id)
	} else {
		if isAdmin(role) {
			orgs, errResp = g.orgUC.ListAllOrganizations(c)
		} else {
			orgs, errResp = g.orgUC.ListMyOrganizations(c, userID)
		}
	}

	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	response := make([]dto.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		response[i] = g.convertToOrgResponse(c, &org)
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[[]dto.OrganizationResponse]{
		Code: netstatus.CodeSuccess,
		Data: response,
	})
}

func (g *OrganizationGateway) ListMyOrganizations(c *gin.Context) {
	userID, _, ok := currentAuth(c)
	if !ok {
		errResp := authRequiredError()
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	orgs, errResp := g.orgUC.ListMyOrganizations(c, userID)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	response := make([]dto.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		response[i] = g.convertToOrgResponse(c, &org)
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[[]dto.OrganizationResponse]{
		Code: netstatus.CodeSuccess,
		Data: response,
	})
}

func (g *OrganizationGateway) GetOrganizationByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID подразделения"})
		return
	}

	org, errResp := g.orgUC.GetOrganization(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.OrganizationResponse]{
		Code: netstatus.CodeSuccess,
		Data: g.convertToOrgResponse(c, org),
	})
}

func (g *OrganizationGateway) GetOrganizationByTag(c *gin.Context) {
	tag := c.Param("tag")
	if tag == "" {
		c.JSON(400, gin.H{"error": "Тег обязателен"})
		return
	}

	org, errResp := g.orgUC.GetOrganizationByTag(c, tag)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.OrganizationResponse]{
		Code: netstatus.CodeSuccess,
		Data: g.convertToOrgResponse(c, org),
	})
}

func (g *OrganizationGateway) UpdateOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID подразделения"})
		return
	}

	if errResp := g.ensureOrganizationManagePermission(c, id); errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	var input dto.UpdateOrganizationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errResp := g.orgUC.UpdateOrganization(c, id, input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *OrganizationGateway) DeleteOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID подразделения"})
		return
	}

	if errResp := g.ensureOrganizationManagePermission(c, id); errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	errResp := g.orgUC.DeleteOrganization(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *OrganizationGateway) AddAccountToOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID подразделения"})
		return
	}

	if errResp := g.ensureOrganizationManagePermission(c, id); errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	var input dto.AddAccountToOrgRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errResp := g.orgUC.AddAccountToOrganization(c, id, input.AccountID)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *OrganizationGateway) RemoveAccountFromOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID подразделения"})
		return
	}

	if errResp := g.ensureOrganizationManagePermission(c, id); errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	var input dto.RemoveAccountFromOrgRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errResp := g.orgUC.RemoveAccountFromOrganization(c, id, input.AccountID)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[any]{
		Code: netstatus.CodeSuccess,
		Data: nil,
	})
}

func (g *OrganizationGateway) ListOrganizationAccounts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Некорректный ID подразделения"})
		return
	}

	if errResp := g.ensureOrganizationManagePermission(c, id); errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	accounts, errResp := g.orgUC.ListOrganizationAccounts(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[[]dto.OrganizationAccountResponse]{
		Code: netstatus.CodeSuccess,
		Data: accounts,
	})
}

func (g *OrganizationGateway) ListAccounts(c *gin.Context) {
	_, role, ok := currentAuth(c)
	if !ok {
		errResp := authRequiredError()
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}
	if !isAdmin(role) {
		errResp := forbiddenError("Список пользователей доступен только администратору")
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	accounts, err := g.userRepo.GetAll(c.Request.Context())
	if err != nil {
		errResp := netsp.BuildError(500, netsp.ErrorDetail{
			Title:    "Не удалось загрузить пользователей",
			Message:  "Не удалось получить список учетных записей",
			Solution: "Повторите попытку позже",
		})
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	response := make([]dto.OrganizationAccountResponse, len(accounts))
	for i, account := range accounts {
		response[i] = dto.OrganizationAccountResponse{
			ID:        uint64(account.ID),
			Email:     account.Email,
			Username:  account.Username,
			Role:      account.Role,
			CreatedAt: account.CreatedAt,
		}
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[[]dto.OrganizationAccountResponse]{
		Code: netstatus.CodeSuccess,
		Data: response,
	})
}

func (g *OrganizationGateway) convertToOrgResponse(c *gin.Context, org *entity.Organization) dto.OrganizationResponse {
	ownerAccount, _ := g.userRepo.GetByID(c.Request.Context(), org.OwnerID)
	ownerDTO := userDTO.AccountResponse{}
	if ownerAccount != nil {
		ownerDTO = userDTO.AccountResponse{
			ID:        uint(ownerAccount.ID),
			Email:     ownerAccount.Email,
			Username:  ownerAccount.Username,
			Role:      ownerAccount.Role,
			CreatedAt: ownerAccount.CreatedAt,
		}
	}

	return dto.OrganizationResponse{
		ID:          org.ID,
		Title:       org.Title,
		Tag:         org.Tag,
		Description: org.Description,
		ImageURL:    org.ImageURL,
		HeaderTitle: org.HeaderTitle,
		Owner:       ownerDTO,
		CreatedAt:   org.CreatedAt,
	}
}

func (g *OrganizationGateway) ensureOrganizationManagePermission(c *gin.Context, orgID uint64) *netsp.Response[netsp.ErrorDetail] {
	userID, role, ok := currentAuth(c)
	if !ok {
		return authRequiredError()
	}
	if isAdmin(role) {
		return nil
	}
	org, errResp := g.orgUC.GetOrganization(c, orgID)
	if errResp != nil {
		return errResp
	}
	if org.OwnerID != userID {
		return forbiddenError("Изменять подразделение может только владелец или администратор")
	}
	return nil
}
