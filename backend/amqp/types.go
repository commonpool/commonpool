package amqp

import (
	"go.uber.org/zap/zapcore"
	"time"
)

type AmqpPublishing struct {
	// Application or exchange specific fields,
	// the headers exchange will inspect this field.
	Headers AmqpArgs

	// Properties
	ContentType     string    // MIME content type
	ContentEncoding string    // MIME content encoding
	DeliveryMode    uint8     // Transient (0 or 1) or Persistent (2)
	Priority        uint8     // 0 to 9
	CorrelationId   string    // correlation identifier
	ReplyTo         string    // address to to reply to (ex: RPC)
	Expiration      string    // message expiration spec
	MessageId       string    // message identifier
	Timestamp       time.Time // message timestamp
	Type            string    // message type name
	UserId          string    // creating user id - ex: "guest"
	AppId           string    // creating application id

	// The application specific payload of the message
	Body []byte
}

type EventType string

const (
	MessageEvent = "message"
)

type EventSubType string

const (
	MessageChanged = "message_changed"
	MessageDeleted = "message_deleted"
)

type Event struct {
	Type      EventType    `json:"type"`
	SubType   EventSubType `json:"subType"`
	Channel   string       `json:"channel"`
	User      string       `json:"user"`
	ID        string       `json:"id"`
	Timestamp string       `json:"timestamp"`
	Text      string       `json:"text"`
}

type EventContainer struct {
	Key   string
	Event Event
}

type AmqpAcknowledger interface {
	Ack(tag uint64, multiple bool) error
	Nack(tag uint64, multiple bool, requeue bool) error
	Reject(tag uint64, requeue bool) error
}

type AmqpDelivery struct {
	Acknowledger AmqpAcknowledger // the channel from which this delivery arrived
	Headers      AmqpArgs         // Application or header exchange table
	// Properties
	ContentType     string    // MIME content type
	ContentEncoding string    // MIME content encoding
	DeliveryMode    uint8     // queue implementation use - non-persistent (1) or persistent (2)
	Priority        uint8     // queue implementation use - 0 to 9
	CorrelationId   string    // application use - correlation identifier
	ReplyTo         string    // application use - address to reply to (ex: RPC)
	Expiration      string    // implementation use - message expiration spec
	MessageId       string    // application use - message identifier
	Timestamp       time.Time // application use - message timestamp
	Type            string    // application use - message type name
	UserId          string    // application use - creating user - should be authenticated user
	AppId           string    // application use - creating application id

	// Valid only with Channel.Consume
	ConsumerTag string

	// Valid only with Channel.Get
	MessageCount uint32

	DeliveryTag uint64
	Redelivered bool
	Exchange    string // basic.publish exchange
	RoutingKey  string // basic.publish routing key

	Body []byte
}

type AmqpArgs map[string]interface{}

func (a AmqpArgs) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	if a == nil {
		return nil
	}

	for s, i := range a {
		err := encoder.AddReflected(s, i)
		if err != nil {
			return err
		}
	}

	return nil
}

var _ zapcore.ObjectMarshaler = AmqpArgs{}
