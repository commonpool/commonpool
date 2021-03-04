package postgres

import (
	"context"
	"github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type EventStoreSuite struct {
	suite.Suite
	testDB     *gorm.DB
	eventStore *PostgresEventStore
	ctx        context.Context
}

func (s *EventStoreSuite) SetupSuite() {
	s.testDB = db.NewTestDb()
	if err := s.testDB.AutoMigrate(&eventstore.StreamEvent{}, &eventstore.Stream{}); err != nil {
		s.T().FailNow()
	}
	s.eventStore = &PostgresEventStore{
		db: s.testDB,
	}
	s.ctx = context.Background()
}

func (s *EventStoreSuite) SetupTest() {
	if err := s.testDB.Delete(&eventstore.StreamEvent{}, "1 = 1").Error; err != nil {
		s.T().Error(err)
		s.T().FailNow()
	}
	if err := s.testDB.Delete(&eventstore.Stream{}, "1 = 1").Error; err != nil {
		s.T().Error(err)
		s.T().FailNow()
	}
}

func (s *EventStoreSuite) TestLoadEventsFromEmptyEventStore() {
	streamKey := eventstore.NewStreamKey("stream", "id1")
	events, err := s.eventStore.Load(s.ctx, streamKey)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), events)
}

func (s *EventStoreSuite) TestSaveEventsShouldStoreEvents() {

	streamKey := eventstore.NewStreamKey("stream", "id1")

	evt1 := eventstore.NewStreamEvent(streamKey, eventstore.NewStreamEventKey("event_type", "evt-1"), "payload")
	evt2 := eventstore.NewStreamEvent(streamKey, eventstore.NewStreamEventKey("event_type", "evt-2"), "payload")

	if !assert.NoError(s.T(), s.eventStore.Save(s.ctx, streamKey, 0, []*eventstore.StreamEvent{evt1, evt2})) {
		return
	}

	loadedEvents, err := s.eventStore.Load(s.ctx, streamKey)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), loadedEvents, 2)
	assert.Equal(s.T(), evt1, loadedEvents[0])
	assert.Equal(s.T(), evt2, loadedEvents[1])
}

func (s *EventStoreSuite) TestSaveShouldThrowWhenEmptyStreamIsNotExpectedVersion() {
	streamKey := eventstore.NewStreamKey("stream", "id1")
	evt := eventstore.NewStreamEvent(streamKey, eventstore.NewStreamEventKey("event_type", "evt-1"), "payload")
	createEvents := []*eventstore.StreamEvent{evt}
	assert.Error(s.T(), s.eventStore.Save(s.ctx, streamKey, 1, createEvents))
}

func (s *EventStoreSuite) TestSaveShouldThrowWhenStreamIsNotExpectedVersion() {
	streamKey := eventstore.NewStreamKey("stream", "id1")
	evt1 := eventstore.NewStreamEvent(streamKey, eventstore.NewStreamEventKey("event_type", "evt-1"), "payload")
	assert.NoError(s.T(), s.eventStore.Save(s.ctx, streamKey, 0, []*eventstore.StreamEvent{evt1}))
	evt2 := eventstore.NewStreamEvent(streamKey, eventstore.NewStreamEventKey("event_type", "evt-2"), "payload")
	assert.Error(s.T(), s.eventStore.Save(s.ctx, streamKey, 0, []*eventstore.StreamEvent{evt2}))
}

func (s *EventStoreSuite) TestSaveShouldThrowWhenEventsHaveSameID() {
	streamKey := eventstore.NewStreamKey("stream", "id1")
	evt1 := eventstore.NewStreamEvent(streamKey, eventstore.NewStreamEventKey("event_type", "evt-1"), "payload")
	assert.NoError(s.T(), s.eventStore.Save(s.ctx, streamKey, 0, []*eventstore.StreamEvent{evt1}))
	evt2 := eventstore.NewStreamEvent(streamKey, eventstore.NewStreamEventKey("event_type", "evt-1"), "payload")
	assert.Error(s.T(), s.eventStore.Save(s.ctx, streamKey, 0, []*eventstore.StreamEvent{evt2}))
}

func (s *EventStoreSuite) TestGetEventsByType() {

	now := time.Now()

	streamKey := eventstore.NewStreamKey("stream", "id1")

	evt1 := eventstore.NewStreamEvent(streamKey, eventstore.NewStreamEventKey("event_type", "evt-1"), "payload", eventstore.NewStreamEventOptions{
		EventTime: now.Add(-4 * time.Hour),
	})
	evt2 := eventstore.NewStreamEvent(streamKey, eventstore.NewStreamEventKey("event_type", "evt-2"), "payload", eventstore.NewStreamEventOptions{
		EventTime: now.Add(-2 * time.Hour),
	})
	evt3 := eventstore.NewStreamEvent(streamKey, eventstore.NewStreamEventKey("event_type", "evt-3"), "payload", eventstore.NewStreamEventOptions{
		EventTime: now.Add(-1 * time.Hour),
	})
	evts := []*eventstore.StreamEvent{evt1, evt2, evt3}
	assert.NoError(s.T(), s.eventStore.Save(s.ctx, streamKey, 0, evts))

	var loaded []*eventstore.StreamEvent
	err := s.eventStore.ReplayEventsByType(s.ctx, []string{"event_type"}, now.Add(-3*time.Hour), func(events []*eventstore.StreamEvent) error {
		for _, streamEvent := range events {
			loaded = append(loaded, streamEvent)
		}
		return nil
	})
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Len(s.T(), loaded, 2)

}

func TestEventStoreSuite(t *testing.T) {
	suite.Run(t, new(EventStoreSuite))
}
