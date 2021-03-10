package eventbus

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/clusterlock"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/commonpool/backend/pkg/mq"
	"time"
)

type CatchUpListener struct {
	eventStore   eventstore.EventStore
	getTimestamp func() time.Time
	amqpClient   mq.Client
	deduplicator EventDeduplicator
	clusterLock  clusterlock.Locker
	lockTTL      time.Duration
	lockOptions  *clusterlock.Options
	initialized  bool
	listener     Listener
	eventMapper  *eventsource.EventMapper
}

func NewCatchUpListener(
	eventStore eventstore.EventStore,
	getTimestamp func() time.Time,
	amqpClient mq.Client,
	deduplicator EventDeduplicator,
	clusterLock clusterlock.Locker,
	lockTTL time.Duration,
	lockOptions *clusterlock.Options,
	eventMapper *eventsource.EventMapper) *CatchUpListener {
	return &CatchUpListener{
		eventStore:   eventStore,
		getTimestamp: getTimestamp,
		amqpClient:   amqpClient,
		deduplicator: deduplicator,
		clusterLock:  clusterLock,
		lockTTL:      lockTTL,
		lockOptions:  lockOptions,
		eventMapper:  eventMapper,
	}
}

type CatchUpListenerFactory func(key string, lockTTL time.Duration) *CatchUpListener

func (c *CatchUpListener) Listen(ctx context.Context, listenerFunc ListenerFunc) error {
	if !c.initialized {
		return fmt.Errorf("not initialized")
	}
	return c.listener.Listen(ctx, listenerFunc)
}

func (c *CatchUpListener) Initialize(ctx context.Context, name string, eventTypes []string) error {
	// 	listener := NewDeduplicateListener(
	// 		c.deduplicator,
	// 		)

	listener := NewLockedListener(
		NewSequenceListener(
			[]Listener{
				NewReplayListener(c.eventStore, c.getTimestamp),
				NewRabbitMqListener(c.amqpClient, c.eventMapper),
			}),
		c.clusterLock,
		c.lockTTL,
		c.lockOptions)
	if err := listener.Initialize(ctx, name, eventTypes); err != nil {
		return err
	}
	c.listener = listener
	c.initialized = true
	return nil
}

var _ Listener = &CatchUpListener{}
