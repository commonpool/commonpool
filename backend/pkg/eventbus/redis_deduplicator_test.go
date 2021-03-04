package eventbus

import (
	"context"
	"crypto/tls"
	"github.com/commonpool/backend/pkg/config"
	"github.com/commonpool/backend/pkg/eventstore"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func getRedisClient() (*redis.Client, error) {

	appConfig, err := config.GetAppConfig(os.LookupEnv, ioutil.ReadFile)
	if err != nil {
		return nil, err
	}

	var redisTlsConfig *tls.Config = nil
	if appConfig.RedisTlsEnabled {
		redisTlsConfig = &tls.Config{
			InsecureSkipVerify: appConfig.RedisTlsSkipVerify,
		}
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:      appConfig.RedisHost + ":" + appConfig.RedisPort,
		Password:  appConfig.RedisPassword,
		DB:        0,
		TLSConfig: redisTlsConfig,
	})

	return redisClient, nil
}

func TestRedisDeduplicator(t *testing.T) {

	redisClient, err := getRedisClient()
	if !assert.NoError(t, err) {
		return
	}

	key := "tests.redis-deduplicator"

	if _, err := redisClient.Del(context.TODO(), key).Result(); !assert.NoError(t, err) {
		return
	}

	d := NewRedisDeduplicator(10, redisClient, key)

	var calls []*eventstore.StreamEvent
	if !assert.NoError(t, d.Deduplicate(context.TODO(), evts(
		evt("t1", "1"),
		evt("t1", "2"),
		evt("t1", "1"),
		evt("t1", "3"),
	), func(evt *eventstore.StreamEvent) error {
		calls = append(calls, evt)
		return nil
	})) {
		return
	}

	assert.Len(t, calls, 3)
	assert.Equal(t, "1", calls[0].EventID)
	assert.Equal(t, "2", calls[1].EventID)
	assert.Equal(t, "3", calls[2].EventID)

	calls = []*eventstore.StreamEvent{}

	if !assert.NoError(t, d.Deduplicate(context.TODO(), evts(
		evt("t1", "1"),
		evt("t1", "2"),
		evt("t1", "1"),
		evt("t1", "3"),
		evt("t1", "4"),
		evt("t1", "5"),
	), func(evt *eventstore.StreamEvent) error {
		calls = append(calls, evt)
		return nil
	})) {
		return
	}

	assert.Len(t, calls, 2)
	assert.Equal(t, "4", calls[0].EventID)
	assert.Equal(t, "5", calls[1].EventID)

}
