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

	search, err := h.resourceStore.Search(resource.Query{
		Type:  resourceType,
		Query: &searchQuery,
		Skip:  skip,
		Take:  take,
	})
	if err != nil {
		return NewErrResponse(c, err)
	}

	var resources = make([]web.Resource, len(search.Items))
	for i, item := range search.Items {
		resources[i] = NewResourceResponse(item)
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
		return NewErrResponse(c, NewError(ErrUuidParseError, ErrUuidParseErrorCode, http.StatusBadRequest))
	}

	resource := model.Resource{}
	if err := h.resourceStore.GetByKey(*resourceKey, &resource); err != nil {
		return NewErrResponse(c, err)
	}

	response := web.GetResourceResponse{
		Resource: NewResourceResponse(resource),
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

	req := web.CreateResourceRequest{}

	if err := c.Bind(&req); err != nil {
		err := NewError(ErrCreateResourceCannotBind, ErrCreateResourceCannotBindCode, http.StatusBadRequest)
		return NewErrorResponse(c, *err)
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	resource := sanitize(req.Resource)

	if isSummaryTooShort(resource.Summary) {
		return newSummaryTooShortError(c)
	}
	if isSummaryTooLong(resource.Summary) {
		return newSummaryTooLongError(c)
	}
	if isDescriptionTooLong(resource.Description) {
		return newDescriptionTooLongError(c)
	}
	if isExchangeValueTooLow(resource) {
		return newExchangeValueTooLowError(c)
	}
	if isExchangeValueTooHigh(resource) {
		return newExchangeValueTooHighError(c)
	}
	if isNecessityLevelTooLow(resource) {
		return newNecessityLevelTooLowError(c)
	}
	if isNecessityLevelTooHigh(resource) {
		return newNecessityLevelTooHighError(c)
	}
	if isTimeSensitivityTooLow(resource) {
		return newTimeSensitivityTooLowError(c)
	}
	if isTimeSensitivityTooHigh(resource) {
		return newTimeSensitivityTooHighError(c)
	}

	res := model.NewResource(
		model.NewResourceKey(),
		resource.Type,
		"author",
		resource.Summary,
		resource.Description,
		model.NewTimeSensitivity(resource.TimeSensitivity),
		model.NewNecessityLevel(resource.NecessityLevel),
		model.NewExchangeValue(resource.ExchangeValue))

	err := h.resourceStore.Create(&res)
	if err != nil {
		// todo
	}

	response := web.CreateResourceResponse{
		Resource: NewResourceResponse(res),
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
		newError := NewError(ErrUpdateResourceCannotBind, ErrUpdateResourceCannotBindCode, http.StatusBadRequest)
		return NewErrorResponse(c, *newError)
	}

	// Validates the CreateResourceRequest
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// Gets the resource id
	resourceKey, err := model.ParseResourceKey(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid resource id")
	}

	// Retrieves the resource
	resToUpdate := model.Resource{}
	if err := h.resourceStore.GetByKey(*resourceKey, &resToUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, "Could not get resource")
	}

	// Validate
	desiredResource := sanitize(req.Resource)

	if isSummaryTooLong(desiredResource.Summary) {
		return newSummaryTooLongError(c)
	}
	if isSummaryTooShort(desiredResource.Summary) {
		return newSummaryTooShortError(c)
	}
	if isDescriptionTooLong(desiredResource.Description) {
		return newDescriptionTooLongError(c)
	}
	if isExchangeValueTooLow(desiredResource) {
		return newExchangeValueTooLowError(c)
	}
	if isExchangeValueTooHigh(desiredResource) {
		return newExchangeValueTooHighError(c)
	}
	if isNecessityLevelTooLow(desiredResource) {
		return newNecessityLevelTooLowError(c)
	}
	if isNecessityLevelTooHigh(desiredResource) {
		return newNecessityLevelTooHighError(c)
	}
	if isTimeSensitivityTooLow(desiredResource) {
		return newTimeSensitivityTooLowError(c)
	}
	if isTimeSensitivityTooHigh(desiredResource) {
		return newTimeSensitivityTooHighError(c)
	}

	resToUpdate.Summary = desiredResource.Summary
	resToUpdate.Description = desiredResource.Description
	resToUpdate.Type = desiredResource.Type
	resToUpdate.ExchangeValue = model.NewExchangeValue(desiredResource.ExchangeValue)
	resToUpdate.NecessityLevel = model.NewNecessityLevel(desiredResource.NecessityLevel)
	resToUpdate.TimeSensitivity = model.NewTimeSensitivity(desiredResource.TimeSensitivity)

	if err := h.resourceStore.Update(&resToUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, "Could not update resource")
	}

	return c.JSON(http.StatusOK, web.GetResourceResponse{
		Resource: NewResourceResponse(resToUpdate),
	})

}

