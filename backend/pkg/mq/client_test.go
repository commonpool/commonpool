package mq

import (
	"context"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestClientConfigBuilder(t *testing.T) {
	args := map[string]interface{}{}

	c := NewConsumerConfig(func(d Delivery) error {
		return nil
	}).
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

	go func() {
		assert.NoError(t, q.Consume(ctx, NewConsumerConfig(func(d Delivery) error {
			t.Logf("received message: %s", string(d.Body))
			done <- struct{}{}
			return nil
		})))
	}()

	time.Sleep(100 * time.Millisecond)

	assert.NoError(t, q.Send(ctx, Message{Body: []byte("hello")}, false, false))

	<-done

}

func TestReconnect(t *testing.T) {

	config := NewRabbitConfig(os.Getenv("AMQP_URL"))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	q := NewRabbitQueue(config, NewQueueConfig().WithName("test-rabbit-client-reconnect"))
	p := NewRabbitPublisher(config)

	assert.NoError(t, q.Connect(ctx))
	assert.NoError(t, p.Start(ctx))

	time.Sleep(100 * time.Millisecond)

	assert.NoError(t, q.connection.Close())
	assert.NoError(t, p.connection.Close())

	time.Sleep(100 * time.Millisecond)

	p.Publish(MessageEnvelope{
		Exchange:  "",
		Key:       q.Config.Name,
		Mandatory: false,
		Immediate: false,
		Msg: amqp.Publishing{
			Body: []byte("hello"),
		},
	})

	time.Sleep(100 * time.Millisecond)

	assert.NoError(t, q.connection.Close())
	assert.NoError(t, p.connection.Close())

	time.Sleep(100 * time.Millisecond)

	done := make(chan struct{})
	messageCount := 0
	go func() {
		err := q.Consume(ctx, NewConsumerConfig(func(d Delivery) error {
			t.Logf("received message: %s", string(d.Body))
			messageCount++
			assert.NoError(t, d.Acknowledger.Ack(d.DeliveryTag, false))
			if messageCount == 2 {
				done <- struct{}{}
			}
			return nil
		}))
		assert.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond)

	assert.NoError(t, q.connection.Close())
	assert.NoError(t, p.connection.Close())

	time.Sleep(100 * time.Millisecond)

	p.Publish(MessageEnvelope{
		Exchange:  "",
		Key:       q.Config.Name,
		Mandatory: false,
		Immediate: false,
		Msg: amqp.Publishing{
			Body: []byte("hello"),
		},
	})

	cancel()

	time.Sleep(1 * time.Second)

	<-done

}
