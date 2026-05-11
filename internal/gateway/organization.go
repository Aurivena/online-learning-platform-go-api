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
)

type OrganizationGateway struct {
	orgUC usecase.OrganizationUseCaseInterface
}

func NewOrganizationGateway(orgUC usecase.OrganizationUseCaseInterface) *OrganizationGateway {
	return &OrganizationGateway{orgUC: orgUC}
}

func (g *OrganizationGateway) CreateOrganization(c *gin.Context) {
	var input dto.CreateOrganizationRequest
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

	org, errResp := g.orgUC.CreateOrganization(c, userID.(uint64), input)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.OrganizationResponse]{
		Code: netstatus.CodeSuccess,
		Data: convertToOrgResponse(org),
	})
}

func (g *OrganizationGateway) ListAllOrganizations(c *gin.Context) {
	orgs, errResp := g.orgUC.ListAllOrganizations(c)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	response := make([]dto.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		response[i] = convertToOrgResponse(&org)
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[[]dto.OrganizationResponse]{
		Code: netstatus.CodeSuccess,
		Data: response,
	})
}

func (g *OrganizationGateway) ListMyOrganizations(c *gin.Context) {
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

	orgs, errResp := g.orgUC.ListMyOrganizations(c, userID.(uint64))
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	response := make([]dto.OrganizationResponse, len(orgs))
	for i, org := range orgs {
		response[i] = convertToOrgResponse(&org)
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[[]dto.OrganizationResponse]{
		Code: netstatus.CodeSuccess,
		Data: response,
	})
}

func (g *OrganizationGateway) GetOrganizationByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid organization ID"})
		return
	}

	org, errResp := g.orgUC.GetOrganization(c, id)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.OrganizationResponse]{
		Code: netstatus.CodeSuccess,
		Data: convertToOrgResponse(org),
	})
}

func (g *OrganizationGateway) GetOrganizationByTag(c *gin.Context) {
	tag := c.Param("tag")
	if tag == "" {
		c.JSON(400, gin.H{"error": "Tag is required"})
		return
	}

	org, errResp := g.orgUC.GetOrganizationByTag(c, tag)
	if errResp != nil {
		netoutput.WriteHTTP(c.Writer, *errResp)
		return
	}

	netoutput.WriteHTTP(c.Writer, netsp.Response[dto.OrganizationResponse]{
		Code: netstatus.CodeSuccess,
		Data: convertToOrgResponse(org),
	})
}

func (g *OrganizationGateway) UpdateOrganization(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid organization ID"})
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
		c.JSON(400, gin.H{"error": "Invalid organization ID"})
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
		c.JSON(400, gin.H{"error": "Invalid organization ID"})
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
		c.JSON(400, gin.H{"error": "Invalid organization ID"})
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

func convertToOrgResponse(org *entity.Organization) dto.OrganizationResponse {
	return dto.OrganizationResponse{
		ID:          org.ID,
		Title:       org.Title,
		Tag:         org.Tag,
		Description: org.Description,
		OwnerID:     org.OwnerID,
		CreatedAt:   org.CreatedAt,
		UpdatedAt:   org.UpdatedAt,
	}
}
