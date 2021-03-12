package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type RabbitMQListener struct {
	mqClient    mq.MqClient
	name        string
	eventTypes  []string
	initialized bool
	eventMapper *eventsource.EventMapper
}

type ListenerFunc func(ctx context.Context, events []eventsource.Event) error

func NewRabbitMqListener(mqClient mq.MqClient, eventMapper *eventsource.EventMapper) *RabbitMQListener {
	return &RabbitMQListener{
		mqClient:    mqClient,
		eventMapper: eventMapper,
	}
}

func (s *RabbitMQListener) Initialize(ctx context.Context, name string, eventTypes []string) error {

	s.initialized = true
	s.name = name
	s.eventTypes = eventTypes

	channel := s.mqClient.NewChannel()
	defer channel.Close()

	if err := channel.Connect(ctx); err != nil {
		return err
	}

	if err := channel.QueueDeclare(ctx, s.name, false, false, false, false, map[string]interface{}{}); err != nil {
		return err
	}

	for _, eventType := range s.eventTypes {
		if err := channel.QueueBind(ctx, s.name, "", "events.routed", false, map[string]interface{}{
			"event_type": eventType,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s *RabbitMQListener) Listen(ctx context.Context, listenerFunc ListenerFunc) error {

	l := logging.WithContext(ctx)
	l = l.Named("RabbitMQListener " + s.name)

	if !s.initialized {
		return fmt.Errorf("not initialized")
	}

	l.Debug("creating channel...")
	queue := s.mqClient.NewQueue(mq.NewQueueConfig().WithName(s.name))
	if err := queue.Connect(ctx); err != nil {
		return err
	}

	errChan := make(chan error)

	if err := queue.Consume(ctx, mq.NewConsumerConfig(func(msg mq.Delivery) error {

		var streamEvent eventstore.StreamEvent

		l.Debug("received event", zap.String("event_type", msg.Type))

		err := json.Unmarshal(msg.Body, &streamEvent)
		if err != nil {
			errChan <- errors.Wrap(err, "could not unmarshal event")
			return nil
		}

		evt, err := s.eventMapper.Map(msg.Type, msg.Body)
		if err != nil {
			errChan <- errors.Wrapf(err, "could not map event with type '%s'", msg.Type)
			return nil
		}

		err = listenerFunc(ctx, []eventsource.Event{evt})
		if err != nil {
			errChan <- errors.Wrap(err, "event listener error")
			return nil
		} else {
			if err := msg.Acknowledger.Ack(msg.DeliveryTag, false); err != nil {
				errChan <- errors.Wrap(err, "could not acknowledge delivery")
				return nil
			}
		}

		return nil

	})); err != nil {
		return err
	}

	for {
		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			return nil
		}
	}

}
