package trading

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOfferApprovalState(t *testing.T) {

	offerKey := keys.GenerateOfferKey()
	groupKey := keys.GenerateGroupKey()
	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")
	userTarget := NewUserTarget(user1Key)
	groupTarget := NewGroupTarget(groupKey)
	resourceKey := keys.GenerateResourceKey()

	offerItem1 := NewCreditTransferItem(offerKey, keys.GenerateOfferItemKey(), userTarget, groupTarget, time.Hour)
	offerItem2 := NewResourceTransferItem(offerKey, keys.GenerateOfferItemKey(), userTarget, resourceKey)

	offerItem1OutboundApproval := NewOfferItemApproval(offerItem1.Key, user1Key, Outbound)
	offerItem1InboundApproval := NewOfferItemApproval(offerItem1.Key, user2Key, Inbound)
	offerItem2OutboundApproval := NewOfferItemApproval(offerItem2.Key, user2Key, Outbound)

	var state = NewOfferApprovalState(
		NewOffer(offerKey, groupKey, user1Key, "", nil),
		NewOfferItemsFrom(
			offerItem1,
			offerItem2,
		),
		NewOfferItemApprovals(
			offerItem1OutboundApproval,
			offerItem1InboundApproval,
			offerItem2OutboundApproval))

	retrievedOfferItem1, err := state.GetOfferItem(offerItem1.Key)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, offerItem1, retrievedOfferItem1)

	retrievedOfferItem2, err := state.GetOfferItem(offerItem2.Key)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, offerItem2, retrievedOfferItem2)

	offerItem1OutboundApproved, err := state.IsOutboundApproved(offerItem1.Key)
	if !assert.NoError(t, err) {
		return
	}
	assert.True(t, offerItem1OutboundApproved)

	offerItem1InboundApproved, err := state.IsInboundApproved(offerItem1.Key)
	if !assert.NoError(t, err) {
		return
	}
	assert.True(t, offerItem1InboundApproved)

	offerItem2OutboundApproved, err := state.IsOutboundApproved(offerItem2.Key)
	if !assert.NoError(t, err) {
		return
	}
	assert.True(t, offerItem2OutboundApproved)

	offerItem2InboundApproved, err := state.IsInboundApproved(offerItem2.Key)
	if !assert.NoError(t, err) {
		return
	}
	assert.False(t, offerItem2InboundApproved)

	approvalsForOfferItem1 := state.GetApprovalsForOfferItem(offerItem1.Key)
	assert.Len(t, approvalsForOfferItem1.Items, 2)
	assert.Equal(t, offerItem1OutboundApproval, approvalsForOfferItem1.Items[0])
	assert.Equal(t, offerItem1InboundApproval, approvalsForOfferItem1.Items[1])

	approvalsForOfferItem2 := state.GetApprovalsForOfferItem(offerItem2.Key)
	assert.Len(t, approvalsForOfferItem2.Items, 1)
	assert.Equal(t, offerItem2OutboundApproval, approvalsForOfferItem2.Items[0])

}
