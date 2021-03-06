package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/resource"
	"github.com/commonpool/backend/pkg/resource/handler"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

var resourceCounter = 1

func (s *IntegrationTestSuite) CreateResource(t *testing.T, ctx context.Context, userSession *models.UserSession, opts ...*handler.CreateResourceRequest) (*handler.CreateResourceResponse, *http.Response) {

	resourceCounter++
	var resourceName = "resource-" + strconv.Itoa(resourceCounter)

	payload := &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			Summary:          resourceName,
			Description:      resourceName,
			Type:             resource.Offer,
			SubType:          resource.ObjectResource,
			ValueInHoursFrom: 0,
			ValueInHoursTo:   0,
			SharedWith:       []handler.InputResourceSharing{},
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
		if option.Resource.SubType != "" && option.Resource.SubType != resource.ObjectResource {
			payload.Resource.SubType = option.Resource.SubType
		}
	}

	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/resources", payload)
	s.server.Router.ServeHTTP(recorder, httpReq)
	response := &handler.CreateResourceResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) SearchResources(t *testing.T, ctx context.Context, userSession *models.UserSession, take int, skip int, query string, resourceType resource.Type, sharedWithGroup *string) (*handler.SearchResourcesResponse, *http.Response) {
	target, _ := url.Parse("/api/v1/resources")
	target.Query().Set("take", strconv.Itoa(take))
	target.Query().Set("skip", strconv.Itoa(skip))
	target.Query().Set("query", query)
	target.Query().Set("type", strconv.Itoa(int(resourceType)))
	if sharedWithGroup != nil {
		target.Query().Set("group_id", *sharedWithGroup)
	}
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodGet, target.String(), nil)
	s.server.Router.ServeHTTP(recorder, httpReq)
	response := &handler.SearchResourcesResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) GetResource(t *testing.T, ctx context.Context, userSession *models.UserSession, resourceKey string) (*handler.GetResourceResponse, *http.Response) {
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/resources/%s", resourceKey), nil)
	s.server.Router.ServeHTTP(recorder, httpReq)
	response := &handler.GetResourceResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) UpdateResource(t *testing.T, ctx context.Context, userSession *models.UserSession, resourceKey string, request *handler.UpdateResourceRequest) (*handler.UpdateResourceResponse, *http.Response) {
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPut, fmt.Sprintf("/api/v1/resources/%s", resourceKey), request)
	s.server.Router.ServeHTTP(recorder, httpReq)
	response := &handler.UpdateResourceResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) TestUserCanCreateResource() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	resp, httpResp := s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			Summary:          "Summary",
			Description:      "Description",
			Type:             resource.Offer,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith:       []handler.InputResourceSharing{},
		},
	})
	if !AssertStatusCreated(s.T(), httpResp) {
		return
	}

	assert.Equal(s.T(), "Summary", resp.Resource.Summary)
	assert.Equal(s.T(), "Description", resp.Resource.Description)
	assert.Equal(s.T(), resource.Offer, resp.Resource.Type)
	assert.Equal(s.T(), 1, resp.Resource.ValueInHoursFrom)
	assert.Equal(s.T(), 3, resp.Resource.ValueInHoursTo)

}

func (s *IntegrationTestSuite) TestUserCanSearchResources() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	_, httpResponse := s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			Summary: "Blabbers",
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	res, httpRes := s.SearchResources(s.T(), ctx, user1, 10, 0, "Blabbers", resource.Offer, nil)
	if !AssertOK(s.T(), httpRes) {
		return
	}

	assert.Equal(s.T(), 10, res.Take)
	assert.Equal(s.T(), 0, res.Skip)
	assert.Equal(s.T(), 1, res.TotalCount)

	if !assert.Len(s.T(), res.Resources, 1) {
		return
	}
	assert.Equal(s.T(), "Blabbers", res.Resources[0].Summary)

}

func (s *IntegrationTestSuite) TestUserCanSearchResourcesWhenNoMatch() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			Summary: "SizzlersBopBiBouWap",
		},
	})

	res, httpRes := s.SearchResources(s.T(), ctx, user1, 10, 0, "ResourceNoMatchQuery", resource.Offer, nil)
	assert.Equal(s.T(), http.StatusOK, httpRes.StatusCode)

	assert.Equal(s.T(), 10, res.Take)
	assert.Equal(s.T(), 0, res.Skip)
	assert.Equal(s.T(), 0, len(res.Resources))
	assert.Equal(s.T(), 0, res.TotalCount)

}

