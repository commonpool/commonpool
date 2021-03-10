package mq

import (
	"context"
	"github.com/commonpool/backend/logging"
	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
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
	return r.channel.Publish(exchange, key, mandatory, immediate, mapMessage(publishing))
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

	time.Sleep(2 * time.Second)

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
			amqpDelivery := mapDelivery(delivery)
			amqpChan <- amqpDelivery
		}
		close(amqpChan)
	}()

	return amqpChan, nil
}

func mapDelivery(delivery amqp.Delivery) Delivery {
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
	return amqpDelivery
}

func mapMessage(publishing Message) amqp.Publishing {
	return amqp.Publishing{
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
	}
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

//

type RabbitConnection struct {
	connection   *amqp.Connection
	config       RabbitConfig
	errorChannel chan *amqp.Error
	retryChannel chan struct{}
	closeChannel chan struct{}
	closed       bool
}

func NewRabbitConnection(rabbitConfig RabbitConfig) *RabbitConnection {
	var connection = &RabbitConnection{
		config:       rabbitConfig,
		errorChannel: make(chan *amqp.Error),
	}
	return connection
}

func (c *RabbitConnection) connect(ctx context.Context) (err error) {
	l := logging.WithContext(ctx).Named("RabbitQueue.connect")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			l.Debug("Connecting to RabbitMQ")
			conn, err := amqp.Dial(c.config.URL)
			if err == nil {
				c.connection = conn
				c.connection.NotifyClose(c.errorChannel)
				l.Debug("Connection established")
				return nil
			}
			l.Error("Connection to rabbitMq failed. Retrying in 1 sec...", zap.Error(err))
			time.Sleep(1 * time.Second)
		}
	}
}

func (q *RabbitConnection) reconnector(ctx context.Context) {
	ctx, l := logging.With(ctx)
	for {
		select {
		case err := <-q.errorChannel:
			if !q.closed {
				if err != nil {
					l.Warn("Reconnecting after connection closed", zap.Error(err))
				} else {
					l.Warn("Reconnecting after connection closed")
				}
				if err := q.connect(ctx); err != nil {
					q.retryChannel <- struct{}{}
				}
			} else {
				return
			}
		case <-ctx.Done():
			if q.closed {
				return
			}
			q.close()
		}
	}
}

func (c *RabbitConnection) close() error {
	if c.closed {
		return nil
	}
	if err := c.connection.Close(); err != nil {
		return err
	}
	c.closed = true
	c.closeChannel <- struct{}{}
	return nil
}

//

type RabbitConfig struct {
	URL string
}

func NewRabbitConfig(url string) RabbitConfig {
	return RabbitConfig{
		URL: url,
	}
}

type QueueConfig struct {
	Name       string
	Key        string
	Exchange   string
	Durable    bool
	NoWait     bool
	Args       Args
	AutoDelete bool
	Exclusive  bool
}

func NewQueueConfig() QueueConfig {
	return QueueConfig{}
}

func (q QueueConfig) WithName(name string) QueueConfig {
	var t = &q
	t.Name = name
	return *t
}
func (q QueueConfig) WithKey(key string) QueueConfig {
	var t = &q
	t.Key = key
	return *t
}
func (q QueueConfig) WithExchange(exchange string) QueueConfig {
	var t = &q
	t.Exchange = exchange
	return *t
}
func (q QueueConfig) WithDurable(durable bool) QueueConfig {
	var t = &q
	t.Durable = durable
	return *t
}
func (q QueueConfig) WithNoWait(noWait bool) QueueConfig {
	var t = &q
	t.NoWait = noWait
	return *t
}
func (q QueueConfig) WithAutoDelete(autoDelete bool) QueueConfig {
	var t = &q
	t.AutoDelete = autoDelete
	return *t
}
func (q QueueConfig) WithExclusive(exclusive bool) QueueConfig {
	var t = &q
	t.Exclusive = exclusive
	return *t
}
func (q QueueConfig) WithArgs(args Args) QueueConfig {
	var t = &q
	t.Args = args
	return *t
}

