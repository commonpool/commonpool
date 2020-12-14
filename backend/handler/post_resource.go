package handler

import (
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/handler"
	resource2 "github.com/commonpool/backend/pkg/resource"
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
	res := resource2.NewResource(
		model.NewResourceKey(uuid.NewV4()),
		newResource.Type,
		newResource.SubType,
		loggedInUser.Subject,
		newResource.Summary,
		newResource.Description,
		newResource.ValueInHoursFrom,
		newResource.ValueInHoursTo,
	)

	createResourceResponse := h.resourceStore.Create(resource2.NewCreateResourceQuery(&res, sharedWithGroupKeys))
	if createResourceResponse.Error != nil {
		return createResourceResponse.Error
	}

	getResourceResponse, err := h.resourceStore.GetByKey(ctx, resource2.NewGetResourceByKeyQuery(res.GetKey()))
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
