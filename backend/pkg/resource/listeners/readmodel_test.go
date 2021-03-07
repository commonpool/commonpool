package listeners

import (
	"context"
	userdomain "github.com/commonpool/backend/pkg/auth/domain"
	db2 "github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/eventstore/postgres"
	groupdomain "github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/resource/queries"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	tradingdomain "github.com/commonpool/backend/pkg/trading/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"math/rand"
	"testing"
	"time"
)

type ReadModelTestSuite struct {
	suite.Suite
	eventStore  eventstore.EventStore
	l           *ResourceReadModelHandler
	ctx         context.Context
	db          *gorm.DB
	getResource *queries.GetResource
	getSharings *queries.GetResourceSharings
}

func (s *ReadModelTestSuite) SetupSuite() {

	// database
	db := db2.NewTestDb()
	s.db = db

	// Create event mapper
	eventMapper := eventsource.NewEventMapper()

	// Register events handled by this view model
	if err := domain.RegisterEvents(eventMapper); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := groupdomain.RegisterEvents(eventMapper); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := userdomain.RegisterEvents(eventMapper); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}

	// create event store
	eventStore := postgres.NewPostgresEventStore(db, eventMapper)
	s.eventStore = eventStore

	// migrate eventstore db
	if err := eventStore.MigrateDatabase(); err != nil {
		s.FailNow(err.Error())
	}

	// create read model handler
	s.l = &ResourceReadModelHandler{
		db: db,
	}

	// migrate db
	if err := s.l.migrateDatabase(); err != nil {
		s.FailNow(err.Error())
	}

	// queries
	s.getResource = queries.NewGetResource(s.db)
	s.getSharings = queries.NewGetResourceSharings(s.db)

	// clean the database
	s.db.Delete(&readmodel.ResourceReadModel{}, "1 = 1")
	s.db.Delete(&readmodel.ResourceSharingReadModel{}, "1 = 1")
	s.db.Delete(&readmodel.ResourceGroupNameReadModel{}, "1 = 1")
	s.db.Delete(&readmodel.ResourceUserNameReadModel{}, "1 = 1")
	s.db.Delete(&eventstore.StreamEvent{}, "1 = 1")
	s.db.Delete(&eventstore.Stream{}, "1 = 1")
}

func (s *ReadModelTestSuite) SetupTest() {
	s.ctx = context.TODO()
}

func (s *ReadModelTestSuite) TestShouldCreateUserWhenUserRegistered() {
	userKey := s.aUserKey()
	streamKey := s.aUserStreamKey(userKey)
	userInfo := s.aUserInfo("TestShouldCreateUserWhenUserRegistered")

	user := userdomain.NewUser(userKey)
	if err := user.DiscoverUser(userInfo); !assert.NoError(s.T(), err) {
		return
	}
	_, err := s.saveAndApplyEvents(streamKey, user.GetVersion(), user.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}
	userRm, err := s.findUserNameReadModel(user.GetKey())
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), readmodel.ResourceUserNameReadModel{
		UserKey:  userKey.String(),
		Username: userInfo.Username,
		Version:  0,
	}, *userRm)
}

func (s *ReadModelTestSuite) TestShouldUpdateUserWhenUserInfoChange() {

	userKey := s.aUserKey()
	streamKey := s.aUserStreamKey(userKey)
	userInfo1 := s.aUserInfo("TestShouldUpdateUserWhenUserInfoChange")
	userInfo2 := s.aUserInfo("TestShouldUpdateUserWhenUserInfoChange-2")

	user := userdomain.NewUser(userKey)
	if err := user.DiscoverUser(userInfo1); !assert.NoError(s.T(), err) {
		return
	}

	if err := user.ChangeUserInfo(userInfo2); !assert.NoError(s.T(), err) {
		return
	}

	_, err := s.saveAndApplyEvents(streamKey, user.GetVersion(), user.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}

	userRm, err := s.findUserNameReadModel(user.GetKey())
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), readmodel.ResourceUserNameReadModel{
		UserKey:  userKey.String(),
		Username: userInfo2.Username,
		Version:  1,
	}, *userRm)
}

