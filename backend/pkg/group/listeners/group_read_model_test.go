package listeners

import (
	"context"
	"github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/eventstore/postgres"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/group/queries"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/group/store"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type GroupReadModelTestSuite struct {
	suite.Suite
	db            *gorm.DB
	l             *GroupReadModelListener
	es            *postgres.PostgresEventStore
	repo          domain.GroupRepository
	getGroup      *queries.GetGroup
	getMembership *queries.GetMembershipReadModel
}

func (s *GroupReadModelTestSuite) SetupSuite() {
	s.db = db.NewTestDb("GroupReadModelTestSuite")
	s.l = &GroupReadModelListener{
		db: s.db,
	}
	if err := s.l.migrateDatabase(); err != nil {
		s.FailNow(err.Error())
	}
	s.getGroup = queries.NewGetGroupReadModel(s.db)
	s.getMembership = queries.NewGetMembership(s.db)

	eventMapper := eventsource.NewEventMapper()
	if err := domain.RegisterEvents(eventMapper); !assert.NoError(s.T(), err) {
		return
	}

	s.es = postgres.NewPostgresEventStore(s.db, eventMapper)

	if err := s.es.MigrateDatabase(); !assert.NoError(s.T(), err) {
		s.FailNow(err.Error())
	}

	s.db.Delete(&readmodels.GroupReadModel{}, "1 = 1")
	s.db.Delete(&readmodels.MembershipReadModel{}, "1 = 1")

	s.repo = store.NewEventSourcedGroupRepository(s.es)
}

func (s *GroupReadModelTestSuite) TestGroupReadModelCreateGroup() {

	gk := keys.GenerateGroupKey()
	g := domain.NewGroup(gk)
	owner := keys.NewUserKey("TestGroupReadModelCreateGroup-Owner")
	ownerMembershipKey := keys.NewMembershipKey(g.GetKey(), owner)

	err := g.CreateGroup(owner, domain.GroupInfo{
		Name:        "TestGroupReadModelCreateGroup-name",
		Description: "TestGroupReadModelCreateGroup-description",
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	err = s.repo.Save(context.TODO(), g)
	if !assert.NoError(s.T(), err) {
		return
	}

	events, err := s.es.Load(context.TODO(), keys.NewStreamKey("group", g.GetKey().String()))
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.Len(s.T(), events, 2) {
		return
	}

	if !assert.NoError(s.T(), s.l.applyEvent(events[0])) {
		return
	}

	group, err := s.getGroup.Get(context.TODO(), g.GetKey())
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), "TestGroupReadModelCreateGroup-name", group.Name)
	assert.Equal(s.T(), "TestGroupReadModelCreateGroup-description", group.Description)
	assert.Equal(s.T(), 0, group.Version)
	assert.Equal(s.T(), "TestGroupReadModelCreateGroup-Owner", group.CreatedBy)

	if !assert.NoError(s.T(), s.l.applyEvent(events[1])) {
		return
	}

	membership, err := s.getMembership.Get(context.TODO(), ownerMembershipKey)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), 1, membership.Version)
	assert.Equal(s.T(), owner, membership.UserKey)
	assert.Equal(s.T(), g.GetKey(), membership.GroupKey)
	assert.Equal(s.T(), true, membership.IsMember)
	assert.Equal(s.T(), true, membership.IsAdmin)
	assert.Equal(s.T(), true, membership.IsOwner)
	assert.Equal(s.T(), true, membership.UserConfirmed)
	assert.NotEqual(s.T(), time.Time{}, *membership.UserConfirmedAt)
	assert.Equal(s.T(), true, membership.GroupConfirmed)
	assert.Equal(s.T(), owner.String(), *membership.GroupConfirmedBy)
	assert.NotEqual(s.T(), time.Time{}, *membership.GroupConfirmedAt)
	assert.Equal(s.T(), "TestGroupReadModelCreateGroup-name", membership.GroupName)

}

func (s *GroupReadModelTestSuite) TestChangeGroupInfo() {

	gk := keys.GenerateGroupKey()
	g := domain.NewGroup(gk)
	owner := keys.NewUserKey("TestGroupReadModelCreateGroup-Owner")
	ownerMembershipKey := keys.NewMembershipKey(g.GetKey(), owner)

	err := g.CreateGroup(owner, domain.GroupInfo{
		Name:        "TestGroupReadModelCreateGroup-name",
		Description: "TestGroupReadModelCreateGroup-description",
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	err = g.ChangeInfo(owner, domain.GroupInfo{
		Name:        "TestGroupReadModelCreateGroup-name-2",
		Description: "TestGroupReadModelCreateGroup-description-2",
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	err = s.repo.Save(context.TODO(), g)
	if !assert.NoError(s.T(), err) {
		return
	}

	events, err := s.es.Load(context.TODO(), keys.NewStreamKey("group", g.GetKey().String()))
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.Len(s.T(), events, 3) {
		return
	}

	if !assert.NoError(s.T(), s.l.applyEvents(events)) {
		return
	}

	group, err := s.getGroup.Get(context.TODO(), g.GetKey())
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), "TestGroupReadModelCreateGroup-name-2", group.Name)
	assert.Equal(s.T(), "TestGroupReadModelCreateGroup-description-2", group.Description)
	assert.Equal(s.T(), 2, group.Version)

	member, err := s.getMembership.Get(context.TODO(), ownerMembershipKey)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), member.GroupName, "TestGroupReadModelCreateGroup-name-2")

}

