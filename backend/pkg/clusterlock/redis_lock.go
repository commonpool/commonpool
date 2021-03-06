package clusterlock

import (
	"context"
	"github.com/bsm/redislock"
	"time"
)

type Locker interface {
	Obtain(ctx context.Context, key string, ttl time.Duration, opt *Options) (Lock, error)
}

type Lock interface {
	Metadata() string
	Key() string
	Token() string
	Release(ctx context.Context) error
	Refresh(ctx context.Context, ttl time.Duration, opt *Options) error
}

type Redis struct {
	client *redislock.Client
}

var _ Locker = &Redis{}

type RedisLock struct {
	lock *redislock.Lock
}

var _ Lock = &RedisLock{}

type RetryStrategy interface {
	NextBackoff() time.Duration
}

type EverySecondRetryStrategy struct {
}

func (e EverySecondRetryStrategy) NextBackoff() time.Duration {
	return time.Millisecond * 10
}

var _ RetryStrategy = &EverySecondRetryStrategy{}

var EverySecond = &EverySecondRetryStrategy{}

type Options struct {
	RetryStrategy RetryStrategy
	Metadata      string
}

func NewRedis(redisLockClient *redislock.Client) *Redis {
	return &Redis{
		client: redisLockClient,
	}
}

func (l *Redis) Obtain(ctx context.Context, key string, ttl time.Duration, opt *Options) (Lock, error) {
	var opts *redislock.Options
	if opt != nil {
		opts = &redislock.Options{
			RetryStrategy: opt.RetryStrategy,
			Metadata:      opt.Metadata,
		}
	}
	lock, err := l.client.Obtain(ctx, key, ttl, opts)
	if err != nil {
		return nil, err
	}
	return &RedisLock{
		lock: lock,
	}, nil
}

func (l *RedisLock) Metadata() string {
	l.lock.Key()
	return l.lock.Metadata()
}

func (l *RedisLock) Key() string {
	return l.lock.Key()
}
func (l *RedisLock) Token() string {
	return l.lock.Token()
}

func (l *RedisLock) Release(ctx context.Context) error {
	return l.lock.Release(ctx)
}

func (l *RedisLock) Refresh(ctx context.Context, ttl time.Duration, opt *Options) error {
	var opts *redislock.Options
	if opt != nil {
		opts = &redislock.Options{
			RetryStrategy: opt.RetryStrategy,
			Metadata:      opt.Metadata,
		}
	}
	return l.lock.Refresh(ctx, ttl, opts)
}
