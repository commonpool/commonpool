package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/resource"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"
	"strings"
)

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
	req := web.CreateResourceRequest{}
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
	loggedInUser, err := auth.GetLoggedInUser(ctx)
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
	res := resourcemodel.NewResource(
		resourcemodel.NewResourceKey(uuid.NewV4()),
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
	return c.JSON(http.StatusCreated, web.CreateResourceResponse{
		Resource: NewResourceResponse(&res, loggedInUser.Username, loggedInUser.Subject, groups),
	})
}
