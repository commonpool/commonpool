package trading

import "encoding/json"

type OfferItemApprovals struct {
	Items []*OfferItemApproval
}

func NewOfferItemApprovals(approvals ...*OfferItemApproval) *OfferItemApprovals {
	return &OfferItemApprovals{
		Items: approvals,
	}
}

func (a OfferItemApprovals) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Items)
}

func (a *OfferItemApprovals) UnmarshalJSON(bytes []byte) error {
	var items []*OfferItemApproval
	if err := json.Unmarshal(bytes, &items); err != nil {
		return err
	}
	a.Items = items
	return nil
}
