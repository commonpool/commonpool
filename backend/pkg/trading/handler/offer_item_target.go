package handler

import (
	"fmt"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource"
)

type OfferItemTarget struct {
	UserID  *string             `json:"userId"`
	GroupID *string             `json:"groupId" validatde:"uuid"`
	Type    resource.TargetType `json:"type"`
}

func MapWebOfferItemTarget(target OfferItemTarget) (*resource.Target, error) {
	if target.Type == resource.UserTarget {
		userKey := keys.NewUserKey(*target.UserID)
		return &resource.Target{
			UserKey:  &userKey,
			GroupKey: nil,
			Type:     resource.UserTarget,
		}, nil
	} else if target.Type == resource.GroupTarget {
		groupKey, err := keys.ParseGroupKey(*target.GroupID)
		if err != nil {
			return nil, err
		}
		return &resource.Target{
			UserKey:  nil,
			GroupKey: &groupKey,
			Type:     resource.GroupTarget,
		}, nil
	}
	return nil, fmt.Errorf("invalid target")
}

func MapOfferItemTarget(target *resource.Target) (*OfferItemTarget, error) {

	if target == nil {
		return nil, nil
	}
	if target.IsForGroup() {
		groupId := target.GetGroupKey().String()
		return &OfferItemTarget{
			UserID:  nil,
			GroupID: &groupId,
			Type:    resource.GroupTarget,
		}, nil

	} else if target.IsForUser() {
		userId := target.GetUserKey().String()
		return &OfferItemTarget{
			UserID:  &userId,
			GroupID: nil,
			Type:    resource.UserTarget,
		}, nil
	} else {
		return nil, fmt.Errorf("unexpected offer item type")
	}

}

func (t OfferItemTarget) Parse() (*resource.Target, error) {
	if t.Type == resource.GroupTarget {
		groupKey, err := keys.ParseGroupKey(*t.GroupID)
		if err != nil {
			return nil, err
		}
		return &resource.Target{
			UserKey:  nil,
			GroupKey: &groupKey,
			Type:     resource.GroupTarget,
		}, nil
	} else if t.Type == resource.UserTarget {
		userKey := keys.NewUserKey(*t.UserID)
		return &resource.Target{
			UserKey:  &userKey,
			GroupKey: nil,
			Type:     resource.UserTarget,
		}, nil
	}
	return nil, fmt.Errorf("unexpected target type: %s", t.Type)
}

func NewWebOfferItemTarget(offerItemTarget *resource.Target) *OfferItemTarget {

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
		Type:    resource.GroupTarget,
	}
}

func NewUserTarget(user string) *OfferItemTarget {
	return &OfferItemTarget{
		UserID:  &user,
		GroupID: nil,
		Type:    resource.UserTarget,
	}
}
