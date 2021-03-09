package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/eventstore/postgres"
	"github.com/commonpool/backend/pkg/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type ReplayListenerTestSuite struct {
	suite.Suite
	ctx        context.Context
	db         *gorm.DB
	eventStore *postgres.PostgresEventStore
}

func (s *ReplayListenerTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = db.NewTestDb("ReplayListenerTestSuite")
	if err := s.db.AutoMigrate(&eventstore.StreamEvent{}, eventstore.Stream{}); err != nil {
		s.T().Fatal(err)
	}
	eventMapper := eventsource.NewEventMapper()
	if err := test.RegisterMockEvents(eventMapper); !assert.NoError(s.T(), err) {
		return
	}
	s.eventStore = postgres.NewPostgresEventStore(s.db, eventMapper)
}

func (s *ReplayListenerTestSuite) SetupTest() {
	s.db.Delete(&eventstore.StreamEvent{}, "1 = 1")
	s.db.Delete(&eventstore.Stream{}, "1 = 1")
}

func (s *ReplayListenerTestSuite) TestReplayListener() {

	ctx, cancel := context.WithTimeout(s.ctx, time.Second*2)
	defer cancel()

	var now = time.Now().UTC()
	listener := NewReplayListener(
		s.eventStore,
		func() time.Time {
			return now.Add(-3 * time.Hour)
		})

	if !assert.NoError(s.T(), listener.Initialize(ctx, "replay-listener", []string{test.MockEventType})) {
		return
	}

	events := test.NewMockEvents(
		test.NewMockEvent("evt1", test.MockEventTime(now.Add(-4*time.Hour))),
		test.NewMockEvent("evt2", test.MockEventTime(now.Add(-2*time.Hour))),
		test.NewMockEvent("evt3", test.MockEventTime(now.Add(-1*time.Hour))),
	)

	_, err := s.eventStore.Save(ctx, test.MockStreamKey, 0, events)
	if !assert.NoError(s.T(), err) {
		return
	}

	var loaded []eventsource.Event
	err = listener.Listen(ctx, func(events []eventsource.Event) error {
		for _, loadedEvent := range events {
			loaded = append(loaded, loadedEvent)
			s.T().Logf("event received: %s", loadedEvent.GetEventID())
		}
		go func() {
			if len(loaded) == 3 {
				time.Sleep(50 * time.Millisecond)
				cancel()
			}
		}()

		return nil
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	if !assert.Len(s.T(), loaded, 2) {
		return
	}

	assert.Equal(s.T(), "evt2", loaded[0].GetEventID())
	assert.Equal(s.T(), "evt3", loaded[1].GetEventID())

}

func TestReplayListener(t *testing.T) {
	suite.Run(t, &ReplayListenerTestSuite{})
}
