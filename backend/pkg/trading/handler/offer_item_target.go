package handler

import (
	"fmt"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
)

type OfferItemTarget struct {
	UserID  *string           `json:"userId"`
	GroupID *string           `json:"groupId" validatde:"uuid"`
	Type    domain.TargetType `json:"type"`
}

func MapWebOfferItemTarget(target OfferItemTarget) (*domain.Target, error) {
	if target.Type == domain.UserTarget {
		userKey := keys.NewUserKey(*target.UserID)
		return &domain.Target{
			UserKey:  &userKey,
			GroupKey: nil,
			Type:     domain.UserTarget,
		}, nil
	} else if target.Type == domain.GroupTarget {
		groupKey, err := keys.ParseGroupKey(*target.GroupID)
		if err != nil {
			return nil, err
		}
		return &domain.Target{
			UserKey:  nil,
			GroupKey: &groupKey,
			Type:     domain.GroupTarget,
		}, nil
	}
	return nil, fmt.Errorf("invalid target")
}

func MapOfferItemTarget(targetType, targetKey string) (*OfferItemTarget, error) {
	if targetType == string(domain.GroupTarget) {
		return &OfferItemTarget{
			UserID:  nil,
			GroupID: &targetKey,
			Type:    domain.GroupTarget,
		}, nil

	} else if targetType == string(domain.UserTarget) {
		return &OfferItemTarget{
			UserID:  &targetKey,
			GroupID: nil,
			Type:    domain.UserTarget,
		}, nil
	} else {
		return nil, fmt.Errorf("unexpected offer item type")
	}

}

func (t OfferItemTarget) Parse() (*domain.Target, error) {
	if t.Type == domain.GroupTarget {
		groupKey, err := keys.ParseGroupKey(*t.GroupID)
		if err != nil {
			return nil, err
		}
		return &domain.Target{
			UserKey:  nil,
			GroupKey: &groupKey,
			Type:     domain.GroupTarget,
		}, nil
	} else if t.Type == domain.UserTarget {
		userKey := keys.NewUserKey(*t.UserID)
		return &domain.Target{
			UserKey:  &userKey,
			GroupKey: nil,
			Type:     domain.UserTarget,
		}, nil
	}
	return nil, fmt.Errorf("unexpected target type: %s", t.Type)
}

func NewWebOfferItemTarget(offerItemTarget *domain.Target) *OfferItemTarget {

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
		Type:    domain.GroupTarget,
	}
}

func NewUserTarget(user string) *OfferItemTarget {
	return &OfferItemTarget{
		UserID:  &user,
		GroupID: nil,
		Type:    domain.UserTarget,
	}
}
