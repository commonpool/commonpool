package mq

import (
	"context"
	"github.com/commonpool/backend/logging"
	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var _ Client = RabbitMqClient{}

type Client interface {
	Shutdown() error
	GetChannel() (Channel, error)
}

type Channel interface {
	Close() error
	ExchangeBind(ctx context.Context, destination string, key string, source string, nowait bool, args Args) error
	ExchangeUnbind(ctx context.Context, destination string, key string, source string, noWait bool, args Args) error
	ExchangeDeclare(ctx context.Context, name string, exchangeType string, durable bool, autoDelete bool, internal bool, nowait bool, args Args) error
	ExchangeDelete(ctx context.Context, name string, ifUnused bool, noWait bool) error
	QueueDeclare(ctx context.Context, name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args Args) error
	QueueBind(ctx context.Context, name string, key string, exchange string, nowait bool, args Args) error
	QueueUnbind(ctx context.Context, name string, key string, exchange string, args Args) error
	QueueDelete(ctx context.Context, name string, ifUnused bool, ifEmpty bool, noWait bool) error
	Consume(ctx context.Context, queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args Args) (<-chan Delivery, error)
	Publish(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, publishing Message) error
}

type Ack interface {
	Ack(tag uint64, multiple bool) error
	Nack(tag uint64, multiple bool, requeue bool) error
	Reject(tag uint64, requeue bool) error
}

type RabbitMqChannel struct {
	channel *amqp.Channel
	msg     <-chan amqp.Delivery
}

func (r RabbitMqChannel) Close() error {
	return r.channel.Close()
}

var _ Channel = &RabbitMqChannel{}

type RabbitMqClient struct {
	connection     *amqp.Connection
	isShuttingDown bool
	channel        *amqp.Channel
}

func (r RabbitMqClient) GetChannel() (Channel, error) {
	ch, err := r.connection.Channel()
	if err != nil {
		return nil, err
	}
	var amqpChannel = &RabbitMqChannel{
		channel: ch,
	}
	return amqpChannel, nil
}

func (r RabbitMqClient) Shutdown() error {
	r.isShuttingDown = true
	if !r.connection.IsClosed() {
		if err := r.connection.Close(); err != nil {
			log.Errorf("could not shut down amqp client: %v", err)
			return err
		}
	}
	return nil
}

func (r RabbitMqChannel) Publish(ctx context.Context, exchange string, key string, mandatory bool, immediate bool, publishing Message) error {
	return r.channel.Publish(exchange, key, mandatory, immediate, amqp.Publishing{
		Headers:         map[string]interface{}(publishing.Headers),
		ContentType:     publishing.ContentType,
		ContentEncoding: publishing.ContentEncoding,
		DeliveryMode:    publishing.DeliveryMode,
		Priority:        publishing.Priority,
		CorrelationId:   publishing.CorrelationId,
		ReplyTo:         publishing.ReplyTo,
		Expiration:      publishing.Expiration,
		MessageId:       publishing.MessageId,
		Timestamp:       publishing.Timestamp,
		Type:            publishing.Type,
		UserId:          publishing.UserId,
		AppId:           publishing.AppId,
		Body:            publishing.Body,
	})
}

func NewRabbitMqClient(ctx context.Context, amqpUrl string) (*RabbitMqClient, error) {

	l := logging.WithContext(ctx)

	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		l.Error("could not connect to RabbitMQ", zap.Error(err))
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		l.Error("could not create channel", zap.Error(err))
		return nil, err
	}

	rabbitMqClient := &RabbitMqClient{
		connection: conn,
		channel:    ch,
	}

	channel, err := rabbitMqClient.GetChannel()
	if err != nil {
		l.Error("could not get channel", zap.Error(err))
		return nil, err
	}

	if err := channel.ExchangeDeclare(ctx, MessagesExchange, amqp.ExchangeHeaders, true, false, false, false, nil); err != nil {
		l.Error("could not create messages exchange", zap.Error(err))
		return nil, err
	}

	if err := channel.ExchangeDeclare(
		ctx, WebsocketMessagesExchange, amqp.ExchangeHeaders, true, false, false, false, nil); err != nil {
		l.Error("could not create messages websocket exchange", zap.Error(err))
		return nil, err
	}

	headers := map[string]interface{}{EventTypeKey: EventTypeMessage}
	if err = channel.ExchangeBind(ctx, WebsocketMessagesExchange, "", MessagesExchange, false, headers); err != nil {
		l.Error("could not bind messages and websocket exchanges", zap.Error(err))
		return nil, err
	}

	return rabbitMqClient, nil

}

