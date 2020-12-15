package model

type Offers struct {
	Items []*Offer
}

func NewOffers(items []*Offer) *Offers {
	return &Offers{
		Items: items,
	}
}

func (o *Offers) GetOfferKeys() *OfferKeys {
	var offerKeys []OfferKey
	for _, offer := range o.Items {
		offerKeys = append(offerKeys, offer.GetKey())
	}
	return NewOfferKeys(offerKeys)
}
