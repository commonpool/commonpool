package domain

import (
	"encoding/json"
	"github.com/commonpool/backend/pkg/keys"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func assertApproved(t *testing.T, offer *Offer, offerItemKey keys.OfferItemKey, direction ApprovalDirection, isApproved bool) {
	isApprovedOutbound, err := offer.IsOfferItemApproved(offerItemKey, direction)
	assert.NoError(t, err)
	assert.Equal(t, isApproved, isApprovedOutbound)
}

var approveAllMatrix PermissionMatrix = func(userKey keys.UserKey, offerItem OfferItem, direction ApprovalDirection) bool {
	return true
}

var denyAllMatrix PermissionMatrix = func(userKey keys.UserKey, offerItem OfferItem, direction ApprovalDirection) bool {
	return false
}

func assertError(t *testing.T, expected string, err error) {
	if assert.Error(t, err) {
		if assert.Equal(t, expected, err.Error()) {
			t.Log(err.Error())
		}
	}
}
func assertChangeCount(t *testing.T, offer *Offer, expected int) int {
	assert.Equal(t, expected, len(offer.changes))
	return len(offer.changes)
}

func TestSubmitOffer(t *testing.T) {

	var offer = NewOffer()
	groupKey := keys.NewGroupKey(uuid.NewV4())
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	err := offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			CreditTransferItem{
				OfferItemKey: offerItemKey,
				Amount:       time.Hour * 2,
				From:         NewUserTarget(user1Key),
				To:           NewUserTarget(user2Key),
			}},
	})

	assert.NoError(t, err)
	assert.Len(t, offer.changes, 1)
	assert.Equal(t, 1, offer.OfferItemCount())

	assertApproved(t, offer, offerItemKey, Outbound, false)
	assertApproved(t, offer, offerItemKey, Inbound, false)

	err = offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix)
	assert.NoError(t, err)
	assert.Len(t, offer.changes, 2)

	assertApproved(t, offer, offerItemKey, Outbound, true)
	assertApproved(t, offer, offerItemKey, Inbound, false)

	err = offer.ApproveOfferItem(user2Key, offerItemKey, Inbound, approveAllMatrix)
	assert.NoError(t, err)
	assert.Len(t, offer.changes, 4)

	assertApproved(t, offer, offerItemKey, Outbound, true)
	assertApproved(t, offer, offerItemKey, Inbound, true)

	assert.Equal(t, Approved, offer.status)

}

func TestCannotApprove(t *testing.T) {

	offer := NewOffer()
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	groupKey := keys.NewGroupKey(uuid.NewV4())
	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	err := offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			CreditTransferItem{
				OfferItemKey: offerItemKey,
				Amount:       time.Hour * 2,
				From:         NewUserTarget(user1Key),
				To:           NewUserTarget(user2Key),
			}},
	})

	assert.NoError(t, err)
	assert.Len(t, offer.changes, 1)
	assert.Equal(t, 1, offer.OfferItemCount())

	err = offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, denyAllMatrix)
	assertError(t, "cannot approve offer item (outbound): user 'user1' is not allowed to do this operation", err)
	assert.Len(t, offer.changes, 1)

}

func TestApproveTwiceIdempotent(t *testing.T) {

	offer := NewOffer()
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	groupKey := keys.NewGroupKey(uuid.NewV4())
	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	err := offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			CreditTransferItem{
				OfferItemKey: offerItemKey,
				Amount:       time.Hour * 2,
				From:         NewUserTarget(user1Key),
				To:           NewUserTarget(user2Key),
			}},
	})

	assert.NoError(t, err)
	assert.Len(t, offer.changes, 1)
	assert.Equal(t, 1, offer.OfferItemCount())

	err = offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix)
	assert.NoError(t, err)
	assert.Len(t, offer.changes, 2)

	err = offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix)
	assert.NoError(t, err)
	assert.Len(t, offer.changes, 2)

}

func TestApproveDeclinedOfferShouldThrow(t *testing.T) {

	offer := NewOffer()
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	groupKey := keys.NewGroupKey(uuid.NewV4())
	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	err := offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			CreditTransferItem{
				OfferItemKey: offerItemKey,
				Amount:       time.Hour * 2,
				From:         NewUserTarget(user1Key),
				To:           NewUserTarget(user2Key),
			}},
	})

	err = offer.DeclineOffer(user1Key)
	assert.NoError(t, err)
	assert.Len(t, offer.changes, 2)

	err = offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix)
	assertError(t, "cannot approve offer item (outbound): offer status must be 'pending' but is 'declined'", err)
	assert.Len(t, offer.changes, 2)

}

