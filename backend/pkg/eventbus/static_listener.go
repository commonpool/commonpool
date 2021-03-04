package eventbus

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/eventstore"
)

type StaticListener struct {
	events      []*eventstore.StreamEvent
	initialized bool
}

func NewStaticListener(events []*eventstore.StreamEvent) *StaticListener {
	return &StaticListener{
		events: events,
	}
}

func (s *StaticListener) Listen(ctx context.Context, listenerFunc ListenerFunc) error {
	if !s.initialized {
		return fmt.Errorf("not initialized")
	}
	return listenerFunc(s.events)
}

func (s *StaticListener) Initialize(ctx context.Context, name string, eventTypes []string) error {
	s.initialized = true
	return nil
}

var _ Listener = &StaticListener{}
