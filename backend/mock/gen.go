package mock

//go:generate moq -pkg mock -out auth_mock.go ../resource Store
//go:generate moq -pkg mock -out resource_store_mock.go ../resource Store
//go:generate moq -pkg mock -out amqp_client_mock.go ../amqp AmqpClient
//go:generate moq -pkg mock -out amqp_channel_mock.go ../amqp AmqpChannel
//go:generate moq -pkg mock -out amcp_ack_mock.go ../amqp AmqpAcknowledger
