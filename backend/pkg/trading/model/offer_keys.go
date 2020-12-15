package model

type OfferKeys struct {
	Items []OfferKey
}

func NewOfferKeys(items []OfferKey) *OfferKeys {
	copied := make([]OfferKey, len(items))
	copy(copied, items)
	return &OfferKeys{
		Items: copied,
	}
}

func (o *OfferKeys) Strings() []string {
	var result []string
	for _, item := range o.Items {
		result = append(result, item.String())
	}
	if result == nil {
		result = []string{}
	}
	return result
}
