package trading

import "github.com/commonpool/backend/pkg/keys"

type Offers struct {
	Items []*Offer
}

func NewOffers(items []*Offer) *Offers {
	return &Offers{
		Items: items,
	}
}

func (o *Offers) GetOfferKeys() *keys.OfferKeys {
	var offerKeys []keys.OfferKey
	for _, offer := range o.Items {
		offerKeys = append(offerKeys, offer.GetKey())
	}
	return keys.NewOfferKeys(offerKeys)
}
