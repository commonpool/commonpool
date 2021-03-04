package domain

type OfferApproved struct {
	Type    OfferEvent `json:"type"`
	Version int        `json:"version"`
}

func NewOfferApproved() *OfferApproved {
	return &OfferApproved{
		Type:    OfferApprovedEvent,
		Version: 1,
	}
}

func (o *OfferApproved) GetType() OfferEvent {
	return o.Type
}

func (o *OfferApproved) GetVersion() int {
	return o.Version
}

var _ Event = &OfferApproved{}
