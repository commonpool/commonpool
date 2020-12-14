package amqp

import (
	"github.com/commonpool/backend/model"
	"go.uber.org/zap/zapcore"
	"time"
)

type Message struct {
	// Application or exchange specific fields,
	// the headers exchange will inspect this field.
	Headers Args

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

func NewMessage() *Message {
	return &Message{}
}

func (p *Message) WithContentType(contentType string) *Message {
	p.ContentType = contentType
	return p
}

func (p *Message) WithType(publishingType string) *Message {
	p.Type = publishingType
	return p
}

func (p *Message) WithJsonContentType() *Message {
	return p.WithContentType("application/json")
}

func (p *Message) WithHeaders(headers Args) *Message {
	p.Headers = headers
	return p
}

func (p *Message) WithTimestamp(timestamp time.Time) *Message {
	p.Timestamp = timestamp
	return p
}

func (p *Message) WithMessageId(id string) *Message {
	p.MessageId = id
	return p
}

func (p *Message) WithBody(body []byte) *Message {
	p.Body = body
	return p
}

func (p *Message) WithJsonBody(json string) *Message {
	return p.WithJsonContentType().WithBody([]byte(json))
}

type EventType string

const (
	NewChatMessage EventType = "chat.message"
)

type EventSubType string

const (
	UserMessage EventSubType = "user"
)

type Event struct {
	Channel   string       `json:"channel"`
	ID        string       `json:"id"`
	SubType   EventSubType `json:"subType"`
	Text      string       `json:"text"`
	Timestamp string       `json:"timestamp"`
	Type      EventType    `json:"type"`
	User      string       `json:"user"`
}

type EventContainer struct {
	Key   string
	Event Event
}

type Delivery struct {
	Acknowledger    Ack       // the channel from which this delivery arrived
	Headers         Args      // Application or header exchange table
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

type Args map[string]interface{}

func NewArgs() Args {
	return Args{}
}

func (a Args) WithChannelKey(channel model.ChannelKey) Args {
	a["channel_id"] = channel.String()
	return a
}

func (a Args) WithEventType(eventType EventType) Args {
	a["event_type"] = eventType
	return a
}

func (a Args) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
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

var _ zapcore.ObjectMarshaler = Args{}
