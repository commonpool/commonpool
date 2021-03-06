package handler

import (
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"
	"strings"
)

type CreateResourceRequest struct {
	Resource CreateResourcePayload `json:"resource"`
}

type CreateResourcePayload struct {
	Summary          string                 `json:"summary" validate:"required,max=100"`
	Description      string                 `json:"description" validate:"required,max=2000"`
	Type             resource.Type          `json:"type" validate:"min=0,max=1"`
	SubType          resource.SubType       `json:"subType"`
	ValueInHoursFrom int                    `json:"valueInHoursFrom" validate:"min=0"`
	ValueInHoursTo   int                    `json:"valueInHoursTo" validate:"min=0"`
	SharedWith       []InputResourceSharing `json:"sharedWith"`
}

type CreateResourceResponse struct {
	Resource Resource `json:"resource"`
}

// CreateResource
// @Summary Creates a resource
// @Description Creates a resource
// @ID createResource
// @Tags resources
// @Accept json
// @Produce json
// @Param resource body web.CreateResourceRequest true "Resource to create"
// @Success 200 {object} web.CreateResourceResponse
// @Failure 400 {object} utils.Error
// @Router /resources [post]
func (h *ResourceHandler) CreateResource(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "CreateResource")

	var err error

	// convert input body
	req := CreateResourceRequest{}
	if err = c.Bind(&req); err != nil {
		return err
	}

	// validating body
	if err = c.Validate(req); err != nil {
		return err
	}

	req.Resource.Summary = strings.TrimSpace(req.Resource.Summary)
	req.Resource.Description = strings.TrimSpace(req.Resource.Description)

	// get logged in user
	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	// getting group keys that resource is shared with
	sharedWithGroupKeys, err, done := h.parseGroupKeys(c, req.Resource.SharedWith)
	if done {
		return err
	}

	err, done = h.ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf(c, loggedInUser.GetUserKey(), sharedWithGroupKeys)
	if done {
		return err
	}

	newResource := req.Resource
	res := resource.NewResource(
		keys.NewResourceKey(uuid.NewV4()),
		newResource.Type,
		newResource.SubType,
		loggedInUser.Subject,
		newResource.Summary,
		newResource.Description,
		newResource.ValueInHoursFrom,
		newResource.ValueInHoursTo,
	)

	if err = h.resourceService.Create(ctx, resource.NewCreateResourceQuery(&res, sharedWithGroupKeys)); err != nil {
		return err
	}

	getResourceResponse, err := h.resourceService.GetByKey(ctx, resource.NewGetResourceByKeyQuery(res.GetKey()))
	if err != nil {
		return err
	}

	groups, err := h.groupService.GetGroupsByKeys(ctx, getResourceResponse.Sharings.GetAllGroupKeys())
	if err != nil {
		return err
	}

	// send response
	return c.JSON(http.StatusCreated, CreateResourceResponse{
		Resource: NewResourceResponse(&res, loggedInUser.Username, loggedInUser.Subject, groups),
	})
}
