package domain

import (
	json2 "encoding/json"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/keys"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ResourceTestSuite struct {
	suite.Suite
	user         keys.UserKey
	user2        keys.UserKey
	userTarget   keys.Target
	resource     *Resource
	resourceInfo ResourceInfo
	group1       keys.GroupKey
	group2       keys.GroupKey
	group3       keys.GroupKey
	resourceKey  keys.ResourceKey
}

func (s *ResourceTestSuite) SetupTest() {
	s.resourceKey = keys.NewResourceKey(uuid.NewV4())
	r := NewResource(s.resourceKey)
	s.user = keys.NewUserKey("user")
	s.user2 = keys.NewUserKey("user2")
	s.userTarget = *keys.NewUserTarget(s.user)
	s.group1 = keys.NewGroupKey(uuid.NewV4())
	s.group2 = keys.NewGroupKey(uuid.NewV4())
	s.group3 = keys.NewGroupKey(uuid.NewV4())
	s.resourceInfo = ResourceInfo{
		Value: ResourceValueEstimation{
			ValueType:         FromToDuration,
			ValueFromDuration: 3 * time.Hour,
			ValueToDuration:   5 * time.Hour,
		},
		Name:         "resource",
		Description:  "description",
		CallType:     Offer,
		ResourceType: ServiceResource,
	}
	err := r.Register(s.user, s.userTarget, ServiceResource, s.resourceInfo, *keys.NewEmptyGroupKeys())
	if !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	s.resource = r
}

func (s *ResourceTestSuite) TestNewResource() {
	r := NewResource(s.resourceKey)
	user := keys.NewUserKey("user")
	err := r.Register(user, *keys.NewUserTarget(user), ServiceResource, s.resourceInfo, *keys.NewGroupKeys([]keys.GroupKey{}))
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Len(s.T(), r.GetChanges(), 1)
}

func (s *ResourceTestSuite) TestNewResourceWithSharings() {
	r := NewResource(s.resourceKey)
	user := keys.NewUserKey("user")
	err := r.Register(user, *keys.NewUserTarget(user), ServiceResource, s.resourceInfo, *keys.NewGroupKeys([]keys.GroupKey{
		s.group1,
	}))

	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Len(s.T(), r.GetChanges(), 2)
	if !assert.Len(s.T(), r.sharings, 1) {
		return
	}
	assert.Equal(s.T(), s.group1, r.sharings[0].GroupKey)

}

func (s *ResourceTestSuite) TestNewResourceShouldFailIfAlreadyCreated() {
	err := s.resource.Register(s.user, s.userTarget, ServiceResource, s.resourceInfo, *keys.NewEmptyGroupKeys())
	if !assert.Error(s.T(), err) {
		return
	}
}

func (s *ResourceTestSuite) TestChangeInfo() {
	newInfo := ResourceInfo{
		Value: ResourceValueEstimation{
			ValueType:         FromToDuration,
			ValueFromDuration: 100 * time.Hour,
			ValueToDuration:   200 * time.Hour,
		},
		Name:         "TestChangeInfo",
		Description:  "TestChangeInfo-description",
		ResourceType: ObjectResource,
		CallType:     Offer,
	}
	err := s.resource.ChangeInfo(s.user2, newInfo)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), newInfo, s.resource.info)
}

func (s *ResourceTestSuite) TestChangeInfoIdempotent() {
	currentChangeCount := len(s.resource.changes)
	err := s.resource.ChangeInfo(s.user2, s.resourceInfo)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Len(s.T(), s.resource.changes, currentChangeCount)
}

func (s *ResourceTestSuite) TestChangeInfoShouldFailIfNotRegistered() {
	r := NewResource(keys.NewResourceKey(uuid.NewV4()))
	err := r.ChangeInfo(s.user2, s.resourceInfo)
	assert.Error(s.T(), err)
}

func (s *ResourceTestSuite) TestChangeInfoShouldFailIfDeleted() {
	if err := s.resource.Delete(s.user2); !assert.NoError(s.T(), err) {
		return
	}
	err := s.resource.ChangeInfo(s.user2, s.resourceInfo)
	assert.Error(s.T(), err)
}

func (s *ResourceTestSuite) TestChangeInfoShouldFailIfNew() {
	r := NewResource(s.resourceKey)
	err := r.ChangeInfo(s.user2, s.resourceInfo)
	assert.Error(s.T(), err)
}

func (s *ResourceTestSuite) TestChangeSharings() {

	err := s.resource.ChangeSharings(s.user2, *keys.NewGroupKeys([]keys.GroupKey{
		s.group1,
	}))
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.Len(s.T(), s.resource.sharings, 1) {
		return
	}

	assert.Equal(s.T(), s.group1, s.resource.sharings[0].GroupKey)

}

func (s *ResourceTestSuite) TestChangeSharingsShouldFailIfNew() {
	r := NewResource(s.resourceKey)
	err := r.ChangeSharings(s.user2, *keys.NewEmptyGroupKeys())
	assert.Error(s.T(), err)
}

func (s *ResourceTestSuite) TestChangeSharingsShouldFailIfDeleted() {
	err := s.resource.Delete(s.user2)
	if !assert.NoError(s.T(), err) {
		return
	}
	err = s.resource.ChangeSharings(s.user2, *keys.NewEmptyGroupKeys())
	assert.Error(s.T(), err)
}

func (s *ResourceTestSuite) TestChangeSharingsDuplicateKeys() {
	err := s.resource.ChangeSharings(s.user2, *keys.NewGroupKeys([]keys.GroupKey{
		s.group1,
		s.group1,
		s.group1,
	}))
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Len(s.T(), s.resource.sharings, 1)
}

