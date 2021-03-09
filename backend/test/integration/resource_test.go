package integration

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/models"
	handler2 "github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/resource/handler"
	"github.com/commonpool/backend/test"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
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

func (s *IntegrationTestSuite) TestUserCanCreateResource() {

	ctx := context.Background()

	_, user1Cli := s.testUserCli(s.T())
	var resource handler.GetResourceResponse
	if !s.NoError(user1Cli.CreateResource(ctx, &handler.CreateResourceRequest{
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
	}, &resource)) {
		return
	}

	s.Equal("Summary", resource.Resource.Name)
	s.Equal("Description", resource.Resource.Description)
	s.Equal(domain.ObjectResource, resource.Resource.ResourceType)
	s.Equal(3*time.Hour, resource.Resource.Value.ValueFromDuration)
	s.Equal(4*time.Hour, resource.Resource.Value.ValueToDuration)

}

func (s *IntegrationTestSuite) TestUserCanSearchResources() {

	ctx := context.Background()
	uid := uuid.NewV4().String()

	_, user1Cli := s.testUserCli(s.T())

	if !s.NoError(user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid)).AsRequest(), &handler.GetResourceResponse{})) {
		return
	}

	response := &handler.SearchResourcesResponse{}
	if !s.NoError(user1Cli.SearchResources(ctx, uid, nil, nil, 0, 10, nil, response)) {
		return
	}

	s.Equal(10, response.Take)
	s.Equal(0, response.Skip)
	if !s.Len(response.Resources, 1) {
		return
	}
	s.Equal(uid, response.Resources[0].Name)
}

func (s *IntegrationTestSuite) TestUserCanSearchResourcesWhenNoMatch() {

	ctx := context.Background()
	uid1 := uuid.NewV4().String()
	uid2 := uuid.NewV4().String()

	_, user1Cli := s.testUserCli(s.T())

	if !s.NoError(user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid1)).AsRequest(), &handler.GetResourceResponse{})) {
		return
	}

	response := &handler.SearchResourcesResponse{}
	if !s.NoError(user1Cli.SearchResources(ctx, uid2, nil, nil, 0, 10, nil, response)) {
		return
	}

	s.Equal(10, response.Take)
	s.Equal(0, response.Skip)
	s.Len(response.Resources, 0)
}

func (s *IntegrationTestSuite) TestUserCanSearchResourcesWithSkip() {

	ctx := context.Background()
	uid := uuid.NewV4().String()

	_, user1Cli := s.testUserCli(s.T())

	if !s.NoError(user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid)).AsRequest(), &handler.GetResourceResponse{})) {
		return
	}

	if !s.NoError(user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid)).AsRequest(), &handler.GetResourceResponse{})) {
		return
	}

	time.Sleep(1 * time.Second)

	response := &handler.SearchResourcesResponse{}
	if !s.NoError(user1Cli.SearchResources(ctx, uid, nil, nil, 1, 10, nil, response)) {
		return
	}

	s.Len(response.Resources, 1)

}

func (s *IntegrationTestSuite) TestUserCanSearchResourcesSharedWithGroup() {

	ctx := context.Background()
	uid := uuid.NewV4().String()

	user1, user1Cli := s.testUserCli(s.T())

	group1 := &handler2.GetGroupResponse{}
	group2 := &handler2.GetGroupResponse{}

	if !s.NoError(s.testGroup2(s.T(), user1, group1)) {
		return
	}

	if !s.NoError(s.testGroup2(s.T(), user1, group2)) {
		return
	}

	time.Sleep(500 * time.Millisecond)

	resourceInGroup1 := &handler.GetResourceResponse{}
	resourceInGroup2 := &handler.GetResourceResponse{}
	resourceInBothGroups := &handler.GetResourceResponse{}
	resourceInNoGroups := &handler.GetResourceResponse{}

	if !s.NoError(user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid), group1).AsRequest(), resourceInGroup1)) {
		return
	}

	time.Sleep(500 * time.Millisecond)

	if !s.NoError(user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid), group2).AsRequest(), resourceInGroup2)) {
		return
	}

	time.Sleep(500 * time.Millisecond)

	if !s.NoError(user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid), group1, group2).AsRequest(), resourceInBothGroups)) {
		return
	}

	time.Sleep(500 * time.Millisecond)

	if !s.NoError(user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid)).AsRequest(), resourceInNoGroups)) {
		return
	}

	time.Sleep(1 * time.Second)

	searchedInGroup1 := &handler.SearchResourcesResponse{}
	if !s.NoError(user1Cli.SearchResources(ctx, uid, nil, nil, 0, 10, group1, searchedInGroup1)) {
		return
	}

	searchedInGroup2 := &handler.SearchResourcesResponse{}
	if !s.NoError(user1Cli.SearchResources(ctx, uid, nil, nil, 0, 10, group2, searchedInGroup2)) {
		return
	}

	s.Contains(searchedInGroup1.GetResourceKeys(), resourceInGroup1.Resource.ResourceKey)
	s.Contains(searchedInGroup1.GetResourceKeys(), resourceInBothGroups.Resource.ResourceKey)
	s.NotContains(searchedInGroup1.GetResourceKeys(), resourceInNoGroups.Resource.ResourceKey)
	s.NotContains(searchedInGroup1.GetResourceKeys(), resourceInGroup2.Resource.ResourceKey)

	s.NotContains(searchedInGroup2.GetResourceKeys(), resourceInGroup1.Resource.ResourceKey)
	s.Contains(searchedInGroup2.GetResourceKeys(), resourceInBothGroups.Resource.ResourceKey)
	s.NotContains(searchedInGroup2.GetResourceKeys(), resourceInNoGroups.Resource.ResourceKey)
	s.Contains(searchedInGroup2.GetResourceKeys(), resourceInGroup2.Resource.ResourceKey)

	if !s.Len(searchedInGroup1.Resources, 2) {
		return
	}
	if !s.Len(searchedInGroup2.Resources, 2) {
		return
	}

}

