package eventbus

import "context"

type Listener interface {
	Listen(ctx context.Context, listenerFunc ListenerFunc) error
	Initialize(ctx context.Context, name string, eventTypes []string) error
}
