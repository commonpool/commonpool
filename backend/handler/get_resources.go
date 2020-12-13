package handler

import (
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/utils"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

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
func (h *Handler) SearchResources(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "SearchResources")

	skip, err := utils.ParseSkip(c)
	if err != nil {
		return err
	}

	take, err := utils.ParseTake(c, 0, 100)
	if err != nil {
		return err
	}

	searchQuery := strings.TrimSpace(c.QueryParam("query"))

	resourceType, err := resource.ParseResourceType(c.QueryParam("type"))
	if err != nil {
		return err
	}

	resourceSubType, err := resource.ParseResourceSubType(c.QueryParam("sub_type"))
	if err != nil {
		return err
	}

	createdBy := c.QueryParam("created_by")

	var groupKey *model.GroupKey
	groupStr := c.QueryParam("group_id")
	if groupStr != "" {
		groupKey2, err := group.ParseGroupKey(groupStr)
		if err != nil {
			return err
		}
		groupKey = &groupKey2
	}

	resourcesQuery := resource.NewSearchResourcesQuery(&searchQuery, resourceType, resourceSubType, skip, take, createdBy, groupKey)
	resources := h.resourceStore.Search(ctx, resourcesQuery)
	if resources.Error != nil {
		return err
	}

	getGroupsResponse, err := h.groupService.GetGroupsByKeys(ctx, resources.Sharings.GetAllGroupKeys())
	if err != nil {
		return err
	}

	groupMap := map[model.GroupKey]*group.Group{}
	for _, g := range getGroupsResponse.Items {
		groupMap[g.GetKey()] = g
	}

	var createdByKeys []model.UserKey
	for _, item := range resources.Resources.Items {
		createdByKeys = append(createdByKeys, item.GetOwnerKey())
	}

	createdByUsers, err := h.authStore.GetByKeys(ctx, createdByKeys)
	if err != nil {
		return err
	}

	resourceItems := resources.Resources.Items
	var resourcesResponse = make([]web.Resource, len(resourceItems))
	for i, item := range resourceItems {

		createdBy, err := createdByUsers.GetUser(model.NewUserKey(item.CreatedBy))
		if err != nil {
			return err
		}

		var groups []*group.Group
		sharings := resources.Sharings.GetSharingsForResource(item.GetKey())
		for _, groupKey := range sharings.GetAllGroupKeys().Items {
			groups = append(groups, groupMap[groupKey])
		}

		resourcesResponse[i] = NewResourceResponse(item, createdBy.Username, createdBy.ID, group.NewGroups(groups))
	}

	// return
	return c.JSON(http.StatusOK, web.SearchResourcesResponse{
		Resources:  resourcesResponse,
		Take:       take,
		Skip:       skip,
		TotalCount: resources.TotalCount,
	})
}
