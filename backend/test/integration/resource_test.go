package integration

import (
	"context"
	handler2 "github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/resource/handler"
	"github.com/commonpool/backend/test"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserCanCreateResource(t *testing.T) {

	ctx := context.Background()

	_, user1Cli := testUserCli(t)
	var resource handler.GetResourceResponse
	if !assert.NoError(t, user1Cli.CreateResource(ctx, &handler.CreateResourceRequest{
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

	assert.Equal(t, "Summary", resource.Resource.Name)
	assert.Equal(t, "Description", resource.Resource.Description)
	assert.Equal(t, domain.ObjectResource, resource.Resource.ResourceType)
	assert.Equal(t, 3*time.Hour, resource.Resource.Value.ValueFromDuration)
	assert.Equal(t, 4*time.Hour, resource.Resource.Value.ValueToDuration)

}

func TestUserCanSearchResources(t *testing.T) {

	ctx := context.Background()
	uid := uuid.NewV4().String()

	_, user1Cli := testUserCli(t)

	if !assert.NoError(t, user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid)).AsRequest(), &handler.GetResourceResponse{})) {
		return
	}

	response := &handler.SearchResourcesResponse{}
	if !assert.NoError(t, user1Cli.SearchResources(ctx, uid, nil, nil, 0, 10, nil, response)) {
		return
	}

	assert.Equal(t, 10, response.Take)
	assert.Equal(t, 0, response.Skip)
	if !assert.Len(t, response.Resources, 1) {
		return
	}
	assert.Equal(t, uid, response.Resources[0].Name)
}

func TestUserCanSearchResourcesWhenNoMatch(t *testing.T) {

	ctx := context.Background()
	uid1 := uuid.NewV4().String()
	uid2 := uuid.NewV4().String()

	_, user1Cli := testUserCli(t)

	if !assert.NoError(t, user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid1)).AsRequest(), &handler.GetResourceResponse{})) {
		return
	}

	response := &handler.SearchResourcesResponse{}
	if !assert.NoError(t, user1Cli.SearchResources(ctx, uid2, nil, nil, 0, 10, nil, response)) {
		return
	}

	assert.Equal(t, 10, response.Take)
	assert.Equal(t, 0, response.Skip)
	assert.Len(t, response.Resources, 0)
}

func TestUserCanSearchResourcesWithSkip(t *testing.T) {

	ctx := context.Background()
	uid := uuid.NewV4().String()

	_, user1Cli := testUserCli(t)

	if !assert.NoError(t, user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid)).AsRequest(), &handler.GetResourceResponse{})) {
		return
	}

	if !assert.NoError(t, user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid)).AsRequest(), &handler.GetResourceResponse{})) {
		return
	}

	time.Sleep(1 * time.Second)

	response := &handler.SearchResourcesResponse{}
	if !assert.NoError(t, user1Cli.SearchResources(ctx, uid, nil, nil, 1, 10, nil, response)) {
		return
	}

	assert.Len(t, response.Resources, 1)

}

func TestUserCanSearchResourcesSharedWithGroup(t *testing.T) {

	ctx := context.Background()
	uid := uuid.NewV4().String()

	user1, user1Cli := testUserCli(t)

	group1 := &handler2.GetGroupResponse{}
	group2 := &handler2.GetGroupResponse{}

	if !assert.NoError(t, testGroup2(t, user1, group1)) {
		return
	}

	if !assert.NoError(t, testGroup2(t, user1, group2)) {
		return
	}

	time.Sleep(500 * time.Millisecond)

	resourceInGroup1 := &handler.GetResourceResponse{}
	resourceInGroup2 := &handler.GetResourceResponse{}
	resourceInBothGroups := &handler.GetResourceResponse{}
	resourceInNoGroups := &handler.GetResourceResponse{}

	if !assert.NoError(t, user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid), group1).AsRequest(), resourceInGroup1)) {
		return
	}

	time.Sleep(500 * time.Millisecond)

	if !assert.NoError(t, user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid), group2).AsRequest(), resourceInGroup2)) {
		return
	}

	time.Sleep(500 * time.Millisecond)

	if !assert.NoError(t, user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid), group1, group2).AsRequest(), resourceInBothGroups)) {
		return
	}

	time.Sleep(500 * time.Millisecond)

	if !assert.NoError(t, user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName(uid)).AsRequest(), resourceInNoGroups)) {
		return
	}

	time.Sleep(1 * time.Second)

	searchedInGroup1 := &handler.SearchResourcesResponse{}
	if !assert.NoError(t, user1Cli.SearchResources(ctx, uid, nil, nil, 0, 10, group1, searchedInGroup1)) {
		return
	}

	searchedInGroup2 := &handler.SearchResourcesResponse{}
	if !assert.NoError(t, user1Cli.SearchResources(ctx, uid, nil, nil, 0, 10, group2, searchedInGroup2)) {
		return
	}

	assert.Contains(t, searchedInGroup1.GetResourceKeys(), resourceInGroup1.Resource.ResourceKey)
	assert.Contains(t, searchedInGroup1.GetResourceKeys(), resourceInBothGroups.Resource.ResourceKey)
	assert.NotContains(t, searchedInGroup1.GetResourceKeys(), resourceInNoGroups.Resource.ResourceKey)
	assert.NotContains(t, searchedInGroup1.GetResourceKeys(), resourceInGroup2.Resource.ResourceKey)

	assert.NotContains(t, searchedInGroup2.GetResourceKeys(), resourceInGroup1.Resource.ResourceKey)
	assert.Contains(t, searchedInGroup2.GetResourceKeys(), resourceInBothGroups.Resource.ResourceKey)
	assert.NotContains(t, searchedInGroup2.GetResourceKeys(), resourceInNoGroups.Resource.ResourceKey)
	assert.Contains(t, searchedInGroup2.GetResourceKeys(), resourceInGroup2.Resource.ResourceKey)

	if !assert.Len(t, searchedInGroup1.Resources, 2) {
		return
	}
	if !assert.Len(t, searchedInGroup2.Resources, 2) {
		return
	}

}