func (s *IntegrationTestSuite) TestUserCanGetResource() {

	ctx := context.Background()

	_, user1Cli := s.testUserCli(s.T())

	var createdResource handler.GetResourceResponse
	if !s.NoError(user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName("TestUserCanGetResource")).AsRequest(), &createdResource)) {
		return
	}

	var gottenResource handler.GetResourceResponse
	if !s.NoError(user1Cli.GetResource(ctx, createdResource, &gottenResource)) {
		return
	}

	s.Equal(createdResource, gottenResource)

}

func (s *IntegrationTestSuite) TestUserCanUpdateResource() {

	ctx := context.Background()
	_, cli := s.testUserCli(s.T())

	first := test.AResourceInfo().WithName("TestUserCanGetResource").WithDescription("first")
	second := first.
		WithName("second-name").
		WithDescription("second-description").
		WithValue(domain.NewResourceValueEstimation().WithHoursFromTo(9, 18))

	var createdResource handler.GetResourceResponse
	if !s.NoError(cli.CreateResource(ctx, handler.NewCreateResourcePayload(first).AsRequest(), &createdResource)) {
		return
	}

	var updatedResource handler.GetResourceResponse
	if !s.NoError(cli.UpdateResource(ctx, createdResource, handler.NewUpdateResourcePayload(second.AsUpdate()).AsRequest(), &updatedResource)) {
		return
	}

	var gottenResource handler.GetResourceResponse
	if !s.NoError(cli.GetResource(ctx, createdResource, &gottenResource)) {
		return
	}

	s.Equal("second-name", gottenResource.Resource.ResourceInfo.Name)
	s.Equal("second-description", gottenResource.Resource.ResourceInfo.Description)
	s.Equal(9*time.Hour, gottenResource.Resource.ResourceInfo.Value.ValueFromDuration)
	s.Equal(18*time.Hour, gottenResource.Resource.ResourceInfo.Value.ValueToDuration)

}

func (s *IntegrationTestSuite) TestUserCanUpdateResourceSharings() {

	ctx := context.Background()
	user, cli := s.testUserCli(s.T())

	group1 := &handler2.GetGroupResponse{}
	group2 := &handler2.GetGroupResponse{}

	if !s.NoError(s.testGroup2(s.T(), user, group1)) {
		return
	}
	if !s.NoError(s.testGroup2(s.T(), user, group2)) {
		return
	}

	var resource handler.GetResourceResponse
	if !s.NoError(cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resource)) {
		return
	}

	s.Len(resource.Resource.Sharings, 0)

	if !s.NoError(cli.UpdateResource(ctx, resource, resource.AsUpdate().WithShared(group1), &resource)) {
		return
	}

	s.Len(resource.Resource.Sharings, 1)

	if !s.NoError(cli.UpdateResource(ctx, resource, resource.AsUpdate().WithShared(group1, group2), &resource)) {
		return
	}

	s.Len(resource.Resource.Sharings, 2)

	if !s.NoError(cli.UpdateResource(ctx, resource, resource.AsUpdate().WithShared(), &resource)) {
		return
	}

	s.Len(resource.Resource.Sharings, 0)
}
