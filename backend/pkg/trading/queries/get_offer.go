package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	readmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"strings"
)

type GetOfferResult struct {
	Offer      *readmodels.DBOfferReadModel
	OfferItems []*readmodels.OfferItemReadModel
}

type GetOffer struct {
	db *gorm.DB
}

func NewGetOffer(db *gorm.DB) *GetOffer {
	return &GetOffer{
		db: db,
	}
}

func (q *GetOffer) Get(ctx context.Context, offerKey keys.OfferKey) (*readmodels.OfferReadModel, error) {
	offers, err := getOffers(ctx, keys.NewOfferKeys([]keys.OfferKey{offerKey}), q.db)
	if err != nil {
		return nil, err
	}
	if len(offers) == 0 {
		return nil, exceptions.ErrOfferNotFound
	}
	return offers[0], nil
}

func getOffers(ctx context.Context, offerKeys *keys.OfferKeys, db *gorm.DB) ([]*readmodels.OfferReadModel, error) {

	if len(offerKeys.Items) == 0 {
		return []*readmodels.OfferReadModel{}, nil
	}

	var sb strings.Builder
	sb.WriteString("offer_key in (")
	var offerKeysLen = len(offerKeys.Items)
	var params = make([]interface{}, offerKeysLen)
	for i := 0; i < offerKeysLen; i++ {
		sb.WriteString("?")
		if i < offerKeysLen-1 {
			sb.WriteString(",")
		}
		params[i] = offerKeys.Items[i]
	}
	sb.WriteString(")")

	var offers []*readmodels.DBOfferReadModel
	var offerItems []*readmodels.OfferItemReadModel
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		query := db.Model(&readmodels.DBOfferReadModel{}).
			Where(sb.String(), params...).
			Find(&offers)
		if err := query.Error; err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		return db.
			Model(&readmodels.OfferItemReadModel{}).
			Where(sb.String(), params...).
			Find(&offerItems).Error
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}
	if len(offers) == 0 {
		return []*readmodels.OfferReadModel{}, nil
	}

	cache := NewReadModelCache()

	for _, offer := range offers {
		cache.processOffer(offer)
	}

	for _, item := range offerItems {
		cache.processOfferItem(item)
	}

	if err := cache.retrieve(ctx, db); err != nil {
		return nil, err
	}

	mappedItems := mapOfferItems(offerItems, cache)
	groupedItems := map[keys.OfferKey][]*readmodels.OfferItemReadModel2{}
	for _, mappedItem := range mappedItems {
		if _, ok := groupedItems[mappedItem.OfferKey]; !ok {
			groupedItems[mappedItem.OfferKey] = []*readmodels.OfferItemReadModel2{}
		}
		groupedItems[mappedItem.OfferKey] = append(groupedItems[mappedItem.OfferKey], mappedItem)
	}

	var mappedOffers []*readmodels.OfferReadModel
	for _, offer := range offers {
		mappedItemsForOffer := groupedItems[offer.OfferKey]
		mappedOffer := readmodels.OfferReadModel{
			OfferReadModelBase: offer.OfferReadModelBase,
			DeclinedBy:         cache.getUserReadModel(offer.DeclinedBy),
			SubmittedBy:        cache.getUserReadModel(offer.SubmittedBy),
			OfferItems:         mappedItemsForOffer,
		}
		mappedOffers = append(mappedOffers, &mappedOffer)
	}

	return mappedOffers, nil

}

func mapOfferItems(offerItems []*readmodels.OfferItemReadModel, cache *readModelCache) []*readmodels.OfferItemReadModel2 {
	var mappedItems = make([]*readmodels.OfferItemReadModel2, len(offerItems))
	for i, item := range offerItems {
		mappedItem := mapOfferItem(item, cache)
		mappedItems[i] = mappedItem
	}
	return mappedItems
}

func mapOfferItem(item *readmodels.OfferItemReadModel, cache *readModelCache) *readmodels.OfferItemReadModel2 {

	mappedItem := &readmodels.OfferItemReadModel2{
		OfferItemReadModelBase: item.OfferItemReadModelBase,
		ApprovedInboundBy:      cache.getUserReadModel(item.ApprovedInboundBy),
		ApprovedOutboundBy:     cache.getUserReadModel(item.ApprovedOutboundBy),
		ServiceGivenBy:         cache.getUserReadModel(item.ServiceGivenBy),
		ServiceReceivedBy:      cache.getUserReadModel(item.ServiceReceivedBy),
		ResourceGivenBy:        cache.getUserReadModel(item.ResourceGivenBy),
		ResourceTakenBy:        cache.getUserReadModel(item.ResourceTakenBy),
		ResourceBorrowedBy:     cache.getUserReadModel(item.ResourceBorrowedBy),
		ResourceLentBy:         cache.getUserReadModel(item.ResourceLentBy),
		BorrowedItemReturnedBy: cache.getUserReadModel(item.BorrowedItemReturnedBy),
		LentItemReceivedBy:     cache.getUserReadModel(item.LentItemReceivedBy),
		From:                   cache.getTargetReadModel(item.From),
		To:                     cache.getTargetReadModel(item.To),
		Resource:               cache.getResource(item.ResourceKey),
	}

	return mappedItem
}

type readModelCache struct {
	allUserKeys     map[keys.UserKey]bool
	allGroupKeys    map[keys.GroupKey]bool
	allResourceKeys map[keys.ResourceKey]bool
	users           []*readmodels.OfferUserReadModel
	groups          []*readmodels.OfferGroupReadModel
	resources       []*readmodels.OfferResourceReadModel
}