func (s *ResourceTestSuite) TestChangeSharingsIdempotent() {

	for i := 0; i < 2; i++ {
		err := s.resource.ChangeSharings(s.user2, *keys.NewGroupKeys([]keys.GroupKey{
			s.group1,
		}))
		if !assert.NoError(s.T(), err) {
			return
		}
	}

	if !assert.Len(s.T(), s.resource.sharings, 1) {
		return
	}

	assert.Equal(s.T(), s.group1, s.resource.sharings[0].GroupKey)
	assert.Len(s.T(), s.resource.changes, 2)

}
func (s *ResourceTestSuite) TestChangeSharingsEvtDifference() {

	err := s.resource.ChangeSharings(s.user2, *keys.NewGroupKeys([]keys.GroupKey{
		s.group1,
		s.group2,
	}))
	if !assert.NoError(s.T(), err) {
		return
	}

	err = s.resource.ChangeSharings(s.user2, *keys.NewGroupKeys([]keys.GroupKey{
		s.group2,
		s.group3,
	}))
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.Len(s.T(), s.resource.sharings, 2) {
		return
	}

	// Assert state
	assert.Equal(s.T(), s.group2, s.resource.sharings[0].GroupKey)
	assert.Equal(s.T(), s.group3, s.resource.sharings[1].GroupKey)
	assert.Len(s.T(), s.resource.changes, 3)

	// Assert 1st resource sharing change
	event, ok := s.resource.changes[1].(ResourceGroupSharingChanged)
	if !assert.Equal(s.T(), true, ok) {
		return
	}

	if !assert.Len(s.T(), event.NewResourceSharings, 2) {
		return
	}
	assert.Equal(s.T(), s.group1, event.NewResourceSharings[0].GroupKey)
	assert.Equal(s.T(), s.group2, event.NewResourceSharings[1].GroupKey)

	// Assert 2nd resource sharing change
	event, ok = s.resource.changes[2].(ResourceGroupSharingChanged)
	if !assert.Equal(s.T(), true, ok) {
		return
	}

	if !assert.Len(s.T(), event.NewResourceSharings, 2) {
		return
	}
	assert.Equal(s.T(), s.group2, event.NewResourceSharings[0].GroupKey)
	assert.Equal(s.T(), s.group3, event.NewResourceSharings[1].GroupKey)
	assert.Equal(s.T(), s.group1, event.OldResourceSharings[0].GroupKey)
	assert.Equal(s.T(), s.group2, event.OldResourceSharings[1].GroupKey)
	assert.Equal(s.T(), s.group3, event.AddedSharings[0].GroupKey)
	assert.Equal(s.T(), s.group1, event.RemovedSharings[0].GroupKey)

}

func (s *ResourceTestSuite) TestDeleteResource() {
	if err := s.resource.Delete(s.user2); !assert.NoError(s.T(), err) {
		return
	}
	assert.True(s.T(), s.resource.isDeleted)
}

func (s *ResourceTestSuite) TestDeleteResourceIdempotent() {
	c := len(s.resource.changes)
	for i := 0; i < 2; i++ {
		if err := s.resource.Delete(s.user2); !assert.NoError(s.T(), err) {
			return
		}
	}
	assert.Len(s.T(), s.resource.changes, c+1)
}

func (s *ResourceTestSuite) TestDeleteResourceShouldFailIfNew() {
	r := NewResource(s.resourceKey)
	err := r.Delete(s.user2)
	assert.Error(s.T(), err)
}

func (s *ResourceTestSuite) TestMarkAsCommitted() {
	c := len(s.resource.GetChanges())
	v := s.resource.GetVersion()
	s.resource.MarkAsCommitted()
	assert.Equal(s.T(), v+c, s.resource.GetVersion())
	assert.Len(s.T(), s.resource.GetChanges(), 0)
}

func (s *ResourceTestSuite) TestFromEvents() {
	r := NewFromEvents(s.resourceKey, s.resource.GetChanges())
	s.resource.MarkAsCommitted()
	assert.Equal(s.T(), s.resource.info, r.info)
	assert.Equal(s.T(), *s.resource, *r)
}

func (s *ResourceTestSuite) TestMapEvents() {

	err := s.resource.ChangeInfo(s.user2, ResourceInfo{
		Value: ResourceValueEstimation{
			ValueType:         FromToDuration,
			ValueFromDuration: 5 * time.Hour,
			ValueToDuration:   22 * time.Hour,
		},
		Name:        "TestMapEvents",
		Description: "TestMapEvents-desc",
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	err = s.resource.ChangeSharings(s.user, *keys.NewGroupKeys([]keys.GroupKey{s.group2}))
	if !assert.NoError(s.T(), err) {
		return
	}

	err = s.resource.Delete(s.user)
	if !assert.NoError(s.T(), err) {
		return
	}

	mapper := eventsource.NewEventMapper()
	err = RegisterEvents(mapper)
	if !assert.NoError(s.T(), err) {
		return
	}

	for i, event := range s.resource.GetChanges() {
		json, err := json2.Marshal(event)
		if !assert.NoError(s.T(), err) {
			return
		}
		evt, err := mapper.Map(event.GetEventType(), json)
		if !assert.NoError(s.T(), err) {
			return
		}
		assert.Equal(s.T(), s.resource.GetChanges()[i], evt)
	}

}

func TestResourceSuite(t *testing.T) {
	suite.Run(t, &ResourceTestSuite{})
}
