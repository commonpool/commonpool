package listeners

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/domain"
	"github.com/commonpool/backend/pkg/auth/queries"
	"github.com/commonpool/backend/pkg/auth/readmodel"
	"github.com/commonpool/backend/pkg/auth/store"
	"github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/eventstore/postgres"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
)

type UserReadModelTestSuite struct {
	suite.Suite
	db      *gorm.DB
	es      *postgres.PostgresEventStore
	repo    *store.EventSourcedUserRepository
	getUser *queries.GetUser
	l       *UserReadModelListener
}

func (s *UserReadModelTestSuite) saveAndLoadEvents(user *domain.User, alsoHandleEvents bool) ([]eventsource.Event, error) {

	err := s.repo.Save(context.TODO(), user)
	if !assert.NoError(s.T(), err) {
		return nil, err
	}

	evts, err := s.es.Load(context.TODO(), eventstore.NewStreamKey("user", user.GetKey().String()))
	if !assert.NoError(s.T(), err) {
		return nil, err
	}

	if alsoHandleEvents {
		err = s.l.handleEvents(evts)
		if !assert.NoError(s.T(), err) {
			return nil, err
		}
	}
	return evts, nil
}

func (s *UserReadModelTestSuite) SetupSuite() {
	s.db = db.NewTestDb()
	evtMapper := eventsource.NewEventMapper()

	if err := domain.RegisterEvents(evtMapper); err != nil {
		s.FailNow(err.Error())
	}

	s.es = postgres.NewPostgresEventStore(s.db, evtMapper)
	if err := s.es.MigrateDatabase(); err != nil {
		s.FailNow(err.Error())
	}
	s.repo = store.NewEventSourcedUserRepository(s.es)
	s.l = &UserReadModelListener{
		db: s.db,
	}
	if err := s.l.migrateDatabase(); err != nil {
		s.FailNow(err.Error())
	}
	s.getUser = queries.NewGetUser(s.db)

	if err := s.db.Delete(&readmodel.UserReadModel{}, "1 = 1").Error; err != nil {
		s.FailNow(err.Error())
	}

	if err := s.db.Delete(&eventstore.StreamEvent{}, "1 = 1").Error; err != nil {
		s.FailNow(err.Error())
	}

	if err := s.db.Delete(&eventstore.Stream{}, "1 = 1").Error; err != nil {
		s.FailNow(err.Error())
	}
}

func (s *UserReadModelTestSuite) TestNewUser() {
	user := domain.New(keys.NewUserKey("TestNewUser"))
	err := user.DiscoverUser(domain.UserInfo{
		Email:    "TestNewUser@example.com",
		Username: "TestNewUser-username",
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	_, err = s.saveAndLoadEvents(user, true)
	if !assert.NoError(s.T(), err) {
		return
	}

	rm, err := s.getUser.Get(user.GetKey())
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), "TestNewUser", rm.UserKey)
	assert.Equal(s.T(), "TestNewUser-username", rm.Username)
	assert.Equal(s.T(), "TestNewUser@example.com", rm.Email)
	assert.Equal(s.T(), 0, rm.Version)

}

func (s *UserReadModelTestSuite) TestChangeUserInfo() {
	user := domain.New(keys.NewUserKey("TestChangeUserInfo"))

	if err := user.DiscoverUser(domain.UserInfo{
		Email:    "TestChangeUserInfo@example.com",
		Username: "TestChangeUserInfo-username",
	}); !assert.NoError(s.T(), err) {
		return
	}

	if err := user.ChangeUserInfo(domain.UserInfo{
		Email:    "TestChangeUserInfo@example.com2",
		Username: "TestChangeUserInfo-username2",
	}); !assert.NoError(s.T(), err) {
		return
	}

	_, err := s.saveAndLoadEvents(user, true)
	if !assert.NoError(s.T(), err) {
		return
	}

	rm, err := s.getUser.Get(user.GetKey())
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Equal(s.T(), "TestChangeUserInfo", rm.UserKey)
	assert.Equal(s.T(), "TestChangeUserInfo-username2", rm.Username)
	assert.Equal(s.T(), "TestChangeUserInfo@example.com2", rm.Email)
	assert.Equal(s.T(), 0, rm.Version)

}

func TestUserReadModel(t *testing.T) {
	suite.Run(t, &UserReadModelTestSuite{})
}
