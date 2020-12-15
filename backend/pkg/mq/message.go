package mq

import "time"

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
