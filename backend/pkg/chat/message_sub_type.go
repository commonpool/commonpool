package chat

type MessageSubType string

const (
	UserMessage MessageSubType = "user"
	BotMessage  MessageSubType = "bot"
)
