package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	resource2 "github.com/commonpool/backend/pkg/resource"
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
			Type:             resource2.Offer,
			SubType:          resource2.ObjectResource,
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
		if option.Resource.Type != resource2.Offer {
			payload.Resource.Type = option.Resource.Type
		}
		if option.Resource.SubType != "" && option.Resource.SubType != resource2.ObjectResource {
			payload.Resource.SubType = option.Resource.SubType
		}
	}

	c, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/resources", payload)
	assert.NoError(t, a.CreateResource(c))
	response := &web.CreateResourceResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response)
}

func SearchResources(t *testing.T, ctx context.Context, userSession *auth.UserSession, take int, skip int, query string, resourceType resource2.Type, sharedWithGroup *string) (*web.SearchResourcesResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodGet, "/api/v1/resources", nil)
	c.QueryParams()["take"] = []string{strconv.Itoa(take)}
	c.QueryParams()["skip"] = []string{strconv.Itoa(skip)}
	c.QueryParams()["query"] = []string{query}
	c.QueryParams()["type"] = []string{strconv.Itoa(int(resourceType))}
	if sharedWithGroup != nil {
		c.QueryParams()["group_id"] = []string{*sharedWithGroup}
	}
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
			Type:             resource2.Offer,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith:       []web.InputResourceSharing{},
		},
	})

	assert.Equal(t, http.StatusCreated, httpResp.StatusCode)

	assert.Equal(t, "Summary", resp.Resource.Summary)
	assert.Equal(t, "Description", resp.Resource.Description)
	assert.Equal(t, resource2.Offer, resp.Resource.Type)
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

	res, httpRes := SearchResources(t, ctx, user1, 10, 0, "Blabbers", resource2.Offer, nil)
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)

	assert.Equal(t, 10, res.Take)
	assert.Equal(t, 0, res.Skip)
	assert.Equal(t, 1, len(res.Resources))
	assert.Equal(t, 1, res.TotalCount)
	assert.Equal(t, "Blabbers", res.Resources[0].Summary)

}

func TestUserCanSearchResourcesWhenNoMatch(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary: "SizzlersBopBiBouWap",
		},
	})

	res, httpRes := SearchResources(t, ctx, user1, 10, 0, "ResourceNoMatchQuery", resource2.Offer, nil)
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)

	assert.Equal(t, 10, res.Take)
	assert.Equal(t, 0, res.Skip)
	assert.Equal(t, 0, len(res.Resources))
	assert.Equal(t, 0, res.TotalCount)

}

func TestUserCanSearchResourcesWithSkip(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary: "ResourceSkip1",
		},
	})
	CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary: "ResourceSkip2",
		},
	})

	res, httpRes := SearchResources(t, ctx, user1, 10, 1, "ResourceSkip", resource2.Offer, nil)
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)

	assert.Equal(t, 10, res.Take)
	assert.Equal(t, 1, res.Skip)
	assert.Equal(t, 1, len(res.Resources))
	assert.Equal(t, 2, res.TotalCount)
	assert.Equal(t, "ResourceSkip2", res.Resources[0].Summary)

}

func TestUserCanSearchResourcesSharedWithGroup(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	group1 := testGroup(t, user1)
	group2 := testGroup(t, user1)

	CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary: "SharedWithGroup",
			SharedWith: []web.InputResourceSharing{
				{
					GroupID: group1.ID,
				},
			},
		},
	})
	CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary: "SharedWithGroup",
			SharedWith: []web.InputResourceSharing{
				{
					GroupID: group2.ID,
				},
			},
		},
	})

	res, httpRes := SearchResources(t, ctx, user1, 10, 0, "SharedWithGroup", resource2.Offer, &group1.ID)
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)

	assert.Equal(t, 10, res.Take)
	assert.Equal(t, 0, res.Skip)
	assert.Equal(t, 1, len(res.Resources))
	assert.Equal(t, 1, res.TotalCount)
	assert.Equal(t, "SharedWithGroup", res.Resources[0].Summary)

}

