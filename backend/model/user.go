package model

import (
	"fmt"
	"github.com/commonpool/backend/pkg/utils"
	"go.uber.org/zap/zapcore"
	"sort"
	"strings"
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

func (k UserKey) String() string {
	return k.subject
}

func (k UserKey) GetExchangeName() string {
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

func (k *UserKeys) Append(key UserKey) *UserKeys {
	newUserKeys := append(k.Items, key)
	return NewUserKeys(newUserKeys)
}

func (k *UserKeys) IsEmpty() bool {
	return k.Items == nil || len(k.Items) == 0
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

func (k *UserKeys) GetChannelKey() (ChannelKey, error) {

	if k == nil || len(k.Items) == 0 {
		err := fmt.Errorf("cannot get conversation channel for 0 participants")
		return ChannelKey{}, err
	}

	var shortUids []string
	for _, userKey := range k.Items {
		sid, err := utils.ShortUuidFromStr(userKey.String())
		if err != nil {
			return ChannelKey{}, err
		}
		shortUids = append(shortUids, sid)
	}
	sort.Strings(shortUids)
	channelId := strings.Join(shortUids, "-")
	channelKey := NewConversationKey(channelId)

	return channelKey, nil

}

type UserReference interface {
	GetUserKey() UserKey
	GetUsername() string
}

func NewEmptyUserKeys() *UserKeys {
	return &UserKeys{
		Items: []UserKey{},
	}
}