func TestDeclineNewOfferShouldThrow(t *testing.T) {
	var offer = NewOffer()
	user1Key := keys.NewUserKey("user1")
	err := offer.DeclineOffer(user1Key)
	assertError(t, "cannot decline offer: offer has not yet been submitted", err)
	assert.Len(t, offer.changes, 0)
}

func TestApproveOfferItemOfNewOfferShouldThrow(t *testing.T) {
	var offer = NewOffer()

	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	user1Key := keys.NewUserKey("user1")
	err := offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix)
	assertError(t, "cannot approve offer item (outbound): offer has not yet been submitted", err)
	assert.Len(t, offer.changes, 0)
}

func TestSubmitTwiceShouldThrow(t *testing.T) {

	offer := NewOffer()
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	groupKey := keys.NewGroupKey(uuid.NewV4())
	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	err := offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			CreditTransferItem{
				OfferItemKey: offerItemKey,
				Amount:       time.Hour * 2,
				From:         NewUserTarget(user1Key),
				To:           NewUserTarget(user2Key),
			}},
	})

	assert.NoError(t, err)
	assert.Len(t, offer.changes, 1)

	err = offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			CreditTransferItem{
				OfferItemKey: offerItemKey,
				Amount:       time.Hour * 2,
				From:         NewUserTarget(user1Key),
				To:           NewUserTarget(user2Key),
			}},
	})

	assertError(t, "cannot submit offer: offer has already been submitted", err)
	assert.Len(t, offer.changes, 1)

}

func TestReceiveServiceCompletesOffer(t *testing.T) {

	changeCount := 0
	offer := NewOffer()
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	groupKey := keys.NewGroupKey(uuid.NewV4())
	resourceKey := keys.NewResourceKey(uuid.NewV4())

	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	err := offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			&ServiceOfferItem{
				OfferItemKey: offerItemKey,
				Duration:     2 * time.Hour,
				From:         NewUserTarget(user1Key),
				To:           NewUserTarget(user2Key),
				ResourceKey:  resourceKey,
			}},
	})
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Inbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+2)

	err = offer.NotifyServiceGiven(user1Key, offerItemKey)
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	err = offer.NotifyServiceReceived(user1Key, offerItemKey)
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+2)

}

func TestGiveServiceCompletesOffer(t *testing.T) {

	changeCount := 0
	offer := NewOffer()
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	groupKey := keys.NewGroupKey(uuid.NewV4())
	resourceKey := keys.NewResourceKey(uuid.NewV4())

	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	err := offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			&ServiceOfferItem{
				OfferItemKey: offerItemKey,
				Duration:     2 * time.Hour,
				From:         NewUserTarget(user1Key),
				To:           NewUserTarget(user2Key),
				ResourceKey:  resourceKey,
			}},
	})
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Inbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+2)

	err = offer.NotifyServiceReceived(user1Key, offerItemKey)
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	err = offer.NotifyServiceGiven(user1Key, offerItemKey)
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+2)

}

func TestReceiveResourceCompletesOffer(t *testing.T) {

	changeCount := 0
	offer := NewOffer()
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	groupKey := keys.NewGroupKey(uuid.NewV4())
	resourceKey := keys.NewResourceKey(uuid.NewV4())

	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	err := offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			&ResourceTransferItem{
				OfferItemKey: offerItemKey,
				To:           NewUserTarget(user2Key),
				ResourceKey:  resourceKey,
			}},
	})
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Inbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+2)

	err = offer.NotifyResourceGiven(user1Key, offerItemKey)
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	err = offer.NotifyResourceReceived(user1Key, offerItemKey)
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+2)

	t.Log(offer)

}

func TestGiveResourceCompletesOffer(t *testing.T) {

	changeCount := 0
	offer := NewOffer()
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	groupKey := keys.NewGroupKey(uuid.NewV4())
	resourceKey := keys.NewResourceKey(uuid.NewV4())

	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	err := offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			&ResourceTransferItem{
				OfferItemKey: offerItemKey,
				To:           NewUserTarget(user2Key),
				ResourceKey:  resourceKey,
			}},
	})
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Inbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+2)

	err = offer.NotifyResourceReceived(user1Key, offerItemKey)
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	err = offer.NotifyResourceGiven(user1Key, offerItemKey)
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+2)

	t.Log(offer)

}