func ParseSkip(c echo.Context) (int, error) {
	skip, err := ParseQueryParamInt(c, "skip", 0)
	if err != nil {
		return 0, NewErrResponse(c, NewError(ErrInvalidSkip, ErrInvalidSkipCode, http.StatusBadRequest))
	}
	if skip < 0 {
		skip = 0
	}
	return skip, nil
}

func ParseTake(c echo.Context, defaultTake int, maxTake int) (int, error) {
	take, err := ParseQueryParamInt(c, "take", defaultTake)
	if err != nil {
		return 0, NewErrResponse(c, NewError(ErrInvalidTake, ErrInvalidTakeCode, http.StatusBadRequest))
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
		return strconv.Atoi(paramAsStr)
	} else {
		return defaultValue, nil
	}
}

func sanitize(resource web.CreateResourcePayload) web.CreateResourcePayload {
	resource.Summary = strings.TrimSpace(resource.Summary)
	resource.Description = strings.TrimSpace(resource.Description)
	return resource
}

func newDescriptionTooLongError(c echo.Context) error {
	return NewErrResponse(c, NewError(
		ErrDescriptionTooLong,
		ErrDescriptionTooLongCode,
		http.StatusBadRequest))
}

func newSummaryTooLongError(c echo.Context) error {
	return NewErrResponse(c, NewError(
		ErrSummaryTooLong,
		ErrSummaryTooLongCode,
		http.StatusBadRequest))
}

func newSummaryTooShortError(c echo.Context) error {
	return NewErrResponse(c, NewError(
		ErrSummaryEmptyOrNull,
		ErrSummaryEmptyOrNullCode,
		http.StatusBadRequest))
}

func isDescriptionTooLong(description string) bool {
	return len(description) > 100
}

func isSummaryTooLong(summary string) bool {
	return len(summary) > 100
}

func isSummaryTooShort(summary string) bool {
	return len(summary) == 0
}

func newExchangeValueTooLowError(c echo.Context) error {
	return NewErrResponse(c, NewError(
		ErrExchangeValueTooLow,
		ErrExchangeValueTooLowCode,
		http.StatusBadRequest))
}

func newExchangeValueTooHighError(c echo.Context) error {
	return NewErrResponse(c, NewError(
		ErrExchangeValueTooHigh,
		ErrExchangeValueTooHighCode,
		http.StatusBadRequest))
}

func isExchangeValueTooLow(resource web.CreateResourcePayload) bool {
	return resource.ExchangeValue < 0
}

func isExchangeValueTooHigh(resource web.CreateResourcePayload) bool {
	return resource.ExchangeValue > 100
}

func newNecessityLevelTooLowError(c echo.Context) error {
	return NewErrResponse(c, NewError(
		ErrNecessityLevelValueTooLow,
		ErrNecessityLevelValueTooLowCode,
		http.StatusBadRequest))
}

func newNecessityLevelTooHighError(c echo.Context) error {
	return NewErrResponse(c, NewError(
		ErrNecessityLevelValueTooHigh,
		ErrNecessityLevelValueTooHighCode,
		http.StatusBadRequest))
}

func isNecessityLevelTooLow(resource web.CreateResourcePayload) bool {
	return resource.NecessityLevel < 0
}

func isNecessityLevelTooHigh(resource web.CreateResourcePayload) bool {
	return resource.NecessityLevel > 100
}

func newTimeSensitivityTooLowError(c echo.Context) error {
	return NewErrResponse(c, NewError(
		ErrTimeSensitivityValueTooLow,
		ErrTimeSensitivityValueTooLowCode,
		http.StatusBadRequest))
}

func newTimeSensitivityTooHighError(c echo.Context) error {
	return NewErrResponse(c, NewError(
		ErrTimeSensitivityValueTooHigh,
		ErrTimeSensitivityValueTooHighCode,
		http.StatusBadRequest))
}

func isTimeSensitivityTooLow(resource web.CreateResourcePayload) bool {
	return resource.TimeSensitivity < 0
}

func isTimeSensitivityTooHigh(resource web.CreateResourcePayload) bool {
	return resource.TimeSensitivity > 100
}

func NewResourceResponse(res model.Resource) web.Resource {
	return web.Resource{
		Id:              res.ID.String(),
		Type:            res.Type,
		Description:     res.Description,
		Summary:         res.Summary,
		TimeSensitivity: res.TimeSensitivity.Value,
		NecessityLevel:  res.NecessityLevel.Value,
		ExchangeValue:   res.ExchangeValue.Value,
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
