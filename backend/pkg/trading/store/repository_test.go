package store

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/pkg/db"
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
	s.db = db.NewTestDb()
	if err := s.db.AutoMigrate(&eventstore.StreamEvent{}, eventstore.Stream{}); err != nil {
		s.T().Fatal(err)
	}
	s.eventStore = postgres.NewPostgresEventStore(s.db)
	s.repository = &EventSourcedOfferRepository{
		eventStore: s.eventStore,
	}
	s.ctx = context.Background()
}

func (s *RepositoryTestSuite) TestSaveOffer() {

	groupKey := keys.NewGroupKey(uuid.NewV4())
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	user1Key := keys.NewUserKey("key1")
	user2Key := keys.NewUserKey("key2")
	user1Target := domain.NewUserTarget(user1Key)
	user2Target := domain.NewUserTarget(user2Key)

	offer := domain.NewOffer()
	assert.NoError(s.T(), offer.Submit(groupKey, domain.OfferItems{
		Items: []domain.OfferItem{
			&domain.CreditTransferItem{
				OfferItemKey: offerItemKey,
				Amount:       time.Hour * 3,
				From:         user1Target,
				To:           user2Target,
			},
		},
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
