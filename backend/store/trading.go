package store

import (
	"fmt"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/trading"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type TradingStore struct {
	db *gorm.DB
}

func (t TradingStore) SaveOffer(offer model.Offer, items []model.OfferItem) error {

	err := t.db.Create(&offer).Error
	if err != nil {
		return err
	}

	userKeys := map[string]bool{}
	for _, item := range items {
		err := t.db.Create(item).Error
		if err != nil {
			return err
		}
		userKeys[item.FromUserID] = true
		userKeys[item.ToUserID] = true
	}

	for userKey := range userKeys {
		decision := model.OfferDecision{
			OfferID:  offer.ID,
			UserID:   userKey,
			Decision: model.PendingDecision,
		}
		err := t.db.Create(decision).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (t TradingStore) GetOffer(key model.OfferKey) (model.Offer, error) {
	var offer model.Offer
	err := t.db.First(&offer, "id = ?", key.ID.String()).Error
	return offer, err
}

func (t TradingStore) GetItems(key model.OfferKey) ([]model.OfferItem, error) {
	var items []model.OfferItem
	err := t.db.Find(&items, "offer_id = ?", key.ID.String()).Error
	return items, err
}

func (t TradingStore) GetOffers(qry trading.GetOffersQuery) (trading.GetOffersResult, error) {

	chain := t.db.Model(model.Offer{})

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
		userQry := t.db.Model(model.OfferDecision{}).Select("offer_id").Where(qryStr, qryParams...)
		chain = chain.Where("id in (?)", userQry)
	}

	if qry.ResourceKey != nil {
		resQry := t.db.
			Model(model.OfferItem{}).
			Select("offer_id").
			Where("offer_type = ? AND resource_id = ?", model.ResourceOffer, qry.ResourceKey.String())
		chain = chain.Where("id in (?)", resQry)
	}

	if qry.Status != nil {
		chain = chain.Where("status = ?", qry.Status)
	}

	var offers []model.Offer
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
			OfferItems:     []model.OfferItem{},
			OfferDecisions: []model.OfferDecision{},
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

	var offerItems []model.OfferItem
	err = t.db.Model(model.OfferItem{}).Where(qryStr, qryParams...).Find(&offerItems).Error
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

	var offerDecisions []model.OfferDecision
	err = t.db.Model(model.OfferDecision{}).Where(qryStr, qryParams...).Find(&offerDecisions).Error
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

func (t TradingStore) GetDecisions(key model.OfferKey) ([]model.OfferDecision, error) {
	var decisions []model.OfferDecision
	err := t.db.Find(&decisions, "offer_id = ?", key.ID.String()).Error
	return decisions, err
}

func (t TradingStore) SaveDecision(key model.OfferKey, user model.UserKey, decision model.Decision) error {
	return t.db.Model(model.OfferDecision{}).
		Where("offer_id = ? and user_id = ?", key.ID.String(), user.String()).
		Update("decision", decision).
		Error
}

func (t TradingStore) CompleteOffer(key model.OfferKey, status model.OfferStatus) error {
	completionTime := time.Now()
	return t.db.Model(model.Offer{}).
		Where("id = ?", key.ID).
		Updates(map[string]interface{}{
			"completed_at": &completionTime,
			"status":       status,
		}).Error

}

var _ trading.Store = TradingStore{}

func NewTradingStore(db *gorm.DB) trading.Store {
	return TradingStore{db: db}
}
