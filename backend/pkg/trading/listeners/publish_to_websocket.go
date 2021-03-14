package listeners

import (
	"context"
	"github.com/commonpool/backend/pkg/eventbus"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/mq"
)

type PublishToWebsocketListener struct {
	mqClient    mq.MqClient
	eventMapper *eventsource.EventMapper
}

func NewPublishToWebsocketListener(
	mqClient mq.MqClient,
	eventMapper *eventsource.EventMapper) *PublishToWebsocketListener {
	return &PublishToWebsocketListener{
		mqClient:    mqClient,
		eventMapper: eventMapper,
	}
}

func (p *PublishToWebsocketListener) Start(ctx context.Context) error {
	listener := eventbus.NewRabbitMqListener(p.mqClient, p.eventMapper)
	if err := listener.Initialize(ctx, "PublishToWebsocketListener", []string{}); err != nil {
		return err
	}
	return listener.Listen(ctx, func(ctx context.Context, events []eventsource.Event) error {
		for _, _ = range events {
		}
		return nil
	})
}
