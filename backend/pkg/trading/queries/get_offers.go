package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	groupreadmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	"gorm.io/gorm"
	"strings"
)

type GetOffers struct {
	db *gorm.DB
}

func NewGetOffers(db *gorm.DB) *GetOffers {
	return &GetOffers{db: db}
}

func (q *GetOffers) Get(ctx context.Context, userKey keys.UserKey) ([]*groupreadmodels.OfferReadModel, error) {

	// find groups where user is admin/owner

	var memberships []*groupreadmodels.OfferUserMembershipReadModel
	err := q.db.Model(&groupreadmodels.OfferUserMembershipReadModel{}).
		Where("user_key = ? and (is_admin = true or is_owner = true)", userKey).
		Find(&memberships).
		Error
	if err != nil {
		return nil, err
	}

	// find offer keys where user can approve something

	var sb strings.Builder
	sb.WriteString("select * from offer_item_read_models where ( from_user_key = ? or to_user_key = ? ")

	var membershipCount = len(memberships)
	var params = make([]interface{}, 2+2*membershipCount)
	params[0] = userKey
	params[1] = userKey

	if membershipCount > 0 {
		sb.WriteString(" or from_group_key in (")
		for i := 0; i < membershipCount; i++ {
			params[i+2] = memberships[i].GroupKey
			sb.WriteString("?")
			if i < membershipCount-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(") or (")
		for i := 0; i < membershipCount; i++ {
			params[i+2+membershipCount] = memberships[i].GroupKey
			sb.WriteString("?")
			if i < membershipCount-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(")")
	}
	sb.WriteString(") group by offer_key")

	type result struct {
		OfferKey keys.OfferKey
	}
	var results []result

	err = q.db.Raw(sb.String(), params...).
		Find(&results).
		Error
	if err != nil {
		return nil, err
	}

	// retrieve offers

	var resultLen = len(results)
	if resultLen == 0 {
		return []*groupreadmodels.OfferReadModel{}, nil
	}

	sb.Reset()
	params = nil
	params = make([]interface{}, resultLen)

	sb.WriteString("offer_key in (")
	for i := 0; i < resultLen; i++ {
		sb.WriteString("?")
		if i < resultLen-1 {
			sb.WriteString(",")
		}
		params[i] = results[i].OfferKey
	}
	sb.WriteString(")")

	q.db.Model(&groupreadmodels.OfferReadModel{})

	var offerKeys = keys.NewOfferKeys([]keys.OfferKey{})
	for _, r := range results {
		offerKeys.Items = append(offerKeys.Items, r.OfferKey)
	}

	return getOffers(ctx, offerKeys, q.db)

}
