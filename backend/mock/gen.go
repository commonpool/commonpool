package mock

//go:generate moq -pkg mock -out auth_mock.go ../auth Store:AuthStore
//go:generate moq -pkg mock -out authenticator.go ../auth Authenticator
//go:generate moq -pkg mock -out resource_store_mock.go ../resource Store:ResourceStore
//go:generate moq -pkg mock -out chat_store.go ../pkg/chat Store:ChatStore
//go:generate moq -pkg mock -out chat_service.go ../pkg/chat Service:ChatService
//go:generate moq -pkg mock -out amqp_client_mock.go ../amqp Client:AmqpClient
//go:generate moq -pkg mock -out amqp_channel_mock.go ../amqp Channel:AmqpChannel
//go:generate moq -pkg mock -out amcp_ack_mock.go ../amqp Ack:AmqpAck
