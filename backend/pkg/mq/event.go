package mq

type Event struct {
	Channel   string       `json:"channel"`
	ID        string       `json:"id"`
	SubType   EventSubType `json:"subType"`
	Text      string       `json:"text"`
	Timestamp string       `json:"timestamp"`
	Type      EventType    `json:"type"`
	User      string       `json:"user"`
}