func (s *ReadModelTestSuite) TestShouldCreateGroupWhenGroupRegistered() {

	groupKey := s.aGroupKey()
	ownerKey := s.aUserKey()
	streamKey := s.aGroupStreamKey(groupKey)
	groupInfo := s.aGroupInfo("TestShouldCreateGroupWhenGroupRegistered")

	group := groupdomain.NewGroup(groupKey)
	if err := group.CreateGroup(ownerKey, groupInfo); !assert.NoError(s.T(), err) {
		return
	}
	_, err := s.saveAndApplyEvents(streamKey, group.GetVersion(), group.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}

	groupRm, err := s.findGroupNameReadModel(groupKey)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), readmodel.ResourceGroupNameReadModel{
		GroupKey:  groupKey.String(),
		GroupName: groupInfo.Name,
		Version:   0,
	}, *groupRm)

}

func (s *ReadModelTestSuite) TestShouldUpdateGroupWhenGroupInfoChanged() {

	ownerKey := s.aUserKey()

	groupKey := s.aGroupKey()
	groupStreamKey := s.aGroupStreamKey(groupKey)
	group := groupdomain.NewGroup(groupKey)
	groupInfo1 := s.aGroupInfo("TestShouldUpdateGroupWhenGroupInfoChanged")
	groupInfo2 := s.aGroupInfo("TestShouldUpdateGroupWhenGroupInfoChanged-2")

	if err := group.CreateGroup(ownerKey, groupInfo1); !assert.NoError(s.T(), err) {
		return
	}

	if err := group.ChangeInfo(ownerKey, groupInfo2); !assert.NoError(s.T(), err) {
		return
	}

	_, err := s.saveAndApplyEvents(groupStreamKey, group.GetVersion(), group.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}

	groupRm, err := s.findGroupNameReadModel(groupKey)
	if !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), readmodel.ResourceGroupNameReadModel{
		GroupKey:  groupKey.String(),
		GroupName: groupInfo2.Name,
		Version:   2,
	}, *groupRm)
}

func (s *ReadModelTestSuite) TestShouldCreateResourceWhenResourceCreated() {

	userKey, userInfo, err := s.createUser("TestShouldCreateResourceWhenResourceCreated-user1")
	if !assert.NoError(s.T(), err) {
		return
	}

	groupKey, groupInfo, err := s.createGroup("TestShouldCreateResourceWhenResourceCreated-group")
	if !assert.NoError(s.T(), err) {
		return
	}

	resourceKey := s.aResourceKey()
	streamKey := s.aResourceStreamKey(resourceKey)
	resource := domain.NewResource(resourceKey)
	resourceInfo := s.aResourceInfo("TestShouldCreateResourceWhenResourceCreated")

	if err := resource.Register(
		userKey,
		*tradingdomain.NewUserTarget(userKey),
		resourceInfo,
		*keys.NewGroupKeys([]keys.GroupKey{groupKey})); !assert.NoError(s.T(), err) {
		return
	}

	evts, err := s.saveAndApplyEvents(streamKey, resource.GetVersion(), resource.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}

	rm, err := s.getResource.Get(s.ctx, resourceKey)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), readmodel.ResourceReadModel{
		ResourceKey:             resourceKey.String(),
		ResourceName:            resourceInfo.Name,
		Description:             resourceInfo.Description,
		CreatedBy:               userKey.String(),
		CreatedByVersion:        0,
		CreatedByName:           userInfo.Username,
		CreatedAt:               evts[0].GetEventTime(),
		UpdatedBy:               userKey.String(),
		UpdatedByVersion:        0,
		UpdatedByName:           userInfo.Username,
		UpdatedAt:               evts[1].GetEventTime(),
		ResourceValueEstimation: resourceInfo.Value,
		GroupSharingCount:       1,
		Version:                 1,
	}, *rm)

	sharings, err := s.getSharings.Get(s.ctx, resourceKey)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Len(s.T(), sharings, 1) {
		return
	}

	assert.Equal(s.T(), readmodel.ResourceSharingReadModel{
		ResourceKey:  resourceKey.String(),
		GroupKey:     groupKey.String(),
		GroupName:    groupInfo.Name,
		Version:      1,
		GroupVersion: 0,
	}, *sharings[0])

}

