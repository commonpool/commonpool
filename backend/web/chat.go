package web

import "time"

type Thread struct {
	TopicID             string    `json:"topicId"`
	RecipientID         string    `json:"recipientId"`
	LastChars           string    `json:"lastChars"`
	HasUnreadMessages   bool      `json:"hasUnreadMessages"`
	LastMessageAt       time.Time `json:"lastMessageAt"`
	LastMessageUsername string    `json:"lastMessageUsername"`
	LastMessageUserId   string    `json:"lastMessageUserId"`
	Title               string    `json:"title"`
}

type Message struct {
	ID             string    `json:"id"`
	TopicID        string    `json:"topicId"`
	SentBy         string    `json:"sentBy"`
	SentByUsername string    `json:"sentByUsername"`
	SentByMe       bool      `json:"sentByMe"`
	SentAt         time.Time `json:"sentAt"`
	Content        string    `json:"content"`
}

type GetLatestThreadsResponse struct {
	Threads []Thread `json:"threads"`
}

type InquireAboutResourceRequest struct {
	Message string `json:"message"`
}

type SendMessageRequest struct {
	Message string `json:"message"`
}

type GetLatestMessageThreadsResponse struct {
	Messages []Message `json:"messages"`
}

type GetTopicMessagesResponse struct {
	Messages []Message `json:"messages"`
}
