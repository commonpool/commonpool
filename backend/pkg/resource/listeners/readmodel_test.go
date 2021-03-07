package listeners

import (
	"context"
	domain3 "github.com/commonpool/backend/pkg/auth/domain"
	db2 "github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/eventstore/postgres"
	domain2 "github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/resource/queries"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	domain4 "github.com/commonpool/backend/pkg/trading/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"strconv"
	"testing"
	"time"
)

type ReadModelTestSuite struct {
	suite.Suite
	eventStore        eventstore.EventStore
	l                 *ResourceReadModelHandler
	resource          *domain.Resource
	resourceKey       keys.ResourceKey
	resourceStreamKey eventstore.StreamKey
	user1StreamKey    eventstore.StreamKey
	user1Key          keys.UserKey
	user2StreamKey    eventstore.StreamKey
	user2Key          keys.UserKey
	group1Key         keys.GroupKey
	group1StreamKey   eventstore.StreamKey
	group2Key         keys.GroupKey
	group2StreamKey   eventstore.StreamKey
	group1Info        domain2.GroupInfo
	group2Info        domain2.GroupInfo
	user1Info         domain3.UserInfo
	user2Info         domain3.UserInfo
	userCounter       int
	groupCounter      int
	ctx               context.Context
	db                *gorm.DB
	getResource       *queries.GetResource
}

func (s *ReadModelTestSuite) SetupSuite() {

	db := db2.NewTestDb()
	s.db = db

	eventMapper := eventsource.NewEventMapper()

	if err := domain.RegisterEvents(eventMapper); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := domain2.RegisterEvents(eventMapper); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}
	if err := domain3.RegisterEvents(eventMapper); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}

	eventStore := postgres.NewPostgresEventStore(db, eventMapper)
	if err := eventStore.MigrateDatabase(); err != nil {
		s.FailNow(err.Error())
	}
	s.eventStore = eventStore

	s.l = &ResourceReadModelHandler{
		db: db,
	}

	if err := s.l.migrateDatabase(); err != nil {
		s.FailNow(err.Error())
	}

	s.getResource = queries.NewGetResource(s.db)

	s.db.Delete(&readmodel.ResourceReadModel{}, "1 = 1")
	s.db.Delete(&readmodel.ResourceSharingReadModel{}, "1 = 1")
	s.db.Delete(&readmodel.ResourceGroupNameReadModel{}, "1 = 1")
	s.db.Delete(&readmodel.ResourceUserNameReadModel{}, "1 = 1")
	s.db.Delete(&eventstore.StreamEvent{}, "1 = 1")
	s.db.Delete(&eventstore.Stream{}, "1 = 1")
}

func (s *ReadModelTestSuite) SetupTest() {
	s.ctx = context.TODO()
	s.user1Key = keys.NewUserKey(uuid.NewV4().String())
	s.user1StreamKey = eventstore.NewStreamKey("user", s.user1Key.String())
	s.user2Key = keys.NewUserKey(uuid.NewV4().String())
	s.user2StreamKey = eventstore.NewStreamKey("user", s.user2Key.String())
	s.group1Key = keys.NewGroupKey(uuid.NewV4())
	s.group1StreamKey = eventstore.NewStreamKey("group", s.group1Key.String())
	s.group2Key = keys.NewGroupKey(uuid.NewV4())
	s.group2StreamKey = eventstore.NewStreamKey("group", s.group2Key.String())
	s.resourceKey = keys.NewResourceKey(uuid.NewV4())
	s.resourceStreamKey = eventstore.NewStreamKey("resource", s.resourceKey.String())
	s.resource = domain.NewResource(s.resourceKey)
	s.user1Info = domain3.UserInfo{
		Email:    "test" + strconv.Itoa(s.userCounter) + "@example.com",
		Username: "user" + strconv.Itoa(s.userCounter),
	}
	s.userCounter++
	s.user1Info = domain3.UserInfo{
		Email:    "test" + strconv.Itoa(s.userCounter) + "@example.com",
		Username: "user" + strconv.Itoa(s.userCounter),
	}
	s.userCounter++
	s.group1Info = domain2.GroupInfo{
		Name:        "group" + strconv.Itoa(s.groupCounter),
		Description: "group" + strconv.Itoa(s.groupCounter),
	}
	s.groupCounter++
	s.group2Info = domain2.GroupInfo{
		Name:        "group" + strconv.Itoa(s.groupCounter),
		Description: "group" + strconv.Itoa(s.groupCounter),
	}
	s.groupCounter++

}

