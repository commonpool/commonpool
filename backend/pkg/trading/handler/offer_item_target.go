package handler

import (
	"fmt"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
)

type OfferItemTarget struct {
	UserID  *string            `json:"userId"`
	GroupID *string            `json:"groupId" validatde:"uuid"`
	Type    trading.TargetType `json:"type"`
}

func MapWebOfferItemTarget(target OfferItemTarget) (*trading.Target, error) {
	if target.Type == trading.UserTarget {
		userKey := keys.NewUserKey(*target.UserID)
		return &trading.Target{
			UserKey:  &userKey,
			GroupKey: nil,
			Type:     trading.UserTarget,
		}, nil
	} else if target.Type == trading.GroupTarget {
		groupKey, err := keys.ParseGroupKey(*target.GroupID)
		if err != nil {
			return nil, err
		}
		return &trading.Target{
			UserKey:  nil,
			GroupKey: &groupKey,
			Type:     trading.GroupTarget,
		}, nil
	}
	return nil, fmt.Errorf("invalid target")
}

func MapOfferItemTarget(target *trading.Target) (*OfferItemTarget, error) {

	if target == nil {
		return nil, nil
	}
	if target.IsForGroup() {
		groupId := target.GetGroupKey().String()
		return &OfferItemTarget{
			UserID:  nil,
			GroupID: &groupId,
			Type:    trading.GroupTarget,
		}, nil

	} else if target.IsForUser() {
		userId := target.GetUserKey().String()
		return &OfferItemTarget{
			UserID:  &userId,
			GroupID: nil,
			Type:    trading.UserTarget,
		}, nil
	} else {
		return nil, fmt.Errorf("unexpected offer item type")
	}

}

func (t OfferItemTarget) Parse() (*trading.Target, error) {
	if t.Type == trading.GroupTarget {
		groupKey, err := keys.ParseGroupKey(*t.GroupID)
		if err != nil {
			return nil, err
		}
		return &trading.Target{
			UserKey:  nil,
			GroupKey: &groupKey,
			Type:     trading.GroupTarget,
		}, nil
	} else if t.Type == trading.UserTarget {
		userKey := keys.NewUserKey(*t.UserID)
		return &trading.Target{
			UserKey:  &userKey,
			GroupKey: nil,
			Type:     trading.UserTarget,
		}, nil
	}
	return nil, fmt.Errorf("unexpected target type: %s", t.Type)
}

func NewWebOfferItemTarget(offerItemTarget *trading.Target) *OfferItemTarget {

	var userId *string = nil
	var groupId *string = nil

	if offerItemTarget.IsForGroup() {
		groupIdStr := offerItemTarget.GroupKey.String()
		groupId = &groupIdStr
	} else if offerItemTarget.IsForUser() {
		userIdStr := offerItemTarget.UserKey.String()
		userId = &userIdStr
	}

	return &OfferItemTarget{
		UserID:  userId,
		GroupID: groupId,
		Type:    offerItemTarget.Type,
	}

}

func NewGroupTarget(group string) *OfferItemTarget {
	return &OfferItemTarget{
		UserID:  nil,
		GroupID: &group,
		Type:    trading.GroupTarget,
	}
}

func NewUserTarget(user string) *OfferItemTarget {
	return &OfferItemTarget{
		UserID:  &user,
		GroupID: nil,
		Type:    trading.UserTarget,
	}
}
