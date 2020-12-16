package handler

import (
	"fmt"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/resource/model"
	"github.com/commonpool/backend/pkg/user/usermodel"
)

type OfferItemTarget struct {
	UserID  *string          `json:"userId"`
	GroupID *string          `json:"groupId" validatde:"uuid"`
	Type    model.TargetType `json:"type"`
}

func MapWebOfferItemTarget(target OfferItemTarget) (*model.Target, error) {
	if target.Type == model.UserTarget {
		userKey := usermodel.NewUserKey(*target.UserID)
		return &model.Target{
			UserKey:  &userKey,
			GroupKey: nil,
			Type:     model.UserTarget,
		}, nil
	} else if target.Type == model.GroupTarget {
		groupKey, err := group.ParseGroupKey(*target.GroupID)
		if err != nil {
			return nil, err
		}
		return &model.Target{
			UserKey:  nil,
			GroupKey: &groupKey,
			Type:     model.GroupTarget,
		}, nil
	}
	return nil, fmt.Errorf("invalid target")
}

func MapOfferItemTarget(target *model.Target) (*OfferItemTarget, error) {

	if target == nil {
		return nil, nil
	}
	if target.IsForGroup() {
		groupId := target.GetGroupKey().String()
		return &OfferItemTarget{
			UserID:  nil,
			GroupID: &groupId,
			Type:    model.GroupTarget,
		}, nil

	} else if target.IsForUser() {
		userId := target.GetUserKey().String()
		return &OfferItemTarget{
			UserID:  &userId,
			GroupID: nil,
			Type:    model.UserTarget,
		}, nil
	} else {
		return nil, fmt.Errorf("unexpected offer item type")
	}

}

func (t OfferItemTarget) Parse() (*model.Target, error) {
	if t.Type == model.GroupTarget {
		groupKey, err := group.ParseGroupKey(*t.GroupID)
		if err != nil {
			return nil, err
		}
		return &model.Target{
			UserKey:  nil,
			GroupKey: &groupKey,
			Type:     model.GroupTarget,
		}, nil
	} else if t.Type == model.UserTarget {
		userKey := usermodel.NewUserKey(*t.UserID)
		return &model.Target{
			UserKey:  &userKey,
			GroupKey: nil,
			Type:     model.UserTarget,
		}, nil
	}
	return nil, fmt.Errorf("unexpected target type: %s", t.Type)
}

func NewWebOfferItemTarget(offerItemTarget *model.Target) *OfferItemTarget {

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
		Type:    model.GroupTarget,
	}
}

func NewUserTarget(user string) *OfferItemTarget {
	return &OfferItemTarget{
		UserID:  &user,
		GroupID: nil,
		Type:    model.UserTarget,
	}
}