func (s *ReadModelTestSuite) TestShouldCreateUserWhenUserRegistered() {
	user := domain3.New(s.user1Key)
	if err := user.DiscoverUser(s.user1Info); !assert.NoError(s.T(), err) {
		return
	}
	evts, err := s.eventStore.Save(s.ctx, s.user1StreamKey, user.GetVersion(), user.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return
	}
	var userRm readmodel.ResourceUserNameReadModel
	if err := s.db.Model(&readmodel.ResourceUserNameReadModel{}).Find(&userRm, "user_key = ?", user.GetKey().String()).Error; !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), readmodel.ResourceUserNameReadModel{
		UserKey:  s.user1Key.String(),
		Username: s.user1Info.Username,
		Version:  0,
	}, userRm)
}

func (s *ReadModelTestSuite) TestShouldUpdateUserWhenUserInfoChange() {
	user := domain3.New(s.user1Key)
	if err := user.DiscoverUser(s.user1Info); !assert.NoError(s.T(), err) {
		return
	}
	if err := user.ChangeUserInfo(domain3.UserInfo{
		Email:    "TestShouldUpdateUserWhenUserInfoChange@example.com",
		Username: "TestShouldUpdateUserWhenUserInfoChange",
	}); !assert.NoError(s.T(), err) {
		return
	}
	evts, err := s.eventStore.Save(s.ctx, s.user1StreamKey, user.GetVersion(), user.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return
	}
	var userRm readmodel.ResourceUserNameReadModel
	if err := s.db.Model(&readmodel.ResourceUserNameReadModel{}).Find(&userRm, "user_key = ?", user.GetKey().String()).Error; !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), readmodel.ResourceUserNameReadModel{
		UserKey:  s.user1Key.String(),
		Username: "TestShouldUpdateUserWhenUserInfoChange",
		Version:  1,
	}, userRm)
}

func (s *ReadModelTestSuite) TestShouldCreateGroupWhenGroupRegistered() {
	group := domain2.NewGroup(s.group1Key)
	if err := group.CreateGroup(s.user1Key, domain2.GroupInfo{
		Name:        "TestShouldCreateGroupWhenGroupRegistered",
		Description: "TestShouldCreateGroupWhenGroupRegistered-description",
	}); !assert.NoError(s.T(), err) {
		return
	}
	evts, err := s.eventStore.Save(s.ctx, s.group1StreamKey, group.GetVersion(), group.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return
	}
	var groupRm readmodel.ResourceGroupNameReadModel
	if err := s.db.Model(&readmodel.ResourceGroupNameReadModel{}).Find(&groupRm).Error; !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), readmodel.ResourceGroupNameReadModel{
		GroupKey:  s.group1Key.String(),
		GroupName: "TestShouldCreateGroupWhenGroupRegistered",
		Version:   0,
	}, groupRm)
}

func (s *ReadModelTestSuite) TestShouldUpdateGroupWhenGroupInfoChanged() {
	group := domain2.NewGroup(s.group1Key)
	if err := group.CreateGroup(s.user1Key, domain2.GroupInfo{
		Name:        "TestShouldUpdateGroupWhenGroupInfoChanged",
		Description: "TestShouldUpdateGroupWhenGroupInfoChanged-description",
	}); !assert.NoError(s.T(), err) {
		return
	}
	if err := group.ChangeInfo(s.user1Key, domain2.GroupInfo{
		Name:        "TestShouldUpdateGroupWhenGroupInfoChanged-2",
		Description: "TestShouldUpdateGroupWhenGroupInfoChanged-2-description",
	}); !assert.NoError(s.T(), err) {
		return
	}
	evts, err := s.eventStore.Save(s.ctx, s.group1StreamKey, group.GetVersion(), group.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return
	}
	var groupRm readmodel.ResourceGroupNameReadModel
	if err := s.db.Model(&readmodel.ResourceGroupNameReadModel{}).Find(&groupRm, "group_key = ?", s.group1Key.String()).Error; !assert.NoError(s.T(), err) {
		return
	}
	assert.Equal(s.T(), readmodel.ResourceGroupNameReadModel{
		GroupKey:  s.group1Key.String(),
		GroupName: "TestShouldUpdateGroupWhenGroupInfoChanged-2",
		Version:   2,
	}, groupRm)
}

