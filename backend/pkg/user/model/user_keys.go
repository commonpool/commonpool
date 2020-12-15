package model

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

func NewEmptyUserKeys() *UserKeys {
	return &UserKeys{
		Items: []UserKey{},
	}
}