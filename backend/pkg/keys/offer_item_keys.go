package keys

type OfferItemKeys struct {
	Items []OfferItemKey
}

func (t *OfferItemKeys) Strings() []string {
	var strings []string
	for _, item := range t.Items {
		strings = append(strings, item.String())
	}
	return strings
}

func NewOfferItemKeys(items []OfferItemKey) *OfferItemKeys {
	copied := make([]OfferItemKey, len(items))
	copy(copied, items)
	return &OfferItemKeys{
		Items: copied,
	}
}

func (t *OfferItemKeys) IsEmpty() bool {
	return len(t.Items) == 0
}
