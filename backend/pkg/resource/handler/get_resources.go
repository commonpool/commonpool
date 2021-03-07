package handler

import (
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/resource/queries"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"github.com/commonpool/backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type SearchResourcesResponse struct {
	TotalCount int        `json:"totalCount"`
	Take       int        `json:"take"`
	Skip       int        `json:"skip"`
	Resources  []Resource `json:"resources"`
}

// SearchResources godoc
// @Summary Searches resources
// @Description Search for resources
// @ID searchResources
// @Tags resources
// @Accept json
// @Produce json
// @Param query query string false "Search text"
// @Param type query string false "Resource type" Enums(0,1)
// @Param created_by query string false "Created by"
// @Param take query int false "Number of resources to take" minimum(0) maximum(100) default(10)
// @Param skip query int false "Number of resources to skip" minimum(0) default(0)
// @Success 200 {object} web.SearchResourcesResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 400 {object} utils.Error
// @Router /resources [get]
func (h *ResourceHandler) SearchResources(c echo.Context) error {

	var (
		err error
	)

	ctx, _ := handler.GetEchoContext(c, "SearchResources")

	skip, err := utils.ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := utils.ParseTake(c, 0, 100)
	if err != nil {
		return err
	}

	searchQuery := strings.TrimSpace(c.QueryParam("query"))

	var resourceType *domain.ResourceType
	resourceTypeStr := c.QueryParam("type")
	if resourceTypeStr != "" {
		resourceTypeValue, err := domain.ParseResourceType(resourceTypeStr)
		if err != nil {
			return err
		}
		resourceType = &resourceTypeValue
	}

	var callType *domain.CallType
	callTypeStr := c.QueryParam("sub_type")
	if callTypeStr != "" {
		callTypeValue, err := domain.ParseCallType(callTypeStr)
		if err != nil {
			return err
		}
		callType = &callTypeValue
	}

	var createdBy *string
	createdByStr := c.QueryParam("created_by")
	if createdByStr != "" {
		createdBy = &createdByStr
	}

	var groupKey *keys.GroupKey
	groupStr := c.QueryParam("group_id")
	if groupStr != "" {
		groupKey2, err := keys.ParseGroupKey(groupStr)
		if err != nil {
			return err
		}
		groupKey = &groupKey2
	}

	resourcesQuery := queries.NewSearchResourcesQuery(&searchQuery, resourceType, callType, skip, take, createdBy, groupKey)
	resources, err := h.searchResources.Get(ctx, resourcesQuery)
	if err != nil {
		return err
	}

	resourceKeys := keys.NewEmptyResourceKeys()
	for _, resource := range resources {
		resourceKey, err := keys.ParseResourceKey(resource.ResourceKey)
		if err != nil {
			return err
		}
		resourceKeys = resourceKeys.Append(resourceKey)
	}

	sharings, err := h.getResourcesSharings.Get(ctx, resourceKeys)
	if err != nil {
		return err
	}

	groupedSharings := map[string][]*readmodel.ResourceSharingReadModel{}
	for _, sharing := range sharings {
		if _, ok := groupedSharings[sharing.ResourceKey]; !ok {
			groupedSharings[sharing.ResourceKey] = []*readmodel.ResourceSharingReadModel{}
		}
		groupedSharings[sharing.ResourceKey] = append(groupedSharings[sharing.ResourceKey], sharing)
	}

	var resourcesResponse = make([]Resource, len(resources))
	for i, resource := range resources {
		sharingsForResource, ok := groupedSharings[resource.ResourceKey]
		if !ok {
			sharingsForResource = []*readmodel.ResourceSharingReadModel{}
		}
		resourcesResponse[i] = NewResourceResponse(resource, sharingsForResource)
	}

	// return
	return c.JSON(http.StatusOK, SearchResourcesResponse{
		Resources: resourcesResponse,
		Take:      take,
		Skip:      skip,
	})
}
