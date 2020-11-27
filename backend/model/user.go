package model

import (
	"go.uber.org/zap/zapcore"
)

type UserKey struct {
	subject string
}

func (k UserKey) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("user_id", k.subject)
	return nil
}

var _ zapcore.ObjectMarshaler = &UserKey{}

func NewUserKey(subject string) UserKey {
	return UserKey{subject: subject}
}

func (k *UserKey) String() string {
	return k.subject
}

func (k *UserKey) GetExchangeName() string {
	return "users." + k.String()
}

type UserKeys struct {
	Items []UserKey
}

func (k *UserKeys) Contains(key UserKey) bool {
	if k.Items == nil {
		return false
	}
	for _, userKey := range k.Items {
		if userKey == key {
			return true
		}
	}
	return false
}

func NewUserKeys(userKeys []UserKey) *UserKeys {
	if userKeys == nil {
		userKeys = []UserKey{}
	}

	var newUserKeys []UserKey
	userKeyMap := map[UserKey]bool{}
	for _, key := range userKeys {
		if _, ok := userKeyMap[key]; ok {
			continue
		}
		userKeyMap[key] = true
		newUserKeys = append(newUserKeys, key)
	}

	return &UserKeys{
		Items: newUserKeys,
	}
}

func (k *UserKeys) Strings() []string {
	strs := make([]string, len(k.Items))
	for i := range k.Items {
		strs[i] = k.Items[i].String()
	}
	return strs
}

type UserReference interface {
	GetUserKey() UserKey
	GetUsername() string
}