func (s *IntegrationTestSuite) TestUserCanSearchResourcesWithSkip() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			Summary: "ResourceSkip1",
		},
	})
	s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			Summary: "ResourceSkip2",
		},
	})

	res, httpRes := s.SearchResources(s.T(), ctx, user1, 10, 1, "ResourceSkip", resource.Offer, nil)
	assert.Equal(s.T(), http.StatusOK, httpRes.StatusCode)

	assert.Equal(s.T(), 10, res.Take)
	assert.Equal(s.T(), 1, res.Skip)
	assert.Equal(s.T(), 1, len(res.Resources))
	assert.Equal(s.T(), 2, res.TotalCount)
	assert.Equal(s.T(), "ResourceSkip2", res.Resources[0].Summary)

}

func (s *IntegrationTestSuite) TestUserCanSearchResourcesSharedWithGroup() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	group1, err := s.testGroup(s.T(), user1)
	if !assert.NoError(s.T(), err) {
		return
	}
	group2, err := s.testGroup(s.T(), user1)
	if !assert.NoError(s.T(), err) {
		return
	}

	s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			Summary: "SharedWithGroup",
			SharedWith: []handler.InputResourceSharing{
				{
					GroupID: group1.ID,
				},
			},
		},
	})
	s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			Summary: "SharedWithGroup",
			SharedWith: []handler.InputResourceSharing{
				{
					GroupID: group2.ID,
				},
			},
		},
	})

	res, httpRes := s.SearchResources(s.T(), ctx, user1, 10, 0, "SharedWithGroup", resource.Offer, &group1.ID)
	assert.Equal(s.T(), http.StatusOK, httpRes.StatusCode)

	assert.Equal(s.T(), 10, res.Take)
	assert.Equal(s.T(), 0, res.Skip)
	assert.Equal(s.T(), 1, res.TotalCount)
	if !assert.Len(s.T(), res.Resources, 1) {
		return
	}
	assert.Equal(s.T(), "SharedWithGroup", res.Resources[0].Summary)

}

func (s *IntegrationTestSuite) TestUserCanSearchResourcesSharedWithMultipleGroups() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	group1, err := s.testGroup(s.T(), user1)
	if !assert.NoError(s.T(), err) {
		return
	}
	group2, err := s.testGroup(s.T(), user1)
	if !assert.NoError(s.T(), err) {
		return
	}

	createResource1, createResource1Http := s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			Summary: "A-4ccf1c0f-d791-437b-becd-8c4592d3bc1d",
			SharedWith: []handler.InputResourceSharing{
				{
					GroupID: group1.ID,
				}, {
					group2.ID,
				},
			},
		},
	})
	assert.Equal(s.T(), http.StatusCreated, createResource1Http.StatusCode)
	createResource2, createResource2Http := s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			Summary: "B-4ccf1c0f-d791-437b-becd-8c4592d3bc1d",
			Type:    resource.Offer,
			SubType: resource.ObjectResource,
			SharedWith: []handler.InputResourceSharing{
				{
					GroupID: group2.ID,
				},
			},
		},
	})
	assert.Equal(s.T(), http.StatusCreated, createResource2Http.StatusCode)

	searchResource1, searchResources1Http := s.SearchResources(s.T(), ctx, user1, 10, 0, "4ccf1c0f-d791-437b-becd-8c4592d3bc1d", resource.Offer, &group1.ID)
	assert.Equal(s.T(), http.StatusOK, searchResources1Http.StatusCode)
	assert.Equal(s.T(), 10, searchResource1.Take)
	assert.Equal(s.T(), 0, searchResource1.Skip)

	assert.Equal(s.T(), 1, searchResource1.TotalCount)
	if !assert.Len(s.T(), searchResource1.Resources, 1) {
		return
	}
	assert.Equal(s.T(), createResource1.Resource.Id, searchResource1.Resources[0].Id)

	searchResource2, searchResources2Http := s.SearchResources(s.T(), ctx, user1, 10, 0, "4ccf1c0f-d791-437b-becd-8c4592d3bc1d", resource.Offer, &group2.ID)
	assert.Equal(s.T(), http.StatusOK, searchResources2Http.StatusCode)
	assert.Equal(s.T(), 10, searchResource2.Take)
	assert.Equal(s.T(), 0, searchResource2.Skip)
	assert.Equal(s.T(), 2, len(searchResource2.Resources))
	assert.Equal(s.T(), 2, searchResource2.TotalCount)
	if !assert.Len(s.T(), searchResource1.Resources, 2) {
		return
	}
	assert.Equal(s.T(), createResource1.Resource.Id, searchResource1.Resources[0].Id)
	assert.Equal(s.T(), createResource2.Resource.Id, searchResource2.Resources[1].Id)

}

func (s *IntegrationTestSuite) TestUserCanGetResource() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	res, _ := s.CreateResource(s.T(), ctx, user1)

	getResource, httpRes := s.GetResource(s.T(), ctx, user1, res.Resource.Id)
	assert.Equal(s.T(), http.StatusOK, httpRes.StatusCode)
	assert.Equal(s.T(), res.Resource.Id, getResource.Resource.Id)

}

