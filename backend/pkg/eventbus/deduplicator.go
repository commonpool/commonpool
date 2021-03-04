package eventbus

import (
	"context"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/go-redis/redis/v8"
)

type EventDeduplicator interface {
	Deduplicate(
		ctx context.Context,
		evts []*eventstore.StreamEvent,
		do func(evt *eventstore.StreamEvent) error) error
}

type MemoryDeduplicator struct {
	bufferSize int
	lastIds    []string
}

func NewMemoryDeduplicator(bufferSize int) *MemoryDeduplicator {
	return &MemoryDeduplicator{
		bufferSize: bufferSize,
	}
}

func (s *MemoryDeduplicator) Deduplicate(
	ctx context.Context,
	evts []*eventstore.StreamEvent,
	do func(evt *eventstore.StreamEvent) error) error {
	for _, evt := range evts {
		found := false
		for _, id := range s.lastIds {
			if id == evt.EventID {
				found = true
				break
			}
		}
		if !found {
			err := do(evt)
			if err != nil {
				return err
			}
			s.lastIds = append(s.lastIds, evt.EventID)
			lastIdLen := len(s.lastIds)
			if lastIdLen > s.bufferSize {
				s.lastIds = s.lastIds[lastIdLen-s.bufferSize : lastIdLen]
			}
		}
	}
	return nil
}

var _ EventDeduplicator = &MemoryDeduplicator{}

type RedisDeduplicator struct {
	bufferSize  int
	redisClient *redis.Client
	key         string
}

func NewRedisDeduplicator(bufferSize int, redisClient *redis.Client, key string) *RedisDeduplicator {
	return &RedisDeduplicator{
		bufferSize:  bufferSize,
		redisClient: redisClient,
		key:         key,
	}
}

func (r RedisDeduplicator) Deduplicate(
	ctx context.Context,
	evts []*eventstore.StreamEvent,
	do func(evt *eventstore.StreamEvent) error) error {

	return r.redisClient.Watch(ctx, func(tx *redis.Tx) error {

		res := tx.LRange(ctx, r.key, -int64(r.bufferSize), -1)
		if res.Err() != nil {
			return res.Err()
		}

		result, err := res.Result()
		if err != nil {
			return err
		}

		resultMap := map[string]bool{}
		for _, key := range result {
			resultMap[key] = true
		}

		for _, evt := range evts {
			if resultMap[evt.EventID] {
				continue
			}

			err := do(evt)
			if err != nil {
				return err
			}

			if err := tx.RPush(ctx, r.key, evt.EventID).Err(); err != nil {
				return err
			}

			if err := tx.LTrim(ctx, r.key, -int64(r.bufferSize), -1).Err(); err != nil {
				return err
			}

			resultMap[evt.EventID] = true

		}

		return nil

	}, r.key)

}

var _ EventDeduplicator = &RedisDeduplicator{}
