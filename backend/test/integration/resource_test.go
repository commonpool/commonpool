package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/resource/handler"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

var resourceCounter = 1

func (s *IntegrationTestSuite) CreateResource(ctx context.Context, userSession *models.UserSession, opts ...*handler.CreateResourceRequest) (*handler.GetResourceResponse, *http.Response) {

	resourceCounter++
	var resourceName = "resource-" + strconv.Itoa(resourceCounter)

	payload := &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			ResourceInfo: domain.ResourceInfo{
				ResourceInfoBase: domain.ResourceInfoBase{
					Name:         resourceName,
					Description:  resourceName + "-description",
					ResourceType: domain.ObjectResource,
					CallType:     domain.Offer,
				},
				Value: domain.ResourceValueEstimation{
					ValueType:         domain.FromToDuration,
					ValueFromDuration: 0,
					ValueToDuration:   0,
				},
			},
			SharedWith: []handler.InputResourceSharing{},
		},
	}

	for _, option := range opts {

		if option.Resource.ResourceInfo.Name != "" {
			payload.Resource.ResourceInfo.Name = option.Resource.ResourceInfo.Name
		}
		if option.Resource.ResourceInfo.Description != "" {
			payload.Resource.ResourceInfo.Description = option.Resource.ResourceInfo.Description
		}
		if option.Resource.SharedWith != nil {
			for _, sharing := range option.Resource.SharedWith {
				payload.Resource.SharedWith = append(payload.Resource.SharedWith, sharing)
			}
		}
		if option.Resource.ResourceInfo.Value.ValueToDuration != 0 {
			payload.Resource.ResourceInfo.Value.ValueToDuration = option.Resource.ResourceInfo.Value.ValueToDuration
		}
		if option.Resource.ResourceInfo.Value.ValueFromDuration != 0 {
			payload.Resource.ResourceInfo.Value.ValueFromDuration = option.Resource.ResourceInfo.Value.ValueFromDuration
		}
		if option.Resource.ResourceInfo.ResourceType != "" {
			payload.Resource.ResourceInfo.ResourceType = option.Resource.ResourceInfo.ResourceType
		}
		if option.Resource.ResourceInfo.CallType != "" {
			payload.Resource.ResourceInfo.CallType = option.Resource.ResourceInfo.CallType
		}
	}

	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/resources", payload)
	s.server.Router.ServeHTTP(recorder, httpReq)
	response := &handler.GetResourceResponse{}
	s.T().Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) SearchResources(
	t *testing.T,
	ctx context.Context,
	userSession *models.UserSession,
	take int,
	skip int,
	query string,
	callType domain.CallType,
	sharedWithGroup *keys.GroupKey) (*handler.SearchResourcesResponse, *http.Response) {
	var sharedStr = ""
	if sharedWithGroup != nil {
		sharedStr = sharedWithGroup.String()
	}
	target, _ := url.Parse(fmt.Sprintf("/api/v1/resources?take=%d&skip=%d&query=%s&call=%s&group_id=%s", take, skip, query, callType, sharedStr))

	httpReq, recorder := NewRequest(ctx, userSession, http.MethodGet, target.String(), nil)
	s.server.Router.ServeHTTP(recorder, httpReq)
	response := &handler.SearchResourcesResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)

}

func (s *IntegrationTestSuite) GetResource(t *testing.T, ctx context.Context, userSession *models.UserSession, resourceKey keys.ResourceKey) (*handler.GetResourceResponse, *http.Response) {
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodGet, fmt.Sprintf("/api/v1/resources/%s", resourceKey.String()), nil)
	s.server.Router.ServeHTTP(recorder, httpReq)
	response := &handler.GetResourceResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) UpdateResource(t *testing.T, ctx context.Context, userSession *models.UserSession, resourceKey keys.ResourceKey, request *handler.UpdateResourceRequest) (*handler.GetResourceResponse, *http.Response) {
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPut, fmt.Sprintf("/api/v1/resources/%s", resourceKey.String()), request)
	s.server.Router.ServeHTTP(recorder, httpReq)
	response := &handler.GetResourceResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) TestUserCanCreateResource() {
	s.T().Parallel()

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	resp, httpResp := s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			ResourceInfo: domain.ResourceInfo{
				ResourceInfoBase: domain.ResourceInfoBase{
					Name:         "Summary",
					Description:  "Description",
					CallType:     domain.Request,
					ResourceType: domain.ObjectResource,
				},
				Value: domain.ResourceValueEstimation{
					ValueType:         domain.FromToDuration,
					ValueFromDuration: time.Hour * 3,
					ValueToDuration:   time.Hour * 4,
				},
			},
			SharedWith: []handler.InputResourceSharing{},
		},
	})
	if !AssertStatusCreated(s.T(), httpResp) {
		return
	}

	assert.Equal(s.T(), "Summary", resp.Resource.Name)
	assert.Equal(s.T(), "Description", resp.Resource.Description)
	assert.Equal(s.T(), domain.ObjectResource, resp.Resource.ResourceType)
	assert.Equal(s.T(), 3*time.Hour, resp.Resource.Value.ValueFromDuration)
	assert.Equal(s.T(), 4*time.Hour, resp.Resource.Value.ValueToDuration)

}