func (s *IntegrationTestSuite) TestUserCanUpdateResource() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	res, _ := s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			Summary:          "Snippers Boop",
			Description:      "Description",
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
		},
	})

	updateResource, httpRes := s.UpdateResource(s.T(), ctx, user1, res.Resource.Id, &handler.UpdateResourceRequest{
		Resource: handler.UpdateResourcePayload{
			Summary:          "New Summary",
			Description:      "New Description",
			ValueInHoursFrom: 5,
			ValueInHoursTo:   10,
		},
	})
	assert.Equal(s.T(), http.StatusOK, httpRes.StatusCode)

	assert.Equal(s.T(), res.Resource.Id, updateResource.Resource.Id)
	assert.Equal(s.T(), "New Summary", updateResource.Resource.Summary)
	assert.Equal(s.T(), "New Description", updateResource.Resource.Description)
	assert.Equal(s.T(), 5, updateResource.Resource.ValueInHoursFrom)
	assert.Equal(s.T(), 10, updateResource.Resource.ValueInHoursTo)

}

func (s *IntegrationTestSuite) TestUserCanUpdateResourceSharings() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	group1, err := s.testGroup(s.T(), user1)
	if !assert.NoError(s.T(), err) {
		return
	}
	group2, err := s.testGroup(s.T(), user1)
	if !assert.NoError(s.T(), err) {
		return
	}

	createResource, createResourceHttp := s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			Summary:          "Snippers Boop",
			Description:      "Description",
			Type:             resource.Offer,
			SubType:          resource.ObjectResource,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith: []handler.InputResourceSharing{
				{GroupID: group1.ID},
				{GroupID: group2.ID},
			},
		},
	})
	if !AssertStatusCreated(s.T(), createResourceHttp) {
		return
	}

	updateResource1, updateResource1Http := s.UpdateResource(s.T(), ctx, user1, createResource.Resource.Id, &handler.UpdateResourceRequest{
		Resource: handler.UpdateResourcePayload{
			Summary:          "New Summary",
			Description:      "New Description",
			ValueInHoursFrom: 5,
			ValueInHoursTo:   10,
			SharedWith: []handler.InputResourceSharing{
				{GroupID: group1.ID},
			},
		},
	})

	if !AssertOK(s.T(), updateResource1Http) {
		return
	}
	if !assert.Len(s.T(), updateResource1.Resource.SharedWith, 1) {
		return
	}

	updateResource2, updateResource2Http := s.UpdateResource(s.T(), ctx, user1, createResource.Resource.Id, &handler.UpdateResourceRequest{
		Resource: handler.UpdateResourcePayload{
			Summary:          "New Summary",
			Description:      "New Description",
			ValueInHoursFrom: 5,
			ValueInHoursTo:   10,
			SharedWith: []handler.InputResourceSharing{
				{GroupID: group1.ID},
				{GroupID: group2.ID},
			},
		},
	})
	if !AssertOK(s.T(), updateResource2Http) {
		return
	}
	if !assert.Len(s.T(), updateResource2.Resource.SharedWith, 1) {
		return
	}
}

func (s *IntegrationTestSuite) TestUserCanShareResourceWithGroup() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	group1, err := s.testGroup(s.T(), user1)
	if !assert.NoError(s.T(), err) {
		return
	}

	res, httpRes := s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			SharedWith: []handler.InputResourceSharing{{GroupID: group1.ID}},
		},
	})

	if !AssertStatusCreated(s.T(), httpRes) {
		return
	}
	if !assert.Len(s.T(), res.Resource.SharedWith, 1) {
		return
	}
	assert.Equal(s.T(), group1.ID, res.Resource.SharedWith[0].GroupID)
	assert.Equal(s.T(), group1.Name, res.Resource.SharedWith[0].GroupName)

}

func (s *IntegrationTestSuite) TestUserCanShareResourceWithMultipleGroups() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	group1, err := s.testGroup(s.T(), user1)
	if !assert.NoError(s.T(), err) {
		return
	}
	group2, err := s.testGroup(s.T(), user1)
	if !assert.NoError(s.T(), err) {
		return
	}
	group3, err := s.testGroup(s.T(), user1)
	if !assert.NoError(s.T(), err) {
		return
	}

	res, httpRes := s.CreateResource(s.T(), ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			SharedWith: []handler.InputResourceSharing{
				{GroupID: group1.ID},
				{GroupID: group2.ID},
				{GroupID: group3.ID},
			},
		},
	})

	if !AssertStatusCreated(s.T(), httpRes) {
		return
	}
	if !assert.Len(s.T(), res.Resource.SharedWith, 3) {
		return
	}

	for _, groupId := range []string{group1.ID, group2.ID, group3.ID} {
		found := false
		for _, sharing := range res.Resource.SharedWith {
			if sharing.GroupID == groupId {
				found = true
				break
			}
		}
		assert.Equal(s.T(), true, found, "resource sharings should contain group id %s", groupId)
	}

}
