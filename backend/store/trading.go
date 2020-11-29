package store

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/trading"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strings"
)

type TradingStore struct {
	db *gorm.DB
}

var _ trading.Store = TradingStore{}

func NewTradingStore(db *gorm.DB) *TradingStore {
	return &TradingStore{db: db}
}

func (t TradingStore) SaveOfferStatus(key model.OfferKey, status trading.OfferStatus) error {
	qry := t.db.
		Model(trading.Offer{}).
		Where("id = ?", key.ID.String()).
		Update("status", status)

	if qry.Error == nil && qry.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return qry.Error
}

func (t TradingStore) GetItem(ctx context.Context, key model.OfferItemKey) (*trading.OfferItem, error) {
	var item trading.OfferItem
	err := t.db.
		Model(trading.OfferItem{}).
		Where("id = ?", key.ID.String()).
		First(&item).
		Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (t TradingStore) ConfirmItemReceived(ctx context.Context, key model.OfferItemKey) error {
	req := t.db.
		Model(trading.OfferItem{}).
		Where("id = ?", key.ID.String()).
		Update("received", true)
	if req.Error == nil && req.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return req.Error
}

func (t TradingStore) ConfirmItemGiven(ctx context.Context, key model.OfferItemKey) error {

	req := t.db.
		Model(trading.OfferItem{}).
		Where("id = ?", key.ID.String()).
		Update("given", true)

	if req.Error == nil && req.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return req.Error
}

func (t TradingStore) SaveOffer(offer trading.Offer, items *trading.OfferItems) error {

	err := t.db.Create(&offer).Error
	if err != nil {
		return err
	}

	userKeys := map[string]bool{}
	for _, item := range items.Items {
		err := t.db.Create(item).Error
		if err != nil {
			return err
		}
		userKeys[item.FromUserID] = true
		userKeys[item.ToUserID] = true
	}

	for userKey := range userKeys {
		decision := trading.OfferDecision{
			OfferID:  offer.ID,
			UserID:   userKey,
			Decision: trading.PendingDecision,
		}
		err := t.db.Create(decision).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (t TradingStore) GetOffer(key model.OfferKey) (trading.Offer, error) {
	var offer trading.Offer
	err := t.db.First(&offer, "id = ?", key.ID.String()).Error
	return offer, err
}

func (t TradingStore) GetItems(key model.OfferKey) (*trading.OfferItems, error) {
	var items []trading.OfferItem
	err := t.db.Find(&items, "offer_id = ?", key.ID.String()).Error
	return trading.NewOfferItems(items), err
}

func (t TradingStore) GetOffers(qry trading.GetOffersQuery) (trading.GetOffersResult, error) {

	chain := t.db.Model(trading.Offer{})

	if len(qry.UserKeys) > 0 {
		userCount := len(qry.UserKeys)
		qryPlaceholders := make([]string, userCount)
		qryParams := make([]interface{}, userCount)
		for i := 0; i < userCount; i++ {
			qryPlaceholders[i] = "?"
			qryParams[i] = qry.UserKeys[i].String()
		}
		qryStr := strings.Join(qryPlaceholders, ",")
		qryStr = "user_id in (" + qryStr + ")"
		userQry := t.db.Model(trading.OfferDecision{}).Select("offer_id").Where(qryStr, qryParams...)
		chain = chain.Where("id in (?)", userQry)
	}

	if qry.ResourceKey != nil {
		resQry := t.db.
			Model(trading.OfferItem{}).
			Select("offer_id").
			Where("offer_type = ? AND resource_id = ?", resource.Offer, qry.ResourceKey.String())
		chain = chain.Where("id in (?)", resQry)
	}

	if qry.Status != nil {
		chain = chain.Where("status = ?", qry.Status)
	}

	var offers []trading.Offer
	err := chain.Order("created_at desc").Find(&offers).Error
	if err != nil {
		return trading.GetOffersResult{}, err
	}

	offerCount := len(offers)

	offersById := map[model.OfferKey]trading.GetOffersResultItem{}
	offerResultItems := make([]trading.GetOffersResultItem, offerCount)

	for i, offer := range offers {
		resultItem := trading.GetOffersResultItem{
			Offer:          offer,
			OfferItems:     []trading.OfferItem{},
			OfferDecisions: []trading.OfferDecision{},
		}
		offerResultItems[i] = resultItem
		offersById[offer.GetKey()] = resultItem
	}

	var offerIds = make([]uuid.UUID, offerCount)
	for i, offer := range offers {
		offerIds[i] = offer.ID
	}

	qryPlaceholders := make([]string, offerCount)
	qryParams := make([]interface{}, offerCount)
	for i := 0; i < offerCount; i++ {
		qryPlaceholders[i] = "?"
		qryParams[i] = offerIds[i]
	}
	qryStr := strings.Join(qryPlaceholders, ",")
	qryStr = "offer_id in (" + qryStr + ")"

	var offerItems []trading.OfferItem
	err = t.db.Model(trading.OfferItem{}).Where(qryStr, qryParams...).Find(&offerItems).Error
	if err != nil {
		return trading.GetOffersResult{}, err
	}

	for _, offerItem := range offerItems {
		resultItem, ok := offersById[offerItem.GetOfferKey()]
		if !ok {
			return trading.GetOffersResult{}, fmt.Errorf("could not find item offer")
		}
		resultItem.OfferItems = append(resultItem.OfferItems, offerItem)
	}

	var offerDecisions []trading.OfferDecision
	err = t.db.Model(trading.OfferDecision{}).Where(qryStr, qryParams...).Find(&offerDecisions).Error
	if err != nil {
		return trading.GetOffersResult{}, err
	}

	for _, offerDecision := range offerDecisions {
		resultItem, ok := offersById[offerDecision.GetOfferKey()]
		if !ok {
			return trading.GetOffersResult{}, fmt.Errorf("could not find decision offer")
		}
		resultItem.OfferDecisions = append(resultItem.OfferDecisions, offerDecision)
	}

	return trading.GetOffersResult{
		Items: offerResultItems,
	}, nil

}

func (t TradingStore) GetDecisions(key model.OfferKey) ([]trading.OfferDecision, error) {
	var decisions []trading.OfferDecision
	err := t.db.Model(trading.OfferDecision{}).Find(&decisions, "offer_id = ?", key.ID.String()).Error
	return decisions, err
}

func (t TradingStore) SaveDecision(key model.OfferKey, user model.UserKey, decision trading.Decision) error {
	return t.db.Model(trading.OfferDecision{}).
		Where("offer_id = ? and user_id = ?", key.ID.String(), user.String()).
		Updates(map[string]interface{}{
			"decision": decision,
		}).
		Error
}

func (t TradingStore) GetTradingHistory(ctx context.Context, ids *model.UserKeys) ([]trading.HistoryEntry, error) {

	var offerItems []trading.OfferItem

	qry := t.db.Model(trading.OfferItem{}).
		Joins("JOIN offers ON offer_items.offer_id = offers.id AND offers.status = ?", trading.CompletedOffer)

	if ids != nil && ids.Items != nil && len(ids.Items) > 0 {
		var placeholders []string
		var params []interface{}
		for _, item := range ids.Items {
			placeholders = append(placeholders, "?")
			params = append(params, item.String())
		}
		qry = qry.Where(fmt.Sprintf("offers_items.user_id in (%s)", strings.Join(placeholders, ",")), params...)
	}

	err := qry.Find(&offerItems).Error

	if err != nil {
		return nil, err
	}

	var tradingHistory []trading.HistoryEntry

	for _, offerItem := range offerItems {
		fromUser := model.NewUserKey(offerItem.FromUserID)
		toUser := model.NewUserKey(offerItem.ToUserID)

		var resourceKey *model.ResourceKey
		if offerItem.ResourceID != nil {
			rk := model.NewResourceKey(offerItem.ID)
			resourceKey = &rk
		}

		tradingHistory = append(tradingHistory, trading.HistoryEntry{
			FromUserID:        fromUser,
			ToUserID:          toUser,
			ResourceID:        resourceKey,
			TimeAmountSeconds: offerItem.OfferedTimeInSeconds,
		})

	}

	return tradingHistory, nil
}