func (s *ReadModelTestSuite) createUser1() error {
	user := domain3.New(s.user1Key)
	if err := user.DiscoverUser(s.user1Info); !assert.NoError(s.T(), err) {
		return err
	}
	evts, err := s.eventStore.Save(s.ctx, s.user1StreamKey, user.GetVersion(), user.GetChanges())
	if !assert.NoError(s.T(), err) {
		return err
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return err
	}
	return nil
}

func (s *ReadModelTestSuite) createUser2() error {
	user := domain3.New(s.user1Key)
	if err := user.DiscoverUser(s.user2Info); !assert.NoError(s.T(), err) {
		return err
	}
	evts, err := s.eventStore.Save(s.ctx, s.user2StreamKey, user.GetVersion(), user.GetChanges())
	if !assert.NoError(s.T(), err) {
		return err
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return err
	}
	return nil
}

func (s *ReadModelTestSuite) createGroup1() error {
	group := domain2.NewGroup(s.group1Key)
	if err := group.CreateGroup(s.user1Key, domain2.GroupInfo{
		Name:        s.group1Info.Name,
		Description: s.group1Info.Name,
	}); !assert.NoError(s.T(), err) {
		return err
	}
	evts, err := s.eventStore.Save(s.ctx, s.group1StreamKey, group.GetVersion(), group.GetChanges())
	if !assert.NoError(s.T(), err) {
		return err
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return err
	}
	return nil
}

func (s *ReadModelTestSuite) createGroup2() error {
	group := domain2.NewGroup(s.group2Key)
	if err := group.CreateGroup(s.user1Key, domain2.GroupInfo{
		Name:        s.group2Info.Name,
		Description: s.group2Info.Name,
	}); !assert.NoError(s.T(), err) {
		return err
	}
	evts, err := s.eventStore.Save(s.ctx, s.group2StreamKey, group.GetVersion(), group.GetChanges())
	if !assert.NoError(s.T(), err) {
		return err
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return err
	}
	return nil
}

func (s *ReadModelTestSuite) TestShouldCreateResourceWhenResourceCreated() {

	if !assert.NoError(s.T(), s.createUser1()) {
		return
	}

	if !assert.NoError(s.T(), s.createGroup1()) {
		return
	}

	if err := s.resource.Register(s.user1Key, *domain4.NewUserTarget(s.user1Key), domain.ServiceResource, domain.ResourceInfo{
		Value: domain.ResourceValueEstimation{
			ValueType:         domain.FromToDuration,
			ValueFromDuration: 2 * time.Hour,
			ValueToDuration:   3 * time.Hour,
		},
		Name:        "TestShouldCreateResourceWhenResourceCreated",
		Description: "TestShouldCreateResourceWhenResourceCreated-description",
	}, *keys.NewGroupKeys([]keys.GroupKey{s.group1Key})); !assert.NoError(s.T(), err) {
		return
	}

	evts, err := s.eventStore.Save(s.ctx, s.resourceStreamKey, s.resource.GetVersion(), s.resource.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return
	}

	rm, err := s.getResource.Get(s.ctx, s.resourceKey)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), readmodel.ResourceReadModel{
		ResourceKey:      s.resourceKey.String(),
		ResourceName:     "TestShouldCreateResourceWhenResourceCreated",
		Description:      "TestShouldCreateResourceWhenResourceCreated-description",
		CreatedBy:        s.user1Key.String(),
		CreatedByVersion: 0,
		CreatedByName:    s.user1Info.Username,
		CreatedAt:        evts[0].GetEventTime(),
		UpdatedBy:        s.user1Key.String(),
		UpdatedByVersion: 0,
		UpdatedByName:    s.user1Info.Username,
		UpdatedAt:        evts[1].GetEventTime(),
		ResourceValueEstimation: domain.ResourceValueEstimation{
			ValueType:         domain.FromToDuration,
			ValueFromDuration: 2 * time.Hour,
			ValueToDuration:   3 * time.Hour,
		},
		GroupSharingCount: 1,
		Version:           1,
	}, *rm)

	var sharings []*readmodel.ResourceSharingReadModel
	if err := s.db.Find(&sharings, "resource_key = ?", s.resourceKey.String()).Error; !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Len(s.T(), sharings, 1) {
		return
	}
	assert.Equal(s.T(), readmodel.ResourceSharingReadModel{
		ResourceKey:  s.resourceKey.String(),
		GroupKey:     s.group1Key.String(),
		GroupName:    s.group1Info.Name,
		Version:      1,
		GroupVersion: 0,
	}, *sharings[0])

}