func (s *IntegrationTestSuite) TestUserCanSearchResources() {
	s.T().Parallel()

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()
	uid := uuid.NewV4()

	ctx := context.Background()

	_, httpResponse := s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			ResourceInfo: domain.ResourceInfo{
				ResourceInfoBase: domain.ResourceInfoBase{
					Name: uid.String() + "-Bla",
				},
			},
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	res, httpRes := s.SearchResources(s.T(), ctx, user1, 10, 0, uid.String(), domain.Offer, nil)
	if !AssertOK(s.T(), httpRes) {
		return
	}

	assert.Equal(s.T(), 10, res.Take)
	assert.Equal(s.T(), 0, res.Skip)

	if !assert.Len(s.T(), res.Resources, 1) {
		return
	}
	assert.Equal(s.T(), uid.String()+"-Bla", res.Resources[0].Name)

}

func (s *IntegrationTestSuite) TestUserCanSearchResourcesWhenNoMatch() {
	s.T().Parallel()

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			ResourceInfo: domain.ResourceInfo{
				ResourceInfoBase: domain.ResourceInfoBase{
					Name: "SnizzledPopPiDouWabWap",
				},
			},
		},
	})

	res, httpRes := s.SearchResources(s.T(), ctx, user1, 10, 0, "ResourceNoMatchQuery", domain.Offer, nil)
	assert.Equal(s.T(), http.StatusOK, httpRes.StatusCode)

	assert.Equal(s.T(), 10, res.Take)
	assert.Equal(s.T(), 0, res.Skip)
	assert.Equal(s.T(), 0, len(res.Resources))

}

func (s *IntegrationTestSuite) TestUserCanSearchResourcesWithSkip() {
	s.T().Parallel()

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()
	uid := uuid.NewV4().String()

	s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			ResourceInfo: domain.ResourceInfo{
				ResourceInfoBase: domain.ResourceInfoBase{
					Name: uid,
				},
			},
		},
	})
	s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			ResourceInfo: domain.ResourceInfo{
				ResourceInfoBase: domain.ResourceInfoBase{
					Name: uid,
				},
			},
		},
	})

	res, httpRes := s.SearchResources(s.T(), ctx, user1, 10, 1, uid, domain.Offer, nil)
	assert.Equal(s.T(), http.StatusOK, httpRes.StatusCode)

	assert.Equal(s.T(), 10, res.Take)
	assert.Equal(s.T(), 1, res.Skip)
	assert.Equal(s.T(), 1, len(res.Resources))
	assert.Equal(s.T(), uid, res.Resources[0].Name)

}

func (s *IntegrationTestSuite) TestUserCanSearchResourcesSharedWithGroup() {
	s.T().Parallel()

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

	uid := uuid.NewV4()

	s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			ResourceInfo: domain.ResourceInfo{
				ResourceInfoBase: domain.ResourceInfoBase{
					Name: uid.String(),
				},
			},
			SharedWith: handler.NewInputResourceSharings().WithGroups(group1.GroupKey),
		},
	})
	s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			ResourceInfo: domain.ResourceInfo{
				ResourceInfoBase: domain.ResourceInfoBase{
					Name: uid.String(),
				},
			},
			SharedWith: handler.NewInputResourceSharings().WithGroups(group2.GroupKey),
		},
	})

	res, httpRes := s.SearchResources(s.T(), ctx, user1, 10, 0, uid.String(), domain.Request, &group1.GroupKey)
	assert.Equal(s.T(), http.StatusOK, httpRes.StatusCode)

	assert.Equal(s.T(), 10, res.Take)
	assert.Equal(s.T(), 0, res.Skip)
	if !assert.Len(s.T(), res.Resources, 1) {
		return
	}
	assert.Equal(s.T(), uid.String(), res.Resources[0].Name)

}

