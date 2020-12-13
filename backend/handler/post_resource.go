package handler

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
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
func (h *Handler) CreateResource(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "CreateResource")

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
	loggedInUserSession := h.authorization.GetAuthUserSession(c)

	// getting group keys that resource is shared with
	sharedWithGroupKeys, err, done := h.parseGroupKeys(c, req.Resource.SharedWith)
	if done {
		return err
	}

	err, done = h.ensureResourceIsSharedWithGroupsTheUserIsActiveMemberOf(c, loggedInUserSession.GetUserKey(), sharedWithGroupKeys)
	if done {
		return err
	}

	newResource := req.Resource
	res := resource.NewResource(
		model.NewResourceKey(uuid.NewV4()),
		newResource.Type,
		newResource.SubType,
		loggedInUserSession.Subject,
		newResource.Summary,
		newResource.Description,
		newResource.ValueInHoursFrom,
		newResource.ValueInHoursTo,
	)

	createResourceResponse := h.resourceStore.Create(resource.NewCreateResourceQuery(&res, sharedWithGroupKeys))
	if createResourceResponse.Error != nil {
		return createResourceResponse.Error
	}

	getResourceResponse, err := h.resourceStore.GetByKey(ctx, resource.NewGetResourceByKeyQuery(res.GetKey()))
	if err != nil {
		return err
	}

	groups, err := h.groupService.GetGroupsByKeys(ctx, getResourceResponse.Sharings.GetAllGroupKeys())
	if err != nil {
		return err
	}

	// send response
	return c.JSON(http.StatusCreated, web.CreateResourceResponse{
		Resource: NewResourceResponse(&res, loggedInUserSession.Username, loggedInUserSession.Subject, groups),
	})
}
