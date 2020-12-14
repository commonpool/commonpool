package web

import (
	"fmt"
	"github.com/commonpool/backend/model"
	trading2 "github.com/commonpool/backend/pkg/trading"
	"time"
)

type Offer struct {
	ID             string               `json:"id"`
	CreatedAt      time.Time            `json:"createdAt"`
	CompletedAt    *time.Time           `json:"completedAt"`
	Status         trading2.OfferStatus `json:"status"`
	AuthorID       string               `json:"authorId"`
	AuthorUsername string               `json:"authorUsername"`
	Items          []*OfferItem         `json:"items"`
	Message        string               `json:"message"`
}

type OfferItem struct {
	ID                          string                 `json:"id"`
	From                        *OfferItemTarget       `json:"from"`
	To                          *OfferItemTarget       `json:"to"`
	Type                        trading2.OfferItemType `json:"type"`
	ResourceId                  *string                `json:"resourceId"`
	Duration                    *int64                 `json:"duration"`
	Amount                      *int64                 `json:"amount"`
	ReceiverApproved            bool                   `json:"receiverApproved"`
	GiverApproved               bool                   `json:"giverApproved"`
	ReceivingApprovers          []string               `json:"receivingApprovers"`
	GivingApprovers             []string               `json:"givingApprovers"`
	ServiceGivenConfirmation    bool                   `json:"serviceGivenConfirmation"`
	ServiceReceivedConfirmation bool                   `json:"serviceReceivedConfirmation"`
	ItemTaken                   bool                   `json:"itemTaken"`
	ItemGiven                   bool                   `json:"itemGiven"`
	ItemReturnedBack            bool                   `json:"itemReturnedBack"`
	ItemReceivedBack            bool                   `json:"itemReceivedBack"`
}

type GetOfferResponse struct {
	Offer *Offer `json:"offer"`
}
type GetOffersResponse struct {
	Offers []Offer `json:"offers"`
}

type SendOfferRequest struct {
	Offer SendOfferPayload `json:"offer" validate:"required"`
}

type SendOfferPayload struct {
	Items   []SendOfferPayloadItem `json:"items" validate:"min=1"`
	GroupID string                 `json:"groupId" validate:"uuid"`
	Message string                 `json:"message"`
}

type OfferItemTarget struct {
	UserID  *string          `json:"userId"`
	GroupID *string          `json:"groupId" validatde:"uuid"`
	Type    model.TargetType `json:"type"`
}

func MapWebOfferItemTarget(target OfferItemTarget) (*model.Target, error) {
	if target.Type == model.UserTarget {
		userKey := model.NewUserKey(*target.UserID)
		return &model.Target{
			UserKey:  &userKey,
			GroupKey: nil,
			Type:     model.UserTarget,
		}, nil
	} else if target.Type == model.GroupTarget {
		groupKey, err := model.ParseGroupKey(*target.GroupID)
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
		groupKey, err := model.ParseGroupKey(*t.GroupID)
		if err != nil {
			return nil, err
		}
		return &model.Target{
			UserKey:  nil,
			GroupKey: &groupKey,
			Type:     model.GroupTarget,
		}, nil
	} else if t.Type == model.UserTarget {
		userKey := model.NewUserKey(*t.UserID)
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

type SendOfferPayloadItem struct {
	Type       trading2.OfferItemType `json:"type"`
	To         OfferItemTarget        `json:"to" validate:"required,uuid"`
	From       *OfferItemTarget       `json:"from" validate:"required,uuid"`
	ResourceId *string                `json:"resourceId" validate:"required,uuid"`
	Duration   *string                `json:"duration"`
	Amount     *string                `json:"amount"`
}

func NewResourceTransferItem(to *OfferItemTarget, resourceId string) *SendOfferPayloadItem {
	return &SendOfferPayloadItem{
		To:         *to,
		Type:       trading2.ResourceTransfer,
		ResourceId: &resourceId,
	}
}

func NewCreditTransferItem(from *OfferItemTarget, to *OfferItemTarget, time time.Duration) *SendOfferPayloadItem {
	seconds := time.String()
	return &SendOfferPayloadItem{
		From:   from,
		To:     *to,
		Type:   trading2.CreditTransfer,
		Amount: &seconds,
	}
}