func TestBorrowItem(t *testing.T) {

	changeCount := 0
	offer := NewOffer()
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	groupKey := keys.NewGroupKey(uuid.NewV4())
	resourceKey := keys.NewResourceKey(uuid.NewV4())

	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	err := offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			&ResourceBorrowItem{
				OfferItemKey: offerItemKey,
				To:           NewUserTarget(user2Key),
				ResourceKey:  resourceKey,
				Duration:     2 * time.Hour,
			}},
	})
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Inbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+2)

	err = offer.NotifyResourceBorrowed(user1Key, offerItemKey)
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	err = offer.NotifyResourceLent(user1Key, offerItemKey)
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	err = offer.NotifyBorrowerReturnedResource(user1Key, offerItemKey)
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	err = offer.NotifyLenderReceivedBackResource(user1Key, offerItemKey)
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+2)

	t.Log(offer)

}

func TestCannotReturnItemBeforeBorrowingIt(t *testing.T) {

	offer := NewOffer()
	changeCount := 0
	groupKey := keys.NewGroupKey(uuid.NewV4())
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	resourceKey := keys.NewResourceKey(uuid.NewV4())
	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	err := offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			&ResourceBorrowItem{
				OfferItemKey: offerItemKey,
				To:           NewUserTarget(user2Key),
				ResourceKey:  resourceKey,
				Duration:     2 * time.Hour,
			}},
	})
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Inbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+2)

	err = offer.NotifyBorrowerReturnedResource(user1Key, offerItemKey)
	assert.Error(t, err)
	changeCount = assertChangeCount(t, offer, changeCount)

	t.Log(offer)

}

func TestCannotReceiveItemBeforeBorrowingIt(t *testing.T) {

	offer := NewOffer()
	changeCount := 0
	groupKey := keys.NewGroupKey(uuid.NewV4())
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	resourceKey := keys.NewResourceKey(uuid.NewV4())

	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	err := offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			&ResourceBorrowItem{
				OfferItemKey: offerItemKey,
				To:           NewUserTarget(user2Key),
				ResourceKey:  resourceKey,
				Duration:     2 * time.Hour,
			}},
	})
	assert.NoError(t, err)
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Inbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+1)

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix))
	changeCount = assertChangeCount(t, offer, changeCount+2)

	err = offer.NotifyLenderReceivedBackResource(user1Key, offerItemKey)
	assert.Error(t, err)
	changeCount = assertChangeCount(t, offer, changeCount)

	t.Log(offer)

}

func TestFromEvents(t *testing.T) {

	offer := NewOffer()
	groupKey := keys.NewGroupKey(uuid.NewV4())
	offerItemKey := keys.NewOfferItemKey(uuid.NewV4())
	resourceKey := keys.NewResourceKey(uuid.NewV4())
	user1Key := keys.NewUserKey("user1")
	user2Key := keys.NewUserKey("user2")

	assert.NoError(t, offer.Submit(groupKey, OfferItems{
		Items: []OfferItem{
			&ResourceBorrowItem{
				OfferItemKey: offerItemKey,
				To:           NewUserTarget(user2Key),
				ResourceKey:  resourceKey,
				Duration:     2 * time.Hour,
			}},
	}))
	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Inbound, approveAllMatrix))

	assert.NoError(t, offer.ApproveOfferItem(user1Key, offerItemKey, Outbound, approveAllMatrix))

	assert.NoError(t, offer.NotifyResourceBorrowed(user1Key, offerItemKey))
	assert.NoError(t, offer.NotifyResourceLent(user1Key, offerItemKey))
	assert.NoError(t, offer.NotifyBorrowerReturnedResource(user1Key, offerItemKey))
	assert.NoError(t, offer.NotifyLenderReceivedBackResource(user1Key, offerItemKey))

	fromEvents := NewFromEvents(offer.key, offer.changes)
	// ignore changes and version differences
	fromEvents.changes = offer.changes
	fromEvents.version = offer.version

	initialJsBytes, _ := json.MarshalIndent(offer, "", "  ")
	fromEventsJsBytes, _ := json.MarshalIndent(fromEvents, "", "  ")

	t.Log(string(initialJsBytes))
	t.Log(string(fromEventsJsBytes))

	assert.Equal(t, string(initialJsBytes), string(fromEventsJsBytes))

}
