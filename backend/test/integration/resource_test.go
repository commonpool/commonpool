package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/web"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strconv"
	"testing"
)

var resourceCounter = 1

func CreateResource(t *testing.T, ctx context.Context, userSession *auth.UserSession, opts ...*web.CreateResourceRequest) (*web.CreateResourceResponse, *http.Response) {

	resourceCounter++
	var resourceName = "resource-" + strconv.Itoa(resourceCounter)

	payload := &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary:          resourceName,
			Description:      resourceName,
			Type:             resource.Offer,
			ValueInHoursFrom: 0,
			ValueInHoursTo:   0,
			SharedWith:       []web.InputResourceSharing{},
		},
	}

	for _, option := range opts {

		if option.Resource.Summary != "" {
			payload.Resource.Summary = option.Resource.Summary
		}
		if option.Resource.Description != "" {
			payload.Resource.Description = option.Resource.Description
		}
		if option.Resource.SharedWith != nil {
			for _, sharing := range option.Resource.SharedWith {
				payload.Resource.SharedWith = append(payload.Resource.SharedWith, sharing)
			}
		}
		if option.Resource.ValueInHoursTo != 0 {
			payload.Resource.ValueInHoursTo = option.Resource.ValueInHoursTo
		}
		if option.Resource.ValueInHoursFrom != 0 {
			payload.Resource.ValueInHoursFrom = option.Resource.ValueInHoursFrom
		}
		if option.Resource.Type != resource.Offer {
			payload.Resource.Type = option.Resource.Type
		}
	}

	c, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/resources", payload)
	assert.NoError(t, a.CreateResource(c))
	response := &web.CreateResourceResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response)
}

func SearchResources(t *testing.T, ctx context.Context, userSession *auth.UserSession, take int, skip int, query string, resourceType resource.Type) (*web.SearchResourcesResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodGet, "/api/v1/resources", nil)
	c.QueryParams()["take"] = []string{strconv.Itoa(take)}
	c.QueryParams()["skip"] = []string{strconv.Itoa(skip)}
	c.QueryParams()["query"] = []string{query}
	c.QueryParams()["type"] = []string{strconv.Itoa(int(resourceType))}
	assert.NoError(t, a.SearchResources(c))
	response := &web.SearchResourcesResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response)
}

func GetResource(t *testing.T, ctx context.Context, userSession *auth.UserSession, resourceKey string) (*web.GetResourceResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/resources/%s", resourceKey), nil)
	c.SetParamNames("id")
	c.SetParamValues(resourceKey)
	assert.NoError(t, a.GetResource(c))
	response := &web.GetResourceResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response)
}

func UpdateResource(t *testing.T, ctx context.Context, userSession *auth.UserSession, resourceKey string, request *web.UpdateResourceRequest) (*web.UpdateResourceResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodPut, fmt.Sprintf("/api/v1/resources/%s", resourceKey), request)
	c.SetParamNames("id")
	c.SetParamValues(resourceKey)
	assert.NoError(t, a.UpdateResource(c))
	response := &web.UpdateResourceResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response)
}

func TestUserCanCreateResource(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	resp, httpResp := CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary:          "Summary",
			Description:      "Description",
			Type:             resource.Offer,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith:       []web.InputResourceSharing{},
		},
	})

	assert.Equal(t, http.StatusCreated, httpResp.StatusCode)

	assert.Equal(t, "Summary", resp.Resource.Summary)
	assert.Equal(t, "Description", resp.Resource.Description)
	assert.Equal(t, resource.Offer, resp.Resource.Type)
	assert.Equal(t, 1, resp.Resource.ValueInHoursFrom)
	assert.Equal(t, 3, resp.Resource.ValueInHoursTo)

}

func TestUserCanSearchResources(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary: "Blabbers",
		},
	})

	res, httpRes := SearchResources(t, ctx, user1, 10, 0, "Blabbers", resource.Offer)
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)

	assert.Equal(t, 10, res.Take)
	assert.Equal(t, 0, res.Skip)
	assert.Equal(t, 1, len(res.Resources))
	assert.Equal(t, 1, res.TotalCount)
	assert.Equal(t, "Blabbers", res.Resources[0].Summary)

}

func TestUserCanGetResource(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	res, _ := CreateResource(t, ctx, user1)

	getResource, httpRes := GetResource(t, ctx, user1, res.Resource.Id)
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)
	assert.Equal(t, res.Resource.Id, getResource.Resource.Id)

}

func TestUserCanUpdateResource(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	res, _ := CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary:          "Snippers Boop",
			Description:      "Description",
			Type:             resource.Offer,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith:       []web.InputResourceSharing{},
		},
	})

	updateResource, httpRes := UpdateResource(t, ctx, user1, res.Resource.Id, &web.UpdateResourceRequest{
		Resource: web.UpdateResourcePayload{
			Summary:          "New Summary",
			Description:      "New Description",
			Type:             resource.Offer,
			ValueInHoursFrom: 5,
			ValueInHoursTo:   10,
			SharedWith:       []web.InputResourceSharing{},
		},
	})
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)

	assert.Equal(t, res.Resource.Id, updateResource.Resource.Id)
	assert.Equal(t, "New Summary", updateResource.Resource.Summary)
	assert.Equal(t, "New Description", updateResource.Resource.Description)
	assert.Equal(t, 5, updateResource.Resource.ValueInHoursFrom)
	assert.Equal(t, 10, updateResource.Resource.ValueInHoursTo)

}

func TestUserCanShareResourceWithGroup(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	createGroupResponse, _ := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "My Group",
		Description: "Nice Group",
	})

	res, httpRes := CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			SharedWith: []web.InputResourceSharing{{GroupID: createGroupResponse.Group.ID}},
		},
	})

	assert.Equal(t, http.StatusCreated, httpRes.StatusCode)
	assert.Equal(t, 1, len(res.Resource.SharedWith))
	assert.Equal(t, createGroupResponse.Group.ID, res.Resource.SharedWith[0].GroupID)
	assert.Equal(t, "My Group", res.Resource.SharedWith[0].GroupName)

}