func (s *IntegrationTestSuite) TestUserCanSearchResourcesSharedWithMultipleGroups() {
	s.T().Parallel()

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

	uid := uuid.NewV4().String()

	createResource1, createResource1Http := s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			ResourceInfo: domain.ResourceInfo{
				ResourceInfoBase: domain.ResourceInfoBase{
					Name: uid,
				},
			},
			SharedWith: handler.NewInputResourceSharings().WithGroups(group1.GroupKey, group2.GroupKey),
		},
	})
	assert.Equal(s.T(), http.StatusCreated, createResource1Http.StatusCode)
	createResource2, createResource2Http := s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			ResourceInfo: domain.ResourceInfo{
				ResourceInfoBase: domain.ResourceInfoBase{
					Name:         uid,
					CallType:     domain.Offer,
					ResourceType: domain.ObjectResource,
				},
			},
			SharedWith: handler.NewInputResourceSharings().WithGroups(group2.GroupKey),
		},
	})
	assert.Equal(s.T(), http.StatusCreated, createResource2Http.StatusCode)

	searchResource1, searchResources1Http := s.SearchResources(s.T(), ctx, user1, 10, 0, uid, domain.Offer, &group1.GroupKey)
	assert.Equal(s.T(), http.StatusOK, searchResources1Http.StatusCode)
	assert.Equal(s.T(), 10, searchResource1.Take)
	assert.Equal(s.T(), 0, searchResource1.Skip)

	if !assert.Len(s.T(), searchResource1.Resources, 1) {
		return
	}
	assert.Equal(s.T(), createResource1.Resource.ResourceKey, searchResource1.Resources[0].ResourceKey)

	searchResource2, searchResources2Http := s.SearchResources(s.T(), ctx, user1, 10, 0, uid, domain.Offer, &group2.GroupKey)
	assert.Equal(s.T(), http.StatusOK, searchResources2Http.StatusCode)
	assert.Equal(s.T(), 10, searchResource2.Take)
	assert.Equal(s.T(), 0, searchResource2.Skip)
	assert.Equal(s.T(), 2, len(searchResource2.Resources))
	if !assert.Len(s.T(), searchResource2.Resources, 2) {
		return
	}
	assert.Equal(s.T(), createResource1.Resource.ResourceKey, searchResource1.Resources[0].ResourceKey)
	assert.Equal(s.T(), createResource2.Resource.ResourceKey, searchResource2.Resources[1].ResourceKey)

}

func (s *IntegrationTestSuite) TestUserCanGetResource() {
	s.T().Parallel()

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	res, _ := s.CreateResource(ctx, user1)

	getResource, httpRes := s.GetResource(s.T(), ctx, user1, res.Resource.ResourceKey)
	assert.Equal(s.T(), http.StatusOK, httpRes.StatusCode)
	assert.Equal(s.T(), res.Resource.ResourceKey, getResource.Resource.ResourceKey)

}

