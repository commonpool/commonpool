package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/clusterlock"
	"time"
)

type LockedListener struct {
	name        string
	lock        clusterlock.Locker
	listener    Listener
	lockTTL     time.Duration
	lockOptions *clusterlock.Options
	initialized bool
}

func NewLockedListener(listener Listener, locker clusterlock.Locker, lockTtl time.Duration, lockOptions *clusterlock.Options) *LockedListener {
	return &LockedListener{
		lock:        locker,
		listener:    listener,
		lockTTL:     lockTtl,
		lockOptions: lockOptions,
	}
}

func (l *LockedListener) Listen(ctx context.Context, listenerFunc ListenerFunc) error {
	lock, err := l.lock.Obtain(ctx, "locks.listeners."+l.name, l.lockTTL, l.lockOptions)
	if err != nil {
		return err
	}
	defer lock.Release(ctx)
	return l.listener.Listen(ctx, listenerFunc)
}

func (l *LockedListener) Initialize(ctx context.Context, name string, eventTypes []string) error {
	l.name = name
	if err := l.listener.Initialize(ctx, name, eventTypes); err != nil {
		return err
	}
	l.initialized = true
	return nil
}

var _ Listener = &LockedListener{}
