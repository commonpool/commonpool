package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/mq"
	"github.com/labstack/gommon/log"
)

type RabbitMQListener struct {
	amqpClient  mq.Client
	name        string
	eventTypes  []string
	initialized bool
}

type ListenerFunc func(events []*eventstore.StreamEvent) error

func NewRabbitMqListener(amqpClient mq.Client) *RabbitMQListener {
	return &RabbitMQListener{
		amqpClient: amqpClient,
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

	if !s.initialized {
		return fmt.Errorf("not initialized")
	}

	channel, err := s.amqpClient.GetChannel()
	if err != nil {
		return err
	}
	defer channel.Close()

	errChan := make(chan error)

	go func() {

		msgs, err := channel.Consume(ctx, s.name, "", false, false, false, false, map[string]interface{}{})
		if err != nil {
			errChan <- err
			return
		}

		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgs:
				var streamEvent eventstore.StreamEvent
				err := json.Unmarshal(msg.Body, &streamEvent)
				if err != nil {
					log.Printf("could not unmarshal event: %v", err)
					continue
				}
				err = listenerFunc([]*eventstore.StreamEvent{&streamEvent})
				if err != nil {
					log.Printf("listener error: %v", err)
					continue
				} else {
					if err := msg.Acknowledger.Ack(msg.DeliveryTag, false); err != nil {
						log.Printf("could not ack delivery: %v", err)
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