func (s *ReadModelTestSuite) TestShouldUpdateResourceWhenResourceInfoChanged() {

	userKey1, userInfo1, err := s.createUser("TestShouldCreateResourceWhenResourceCreated-user1")
	if !assert.NoError(s.T(), err) {
		return
	}

	userKey2, userInfo2, err := s.createUser("TestShouldCreateResourceWhenResourceCreated-user1")
	if !assert.NoError(s.T(), err) {
		return
	}

	groupKey, _, err := s.createGroup("TestShouldCreateResourceWhenResourceCreated-group")
	if !assert.NoError(s.T(), err) {
		return
	}

	resourceKey := s.aResourceKey()
	streamKey := s.aResourceStreamKey(resourceKey)
	resource := domain.NewResource(resourceKey)
	resourceInfo := s.aResourceInfo("TestShouldCreateResourceWhenResourceCreated")

	if err := resource.Register(
		userKey1, *tradingdomain.NewUserTarget(userKey1),
		resourceInfo,
		*keys.NewGroupKeys([]keys.GroupKey{groupKey})); !assert.NoError(s.T(), err) {
		return
	}

	resourceInfo2 := s.aResourceInfo("TestShouldUpdateResourceWhenResourceInfoChanged2")
	if err := resource.ChangeInfo(userKey2, resourceInfo2); !assert.NoError(s.T(), err) {
		return
	}

	evts, err := s.saveAndApplyEvents(streamKey, resource.GetVersion(), resource.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}

	rm, err := s.getResource.Get(s.ctx, resourceKey)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), readmodel.ResourceReadModel{
		ResourceKey:             resourceKey.String(),
		ResourceName:            resourceInfo2.Name,
		Description:             resourceInfo2.Description,
		CreatedBy:               userKey1.String(),
		CreatedByVersion:        0,
		CreatedByName:           userInfo1.Username,
		CreatedAt:               evts[0].GetEventTime(),
		UpdatedBy:               userKey2.String(),
		UpdatedByVersion:        0,
		UpdatedByName:           userInfo2.Username,
		UpdatedAt:               evts[2].GetEventTime(),
		ResourceValueEstimation: resourceInfo2.Value,
		GroupSharingCount:       1,
		Version:                 1,
	}, *rm)

}

func (s *ReadModelTestSuite) TestShouldDeleteReadModelWhenResourceDeleted() {

	user1Key := s.aUserKey()
	groupKey := s.aGroupKey()

	resourceKey := s.aResourceKey()
	streamKey := s.aResourceStreamKey(resourceKey)
	resource := domain.NewResource(resourceKey)
	resourceInfo := s.aResourceInfo("TestShouldCreateResourceWhenResourceCreated")

	if err := resource.Register(
		user1Key,
		*tradingdomain.NewUserTarget(user1Key),
		resourceInfo,
		*keys.NewGroupKeys([]keys.GroupKey{groupKey})); !assert.NoError(s.T(), err) {
		return
	}

	if err := resource.Delete(user1Key); !assert.NoError(s.T(), err) {
		return
	}

	_, err := s.saveAndApplyEvents(streamKey, resource.GetVersion(), resource.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}

	_, err = s.getResource.Get(s.ctx, resourceKey)
	if !assert.Error(s.T(), err) {
		return
	}
}