func TestUserCanGetResource(t *testing.T) {

	ctx := context.Background()

	_, user1Cli := testUserCli(t)

	var createdResource handler.GetResourceResponse
	if !assert.NoError(t, user1Cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo().WithName("TestUserCanGetResource")).AsRequest(), &createdResource)) {
		return
	}

	var gottenResource handler.GetResourceResponse
	if !assert.NoError(t, user1Cli.GetResource(ctx, createdResource, &gottenResource)) {
		return
	}

	assert.Equal(t, createdResource, gottenResource)

}

func TestUserCanUpdateResource(t *testing.T) {

	ctx := context.Background()
	_, cli := testUserCli(t)

	first := test.AResourceInfo().WithName("TestUserCanGetResource").WithDescription("first")
	second := first.
		WithName("second-name").
		WithDescription("second-description").
		WithValue(domain.NewResourceValueEstimation().WithHoursFromTo(9, 18))

	var createdResource handler.GetResourceResponse
	if !assert.NoError(t, cli.CreateResource(ctx, handler.NewCreateResourcePayload(first).AsRequest(), &createdResource)) {
		return
	}

	var updatedResource handler.GetResourceResponse
	if !assert.NoError(t, cli.UpdateResource(ctx, createdResource, handler.NewUpdateResourcePayload(second.AsUpdate()).AsRequest(), &updatedResource)) {
		return
	}

	var gottenResource handler.GetResourceResponse
	if !assert.NoError(t, cli.GetResource(ctx, createdResource, &gottenResource)) {
		return
	}

	assert.Equal(t, "second-name", gottenResource.Resource.ResourceInfo.Name)
	assert.Equal(t, "second-description", gottenResource.Resource.ResourceInfo.Description)
	assert.Equal(t, 9*time.Hour, gottenResource.Resource.ResourceInfo.Value.ValueFromDuration)
	assert.Equal(t, 18*time.Hour, gottenResource.Resource.ResourceInfo.Value.ValueToDuration)

}

func TestUserCanUpdateResourceSharings(t *testing.T) {

	ctx := context.Background()
	user, cli := testUserCli(t)

	group1 := &handler2.GetGroupResponse{}
	group2 := &handler2.GetGroupResponse{}

	if !assert.NoError(t, testGroup2(t, user, group1)) {
		return
	}
	if !assert.NoError(t, testGroup2(t, user, group2)) {
		return
	}

	var resource handler.GetResourceResponse
	if !assert.NoError(t, cli.CreateResource(ctx, handler.NewCreateResourcePayload(test.AResourceInfo()).AsRequest(), &resource)) {
		return
	}

	assert.Len(t, resource.Resource.Sharings, 0)

	if !assert.NoError(t, cli.UpdateResource(ctx, resource, resource.AsUpdate().WithShared(group1), &resource)) {
		return
	}

	assert.Len(t, resource.Resource.Sharings, 1)

	if !assert.NoError(t, cli.UpdateResource(ctx, resource, resource.AsUpdate().WithShared(group1, group2), &resource)) {
		return
	}

	assert.Len(t, resource.Resource.Sharings, 2)

	if !assert.NoError(t, cli.UpdateResource(ctx, resource, resource.AsUpdate().WithShared(), &resource)) {
		return
	}

	assert.Len(t, resource.Resource.Sharings, 0)
}
