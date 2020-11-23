package model

type UserKey struct {
	subject string
}

func NewUserKey(subject string) UserKey {
	return UserKey{subject: subject}
}

func (k *UserKey) String() string {
	return k.subject
}

type UserKeys struct {
	Items []UserKey
}

func NewUserKeys(userKeys []UserKey) UserKeys {
	if userKeys == nil {
		userKeys = []UserKey{}
	}
	return UserKeys{
		Items: userKeys,
	}
}