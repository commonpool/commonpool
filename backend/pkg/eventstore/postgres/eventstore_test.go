package postgres

import (
	"context"
	"github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/test"
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
	eventMapper := eventsource.NewEventMapper()
	if err := test.RegisterMockEvents(eventMapper); err != nil {
		s.FailNow(err.Error())
	}
	s.eventStore = &PostgresEventStore{
		db:          s.testDB,
		eventMapper: eventMapper,
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
	streamKey := keys.NewStreamKey(test.MockAggregateType, "mock-id")
	events, err := s.eventStore.Load(s.ctx, streamKey)
	assert.NoError(s.T(), err)
	assert.Empty(s.T(), events)
}

func (s *EventStoreSuite) TestSaveEventsShouldSetCorrelationID() {

	streamKey := keys.NewStreamKey(test.MockAggregateType, "mock-id")

	events := test.NewMockEvents(
		test.NewMockEvent("1"),
		test.NewMockEvent("2"),
	)

	_, err := s.eventStore.Save(s.ctx, streamKey, 0, events)
	if !assert.NoError(s.T(), err) {
		return
	}

	loadedEvents, err := s.eventStore.Load(s.ctx, streamKey)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), loadedEvents, 2)
	assert.NotEmpty(s.T(), loadedEvents[0].GetCorrelationID())
	assert.NotEmpty(s.T(), loadedEvents[1].GetCorrelationID())
	assert.Equal(s.T(), loadedEvents[0].GetCorrelationID(), loadedEvents[1].GetCorrelationID())
}
func (s *EventStoreSuite) TestSaveEventsShouldSetEventID() {

	streamKey := keys.NewStreamKey(test.MockAggregateType, "mock-id")

	events := test.NewMockEvents(
		test.NewMockEvent("1"),
		test.NewMockEvent("2"),
	)

	_, err := s.eventStore.Save(s.ctx, streamKey, 0, events)
	if !assert.NoError(s.T(), err) {
		return
	}

	loadedEvents, err := s.eventStore.Load(s.ctx, streamKey)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), loadedEvents, 2)
	assert.NotEmpty(s.T(), loadedEvents[0].GetEventID())
	assert.NotEmpty(s.T(), loadedEvents[1].GetEventID())
	assert.NotEqual(s.T(), loadedEvents[0].GetEventID(), loadedEvents[1].GetEventID())
}

func (s *EventStoreSuite) TestSaveEventsShouldSetSequenceNoForNewStreams() {
	streamKey := keys.NewStreamKey(test.MockAggregateType, "mock-id")
	_, err := s.eventStore.Save(s.ctx, streamKey, 0, test.NewMockEvents(
		test.NewMockEvent("1"),
		test.NewMockEvent("2"),
	))
	if !assert.NoError(s.T(), err) {
		return
	}
	loadedEvents, err := s.eventStore.Load(s.ctx, streamKey)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), loadedEvents, 2)
	assert.Equal(s.T(), 0, loadedEvents[0].GetSequenceNo())
	assert.Equal(s.T(), 1, loadedEvents[1].GetSequenceNo())
}

func (s *EventStoreSuite) TestSaveEventsShouldSetSequenceNoForExistingStreams() {
	streamKey := keys.NewStreamKey(test.MockAggregateType, "mock-id")

	_, err := s.eventStore.Save(s.ctx, streamKey, 0, test.NewMockEvents(
		test.NewMockEvent("1"),
		test.NewMockEvent("2"),
	))
	if !assert.NoError(s.T(), err) {
		return
	}

	_, err = s.eventStore.Save(s.ctx, streamKey, 2, test.NewMockEvents(
		test.NewMockEvent("3"),
	))
	if !assert.NoError(s.T(), err) {
		return
	}

	loadedEvents, err := s.eventStore.Load(s.ctx, streamKey)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), loadedEvents, 3)
	assert.Equal(s.T(), 0, loadedEvents[0].GetSequenceNo())
	assert.Equal(s.T(), 1, loadedEvents[1].GetSequenceNo())
	assert.Equal(s.T(), 2, loadedEvents[2].GetSequenceNo())
}

func (s *EventStoreSuite) TestSaveShouldThrowWhenEmptyStreamIsNotExpectedVersion() {
	streamKey := keys.NewStreamKey(test.MockAggregateType, "mock-id")
	events := test.NewMockEvents(
		test.NewMockEvent("1"),
		test.NewMockEvent("2"),
	)
	_, err := s.eventStore.Save(s.ctx, streamKey, 1, events)
	assert.Error(s.T(), err)
}

func (s *EventStoreSuite) TestSaveShouldThrowWhenStreamIsNotExpectedVersion() {
	streamKey := keys.NewStreamKey(test.MockAggregateType, "mock-id")
	evt1 := test.NewMockEvent("1")
	evt2 := test.NewMockEvent("2")
	_, err := s.eventStore.Save(s.ctx, streamKey, 0, test.NewMockEvents(evt1))
	assert.NoError(s.T(), err)
	_, err = s.eventStore.Save(s.ctx, streamKey, 0, test.NewMockEvents(evt2))
	assert.Error(s.T(), err)
}

func (s *EventStoreSuite) TestSaveShouldThrowWhenEventsHaveSameID() {
	streamKey := keys.NewStreamKey(test.MockAggregateType, "mock-id")
	evt1 := test.NewMockEvent("1")
	evt2 := test.NewMockEvent("2")
	_, err := s.eventStore.Save(s.ctx, streamKey, 0, test.NewMockEvents(evt1))
	assert.NoError(s.T(), err)
	_, err = s.eventStore.Save(s.ctx, streamKey, 0, test.NewMockEvents(evt2))
	assert.Error(s.T(), err)
}

func (s *EventStoreSuite) TestGetEventsByType() {

	now := time.Now()

	streamKey := keys.NewStreamKey(test.MockAggregateType, "mock-id")

	evt1 := test.NewMockEvent("1")
	evt2 := test.NewMockEvent("2")
	evt3 := test.NewMockEvent("3")
	evt1.EventTime = now.Add(-4 * time.Hour)
	evt2.EventTime = now.Add(-2 * time.Hour)
	evt3.EventTime = now.Add(-1 * time.Hour)
	events := test.NewMockEvents(
		evt1,
		evt2,
		evt3,
	)

	_, err := s.eventStore.Save(s.ctx, streamKey, 0, events)
	assert.NoError(s.T(), err)

	var loaded []eventsource.Event
	err = s.eventStore.ReplayEventsByType(s.ctx, []string{test.MockEventType}, now.Add(-3*time.Hour), func(events []eventsource.Event) error {
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
