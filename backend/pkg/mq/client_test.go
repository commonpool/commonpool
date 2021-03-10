package mq

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestClientConfigBuilder(t *testing.T) {
	args := map[string]interface{}{}

	c := NewConsumerConfig(func(d Delivery) {}).
		WithName("name").
		WithArgs(args).
		WithAutoAck(true).
		WithExclusive(true).
		WithNoWait(true)

	assert.Equal(t, "name", c.ConsumerName)
	assert.NotNil(t, c.Consume)
	assert.Equal(t, args, map[string]interface{}(c.Args))
	assert.Equal(t, true, c.AutoAck)
	assert.Equal(t, true, c.Exclusive)
	assert.Equal(t, true, c.NoWait)
}

func TestRabbitClient(t *testing.T) {

	config := NewRabbitConfig(os.Getenv("AMQP_URL"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	q := NewRabbitQueue(config, NewQueueConfig().WithName("test-rabbit-client"))

	assert.NoError(t, q.Connect(ctx))

	done := make(chan struct{})

	assert.NoError(t, q.Consume(ctx, NewConsumerConfig(func(d Delivery) {
		t.Logf("received message: %s", string(d.Body))
		done <- struct{}{}
	})))

	assert.NoError(t, q.Send(ctx, Message{Body: []byte("hello")}, false, false))

	<-done

}

func TestReconnect(t *testing.T) {

	config := NewRabbitConfig(os.Getenv("AMQP_URL"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	q := NewRabbitQueue(config, NewQueueConfig().WithName("test-rabbit-client-reconnect"))

	assert.NoError(t, q.Connect(ctx))

	done := make(chan struct{})

	messageCount := 0
	assert.NoError(t, q.Consume(ctx, NewConsumerConfig(func(d Delivery) {
		t.Logf("received message: %s", string(d.Body))
		messageCount++
		if messageCount == 2 {
			done <- struct{}{}
		}
	})))

	assert.NoError(t, q.Send(ctx, Message{Body: []byte("hello")}, false, false))

	time.Sleep(100 * time.Millisecond)

	assert.NoError(t, q.connection.Close())

	time.Sleep(100 * time.Millisecond)

	assert.NoError(t, q.Send(ctx, Message{Body: []byte("hello")}, false, false))

	<-done

}