func (r RabbitMqChannel) Consume(ctx context.Context, queue string, consumer string, autoAck bool, exclusive bool, noLocal bool, noWait bool, args Args) (<-chan Delivery, error) {

	l := logging.WithContext(ctx).With(
		zap.String("queue", queue),
		zap.String("consumer", consumer),
		zap.Bool("autoAck", autoAck),
		zap.Bool("exclusive", exclusive),
		zap.Bool("noLocal", noLocal),
		zap.Bool("noWait", noWait),
		zap.Object("args", args))

	ch, err := r.channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, map[string]interface{}(args))
	if err != nil {
		l.Error("could not consume queue", zap.Error(err))
	}

	amqpChan := make(chan Delivery)

	go func() {
		for delivery := range ch {
			var amqpDelivery = Delivery{
				Acknowledger:    delivery.Acknowledger,
				Headers:         map[string]interface{}(delivery.Headers),
				ContentType:     delivery.ContentType,
				ContentEncoding: delivery.ContentEncoding,
				DeliveryMode:    delivery.DeliveryMode,
				Priority:        delivery.Priority,
				CorrelationId:   delivery.CorrelationId,
				ReplyTo:         delivery.ReplyTo,
				Expiration:      delivery.Expiration,
				MessageId:       delivery.MessageId,
				Timestamp:       delivery.Timestamp,
				Type:            delivery.Type,
				UserId:          delivery.UserId,
				AppId:           delivery.AppId,
				ConsumerTag:     delivery.ConsumerTag,
				MessageCount:    delivery.MessageCount,
				DeliveryTag:     delivery.DeliveryTag,
				Redelivered:     delivery.Redelivered,
				Exchange:        delivery.Exchange,
				RoutingKey:      delivery.RoutingKey,
				Body:            delivery.Body,
			}
			amqpChan <- amqpDelivery
		}
		close(amqpChan)
	}()

	return amqpChan, nil
}

func (r RabbitMqChannel) ExchangeDeclare(ctx context.Context, name string, exchangeType string, durable bool, autoDelete bool, internal bool, nowait bool, args Args) error {

	l := logging.WithContext(ctx).With(
		zap.String("exchange_name", name),
		zap.String("exchange_type", exchangeType),
		zap.Bool("durable", durable),
		zap.Bool("autoDelete", autoDelete),
		zap.Bool("internal", internal),
		zap.Bool("nowait", nowait),
		zap.Object("args", args))

	if err := r.channel.ExchangeDeclare(name, exchangeType, durable, autoDelete, internal, nowait, map[string]interface{}(args)); err != nil {
		l.Error("could not declare exchange", zap.Error(err))
	}

	return nil

}

func (r RabbitMqChannel) ExchangeBind(ctx context.Context, destination string, key string, source string, nowait bool, args Args) error {

	l := logging.WithContext(ctx).With(
		zap.String("source", source),
		zap.String("destination", destination),
		zap.String("binding_key", key),
		zap.Bool("nowait", nowait),
		zap.Object("args", args))

	if err := r.channel.ExchangeBind(destination, key, source, nowait, map[string]interface{}(args)); err != nil {
		l.Error("could not bind exchanges")
		return err
	}

	return nil

}

func (r RabbitMqChannel) ExchangeUnbind(ctx context.Context, destination string, key string, source string, noWait bool, args Args) error {

	l := logging.WithContext(ctx).With(
		zap.String("destination", destination),
		zap.String("key", key),
		zap.String("source", source),
		zap.Bool("noWait", noWait),
		zap.Object("args", args))

	if err := r.channel.ExchangeUnbind(destination, key, source, noWait, map[string]interface{}(args)); err != nil {
		l.Error("could not unbind exchanges")
		return err
	}

	return nil

}

func (r RabbitMqChannel) ExchangeDelete(
	ctx context.Context,
	name string,
	ifUnused bool,
	noWait bool,
) error {

	l := logging.WithContext(ctx).With(
		zap.String("name", name),
		zap.Bool("if_unused", ifUnused),
		zap.Bool("noWait", noWait))

	if err := r.channel.ExchangeDelete(name, ifUnused, noWait); err != nil {
		l.Error("could not delete exchange")
		return err
	}

	return nil

}

func (r RabbitMqChannel) QueueDeclare(ctx context.Context, name string, durable bool, autoDelete bool, exclusive bool, noWait bool, args Args) error {

	l := logging.WithContext(ctx).With(
		zap.String("name", name),
		zap.Bool("durable", durable),
		zap.Bool("autoDelete", autoDelete),
		zap.Bool("exclusive", exclusive),
		zap.Bool("noWait", noWait),
		zap.Object("args", args))

	if _, err := r.channel.QueueDeclare(name, durable, autoDelete, exclusive, noWait, map[string]interface{}(args)); err != nil {
		l.Error("could not declare queue")
		return err
	}

	return nil

}

func (r RabbitMqChannel) QueueBind(ctx context.Context, name string, key string, exchange string, noWait bool, args Args) error {

	l := logging.WithContext(ctx).With(
		zap.String("name", name),
		zap.String("binding_key", key),
		zap.String("exchange", exchange),
		zap.Bool("noWait", noWait),
		zap.Object("args", args))

	if err := r.channel.QueueBind(name, key, exchange, noWait, map[string]interface{}(args)); err != nil {
		l.Error("could not bind queue")
		return err
	}

	return nil

}

func (r RabbitMqChannel) QueueUnbind(ctx context.Context, name string, key string, exchange string, args Args) error {

	l := logging.WithContext(ctx).With(
		zap.String("name", name),
		zap.String("binding_key", key),
		zap.String("exchange", exchange),
		zap.Object("args", args))

	if err := r.channel.QueueUnbind(name, key, exchange, map[string]interface{}(args)); err != nil {
		l.Error("could not unbind queue")
		return err
	}

	return nil

}

func (r RabbitMqChannel) QueueDelete(ctx context.Context, name string, ifUnused bool, ifEmpty bool, noWait bool) error {

	l := logging.WithContext(ctx).With(
		zap.String("name", name),
		zap.Bool("if_unused", ifUnused),
		zap.Bool("if_empty", ifEmpty),
		zap.Bool("no_wait", noWait))

	if _, err := r.channel.QueueDelete(name, ifUnused, ifEmpty, noWait); err != nil {
		l.Error("could not delete queue")
		return err
	}

	return nil

}
