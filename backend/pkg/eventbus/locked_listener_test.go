package eventbus

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/commonpool/backend/pkg/clusterlock"
	"github.com/commonpool/backend/pkg/eventsource"
	"github.com/commonpool/backend/pkg/test"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"strconv"
	"testing"
	"time"
)

func TestLockedListener(t *testing.T) {

	redisClient, err := getRedisClient()
	if !assert.NoError(t, err) {
		return
	}

	lock := clusterlock.NewRedis(redislock.New(redisClient))

	var listeners []Listener = []Listener{}
	var listenerCount = 30
	var eventCount = 30

	for i := 0; i < listenerCount; i++ {
		events := test.NewMockEvents()
		for j := 0; j < eventCount; j++ {
			events = append(events, test.NewMockEvent(strconv.Itoa(j+i*eventCount)))
		}
		staticListener := NewStaticListener(events)
		lockedListener := NewLockedListener(staticListener, lock, time.Second*5, &clusterlock.Options{
			RetryStrategy: clusterlock.EverySecond,
		})
		listeners = append(listeners, lockedListener)
	}

	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	var callIds []string

	runListener := func(i int, listener Listener, initialDelay time.Duration) {
		g.Go(func() error {
			time.Sleep(initialDelay)

			if err := listener.Initialize(ctx, "listener", []string{"bla"}); !assert.NoError(t, err) {
				return err
			}
			if err := listener.Listen(ctx, func(events []eventsource.Event) error {
				t.Logf("entering listener %d", i)
				defer t.Logf(" exiting listener %d", i)
				for _, event := range events {
					callIds = append(callIds, event.GetEventID())
				}
				return nil
			}); !assert.NoError(t, err) {
				return err
			}
			return nil
		})
	}

	for i := 0; i < listenerCount; i++ {
		runListener(i, listeners[i], time.Millisecond*time.Duration(i)*20)
	}

	g.Wait()

	assert.Len(t, callIds, listenerCount*eventCount)
	for i, _ := range callIds {
		assert.Equal(t, strconv.Itoa(i), callIds[i])
	}

	<-ctx.Done()

}