func (q QueueConfig) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("name", q.Name)
	encoder.AddString("key", q.Key)
	encoder.AddBool("durable", q.Durable)
	encoder.AddBool("auto_delete", q.AutoDelete)
	encoder.AddBool("exclusive", q.Exclusive)
	encoder.AddString("exchange", q.Exchange)
	encoder.AddBool("no_wait", q.NoWait)
	return encoder.AddObject("args", q.Args)
}

type RabbitQueue struct {
	c            *RabbitConnection
	RabbitConfig RabbitConfig
	Config       QueueConfig
	connection   *amqp.Connection
	errorChannel chan *amqp.Error
	channel      *amqp.Channel
	consumers    []ConsumerConfig
	consume      func(delivery Delivery)
	closed       bool
	l            *zap.Logger
}

func NewRabbitQueue(rabbitConfig RabbitConfig, queueConfig QueueConfig) *RabbitQueue {
	q := new(RabbitQueue)
	q.RabbitConfig = rabbitConfig
	q.Config = queueConfig
	q.c = NewRabbitConnection(rabbitConfig)
	return q
}

func (q *RabbitQueue) Connect(ctx context.Context) error {
	err := q.connect(logging.NewContext(ctx, zap.Object("config", q.Config)))
	if err == nil {
		go q.reconnector(ctx)
	}
	return err
}

func (q *RabbitQueue) connect(ctx context.Context) error {
	l := logging.WithContext(ctx).Named("RabbitQueue.connect")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			l.Debug("Connecting to RabbitMQ", zap.Object("config", q.Config))
			conn, err := amqp.Dial(q.RabbitConfig.URL)
			if err == nil {
				q.connection = conn
				q.errorChannel = make(chan *amqp.Error)
				q.connection.NotifyClose(q.errorChannel)
				l.Debug("Connection established", zap.Object("config", q.Config))

				if err = q.openChannel(ctx); err == nil {
					if err = q.declareQueue(ctx); err == nil {
						return nil
					}
				}
			}
			l.Error("Connection to rabbitMq failed. Retrying in 1 sec...", zap.Error(err))
			time.Sleep(1 * time.Second)
		}
	}
}

func (q *RabbitQueue) Consume(ctx context.Context, consumerConfig ConsumerConfig) error {
	ctx, _ = logging.With(ctx, zap.Object("queue_config", q.Config), zap.Object("consumer_config", consumerConfig))
	deliveries, err := q.registerConsumer(ctx, consumerConfig)
	if err == nil {
		q.executeMessageConsumer(consumerConfig, deliveries, false)
	}
	return nil
}

func (q *RabbitQueue) Send(ctx context.Context, message Message, mandatory bool, immediate bool) error {
	l := logging.WithContext(ctx).Named("RabbitQueue.Send")
	l.Debug("sending message")
	return q.channel.Publish(q.Config.Exchange, q.Config.Name, mandatory, immediate, mapMessage(message))
}

func (q *RabbitQueue) Close() {
	q.closed = true
	q.channel.Close()
	q.connection.Close()
}

func (q *RabbitQueue) openChannel(ctx context.Context) error {
	l := logging.WithContext(ctx).Named("RabbitQueue.openChannel")
	channel, err := q.connection.Channel()
	if err != nil {
		l.Error("could not open channel", zap.Error(err))
		return err
	}
	q.channel = channel
	return nil
}

func (q *RabbitQueue) declareQueue(ctx context.Context) error {
	l := logging.WithContext(ctx).Named("RabbitQueue.declareQueue")
	_, err := q.channel.QueueDeclare(q.Config.Name, q.Config.Durable, q.Config.AutoDelete, q.Config.Exclusive, q.Config.NoWait, map[string]interface{}(q.Config.Args))
	if err != nil {
		l.Error("could not declare queue", zap.Error(err))
	}
	return err
}

