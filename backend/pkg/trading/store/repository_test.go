package store

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/eventstore/postgres"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type RepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	eventStore eventstore.EventStore
	repository *EventSourcedOfferRepository
	ctx        context.Context
}

func (s *RepositoryTestSuite) SetupSuite() {
	s.db = db.NewTestDb("RepositoryTestSuite")
	if err := s.db.AutoMigrate(&eventstore.StreamEvent{}, eventstore.Stream{}); err != nil {
		s.T().Fatal(err)
	}
	eventMapper := eventsource.NewEventMapper()
	if err := domain.RegisterEvents(eventMapper); !assert.NoError(s.T(), err) {
		return
	}
	s.eventStore = postgres.NewPostgresEventStore(s.db, eventMapper)
	s.repository = NewEventSourcedOfferRepository(s.eventStore)
	s.ctx = context.Background()
}

func (s *RepositoryTestSuite) TestSaveOffer() {

	groupKey := keys.NewGroupKey(uuid.NewV4())
	user1Key := keys.NewUserKey("key1")
	user2Key := keys.NewUserKey("key2")
	offerItemKey := keys.GenerateOfferItemKey()

	offer := domain.NewOffer(keys.GenerateOfferKey())
	assert.NoError(s.T(), offer.Submit(user1Key, groupKey, []domain.SubmitOfferItem{
		domain.NewCreditTransferItemInput(offerItemKey, user1Key, user2Key, time.Hour*3),
	}))

	assert.NoError(s.T(), s.repository.Save(s.ctx, offer))

	loaded, err := s.repository.Load(s.ctx, offer.GetKey())
	assert.NoError(s.T(), err)

	offerJs, _ := json.Marshal(offer)
	loadedJs, _ := json.Marshal(loaded)

	assert.Equal(s.T(), string(offerJs), string(loadedJs))

}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
