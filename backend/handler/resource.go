package handler

import (
	. "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

// SearchResources godoc
// @Summary Searches resources
// @Description Search for resources
// @ID searchResources
// @Tags resources
// @Accept json
// @Produce json
// @Param query query string true "Search text"
// @Param type query string true "Resource type" Enums(0,1)
// @Param created_by query string true "Created by"
// @Param take query int false "Number of resources to take" minimum(0) maximum(100) default(10)
// @Param skip query int false "Number of resources to skip" minimum(0) default(0)
// @Success 200 {object} web.SearchResourcesResponse
// @Failure 400 {object} utils.Error
// @Router /resources [get]
func (h *Handler) SearchResources(c echo.Context) error {

	skip, err := ParseSkip(c)
	if err != nil {
		return NewErrResponse(c, err)
	}

	take, err := ParseTake(c, 0, 100)
	if err != nil {
		return NewErrResponse(c, err)
	}

	searchQuery := strings.TrimSpace(c.QueryParam("query"))
	resourceType, err := model.ParseResourceType(c.QueryParam("type"))
	if err != nil {
		return NewErrResponse(c, err)
	}

	createdBy := c.QueryParam("created_by")

	search, err := h.resourceStore.Search(resource.Query{
		Type:      resourceType,
		Query:     &searchQuery,
		CreatedBy: createdBy,
		Skip:      skip,
		Take:      take,
	})
	if err != nil {
		return NewErrResponse(c, err)
	}

	var resources = make([]web.Resource, len(search.Items))
	for i, item := range search.Items {
		createdBy := &model.User{}
		err = h.authStore.GetByKey(model.NewUserKey(item.CreatedBy), createdBy)
		if err != nil {
			// todo
		}
		resources[i] = NewResourceResponse(item, createdBy)
	}
	response := web.SearchResourcesResponse{
		Resources:  resources,
		Take:       take,
		Skip:       skip,
		TotalCount: search.TotalCount,
	}

	return c.JSON(http.StatusOK, response)
}

// GetResource
// @Summary Gets a single resource
// @Description Gets a resource by id
// @ID getResource
// @Tags resources
// @Accept json
// @Produce json
// @Param id path string true "Resource id" format(uuid)
// @Success 200 {object} web.GetResourceResponse
// @Failure 400 {object} utils.Error
// @Router /resources/:id [get]
func (h *Handler) GetResource(c echo.Context) error {

	resourceKey, err := model.ParseResourceKey(c.Param("id"))
	if err != nil {
		return NewErrResponse(c, err)
	}

	res := model.Resource{}
	if err := h.resourceStore.GetByKey(*resourceKey, &res); err != nil {
		return NewErrResponse(c, err)
	}

	createdBy := &model.User{}
	err = h.authStore.GetByKey(model.NewUserKey(res.CreatedBy), createdBy)
	if err != nil {
		return NewErrResponse(c, err)
	}

	response := web.GetResourceResponse{
		Resource: NewResourceResponse(res, createdBy),
	}

	return c.JSON(http.StatusOK, response)
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
func (h *Handler) CreateResource(c echo.Context) error {

	var err error

	req := web.CreateResourceRequest{}

	if err = c.Bind(&req); err != nil {
		response := ErrCreateResourceBadRequest(err)
		return NewErrResponse(c, &response)
	}

	if err = c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	sanitized := sanitizeResource(req.Resource)

	err = h.validateResource(c, sanitized)
	if err != nil {
		return NewErrResponse(c, err)
	}

	subject := h.authorization.GetAuthUserSession(c).Subject

	res := model.NewResource(
		model.NewResourceKey(),
		sanitized.Type,
		subject,
		sanitized.Summary,
		sanitized.Description,
		sanitized.ValueInHoursFrom,
		sanitized.ValueInHoursTo,
	)

	err = h.resourceStore.Create(&res)
	if err != nil {
		return NewErrResponse(c, err)
	}

	createdBy := &model.User{}
	err = h.authStore.GetByKey(model.NewUserKey(res.CreatedBy), createdBy)
	if err != nil {
		return NewErrResponse(c, err)
	}

	response := web.CreateResourceResponse{
		Resource: NewResourceResponse(res, createdBy),
	}
	return c.JSON(http.StatusCreated, response)
}

// UpdateResource
// @Summary Updates a resource
// @Description Updates a resource
// @ID updateResource
// @Tags resources
// @Accept json
// @Produce json
// @Param id path string true "Resource id" format(uuid)
// @Param resource body web.UpdateResourceRequest true "Resource to create"
// @Success 200 {object} web.UpdateResourceResponse
// @Failure 400 {object} utils.Error
// @Router /resources [put]
func (h *Handler) UpdateResource(c echo.Context) error {
	req := web.CreateResourceRequest{}

	// Binds the request payload to the web.CreateResourceRequest instance
	if err := c.Bind(&req); err != nil {
		response := ErrUpdateResourceBadRequest(err)
		return NewErrResponse(c, &response)
	}

	// Validates the CreateResourceRequest
	if err := c.Validate(req); err != nil {
		response := ErrValidation(err.Error())
		return NewErrResponse(c, &response)
	}

	// Gets the resource id
	resourceKey, err := model.ParseResourceKey(c.Param("id"))
	if err != nil {
		response := ErrInvalidResourceKey(err.Error())
		return NewErrResponse(c, &response)
	}

	// Retrieves the resource
	resToUpdate := model.Resource{}
	if err := h.resourceStore.GetByKey(*resourceKey, &resToUpdate); err != nil {
		return NewErrResponse(c, err)
	}

	// Validate
	sanitized := sanitizeResource(req.Resource)

	err = h.validateResource(c, sanitized)
	if err != nil {
		return NewErrResponse(c, err)
	}

	resToUpdate.Summary = sanitized.Summary
	resToUpdate.Description = sanitized.Description
	resToUpdate.Type = sanitized.Type
	resToUpdate.ValueInHoursFrom = sanitized.ValueInHoursFrom
	resToUpdate.ValueInHoursTo = sanitized.ValueInHoursTo

	if err := h.resourceStore.Update(&resToUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, "Could not update resource")
	}

	createdBy := &model.User{}
	err = h.authStore.GetByKey(model.NewUserKey(resToUpdate.CreatedBy), createdBy)
	if err != nil {
		// todo
	}

	return c.JSON(http.StatusOK, web.GetResourceResponse{
		Resource: NewResourceResponse(resToUpdate, createdBy),
	})

}

func ParseSkip(c echo.Context) (int, error) {
	skip, err := ParseQueryParamInt(c, "skip", 0)
	if err != nil {
		response := ErrParseSkip(err.Error())
		return 0, &response
	}
	if skip < 0 {
		skip = 0
	}
	return skip, nil
}

func ParseTake(c echo.Context, defaultTake int, maxTake int) (int, error) {
	take, err := ParseQueryParamInt(c, "take", defaultTake)
	if err != nil {
		response := ErrParseTake(err.Error())
		return 0, &response
	}
	if take < 0 {
		take = 0
	}
	if take > maxTake {
		take = maxTake
	}
	return take, nil
}

func ParseQueryParamInt(c echo.Context, paramName string, defaultValue int) (int, error) {
	paramAsStr := c.QueryParam(paramName)
	if paramAsStr != "" {
		int, err := strconv.Atoi(paramAsStr)
		if err != nil {
			response := ErrCannotConvertToInt(paramAsStr, err.Error())
			return 0, &response
		}
		return int, nil
	} else {
		return defaultValue, nil
	}
}

func sanitizeResource(resource web.CreateResourcePayload) web.CreateResourcePayload {
	resource.Summary = strings.TrimSpace(resource.Summary)
	resource.Description = strings.TrimSpace(resource.Description)
	return resource
}

func (h *Handler) validateResource(c echo.Context, sanitized web.CreateResourcePayload) error {
	var err error
	err = h.checkSummaryNotTooShort(c, sanitized.Summary)
	if err != nil {
		return err
	}
	err = h.checkSummaryNotTooLong(c, sanitized.Summary)
	if err != nil {
		return err
	}
	err = h.checkDescriptionNotTooLong(c, sanitized.Description)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) checkSummaryNotTooShort(c echo.Context, summary string) error {
	var err error = nil
	if len(summary) < 5 {
		response := ErrValidation("summary is too short")
		err = &response
		return err
	}
	return nil
}

func (h *Handler) checkSummaryNotTooLong(c echo.Context, summary string) error {
	if len(summary) > 100 {
		response := ErrValidation("summary is too long")
		return &response
	}
	return nil
}

func (h *Handler) checkDescriptionNotTooLong(c echo.Context, description string) error {
	if len(description) > 100 {
		response := ErrValidation("description is too long")
		return &response
	}
	return nil
}

func NewResourceResponse(res model.Resource, usr *model.User) web.Resource {
	return web.Resource{
		Id:               res.ID.String(),
		Type:             res.Type,
		Description:      res.Description,
		Summary:          res.Summary,
		CreatedBy:        usr.Username,
		CreatedById:      usr.ID,
		CreatedAt:        res.CreatedAt,
		ValueInHoursFrom: res.ValueInHoursFrom,
		ValueInHoursTo:   res.ValueInHoursTo,
	}
}

func NewErrResponse(c echo.Context, err error) error {
	res, ok := err.(*ErrorResponse)
	if !ok {
		statusCode := http.StatusInternalServerError
		return c.JSON(statusCode, NewError("Server error", "", statusCode))
	}
	return c.JSON(res.StatusCode, res)
}

func NewErrorResponse(c echo.Context, err ErrorResponse) error {
	return c.JSON(err.StatusCode, err)
}
