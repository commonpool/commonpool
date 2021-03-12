package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"gorm.io/gorm"
	"strings"
)

type GetOffersGroupKeys struct {
	db *gorm.DB
}

func NewGetOffersGroupKeys(db *gorm.DB) *GetOffersGroupKeys {
	return &GetOffersGroupKeys{db: db}
}

type offerGroupKeys struct {
	OfferKey keys.OfferKey
	GroupKey keys.GroupKey
}

func (q *GetOffersGroupKeys) Get(ctx context.Context, offerKeys *keys.OfferKeys) (map[keys.OfferKey]keys.GroupKey, error) {

	if len(offerKeys.Items) == 0 {
		return map[keys.OfferKey]keys.GroupKey{}, nil
	}

	var offersSb strings.Builder
	var offersParams []interface{}
	offersSb.WriteString("offer_key in (")
	for i, offerKey := range offerKeys.Items {
		offersSb.WriteString("?")
		offersParams = append(offersParams, offerKey.String())
		if i < len(offerKeys.Items)-1 {
			offersSb.WriteString(",")
		}
	}
	offersSb.WriteString(")")

	var result []*offerGroupKeys

	if err := q.db.Raw("select offer_key, group_key from offer_read_models where "+offersSb.String(), offersParams...).Find(&result).Error; err != nil {
		return nil, err
	}

	var resultMap = map[keys.OfferKey]keys.GroupKey{}
	for _, resultItem := range result {
		resultMap[resultItem.OfferKey] = resultItem.GroupKey
	}

	return resultMap, nil

}
