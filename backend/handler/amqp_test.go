package handler

import (
	"context"
	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqptest"
	"github.com/NeowayLabs/wabbit/amqptest/server"
	"github.com/commonpool/backend/amqp"
)

type WabbitMqClient struct {
	Server     *server.AMQPServer
	Connection *amqptest.Conn
}

var _ amqp.AmqpClient = WabbitMqClient{}

type WabbitMqChannel struct {
	channel wabbit.Channel
}

func (w WabbitMqChannel) Close() error {
	panic("implement me")
}

func (w WabbitMqChannel) ExchangeBind(ctx context.Context, destination string, key string, source string, nowait bool, args amqp.AmqpArgs) error {

	return w.channel.Bind

	panic("implement me")
}

func (w WabbitMqChannel) ExchangeUnbind(ctx context.Context, destination string, key string, source string, noWait bool, args amqp.AmqpArgs) error {
	panic("implement me")
}

func (w WabbitMqChannel) ExchangeDeclare(ctx context.Context, name string, exchangeType string, durable bool, autoDelete bool, internal bool, nowait bool, args amqp.AmqpArgs) error {
	panic("implement me")
}

func (w WabbitMqChannel) ExchangeDelete(ctx context.Context, name string, ifUnused bool, noWait bool) error {
	panic("implement me")
}

func (w WabbitMqChannel) QueueDeclare(ctx context.Context, name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args amqp.AmqpArgs) error {
	panic("implement me")
}

func (w WabbitMqChannel) QueueBind(ctx context.Context, name string, key string, exchange string, nowait bool, args amqp.AmqpArgs) error {
	panic("implement me")
}

func (w WabbitMqChannel) QueueUnbind(ctx context.Context, name string, key string, exchange string, args amqp.AmqpArgs) error {
	panic("implement me")
}

func (w WabbitMqChannel) QueueDelete(ctx context.Context, name string, ifUnused bool, ifEmpty bool, noWait bool) error {
	panic("implement me")
}

func (w WabbitMqChannel) Consume(ctx context.Context, queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args amqp.AmqpArgs) (<-chan amqp.AmqpDelivery, error) {
	panic("implement me")
}

func (w WabbitMqChannel) Publish(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, publishing amqp.AmqpPublishing) error {
	panic("implement me")
}

var _ amqp.AmqpChannel = WabbitMqChannel{}

func (w WabbitMqClient) Shutdown() error {
	panic("implement me")
}

func (w WabbitMqClient) GetChannel() (amqp.AmqpChannel, error) {
	panic("implement me")
}

func NewWabbitMqClient(server *server.AMQPServer) (*amqp.AmqpClient, error) {
	mockConn, err := amqptest.Dial("amqp://localhost:5672/%2f")
	if err != nil {
		return nil, err
	}
	wabbitMqClient := WabbitMqClient{
		Server:     server,
		Connection: mockConn,
	}
	amqpClient := amqp.AmqpClient(wabbitMqClient)
	return &amqpClient, nil
}
