package eventbus

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/eventsource"
)

type StaticListener struct {
	events      []eventsource.Event
	initialized bool
}

func NewStaticListener(events []eventsource.Event) *StaticListener {
	return &StaticListener{
		events: events,
	}
}

func (s *StaticListener) Listen(ctx context.Context, listenerFunc ListenerFunc) error {
	if !s.initialized {
		return fmt.Errorf("not initialized")
	}
	return listenerFunc(ctx, s.events)
}

func (s *StaticListener) Initialize(ctx context.Context, name string, eventTypes []string) error {
	s.initialized = true
	return nil
}

var _ Listener = &StaticListener{}