func (s *GroupReadModelTestSuite) TestInviteUser() {

	gk := keys.GenerateGroupKey()
	g := domain.NewGroup(gk)
	owner := keys.NewUserKey("TestInviteUser-Owner")
	user := keys.NewUserKey("TestInviteUser-User")

	err := g.CreateGroup(owner, domain.GroupInfo{
		Name:        "TestInviteUser",
		Description: "TestInviteUser",
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.NoError(s.T(), g.JoinGroup(owner, user)) {
		return
	}

	err = s.repo.Save(context.TODO(), g)
	if !assert.NoError(s.T(), err) {
		return
	}

	events, err := s.es.Load(context.TODO(), keys.NewStreamKey("group", g.GetKey().String()))
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.NoError(s.T(), s.l.applyEvents(events)) {
		return
	}

	userMembership, err := s.getMembership.Get(context.TODO(), keys.NewMembershipKey(g.GetKey(), user))
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), owner.String(), *userMembership.GroupConfirmedBy)
	assert.NotNil(s.T(), userMembership.GroupConfirmedAt)
	assert.True(s.T(), userMembership.GroupConfirmed)
	assert.False(s.T(), userMembership.UserConfirmed)
	assert.Nil(s.T(), userMembership.UserConfirmedAt)
	assert.False(s.T(), userMembership.IsOwner)
	assert.False(s.T(), userMembership.IsAdmin)
	assert.False(s.T(), userMembership.IsMember)
	assert.Equal(s.T(), g.GetKey(), userMembership.GroupKey)
	assert.Equal(s.T(), user, userMembership.UserKey)
	assert.Equal(s.T(), 2, userMembership.Version)
}

func (s *GroupReadModelTestSuite) TestUserAcceptedInvitation() {

	gk := keys.GenerateGroupKey()
	g := domain.NewGroup(gk)
	owner := keys.NewUserKey("TestUserAcceptedInvitation-Owner")
	user := keys.NewUserKey("TestUserAcceptedInvitation-User")

	err := g.CreateGroup(owner, domain.GroupInfo{
		Name:        "TestUserAcceptedInvitation",
		Description: "TestUserAcceptedInvitation",
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.NoError(s.T(), g.JoinGroup(owner, user)) {
		return
	}

	if !assert.NoError(s.T(), g.JoinGroup(user, user)) {
		return
	}

	err = s.repo.Save(context.TODO(), g)
	if !assert.NoError(s.T(), err) {
		return
	}

	events, err := s.es.Load(context.TODO(), keys.NewStreamKey("group", g.GetKey().String()))
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.NoError(s.T(), s.l.applyEvents(events)) {
		return
	}

	userMembership, err := s.getMembership.Get(context.TODO(), keys.NewMembershipKey(g.GetKey(), user))
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), owner.String(), *userMembership.GroupConfirmedBy)
	assert.NotNil(s.T(), userMembership.GroupConfirmedAt)
	assert.True(s.T(), userMembership.GroupConfirmed)
	assert.True(s.T(), userMembership.UserConfirmed)
	assert.NotNil(s.T(), userMembership.UserConfirmedAt)
	assert.False(s.T(), userMembership.IsOwner)
	assert.False(s.T(), userMembership.IsAdmin)
	assert.True(s.T(), userMembership.IsMember)
	assert.Equal(s.T(), g.GetKey(), userMembership.GroupKey)
	assert.Equal(s.T(), user, userMembership.UserKey)
	assert.Equal(s.T(), 3, userMembership.Version)

}

