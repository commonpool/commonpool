package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/eventstore/postgres"
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
	s.db = db.NewTestDb()
	if err := s.db.AutoMigrate(&eventstore.StreamEvent{}, eventstore.Stream{}); err != nil {
		s.T().Fatal(err)
	}
	s.eventStore = postgres.NewPostgresEventStore(s.db)
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

	if !assert.NoError(s.T(), listener.Initialize(ctx, "replay-listener", []string{"test-replay-listener-evt"})) {
		return
	}

	streamKey := eventstore.NewStreamKey("stream", "id1")

	evt1 := eventstore.NewStreamEvent(
		streamKey,
		eventstore.NewStreamEventKey("test-replay-listener-evt", "evt-1"),
		"payload",
		eventstore.NewStreamEventOptions{
			EventTime: now.Add(-4 * time.Hour),
		})

	evt2 := eventstore.NewStreamEvent(
		streamKey,
		eventstore.NewStreamEventKey("test-replay-listener-evt", "evt-2"),
		"payload",
		eventstore.NewStreamEventOptions{
			EventTime: now.Add(-2 * time.Hour),
		})

	evt3 := eventstore.NewStreamEvent(
		streamKey,
		eventstore.NewStreamEventKey("test-replay-listener-evt", "evt-3"),
		"payload",
		eventstore.NewStreamEventOptions{
			EventTime: now.Add(-1 * time.Hour),
		})

	assert.NoError(s.T(), s.eventStore.Save(ctx, streamKey, 0, []*eventstore.StreamEvent{evt1, evt2, evt3}))

	var loaded []*eventstore.StreamEvent
	err := listener.Listen(ctx, func(events []*eventstore.StreamEvent) error {
		for _, loadedEvent := range events {
			loaded = append(loaded, loadedEvent)
			s.T().Logf("event received: %s", loadedEvent.EventID)
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

	assert.Len(s.T(), loaded, 2)

	assert.Equal(s.T(), "evt-2", loaded[0].EventID)
	assert.Equal(s.T(), "evt-3", loaded[1].EventID)

}

func TestReplayListener(t *testing.T) {
	suite.Run(t, &ReplayListenerTestSuite{})
}