func NewReadModelCache() *readModelCache {
	return &readModelCache{
		allUserKeys:     map[keys.UserKey]bool{},
		allGroupKeys:    map[keys.GroupKey]bool{},
		allResourceKeys: map[keys.ResourceKey]bool{},
		users:           []*readmodels.OfferUserReadModel{},
		groups:          []*readmodels.OfferGroupReadModel{},
		resources:       []*readmodels.OfferResourceReadModel{},
	}
}

func (c *readModelCache) processOffer(offer *readmodels.DBOfferReadModel) {
	c.addUser(offer.SubmittedBy)
	c.addUser(offer.DeclinedBy)
}

func (c *readModelCache) processOfferItem(offerItem *readmodels.OfferItemReadModel) {
	c.addUser(offerItem.From.UserKey)
	c.addUser(offerItem.To.UserKey)
	c.addUser(offerItem.ApprovedInboundBy)
	c.addUser(offerItem.ApprovedOutboundBy)
	c.addUser(offerItem.ResourceGivenBy)
	c.addUser(offerItem.ResourceTakenBy)
	c.addUser(offerItem.ResourceBorrowedBy)
	c.addUser(offerItem.ResourceLentBy)
	c.addUser(offerItem.BorrowedItemReturnedBy)
	c.addUser(offerItem.LentItemReceivedBy)
	c.addUser(offerItem.ServiceGivenBy)
	c.addUser(offerItem.ServiceReceivedBy)
	c.addGroup(offerItem.From.GroupKey)
	c.addGroup(offerItem.To.GroupKey)
	c.addResource(offerItem.ResourceKey)
}

func (c *readModelCache) addUser(userKey *keys.UserKey) {
	if userKey == nil {
		return
	}
	c.allUserKeys[*userKey] = true
}

func (c *readModelCache) addGroup(groupKey *keys.GroupKey) {
	if groupKey == nil {
		return
	}
	c.allGroupKeys[*groupKey] = true
}

func (c *readModelCache) addResource(resourceKey *keys.ResourceKey) {
	if resourceKey == nil {
		return
	}
	c.allResourceKeys[*resourceKey] = true
}

func (c *readModelCache) getUserReadModel(userKey *keys.UserKey) *readmodels.OfferUserReadModel {
	if userKey == nil {
		return nil
	}
	for _, user := range c.users {
		if user.UserKey == *userKey {
			return user
		}
	}
	return &readmodels.OfferUserReadModel{
		UserKey: *userKey,
		Version: -1,
	}
}

func (c *readModelCache) getGroupReadModel(groupKey *keys.GroupKey) *readmodels.OfferGroupReadModel {
	if groupKey == nil {
		return nil
	}
	for _, group := range c.groups {
		if group.GroupKey == *groupKey {
			return group
		}
	}
	return &readmodels.OfferGroupReadModel{
		GroupKey: *groupKey,
		Version:  -1,
	}
}

func (c *readModelCache) getResource(resourceKey *keys.ResourceKey) *readmodels.OfferResourceReadModel {
	if resourceKey == nil {
		return nil
	}
	for _, resource := range c.resources {
		if resource.ResourceKey == *resourceKey {
			return resource
		}
	}
	return nil
}

func (c *readModelCache) getTargetReadModel(target *keys.Target) *readmodels.OfferItemTargetReadModel {
	if target == nil {
		return nil
	}
	if target.GroupKey == nil && target.UserKey == nil {
		return nil
	}
	var userName *string
	var groupName *string
	var userVersion *int
	var groupVersion *int
	if target.UserKey != nil {
		user := c.getUserReadModel(target.UserKey)
		if user != nil {
			userName = &user.Username
			userVersion = &user.Version
		}
	} else if target.GroupKey != nil {
		group := c.getGroupReadModel(target.GroupKey)
		if group != nil {
			groupName = &group.GroupName
			groupVersion = &group.Version
		}
	}
	result := &readmodels.OfferItemTargetReadModel{
		Target:       *target,
		GroupName:    groupName,
		UserName:     userName,
		UserVersion:  userVersion,
		GroupVersion: groupVersion,
	}
	result.Name = result.GetName()
	return result
}

func (c *readModelCache) retrieve(ctx context.Context, db *gorm.DB) error {

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if len(c.allGroupKeys) != 0 {
			var groupKeyParams []interface{}
			var groupKeyPlaceholders []string
			for key, _ := range c.allGroupKeys {
				groupKeyParams = append(groupKeyParams, key)
				groupKeyPlaceholders = append(groupKeyPlaceholders, "?")
			}
			err := db.
				Model(&readmodels.OfferGroupReadModel{}).
				Where("group_key in ("+strings.Join(groupKeyPlaceholders, ",")+")", groupKeyParams...).
				Find(&c.groups).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	g.Go(func() error {
		if len(c.allUserKeys) != 0 {
			var userKeyParams []interface{}
			var userKeyPlaceholders []string
			for key, _ := range c.allUserKeys {
				userKeyParams = append(userKeyParams, key)
				userKeyPlaceholders = append(userKeyPlaceholders, "?")
			}
			err := db.Where("user_key in ("+strings.Join(userKeyPlaceholders, ",")+")", userKeyParams...).
				Find(&c.users).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	g.Go(func() error {
		if len(c.allResourceKeys) != 0 {
			var resParams []interface{}
			var resPlaceholders []string
			for key, _ := range c.allResourceKeys {
				resParams = append(resParams, key)
				resPlaceholders = append(resPlaceholders, "?")
			}
			err := db.Where("resource_key in ("+strings.Join(resPlaceholders, ",")+")", resParams...).
				Find(&c.resources).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	return g.Wait()

}
