package chat

import (
	"encoding/json"
	"github.com/commonpool/backend/pkg/mq"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
	"time"
)

type Message struct {
	Key            MessageKey
	ChannelKey     ChannelKey
	MessageType    MessageType
	MessageSubType MessageSubType
	SentBy         MessageSender
	SentAt         time.Time
	Text           string
	Blocks         []Block
	Attachments    []Attachment
	VisibleToUser  *usermodel.UserKey
}

func (m *Message) AsAmqpMessage() (*mq.Message, error) {

	evt := mq.Event{
		Type:      mq.NewChatMessage,
		SubType:   mq.UserMessage,
		Channel:   m.ChannelKey.String(),
		User:      m.SentBy.UserKey.String(),
		ID:        m.Key.String(),
		Timestamp: m.SentAt.String(),
		Text:      m.Text,
	}

	js, err := json.Marshal(evt)
	if err != nil {
		return nil, err
	}

	headers := mq.NewArgs().
		With(mq.ChannelIdArg, m.ChannelKey.String()).
		WithEventType(mq.NewChatMessage)

	return mq.NewMessage().
			WithType(string(evt.Type)).
			WithHeaders(headers).
			WithJsonBody(string(js)).
			WithTimestamp(m.SentAt).
			WithMessageId(m.Key.String()),
		nil

}
