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
