package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/db"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/eventstore/postgres"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"os"
	"testing"
	"time"
)

type RabbitListenerTestSuite struct {
	suite.Suite
	ctx           context.Context
	amqpClient    *mq.RabbitMqClient
	AmqpPublisher *AmqpPublisher
	db            *gorm.DB
	eventStore    *postgres.PostgresEventStore
}

func (s *RabbitListenerTestSuite) SetupSuite() {
	s.ctx = context.Background()
	s.db = db.NewTestDb()
	if err := s.db.AutoMigrate(&eventstore.StreamEvent{}, eventstore.Stream{}); err != nil {
		s.T().Fatal(err)
	}
	s.eventStore = postgres.NewPostgresEventStore(s.db)
	amqpClient, err := mq.NewRabbitMqClient(s.ctx, os.Getenv("AMQP_URL"))
	if err != nil {
		s.T().Fatal(err)
	}
	s.amqpClient = amqpClient
	s.AmqpPublisher = NewAmqpPublisher(s.amqpClient)
	err = s.AmqpPublisher.Init(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *RabbitListenerTestSuite) SetupTest() {
	s.db.Delete(&eventstore.StreamEvent{}, "1 = 1")
	s.db.Delete(&eventstore.Stream{}, "1 = 1")
}

func evt(eventType string, id string) *eventstore.StreamEvent {
	streamKey := eventstore.NewStreamKey("test-stream", "1")
	streamEventKey := eventstore.NewStreamEventKey(eventType, id)
	return eventstore.NewStreamEvent(
		streamKey,
		streamEventKey,
		`{"type":"hello"}`)
}

func anEvent(eventType string) *eventstore.StreamEvent {
	streamKey := eventstore.NewStreamKey("test-stream", "1")
	streamEventKey := eventstore.NewStreamEventKey(eventType, "1")
	return eventstore.NewStreamEvent(
		streamKey,
		streamEventKey,
		`{"type":"hello"}`)
}

func (s *RabbitListenerTestSuite) TestSubscriberIsCalled() {

	ctx, cancel := context.WithCancel(s.ctx)
	ctx, cancel = context.WithTimeout(ctx, time.Millisecond*5000)
	defer cancel()

	sub := NewRabbitMqListener(s.amqpClient)
	if assert.NoError(s.T(), sub.Initialize(ctx, "test-event-subscriber", []string{"event-type-1"})) {
		return
	}

	subscriberCalled := false
	go func() {
		err := sub.Listen(ctx, func(events []*eventstore.StreamEvent) error {
			s.T().Log("subscriber called")
			subscriberCalled = true
			cancel()
			return nil
		})
		assert.NoError(s.T(), err)
	}()

	go func() {
		err := s.AmqpPublisher.PublishEvents(ctx, []*eventstore.StreamEvent{anEvent("event-type-1")})
		if err != nil {
			s.T().Fatal(err)
		}
	}()

	<-ctx.Done()

	assert.True(s.T(), subscriberCalled)

}

func (s *RabbitListenerTestSuite) TestMessagesPersisted() {

	sub := NewRabbitMqListener(s.amqpClient)

	ctx1, cancel1 := context.WithTimeout(s.ctx, time.Millisecond*5000)
	defer cancel1()

	if !assert.NoError(s.T(), sub.Initialize(ctx1, "test-messages-persisted", []string{"event-type-2"})) {
		return
	}

	go func() {
		err := sub.Listen(ctx1, func(events []*eventstore.StreamEvent) error {
			return nil
		})
		assert.NoError(s.T(), err)
	}()

	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel1()
	}()

	<-ctx1.Done()

	s.T().Log("Context 1 is done")

	s.T().Log("Publishing event")
	err := s.AmqpPublisher.PublishEvents(s.ctx, []*eventstore.StreamEvent{anEvent("event-type-2")})
	if err != nil {
		s.T().Fatal(err)
	}

	ctx2, cancel2 := context.WithTimeout(s.ctx, time.Millisecond*5000)

	called := false

	go func() {
		err := sub.Listen(ctx2, func(events []*eventstore.StreamEvent) error {
			s.T().Log("Event received")
			called = true
			cancel2()
			return nil
		})
		assert.NoError(s.T(), err)
	}()

	<-ctx2.Done()

	assert.True(s.T(), called)
}

func TestRabbitListener(t *testing.T) {
	suite.Run(t, &RabbitListenerTestSuite{})
}
