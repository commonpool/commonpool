package model

type TopicKey struct {
	ID string
}

func (tk *TopicKey) String() string {
	return tk.ID
}

func NewTopicKey(key string) TopicKey {
	return TopicKey{
		ID: key,
	}
}
