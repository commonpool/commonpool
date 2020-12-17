package service

import (
	"context"
	"encoding/json"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/mq"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func (s *serviceTestSuite) TestSendMessage() {
	ctx := context.TODO()

	auth.SetContextAuthenticatedUser(ctx, "username", "user", "user@email.com")
	channelKey := keys.NewChannelKey("channel-id")
	messageKey := keys.NewMessageKey(uuid.FromStringOrNil("1370bb5e-4310-4d79-95f7-3923ba3f552a"))
	timestamp := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	message := &chat.Message{
		Key:            messageKey,
		ChannelKey:     channelKey,
		MessageType:    chat.NormalMessage,
		MessageSubType: chat.UserMessage,
		SentBy: chat.MessageSender{
			Type:     chat.UserMessageSender,
			UserKey:  keys.NewUserKey("user"),
			Username: "username",
		},
		SentAt:        timestamp,
		Text:          "Hello",
		Blocks:        []chat.Block{},
		Attachments:   []chat.Attachment{},
		VisibleToUser: nil,
	}

	err := s.Service.SendMessage(ctx, message)
	if !assert.NoError(s.T(), err) {
		return
	}

	assert.Len(s.T(), s.ChatStore.SaveMessageCalls(), 1)
	assert.Len(s.T(), s.AmqpChannel.PublishCalls(), 1)
	assert.Equal(s.T(), "Hello", message.Text)
	assert.Equal(s.T(), mq.MessagesExchange, s.AmqpChannel.PublishCalls()[0].Exchange)

	expectedEvent := mq.Message{
		Headers: mq.NewArgs().
			With(mq.ChannelIdArg, channelKey.String()).
			WithEventType(mq.NewChatMessage),
		ContentType:     "application/json",
		ContentEncoding: "",
		DeliveryMode:    0,
		Priority:        0,
		CorrelationId:   "",
		ReplyTo:         "",
		Expiration:      "",
		MessageId:       messageKey.String(),
		Timestamp:       timestamp,
		Type:            mq.EventTypeMessage,
		UserId:          "",
		AppId:           "",
		Body: []byte(CompactJson(s.T(), `{
		"channel"   : "channel-id",
		"type"      : "chat.message",
		"subType"   : "user",
		"user"      : "user",
		"id"        : "1370bb5e-4310-4d79-95f7-3923ba3f552a",
		"timestamp" : "2020-01-01 00:00:00 +0000 UTC",
		"text"      : "Hello"
	}`)),
	}
	assert.Equal(s.T(), expectedEvent, s.AmqpChannel.PublishCalls()[0].Publishing)
}

func CompactJson(t *testing.T, js string) string {
	var o2 interface{}
	var err error
	err = json.Unmarshal([]byte(js), &o2)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := json.Marshal(o2)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}
