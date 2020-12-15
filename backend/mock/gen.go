package mock

//go:generate moq -pkg mock -out auth_mock.go ../pkg/user Store:UserStore
//go:generate moq -pkg mock -out authenticator.go ../pkg/auth Authenticator
//go:generate moq -pkg mock -out resource_store_mock.go ../pkg/resource Store:ResourceStore
//go:generate moq -pkg mock -out chat_store.go ../pkg/chat Store:ChatStore
//go:generate moq -pkg mock -out chat_service.go ../pkg/chat Service:ChatService
//go:generate moq -pkg mock -out amqp_client_mock.go ../pkg/mq Client:AmqpClient
//go:generate moq -pkg mock -out amqp_channel_mock.go ../pkg/mq Channel:AmqpChannel
//go:generate moq -pkg mock -out amcp_ack_mock.go ../pkg/mq Ack:AmqpAck