func (s *IntegrationTestSuite) TestUserCanUpdateResource() {
	s.T().Parallel()

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	res, _ := s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			ResourceInfo: domain.ResourceInfo{
				ResourceInfoBase: domain.ResourceInfoBase{
					Name:        "Snippers Boop",
					Description: "Description",
				},
				Value: domain.ResourceValueEstimation{
					ValueFromDuration: time.Hour * 1,
					ValueToDuration:   time.Hour * 3,
				},
			},
		},
	})

	updateResource, httpRes := s.UpdateResource(s.T(), ctx, user1, res.Resource.ResourceKey, &handler.UpdateResourceRequest{
		Resource: handler.UpdateResourcePayload{
			ResourceInfo: domain.ResourceInfoUpdate{
				Name:        "New Summary",
				Description: "New Description",
				Value: domain.ResourceValueEstimation{
					ValueFromDuration: time.Hour * 5,
					ValueToDuration:   time.Hour * 10,
				},
			},
			SharedWith: nil,
		},
	})
	assert.Equal(s.T(), http.StatusOK, httpRes.StatusCode)

	assert.Equal(s.T(), res.Resource.ResourceKey, updateResource.Resource.ResourceKey)
	assert.Equal(s.T(), "New Summary", updateResource.Resource.Name)
	assert.Equal(s.T(), "New Description", updateResource.Resource.Description)
	assert.Equal(s.T(), 5*time.Hour, updateResource.Resource.Value.ValueFromDuration)
	assert.Equal(s.T(), 10*time.Hour, updateResource.Resource.Value.ValueToDuration)

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

	resourceInfo := domain.ResourceInfo{
		ResourceInfoBase: domain.ResourceInfoBase{
			Name:         "Snippers Boop",
			Description:  "Description",
			ResourceType: domain.ObjectResource,
			CallType:     domain.Request,
		},
		Value: domain.ResourceValueEstimation{
			ValueFromDuration: time.Hour * 1,
			ValueToDuration:   time.Hour * 3,
		},
	}

	createResource, createResourceHttp := s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			ResourceInfo: resourceInfo,
			SharedWith:   handler.NewInputResourceSharings().WithGroups(group1.GroupKey, group2.GroupKey),
		},
	})
	if !AssertStatusCreated(s.T(), createResourceHttp) {
		return
	}

	updateResource1, updateResource1Http := s.UpdateResource(s.T(), ctx, user1, createResource.Resource.ResourceKey, &handler.UpdateResourceRequest{
		Resource: handler.UpdateResourcePayload{
			ResourceInfo: resourceInfo.AsUpdate(),
			SharedWith:   handler.NewInputResourceSharings().WithGroups(group1.GroupKey),
		},
	})

	if !AssertOK(s.T(), updateResource1Http) {
		return
	}
	if !assert.Len(s.T(), updateResource1.Resource.Sharings, 1) {
		return
	}

	updateResource2, updateResource2Http := s.UpdateResource(s.T(), ctx, user1, createResource.Resource.ResourceKey, &handler.UpdateResourceRequest{
		Resource: handler.UpdateResourcePayload{
			ResourceInfo: resourceInfo.AsUpdate(),
			SharedWith:   handler.NewInputResourceSharings().WithGroups(group1.GroupKey, group2.GroupKey),
		},
	})
	if !AssertOK(s.T(), updateResource2Http) {
		return
	}
	if !assert.Len(s.T(), updateResource2.Resource.Sharings, 2) {
		return
	}

	updateResource3, updateResource3Http := s.UpdateResource(s.T(), ctx, user1, createResource.Resource.ResourceKey, &handler.UpdateResourceRequest{
		Resource: handler.UpdateResourcePayload{
			ResourceInfo: resourceInfo.AsUpdate(),
			SharedWith:   handler.NewInputResourceSharings().WithGroups(),
		},
	})
	if !AssertOK(s.T(), updateResource3Http) {
		return
	}
	if !assert.Len(s.T(), updateResource3.Resource.Sharings, 0) {
		return
	}
}

func (s *IntegrationTestSuite) TestUserCanShareResourceWithGroup() {
	s.T().Parallel()

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	ctx := context.Background()

	group1, err := s.testGroup(s.T(), user1)
	if !assert.NoError(s.T(), err) {
		return
	}

	res, httpRes := s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			SharedWith: handler.NewInputResourceSharings().WithGroups(group1.GroupKey),
		},
	})

	if !AssertStatusCreated(s.T(), httpRes) {
		return
	}
	if !assert.Len(s.T(), res.Resource.Sharings, 1) {
		return
	}
	assert.Equal(s.T(), group1.GroupKey, res.Resource.Sharings[0].GroupKey)
	assert.Equal(s.T(), group1.Name, res.Resource.Sharings[0].GroupName)

}

func (s *IntegrationTestSuite) TestUserCanShareResourceWithMultipleGroups() {
	s.T().Parallel()

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

	res, httpRes := s.CreateResource(ctx, user1, &handler.CreateResourceRequest{
		Resource: handler.CreateResourcePayload{
			SharedWith: handler.NewInputResourceSharings().WithGroups(group1.GroupKey, group2.GroupKey, group3.GroupKey),
		},
	})

	if !AssertStatusCreated(s.T(), httpRes) {
		return
	}
	if !assert.Len(s.T(), res.Resource.Sharings, 3) {
		return
	}

	for _, groupId := range []keys.GroupKey{group1.GroupKey, group2.GroupKey, group3.GroupKey} {
		found := false
		for _, sharing := range res.Resource.Sharings {
			if sharing.GroupKey == groupId {
				found = true
				break
			}
		}
		assert.Equal(s.T(), true, found, "resource sharings should contain group id %s", groupId)
	}

}