func (s *ReadModelTestSuite) TestShouldUpdateResourceWhenResourceInfoChanged() {

	if !assert.NoError(s.T(), s.createUser1()) {
		return
	}

	if !assert.NoError(s.T(), s.createGroup1()) {
		return
	}

	if err := s.resource.Register(s.user1Key, *domain4.NewUserTarget(s.user1Key), domain.ServiceResource, domain.ResourceInfo{
		Value: domain.ResourceValueEstimation{
			ValueType:         domain.FromToDuration,
			ValueFromDuration: 2 * time.Hour,
			ValueToDuration:   3 * time.Hour,
		},
		Name:        "TestShouldUpdateResourceWhenResourceInfoChanged",
		Description: "TestShouldUpdateResourceWhenResourceInfoChanged-description",
	}, *keys.NewGroupKeys([]keys.GroupKey{s.group1Key})); !assert.NoError(s.T(), err) {
		return
	}

	if err := s.resource.ChangeInfo(s.user2Key, domain.ResourceInfo{
		Value: domain.ResourceValueEstimation{
			ValueType:         domain.FromToDuration,
			ValueFromDuration: 2 * time.Hour,
			ValueToDuration:   3 * time.Hour,
		},
		Name:        "TestShouldUpdateResourceWhenResourceInfoChanged-2",
		Description: "TestShouldUpdateResourceWhenResourceInfoChanged-2-description",
	}); !assert.NoError(s.T(), err) {
		return
	}

	evts, err := s.eventStore.Save(s.ctx, s.resourceStreamKey, s.resource.GetVersion(), s.resource.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return
	}

	rm, err := s.getResource.Get(s.ctx, s.resourceKey)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), readmodel.ResourceReadModel{
		ResourceKey:      s.resourceKey.String(),
		ResourceName:     "TestShouldUpdateResourceWhenResourceInfoChanged-2",
		Description:      "TestShouldUpdateResourceWhenResourceInfoChanged-2-description",
		CreatedBy:        s.user1Key.String(),
		CreatedByVersion: 0,
		CreatedByName:    s.user1Info.Username,
		CreatedAt:        evts[0].GetEventTime(),
		UpdatedBy:        s.user2Key.String(),
		UpdatedByVersion: 0,
		UpdatedByName:    s.user1Info.Username,
		UpdatedAt:        evts[2].GetEventTime(),
		ResourceValueEstimation: domain.ResourceValueEstimation{
			ValueType:         domain.FromToDuration,
			ValueFromDuration: 2 * time.Hour,
			ValueToDuration:   3 * time.Hour,
		},
		GroupSharingCount: 1,
		Version:           1,
	}, *rm)

}

