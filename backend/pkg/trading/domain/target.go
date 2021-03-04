package domain

import "github.com/commonpool/backend/pkg/keys"

type OfferItemTargetType string

const (
	GroupTarget OfferItemTargetType = "group"
	UserTarget  OfferItemTargetType = "user"
)

type OfferItemTarget struct {
	Type     OfferItemTargetType `json:"type"`
	UserKey  *keys.UserKey       `json:"user_key,omitempty"`
	GroupKey *keys.GroupKey      `json:"group_key,omitempty"`
}

func NewGroupTarget(groupKey keys.GroupKey) *OfferItemTarget {
	return &OfferItemTarget{
		Type:     GroupTarget,
		GroupKey: &groupKey,
	}
}

func NewUserTarget(userKey keys.UserKey) *OfferItemTarget {
	return &OfferItemTarget{
		Type:    UserTarget,
		UserKey: &userKey,
	}
}
