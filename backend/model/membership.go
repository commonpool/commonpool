package model

import "go.uber.org/zap/zapcore"

type MembershipKey struct {
	UserKey  UserKey
	GroupKey GroupKey
}

func (m MembershipKey) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	err := encoder.AddObject("user", m.UserKey)
	if err != nil {
		return err
	}
	return encoder.AddObject("group", m.GroupKey)
}

var _ zapcore.ObjectMarshaler = MembershipKey{}

func NewMembershipKey(groupKey GroupKey, userKey UserKey) MembershipKey {
	return MembershipKey{
		UserKey:  userKey,
		GroupKey: groupKey,
	}
}
