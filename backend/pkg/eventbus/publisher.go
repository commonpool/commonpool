package eventbus

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/mq"
)

type EventPublisher interface {
	PublishEvents(ctx context.Context, events []eventsource.Event) error
}

type AmqpPublisher struct {
	ctx        context.Context
	amqpClient mq.Client
}

var _ EventPublisher = &AmqpPublisher{}

func NewAmqpPublisher(amqpClient mq.Client) *AmqpPublisher {
	return &AmqpPublisher{
		amqpClient: amqpClient,
	}
}

func (a *AmqpPublisher) Init(ctx context.Context) error {
	channel, err := a.amqpClient.GetChannel()
	if err != nil {
		return err
	}
	if err := channel.ExchangeDeclare(ctx, "events.fanout", "fanout", true, false, false, false, map[string]interface{}{}); err != nil {
		return err
	}

	if err := channel.ExchangeDeclare(ctx, "events.routed", "headers", true, false, false, false, map[string]interface{}{}); err != nil {
		return err
	}

	if err := channel.ExchangeBind(ctx, "events.routed", "", "events.fanout", false, map[string]interface{}{}); err != nil {
		return err
	}

	return nil
}

func (a *AmqpPublisher) PublishEvents(ctx context.Context, events []eventsource.Event) error {

	channel, err := a.amqpClient.GetChannel()
	if err != nil {
		return err
	}

	for _, event := range events {
		evtBody, err := json.Marshal(event)
		if err != nil {
			return err
		}

		if err := channel.Publish(ctx, "events.fanout", "", true, false, mq.Message{
			Headers: map[string]interface{}{
				"event_type":     event.GetEventType(),
				"aggregate_id":   event.GetAggregateID(),
				"aggregate_type": event.GetAggregateType(),
				"correlation_id": event.GetCorrelationID(),
				"sequence_no":    event.GetSequenceNo(),
				"event_version":  event.GetEventVersion(),
			},
			ContentType:     "application/json",
			ContentEncoding: "utf8",
			DeliveryMode:    2,
			Priority:        0,
			CorrelationId:   event.GetCorrelationID(),
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       event.GetEventID(),
			Timestamp:       event.GetEventTime(),
			Type:            event.GetEventType(),
			UserId:          "",
			AppId:           "",
			Body:            evtBody,
		}); err != nil {
			return err
		}
	}

	return nil
}
