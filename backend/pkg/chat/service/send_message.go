package service

import (
	"context"
	"github.com/commonpool/backend/pkg/chat"
	"github.com/commonpool/backend/pkg/mq"
)

func (c ChatService) SendMessage(ctx context.Context, message *chat.Message) error {

	err := c.chatStore.SaveMessage(ctx, message)

	amqpChannel, err := c.amqpClient.GetChannel()
	if err != nil {
		return err
	}
	defer amqpChannel.Close()

	publishing, err := message.AsAmqpMessage()
	if err != nil {
		return err
	}

	return amqpChannel.Publish(ctx, mq.MessagesExchange, "", false, false, *publishing)

}
