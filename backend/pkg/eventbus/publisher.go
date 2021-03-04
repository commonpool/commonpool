package eventbus

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/mq"
)

type EventPublisher interface {
	PublishEvents(ctx context.Context, events []*eventstore.StreamEvent) error
}

type AmqpPublisher struct {
	ctx        context.Context
	amqpClient mq.Client
}

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

func (a *AmqpPublisher) PublishEvents(ctx context.Context, events []*eventstore.StreamEvent) error {

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
				"event_type":     event.EventType,
				"stream_id":      event.StreamID,
				"stream_type":    event.StreamType,
				"correlation_id": event.CorrelationID,
				"sequence_no":    event.SequenceNo,
				"event_version":  event.EventVersion,
			},
			ContentType:     "application/json",
			ContentEncoding: "utf8",
			DeliveryMode:    2,
			Priority:        0,
			CorrelationId:   event.CorrelationID,
			ReplyTo:         "",
			Expiration:      "",
			MessageId:       event.EventID,
			Timestamp:       event.EventTime,
			Type:            event.EventType,
			UserId:          "",
			AppId:           "",
			Body:            evtBody,
		}); err != nil {
			return err
		}
	}

	return nil
}
