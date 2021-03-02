package mock

//go:generate moq -pkg mock -out user_store.go ../pkg/user Store:UserStore
//go:generate moq -pkg mock -out user_service.go ../pkg/user Service:UserService
//go:generate moq -pkg mock -out authenticator.go ../pkg/auth Authenticator
//go:generate moq -pkg mock -out resource_service.go ../pkg/resource Service:ResourceService
//go:generate moq -pkg mock -out resource_store.go ../pkg/resource Store:ResourceStore
//go:generate moq -pkg mock -out trading_store.go ../pkg/trading Store:TradingStore
//go:generate moq -pkg mock -out trading_service.go ../pkg/trading Service:TradingService
//go:generate moq -pkg mock -out group_service.go ../pkg/group Service:GroupService
//go:generate moq -pkg mock -out group_store.go ../pkg/group Store:GroupStore
//go:generate moq -pkg mock -out chat_store.go ../pkg/chat Store:ChatStore
//go:generate moq -pkg mock -out chat_service.go ../pkg/chat Service:ChatService
//go:generate moq -pkg mock -out amqp_client.go ../pkg/mq Client:AmqpClient
//go:generate moq -pkg mock -out amqp_channel.go ../pkg/mq Channel:AmqpChannel
//go:generate moq -pkg mock -out amcp_ack.go ../pkg/mq Ack:AmqpAck
//go:generate moq -pkg mock -out offer_approvers.go ../pkg/trading Approvers:OfferApprovers
