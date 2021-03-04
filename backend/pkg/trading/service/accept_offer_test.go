package service

import (
	"context"
	"github.com/commonpool/backend/mock"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/pkg/user"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"time"
)

func (t *tradingTestSuite) TestAcceptOffer() {

	offerKey := keys.NewOfferKey(uuid.NewV4())
	groupKey := keys.NewGroupKey(uuid.NewV4())
	userKey := keys.NewUserKey("user")
	userTarget := trading.NewUserTarget(userKey)
	groupTarget := trading.NewGroupTarget(groupKey)
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	offer := &trading.Offer{
		Key:      offerKey,
		GroupKey: groupKey,
	}
	offerItems := trading.NewOfferItems([]trading.OfferItem{
		trading.NewCreditTransferItem(offerKey, offerItemKey, userTarget, groupTarget, time.Hour, trading.NewCreditTransferItemOptions{
			GiverAccepted:    false,
			ReceiverAccepted: false,
		}),
	})

	approvers := &mock.OfferApprovers{
		HasAnyOfferItemsToApproveFunc: func(userKey keys.UserKey) bool {
			return true
		},
		GetOutboundOfferItemsFunc: func(userKey keys.UserKey) *keys.OfferItemKeys {
			return keys.NewOfferItemKeys([]keys.OfferItemKey{offerItemKey})
		},
		GetInboundOfferItemsFunc: func(userKey keys.UserKey) *keys.OfferItemKeys {
			return keys.NewOfferItemKeys([]keys.OfferItemKey{})
		},
		AllUserKeysFunc: func() *keys.UserKeys {
			return keys.NewUserKeys([]keys.UserKey{
				userKey,
			})
		},
	}

	tradingStore := &mock.TradingStore{
		GetOfferFunc: func(key keys.OfferKey) (*trading.Offer, error) { return offer, nil },
		GetOfferItemsForOfferFunc: func(key keys.OfferKey) (*trading.OfferItems, error) {
			return offerItems, nil
		},
		FindApproversForOfferFunc: func(offerKey keys.OfferKey) (trading.Approvers, error) {
			return approvers, nil
		},
		MarkOfferItemsAsAcceptedFunc: func(ctx context.Context, approvedBy keys.UserKey, approvedByGiver *keys.OfferItemKeys, approvedByReceiver *keys.OfferItemKeys) error {
			return nil
		},
	}

	userStore := &mock.UserStore{
		GetByKeysFunc: func(ctx context.Context, userKeys *keys.UserKeys) (*user.Users, error) {
			return user.NewUsers([]*user.User{
				{
					ID: userKey.String(),
				},
			}), nil
		},
	}

	tradingService := &TradingService{
		tradingStore: tradingStore,
		userStore:    userStore,
	}

	ctx := context.TODO()
	ctx = auth.SetContextAuthenticatedUser(ctx, "user", "username", "user@email.com")

	err := tradingService.AcceptOffer(ctx, offerKey)

	assert.NoError(t.T(), err)
	assert.Len(t.T(), tradingStore.MarkOfferItemsAsAcceptedCalls(), 1)

}
