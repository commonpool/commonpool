package model

type ThreadKey struct {
	TopicKey TopicKey
	UserKey  UserKey
}

func NewThreadKey(topic TopicKey, userKey UserKey) ThreadKey {
	return ThreadKey{
		TopicKey: topic,
		UserKey:  userKey,
	}
}
