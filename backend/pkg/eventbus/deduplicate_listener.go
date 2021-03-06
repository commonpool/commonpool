package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/eventsource"
)

type DeduplicateListener struct {
	listener     Listener
	initialized  bool
	deduplicator EventDeduplicator
}

func NewDeduplicateListener(deduplicator EventDeduplicator, listener Listener) *DeduplicateListener {
	return &DeduplicateListener{
		listener:     listener,
		deduplicator: deduplicator,
	}
}

func (s *DeduplicateListener) Listen(ctx context.Context, listenerFunc ListenerFunc) error {
	return s.listener.Listen(ctx, func(events []eventsource.Event) error {
		return s.deduplicator.Deduplicate(ctx, events, func(evt eventsource.Event) error {
			return listenerFunc([]eventsource.Event{evt})
		})
	})
}

func (s *DeduplicateListener) Initialize(ctx context.Context, name string, eventTypes []string) error {
	if err := s.listener.Initialize(ctx, name, eventTypes); err != nil {
		return err
	}
	s.initialized = true
	return nil
}

var _ Listener = &DeduplicateListener{}