func TestUserCanSearchResourcesSharedWithMultipleGroups(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	group1 := testGroup(t, user1)
	group2 := testGroup(t, user1)

	createResource1, createResource1Http := CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary: "A-4ccf1c0f-d791-437b-becd-8c4592d3bc1d",
			SharedWith: []web.InputResourceSharing{
				{
					GroupID: group1.ID,
				}, {
					group2.ID,
				},
			},
		},
	})
	assert.Equal(t, http.StatusCreated, createResource1Http.StatusCode)
	createResource2, createResource2Http := CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary: "B-4ccf1c0f-d791-437b-becd-8c4592d3bc1d",
			Type:    resource2.Offer,
			SubType: resource2.ObjectResource,
			SharedWith: []web.InputResourceSharing{
				{
					GroupID: group2.ID,
				},
			},
		},
	})
	assert.Equal(t, http.StatusCreated, createResource2Http.StatusCode)

	searchResource1, searchResources1Http := SearchResources(t, ctx, user1, 10, 0, "4ccf1c0f-d791-437b-becd-8c4592d3bc1d", resource2.Offer, &group1.ID)
	assert.Equal(t, http.StatusOK, searchResources1Http.StatusCode)
	assert.Equal(t, 10, searchResource1.Take)
	assert.Equal(t, 0, searchResource1.Skip)
	assert.Equal(t, 1, len(searchResource1.Resources))
	assert.Equal(t, 1, searchResource1.TotalCount)
	assert.Equal(t, createResource1.Resource.Id, searchResource1.Resources[0].Id)

	searchResource2, searchResources2Http := SearchResources(t, ctx, user1, 10, 0, "4ccf1c0f-d791-437b-becd-8c4592d3bc1d", resource2.Offer, &group2.ID)
	assert.Equal(t, http.StatusOK, searchResources2Http.StatusCode)
	assert.Equal(t, 10, searchResource2.Take)
	assert.Equal(t, 0, searchResource2.Skip)
	assert.Equal(t, 2, len(searchResource2.Resources))
	assert.Equal(t, 2, searchResource2.TotalCount)
	assert.Equal(t, createResource1.Resource.Id, searchResource1.Resources[0].Id)
	assert.Equal(t, createResource2.Resource.Id, searchResource2.Resources[1].Id)

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
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
		},
	})

	updateResource, httpRes := UpdateResource(t, ctx, user1, res.Resource.Id, &web.UpdateResourceRequest{
		Resource: web.UpdateResourcePayload{
			Summary:          "New Summary",
			Description:      "New Description",
			ValueInHoursFrom: 5,
			ValueInHoursTo:   10,
		},
	})
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)

	assert.Equal(t, res.Resource.Id, updateResource.Resource.Id)
	assert.Equal(t, "New Summary", updateResource.Resource.Summary)
	assert.Equal(t, "New Description", updateResource.Resource.Description)
	assert.Equal(t, 5, updateResource.Resource.ValueInHoursFrom)
	assert.Equal(t, 10, updateResource.Resource.ValueInHoursTo)

}

func TestUserCanUpdateResourceSharings(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	group1 := testGroup(t, user1)
	group2 := testGroup(t, user1)

	createResource, createResourceHttp := CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary:          "Snippers Boop",
			Description:      "Description",
			Type:             resource2.Offer,
			SubType:          resource2.ObjectResource,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith: []web.InputResourceSharing{
				{GroupID: group1.ID},
				{GroupID: group2.ID},
			},
		},
	})
	assert.Equal(t, http.StatusCreated, createResourceHttp.StatusCode)

	updateResource1, updateResource1Http := UpdateResource(t, ctx, user1, createResource.Resource.Id, &web.UpdateResourceRequest{
		Resource: web.UpdateResourcePayload{
			Summary:          "New Summary",
			Description:      "New Description",
			ValueInHoursFrom: 5,
			ValueInHoursTo:   10,
			SharedWith: []web.InputResourceSharing{
				{GroupID: group1.ID},
			},
		},
	})
	assert.Equal(t, http.StatusOK, updateResource1Http.StatusCode)
	assert.Equal(t, 1, len(updateResource1.Resource.SharedWith))

	updateResource2, updateResource2Http := UpdateResource(t, ctx, user1, createResource.Resource.Id, &web.UpdateResourceRequest{
		Resource: web.UpdateResourcePayload{
			Summary:          "New Summary",
			Description:      "New Description",
			ValueInHoursFrom: 5,
			ValueInHoursTo:   10,
			SharedWith: []web.InputResourceSharing{
				{GroupID: group1.ID},
				{GroupID: group2.ID},
			},
		},
	})
	assert.Equal(t, http.StatusOK, updateResource2Http.StatusCode)
	assert.Equal(t, 2, len(updateResource2.Resource.SharedWith))

}

func TestUserCanShareResourceWithGroup(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	group1 := testGroup(t, user1)

	res, httpRes := CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			SharedWith: []web.InputResourceSharing{{GroupID: group1.ID}},
		},
	})

	assert.Equal(t, http.StatusCreated, httpRes.StatusCode)
	assert.Equal(t, 1, len(res.Resource.SharedWith))
	assert.Equal(t, group1.ID, res.Resource.SharedWith[0].GroupID)
	assert.Equal(t, group1.Name, res.Resource.SharedWith[0].GroupName)

}

func TestUserCanShareResourceWithMultipleGroups(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	group1 := testGroup(t, user1)
	group2 := testGroup(t, user1)
	group3 := testGroup(t, user1)

	res, httpRes := CreateResource(t, ctx, user1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			SharedWith: []web.InputResourceSharing{
				{GroupID: group1.ID},
				{GroupID: group2.ID},
				{GroupID: group3.ID},
			},
		},
	})

	assert.Equal(t, http.StatusCreated, httpRes.StatusCode)
	assert.Equal(t, 3, len(res.Resource.SharedWith))

	for _, groupId := range []string{group1.ID, group2.ID, group3.ID} {
		found := false
		for _, sharing := range res.Resource.SharedWith {
			if sharing.GroupID == groupId {
				found = true
				break
			}
		}
		assert.Equal(t, true, found, "resource sharings should contain group id %s", groupId)
	}

}
