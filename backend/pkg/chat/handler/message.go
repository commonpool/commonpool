package handler

import (
	"github.com/commonpool/backend/pkg/chat"
	"time"
)

type Message struct {
	ID             string              `json:"id"`
	ChannelID      string              `json:"channelId"`
	MessageType    chat.MessageType    `json:"messageType"`
	MessageSubType chat.MessageSubType `json:"messageSubType"`
	SentById       string              `json:"sentById"`
	SentByUsername string              `json:"sentByUsername"`
	SentAt         time.Time           `json:"sentAt"`
	Text           string              `json:"text"`
	Blocks         []chat.Block        `json:"blocks"`
	Attachments    []chat.Attachment   `json:"attachments"`
	VisibleToUser  *string             `json:"visibleToUser"`
}

func MapMessage(message *chat.Message) *Message {
	var visibleToUser *string = nil
	if message.VisibleToUser != nil {
		visibleToUserStr := message.VisibleToUser.String()
		visibleToUser = &visibleToUserStr
	}
	return &Message{
		ID:             message.Key.String(),
		ChannelID:      message.ChannelKey.String(),
		MessageType:    message.MessageType,
		MessageSubType: message.MessageSubType,
		SentById:       message.SentBy.UserKey.String(),
		SentByUsername: message.SentBy.Username,
		SentAt:         message.SentAt,
		Text:           message.Text,
		Blocks:         message.Blocks,
		Attachments:    message.Attachments,
		VisibleToUser:  visibleToUser,
	}
}