func (s *ReadModelTestSuite) TestShouldDeleteReadModelWhenResourceDeleted() {

	if err := s.resource.Register(s.user1Key, *domain4.NewUserTarget(s.user1Key), domain.ServiceResource, domain.ResourceInfo{
		Value: domain.ResourceValueEstimation{
			ValueType:         domain.FromToDuration,
			ValueFromDuration: 2 * time.Hour,
			ValueToDuration:   3 * time.Hour,
		},
		Name:        "TestShouldDeleteReadModelWhenResourceDeleted",
		Description: "TestShouldDeleteReadModelWhenResourceDeleted-description",
	}, *keys.NewGroupKeys([]keys.GroupKey{s.group1Key})); !assert.NoError(s.T(), err) {
		return
	}

	if err := s.resource.Delete(s.user1Key); !assert.NoError(s.T(), err) {
		return
	}

	evts, err := s.eventStore.Save(s.ctx, s.resourceStreamKey, s.resource.GetVersion(), s.resource.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return
	}

	_, err = s.getResource.Get(s.ctx, s.resourceKey)
	if !assert.Error(s.T(), err) {
		return
	}
}

func (s *ReadModelTestSuite) TestShouldUpdateSharings() {

	if !assert.NoError(s.T(), s.createUser1()) {
		return
	}

	if !assert.NoError(s.T(), s.createGroup1()) {
		return
	}

	if !assert.NoError(s.T(), s.createGroup2()) {
		return
	}

	if err := s.resource.Register(s.user1Key, *domain4.NewUserTarget(s.user1Key), domain.ServiceResource, domain.ResourceInfo{
		Value: domain.ResourceValueEstimation{
			ValueType:         domain.FromToDuration,
			ValueFromDuration: 2 * time.Hour,
			ValueToDuration:   3 * time.Hour,
		},
		Name:        "TestShouldCreateResourceWhenResourceCreated",
		Description: "TestShouldCreateResourceWhenResourceCreated-description",
	}, *keys.NewGroupKeys([]keys.GroupKey{s.group1Key})); !assert.NoError(s.T(), err) {
		return
	}

	if err := s.resource.ChangeSharings(s.user1Key, *keys.NewGroupKeys([]keys.GroupKey{s.group2Key})); !assert.NoError(s.T(), err) {
		return
	}

	evts, err := s.eventStore.Save(s.ctx, s.resourceStreamKey, s.resource.GetVersion(), s.resource.GetChanges())
	if !assert.NoError(s.T(), err) {
		return
	}
	if err := s.l.handleEvents(evts); !assert.NoError(s.T(), err) {
		return
	}

	rm, err := s.getResource.Get(s.ctx, s.resourceKey)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), readmodel.ResourceReadModel{
		ResourceKey:      s.resourceKey.String(),
		ResourceName:     "TestShouldCreateResourceWhenResourceCreated",
		Description:      "TestShouldCreateResourceWhenResourceCreated-description",
		CreatedBy:        s.user1Key.String(),
		CreatedByVersion: 0,
		CreatedByName:    s.user1Info.Username,
		CreatedAt:        evts[0].GetEventTime(),
		UpdatedBy:        s.user1Key.String(),
		UpdatedByVersion: 0,
		UpdatedByName:    s.user1Info.Username,
		UpdatedAt:        evts[2].GetEventTime(),
		ResourceValueEstimation: domain.ResourceValueEstimation{
			ValueType:         domain.FromToDuration,
			ValueFromDuration: 2 * time.Hour,
			ValueToDuration:   3 * time.Hour,
		},
		GroupSharingCount: 1,
		Version:           2,
	}, *rm)

	var sharings []*readmodel.ResourceSharingReadModel
	if err := s.db.Find(&sharings, "resource_key = ?", s.resourceKey.String()).Error; !assert.NoError(s.T(), err) {
		return
	}
	if !assert.Len(s.T(), sharings, 1) {
		return
	}
	assert.Equal(s.T(), readmodel.ResourceSharingReadModel{
		ResourceKey:  s.resourceKey.String(),
		GroupKey:     s.group2Key.String(),
		GroupName:    s.group2Info.Name,
		Version:      2,
		GroupVersion: 0,
	}, *sharings[0])

}

func TestReadModelSuite(t *testing.T) {
	suite.Run(t, &ReadModelTestSuite{})
}
