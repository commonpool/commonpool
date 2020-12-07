package trading

import (
	"github.com/commonpool/backend/model"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestOfferItem struct {
	User1Key model.UserKey
	User2Key model.UserKey
	User3Key model.UserKey
	OfferKey model.OfferKey

	Resource1Key  model.ResourceKey
	Resource2Key  model.ResourceKey
	Resource3Key  model.ResourceKey
	ResourceItem1 OfferItem
	ResourceItem2 OfferItem
	ResourceItem3 OfferItem

	TimeItem1 OfferItem
	TimeItem2 OfferItem
}

func NewTestOfferItem() TestOfferItem {

	user1Key := model.NewUserKey("user1")
	user2Key := model.NewUserKey("user2")
	user3Key := model.NewUserKey("user3")
	offerKey := model.NewOfferKey(uuid.NewV4())
	resource1Key := model.NewResourceKey(uuid.NewV4())
	resource2Key := model.NewResourceKey(uuid.NewV4())
	resource3Key := model.NewResourceKey(uuid.NewV4())
	resourceItem1 := model.NewOfferItemKey(uuid.NewV4())
	resourceItem2 := model.NewOfferItemKey(uuid.NewV4())
	resourceItem3 := model.NewOfferItemKey(uuid.NewV4())
	timeItem1 := model.NewOfferItemKey(uuid.NewV4())
	timeItem2 := model.NewOfferItemKey(uuid.NewV4())

	timeInSeconds1 := int64((time.Hour * 4).Hours())
	timeInSeconds2 := int64((time.Hour * 4).Hours())
	return TestOfferItem{
		User1Key:     user1Key,
		User2Key:     user2Key,
		OfferKey:     offerKey,
		Resource1Key: resource1Key,
		Resource2Key: resource2Key,
		ResourceItem1: OfferItem{
			Key:         resourceItem1.ID,
			ResourceKey: &resource1Key.ID,
			From:        user1Key.String(),
			To:          user2Key.String(),
			ItemType:    ResourceItem,
			OfferKey:    offerKey.ID,
		},
		ResourceItem2: OfferItem{
			Key:         resourceItem2.ID,
			ResourceKey: &resource1Key.ID,
			From:        user2Key.String(),
			To:          user1Key.String(),
			ItemType:    ResourceItem,
			OfferKey:    offerKey.ID,
		},
		ResourceItem3: OfferItem{
			Key:         resourceItem3.ID,
			OfferKey:    offerKey.ID,
			ItemType:    ResourceItem,
			From:        user1Key.String(),
			To:          user3Key.String(),
			ResourceKey: &resource3Key.ID,
		},
		TimeItem1: OfferItem{
			Key:                  timeItem1.ID,
			OfferedTimeInSeconds: &timeInSeconds1,
			From:                 user2Key.String(),
			To:                   user1Key.String(),
			ItemType:             TimeItem,
			OfferKey:             offerKey.ID,
		},
		TimeItem2: OfferItem{
			Key:                  timeItem2.ID,
			OfferedTimeInSeconds: &timeInSeconds2,
			From:                 user2Key.String(),
			To:                   user1Key.String(),
			ItemType:             TimeItem,
			OfferKey:             offerKey.ID,
		},
	}

}

func TestNewOfferItemsDuplicate(t *testing.T) {
	test := NewTestOfferItem()
	i := NewOfferItems([]OfferItem{test.TimeItem1, test.TimeItem1})
	assert.Equal(t, 1, len(i.Items))
}

func TestGetUsers(t *testing.T) {
	test := NewTestOfferItem()
	i := NewOfferItems([]OfferItem{test.TimeItem1, test.TimeItem2, test.ResourceItem1, test.ResourceItem2})
	u := i.GetUserKeys()

	assert.True(t, u.Contains(test.User1Key))
	assert.True(t, u.Contains(test.User2Key))
	assert.Equal(t, 2, len(u.Items))
}

func TestGetItemsForUser(t *testing.T) {
	test := NewTestOfferItem()
	i := NewOfferItems([]OfferItem{
		test.TimeItem1,
		test.TimeItem2,
		test.ResourceItem1,
		test.ResourceItem2,
		test.ResourceItem3,
	})
	assert.Equal(t, 5, len(i.GetOfferItemsReceivedByUser(test.User1Key).Items))
	assert.Equal(t, 4, len(i.GetOfferItemsReceivedByUser(test.User2Key).Items))
	assert.Equal(t, 0, len(i.GetOfferItemsReceivedByUser(model.NewUserKey("abc")).Items))
}
