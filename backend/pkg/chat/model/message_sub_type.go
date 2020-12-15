package model

type MessageSubType string

const (
	UserMessage MessageSubType = "user"
	BotMessage  MessageSubType = "bot"
)
