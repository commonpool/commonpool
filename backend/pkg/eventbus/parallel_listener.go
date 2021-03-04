package eventbus

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type ParallelListener struct {
	listeners []Listener
}

func NewParallelListener(listeners []Listener) *ParallelListener {
	return &ParallelListener{
		listeners: listeners,
	}
}

func (s *ParallelListener) Listen(ctx context.Context, listenerFunc ListenerFunc) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, listener := range s.listeners {
		listener := listener
		g.Go(func() error {
			return listener.Listen(ctx, listenerFunc)
		})
	}
	return g.Wait()
}

func (s *ParallelListener) Initialize(ctx context.Context, name string, eventTypes []string) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, listener := range s.listeners {
		listener := listener
		g.Go(func() error {
			return listener.Initialize(ctx, name, eventTypes)
		})
	}
	return g.Wait()
}

var _ Listener = &ParallelListener{}