func (s *GroupReadModelTestSuite) TestJoinGroup() {

	gk := keys.GenerateGroupKey()
	g := domain.NewGroup(gk)
	owner := keys.NewUserKey("TestJoinGroup-Owner")
	user := keys.NewUserKey("TestJoinGroup-User")

	err := g.CreateGroup(owner, domain.GroupInfo{
		Name:        "TestJoinGroup",
		Description: "TestJoinGroup",
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.NoError(s.T(), g.JoinGroup(user, user)) {
		return
	}

	err = s.repo.Save(context.TODO(), g)
	if !assert.NoError(s.T(), err) {
		return
	}

	events, err := s.es.Load(context.TODO(), keys.NewStreamKey("group", g.GetKey().String()))
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.NoError(s.T(), s.l.applyEvents(events)) {
		return
	}

	userMembership, err := s.getMembership.Get(context.TODO(), keys.NewMembershipKey(g.GetKey(), user))
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Nil(s.T(), userMembership.GroupConfirmedBy)
	assert.Nil(s.T(), userMembership.GroupConfirmedAt)
	assert.False(s.T(), userMembership.GroupConfirmed)
	assert.True(s.T(), userMembership.UserConfirmed)
	assert.NotNil(s.T(), userMembership.UserConfirmedAt)
	assert.False(s.T(), userMembership.IsOwner)
	assert.False(s.T(), userMembership.IsAdmin)
	assert.False(s.T(), userMembership.IsMember)
	assert.Equal(s.T(), g.GetKey(), userMembership.GroupKey)
	assert.Equal(s.T(), user, userMembership.UserKey)
	assert.Equal(s.T(), 2, userMembership.Version)

}

func (s *GroupReadModelTestSuite) TestGroupAcceptedInvitation() {

	gk := keys.GenerateGroupKey()
	g := domain.NewGroup(gk)
	owner := keys.NewUserKey("TestGroupAcceptedInvitation-Owner")
	user := keys.NewUserKey("TestGroupAcceptedInvitation-User")

	err := g.CreateGroup(owner, domain.GroupInfo{
		Name:        "TestGroupAcceptedInvitation",
		Description: "TestGroupAcceptedInvitation",
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.NoError(s.T(), g.JoinGroup(user, user)) {
		return
	}

	if !assert.NoError(s.T(), g.JoinGroup(owner, user)) {
		return
	}

	err = s.repo.Save(context.TODO(), g)
	if !assert.NoError(s.T(), err) {
		return
	}

	events, err := s.es.Load(context.TODO(), keys.NewStreamKey("group", g.GetKey().String()))
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.NoError(s.T(), s.l.applyEvents(events)) {
		return
	}

	userMembership, err := s.getMembership.Get(context.TODO(), keys.NewMembershipKey(g.GetKey(), user))
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), owner.String(), *userMembership.GroupConfirmedBy)
	assert.NotNil(s.T(), userMembership.GroupConfirmedAt)
	assert.True(s.T(), userMembership.GroupConfirmed)
	assert.True(s.T(), userMembership.UserConfirmed)
	assert.NotNil(s.T(), userMembership.UserConfirmedAt)
	assert.False(s.T(), userMembership.IsOwner)
	assert.False(s.T(), userMembership.IsAdmin)
	assert.True(s.T(), userMembership.IsMember)
	assert.Equal(s.T(), g.GetKey(), userMembership.GroupKey)
	assert.Equal(s.T(), user, userMembership.UserKey)
	assert.Equal(s.T(), 3, userMembership.Version)

}

func (s *GroupReadModelTestSuite) TestPermissionChanged() {

	gk := keys.GenerateGroupKey()
	g := domain.NewGroup(gk)
	owner := keys.NewUserKey("TestPermissionChanged-Owner")
	user := keys.NewUserKey("TestPermissionChanged-User")

	err := g.CreateGroup(owner, domain.GroupInfo{
		Name:        "TestPermissionChanged",
		Description: "TestPermissionChanged",
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.NoError(s.T(), g.JoinGroup(user, user)) {
		return
	}

	if !assert.NoError(s.T(), g.JoinGroup(owner, user)) {
		return
	}

	if !assert.NoError(s.T(), g.AssignPermission(owner, user, domain.Admin)) {
		return
	}

	err = s.repo.Save(context.TODO(), g)
	if !assert.NoError(s.T(), err) {
		return
	}

	events, err := s.es.Load(context.TODO(), keys.NewStreamKey("group", g.GetKey().String()))
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.NoError(s.T(), s.l.applyEvents(events)) {
		return
	}

	userMembership, err := s.getMembership.Get(context.TODO(), keys.NewMembershipKey(g.GetKey(), user))
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.True(s.T(), userMembership.IsAdmin)
	assert.Equal(s.T(), 4, userMembership.Version)

}

func (s *GroupReadModelTestSuite) TestLeftGroup() {

	gk := keys.GenerateGroupKey()
	g := domain.NewGroup(gk)
	owner := keys.NewUserKey("TestLeftGroup-Owner")
	user := keys.NewUserKey("TestLeftGroup-User")

	err := g.CreateGroup(owner, domain.GroupInfo{
		Name:        "TestLeftGroup",
		Description: "TestLeftGroup",
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.NoError(s.T(), g.JoinGroup(user, user)) {
		return
	}

	if !assert.NoError(s.T(), g.JoinGroup(owner, user)) {
		return
	}

	if !assert.NoError(s.T(), g.CancelMembership(user, user)) {
		return
	}

	err = s.repo.Save(context.TODO(), g)
	if !assert.NoError(s.T(), err) {
		return
	}

	events, err := s.es.Load(context.TODO(), keys.NewStreamKey("group", g.GetKey().String()))
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.NoError(s.T(), s.l.applyEvents(events)) {
		return
	}

	_, err = s.getMembership.Get(context.TODO(), keys.NewMembershipKey(g.GetKey(), user))
	assert.Error(s.T(), err)

}

func TestGroupReadModel(t *testing.T) {
	suite.Run(t, &GroupReadModelTestSuite{})
}
