package mock

//go:generate moq -pkg mock -out amqp_client/amqp_client.go ../pkg/mq Client:AmqpClient
//go:generate moq -pkg mock -out amqp_channel/amqp_channel.go ../pkg/mq Channel:AmqpChannel
//go:generate moq -pkg mock -out amqp_ack/amcp_ack.go ../pkg/mq Ack:AmqpAck
//go:generate moq -pkg mock -out chat_store/chat_store.go ../pkg/chat/store Store:ChatStore
//go:generate moq -pkg mock -out chat_service/chat_service.go ../pkg/chat/service Service:ChatService
//go:generate moq -pkg mock -out authenticator/authenticator.go ../pkg/auth/authenticator Authenticator