func (s *ReadModelTestSuite) TestShouldUpdateSharings() {

	userKey, userInfo, err := s.createUser("TestShouldCreateResourceWhenResourceCreated-user1")
	if !assert.NoError(s.T(), err) {
		return
	}

	groupKey1, _, err := s.createGroup("TestShouldCreateResourceWhenResourceCreated-group")
	if !assert.NoError(s.T(), err) {
		return
	}

	groupKey2, groupInfo2, err := s.createGroup("TestShouldCreateResourceWhenResourceCreated-group")
	if !assert.NoError(s.T(), err) {
		return
	}

	resourceKey := s.aResourceKey()
	streamKey := s.aResourceStreamKey(resourceKey)
	resource := domain.NewResource(resourceKey)
	resourceInfo := s.aResourceInfo("TestShouldDeleteReadModelWhenResourceDeleted")

	if err := resource.Register(
		userKey,
		*tradingdomain.NewUserTarget(userKey),
		resourceInfo,
		*keys.NewGroupKeys([]keys.GroupKey{groupKey1})); !assert.NoError(s.T(), err) {
		return
	}

	if err := resource.ChangeSharings(userKey, *keys.NewGroupKeys([]keys.GroupKey{groupKey2})); !assert.NoError(s.T(), err) {
		return
	}

	evts, err := s.saveAndApplyEvents(streamKey, resource.GetVersion(), resource.GetChanges())

	rm, err := s.getResource.Get(s.ctx, resourceKey)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), readmodel.ResourceReadModel{
		ResourceKey:             resourceKey.String(),
		ResourceName:            resourceInfo.Name,
		Description:             resourceInfo.Description,
		CreatedBy:               userKey.String(),
		CreatedByVersion:        0,
		CreatedByName:           userInfo.Username,
		CreatedAt:               evts[0].GetEventTime(),
		UpdatedBy:               userKey.String(),
		UpdatedByVersion:        0,
		UpdatedByName:           userInfo.Username,
		UpdatedAt:               evts[2].GetEventTime(),
		ResourceValueEstimation: resourceInfo.Value,
		GroupSharingCount:       1,
		Version:                 2,
	}, *rm)

	sharings, err := s.getSharings.Get(s.ctx, resourceKey)
	if !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Len(s.T(), sharings, 1) {
		return
	}
	assert.Equal(s.T(), readmodel.ResourceSharingReadModel{
		ResourceKey:  resourceKey.String(),
		GroupKey:     groupKey2.String(),
		GroupName:    groupInfo2.Name,
		Version:      2,
		GroupVersion: 0,
	}, *sharings[0])

}

func (s *ReadModelTestSuite) aGroupKey() keys.GroupKey {
	return keys.NewGroupKey(uuid.NewV4())
}

func (s *ReadModelTestSuite) aUserKey() keys.UserKey {
	return keys.NewUserKey(uuid.NewV4().String())
}

func (s *ReadModelTestSuite) aUserStreamKey(userKey keys.UserKey) keys.StreamKey {
	return keys.NewStreamKey("user", userKey.String())
}

func (s *ReadModelTestSuite) aResourceStreamKey(resourceKey keys.ResourceKey) keys.StreamKey {
	return keys.NewStreamKey("resource", resourceKey.String())
}

func (s *ReadModelTestSuite) aGroupStreamKey(groupKey keys.GroupKey) keys.StreamKey {
	return keys.NewStreamKey("group", groupKey.String())
}

func (s *ReadModelTestSuite) aResourceKey() keys.ResourceKey {
	return keys.NewResourceKey(uuid.NewV4())
}

func (s *ReadModelTestSuite) aUserInfo(prefix string) userdomain.UserInfo {
	return userdomain.UserInfo{
		Email:    prefix + "@example.com",
		Username: prefix + "-user",
	}
}

func (s *ReadModelTestSuite) aGroupInfo(prefix string) groupdomain.GroupInfo {
	return groupdomain.GroupInfo{
		Name:        prefix,
		Description: prefix + "-description",
	}
}