type ConsumerConfig struct {
	ConsumerName string
	AutoAck      bool
	Exclusive    bool
	NoLocal      bool
	NoWait       bool
	Args         Args
	Consume      func(d Delivery)
}

func NewConsumerConfig(consume func(d Delivery)) ConsumerConfig {
	return ConsumerConfig{
		Consume: consume,
	}
}

func (c ConsumerConfig) WithName(name string) ConsumerConfig {
	var r = &c
	r.ConsumerName = name
	return *r
}

func (c ConsumerConfig) WithAutoAck(autoAck bool) ConsumerConfig {
	var r = &c
	r.AutoAck = autoAck
	return *r
}

func (c ConsumerConfig) WithExclusive(exclusive bool) ConsumerConfig {
	var r = &c
	r.Exclusive = exclusive
	return *r
}
func (c ConsumerConfig) WithNoWait(noWait bool) ConsumerConfig {
	var r = &c
	r.NoWait = noWait
	return *r
}
func (c ConsumerConfig) WithArgs(args Args) ConsumerConfig {
	var r = &c
	r.Args = args
	return *r
}
func (c ConsumerConfig) WithConsume(consume func(d Delivery)) ConsumerConfig {
	var r = &c
	r.Consume = consume
	return *r
}

func (c ConsumerConfig) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("consumer_name", c.ConsumerName)
	encoder.AddBool("auto_ack", c.AutoAck)
	encoder.AddBool("exclusive", c.Exclusive)
	encoder.AddBool("no_local", c.NoLocal)
	encoder.AddBool("no_wait", c.NoWait)
	encoder.AddBool("no_wait", c.NoWait)
	return encoder.AddObject("args", c.Args)
}

func (q *RabbitQueue) registerConsumer(ctx context.Context, consumerConfig ConsumerConfig) (<-chan amqp.Delivery, error) {
	l := logging.WithContext(ctx).Named("RabbitQueue.registerConsumer")
	l.Debug("registering consumer")
	msgs, err := q.channel.Consume(
		q.Config.Name,
		consumerConfig.ConsumerName,
		consumerConfig.AutoAck,
		consumerConfig.Exclusive,
		consumerConfig.NoLocal,
		consumerConfig.NoWait,
		map[string]interface{}(consumerConfig.Args),
	)
	if err != nil {
		l.Error("could not register consumer", zap.Error(err))
	} else {
		l.Debug("successfully registered consumer")
	}
	return msgs, err
}

func (q *RabbitQueue) recoverConsumers(ctx context.Context) {
	ctx, l := logging.With(ctx)
	l = l.Named("RabbitQueue.recoverConsumers")

	for i := range q.consumers {
		consumer := q.consumers[i]
		l = l.With(zap.Object("consumer_config", consumer))
		l.Debug("recovering consumer")
		msgs, err := q.registerConsumer(ctx, consumer)
		if err != nil {
			q.executeMessageConsumer(consumer, msgs, true)
		}
	}

}

func (q *RabbitQueue) executeMessageConsumer(consumerConfig ConsumerConfig, deliveries <-chan amqp.Delivery, isRecovery bool) {
	if !isRecovery {
		q.consumers = append(q.consumers, consumerConfig)
	}
	go func() {
		for delivery := range deliveries {
			d := mapDelivery(delivery)
			consumerConfig.Consume(d)
		}
	}()
}

func (q *RabbitQueue) reconnector(ctx context.Context) {
	ctx, l := logging.With(ctx)
	for {
		select {
		case err := <-q.errorChannel:
			if !q.closed {
				if err != nil {
					l.Warn("Reconnecting after connection closed", zap.Error(err))
				} else {
					l.Warn("Reconnecting after connection closed")
				}

				if err := q.connect(ctx); err != nil {
					q.recoverConsumers(ctx)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
