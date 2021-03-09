package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/mq"
	"go.uber.org/zap"
)

type RabbitMQListener struct {
	amqpClient  mq.Client
	name        string
	eventTypes  []string
	initialized bool
	eventMapper *eventsource.EventMapper
}

type ListenerFunc func(ctx context.Context, events []eventsource.Event) error

func NewRabbitMqListener(amqpClient mq.Client, eventMapper *eventsource.EventMapper) *RabbitMQListener {
	return &RabbitMQListener{
		amqpClient:  amqpClient,
		eventMapper: eventMapper,
	}
}

func (s *RabbitMQListener) Initialize(ctx context.Context, name string, eventTypes []string) error {

	s.initialized = true
	s.name = name
	s.eventTypes = eventTypes

	channel, err := s.amqpClient.GetChannel()
	if err != nil {
		return err
	}
	defer channel.Close()

	if err := channel.QueueDeclare(ctx, s.name, true, false, false, false, map[string]interface{}{}); err != nil {
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
	channel, err := s.amqpClient.GetChannel()
	if err != nil {
		return err
	}
	l.Debug("creating channel... done!")
	defer channel.Close()

	errChan := make(chan error)

	go func() {

		l.Debug("consuming RabbitMQ messages...")
		msgs, err := channel.Consume(ctx, s.name, "", false, false, false, false, map[string]interface{}{})
		if err != nil {
			errChan <- err
			return
		}
		l.Debug("consuming RabbitMQ messages... consuming!")

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgs:
				var streamEvent eventstore.StreamEvent

				l.Debug("received event", zap.String("event_type", msg.Type))

				err := json.Unmarshal(msg.Body, &streamEvent)
				if err != nil {
					l.Error("could not unmarshal event", zap.Error(err))
					continue
				}

				evt, err := s.eventMapper.Map(msg.Type, msg.Body)
				if err != nil {
					l.Error("could not map event with type", zap.String("event_type", msg.Type), zap.Error(err))
					continue
				}

				err = listenerFunc(ctx, []eventsource.Event{evt})
				if err != nil {
					l.Error("listener error", zap.Error(err), zap.String("event_type", msg.Type))
					continue
				} else {
					if err := msg.Acknowledger.Ack(msg.DeliveryTag, false); err != nil {
						l.Error("could not ack delivery", zap.Error(err))
						continue
					}
				}
			}
		}

	}()

	for {
		select {
		case err := <-errChan:
			return err
		case <-ctx.Done():
			return nil
		}
	}

}
