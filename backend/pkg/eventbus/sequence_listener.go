package eventbus

import (
	"context"
	"fmt"
)

type SequenceListener struct {
	listeners   []Listener
	initialized bool
}

func NewSequenceListener(listeners []Listener) *SequenceListener {
	return &SequenceListener{
		listeners: listeners,
	}
}

func (s *SequenceListener) Listen(ctx context.Context, listenerFunc ListenerFunc) error {
	if !s.initialized {
		return fmt.Errorf("not initialized")
	}
	for _, listener := range s.listeners {
		err := listener.Listen(ctx, listenerFunc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SequenceListener) Initialize(ctx context.Context, name string, eventTypes []string) error {
	for _, listener := range s.listeners {
		if err := listener.Initialize(ctx, name, eventTypes); err != nil {
			return err
		}
	}
	s.initialized = true
	return nil
}

var _ Listener = &SequenceListener{}