func (s *ReadModelTestSuite) aResourceInfo(prefix string) domain.ResourceInfo {
	valueEstimation := domain.ResourceValueEstimation{
		ValueType:         domain.FromToDuration,
		ValueFromDuration: time.Duration(rand.Intn(10)) * time.Hour,
		ValueToDuration:   time.Duration(rand.Intn(10)+10) * time.Hour,
	}
	return domain.ResourceInfo{
		Value:        valueEstimation,
		Name:         prefix,
		Description:  prefix + "-description",
		CallType:     domain.Offer,
		ResourceType: domain.ObjectResource,
	}
}

func (s *ReadModelTestSuite) saveAndApplyEvents(streamKey keys.StreamKey, expectedRevision int, changes []eventsource.Event) ([]eventsource.Event, error) {
	evts, err := s.eventStore.Save(s.ctx, streamKey, expectedRevision, changes)
	if err != nil {
		return nil, err
	}
	if err := s.l.handleEvents(evts); err != nil {
		return nil, err
	}
	return evts, err
}

func (s *ReadModelTestSuite) findUserNameReadModel(userKey keys.UserKey) (*readmodel.ResourceUserNameReadModel, error) {
	var userRm readmodel.ResourceUserNameReadModel
	if err := s.db.Model(&readmodel.ResourceUserNameReadModel{}).Find(&userRm, "user_key = ?", userKey.String()).Error; !assert.NoError(s.T(), err) {
		return nil, err
	}
	return &userRm, nil
}

func (s *ReadModelTestSuite) findGroupNameReadModel(groupKey keys.GroupKey) (*readmodel.ResourceGroupNameReadModel, error) {
	var groupRm readmodel.ResourceGroupNameReadModel
	if err := s.db.Model(&readmodel.ResourceGroupNameReadModel{}).Find(&groupRm, "group_key = ?", groupKey.String()).Error; !assert.NoError(s.T(), err) {
		return nil, err
	}
	return &groupRm, nil
}

func (s *ReadModelTestSuite) createUser(prefix string) (keys.UserKey, userdomain.UserInfo, error) {
	userKey := s.aUserKey()
	userInfo := s.aUserInfo(prefix)
	user := userdomain.NewUser(userKey)
	streamKey := s.aUserStreamKey(userKey)

	if err := user.DiscoverUser(userInfo); !assert.NoError(s.T(), err) {
		return keys.UserKey{}, userdomain.UserInfo{}, err
	}
	evts, err := s.eventStore.Save(s.ctx, streamKey, user.GetVersion(), user.GetChanges())
	if !assert.NoError(s.T(), err) {
		return keys.UserKey{}, userdomain.UserInfo{}, err
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return keys.UserKey{}, userdomain.UserInfo{}, err
	}
	return userKey, userInfo, nil
}

func (s *ReadModelTestSuite) createGroup(prefix string) (keys.GroupKey, groupdomain.GroupInfo, error) {
	groupKey := s.aGroupKey()
	groupStreamKey := s.aGroupStreamKey(groupKey)
	groupInfo := s.aGroupInfo(prefix)
	group := groupdomain.NewGroup(groupKey)

	if err := group.CreateGroup(s.aUserKey(), groupdomain.GroupInfo{
		Name:        groupInfo.Name,
		Description: groupInfo.Name,
	}); !assert.NoError(s.T(), err) {
		return keys.GroupKey{}, groupdomain.GroupInfo{}, err
	}
	evts, err := s.eventStore.Save(s.ctx, groupStreamKey, group.GetVersion(), group.GetChanges())
	if !assert.NoError(s.T(), err) {
		return keys.GroupKey{}, groupdomain.GroupInfo{}, err
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return keys.GroupKey{}, groupdomain.GroupInfo{}, err
	}
	return groupKey, groupInfo, nil
}

func TestReadModelSuite(t *testing.T) {
	suite.Run(t, &ReadModelTestSuite{})
}
