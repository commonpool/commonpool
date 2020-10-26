package model

type MembershipKey struct {
	UserKey  UserKey
	GroupKey GroupKey
}

func NewMembershipKey(groupKey GroupKey, userKey UserKey) MembershipKey {
	return MembershipKey{
		UserKey:  userKey,
		GroupKey: groupKey,
	}
}
