package domain

type OfferCompleted struct {
	Type    OfferEvent `json:"type"`
	Version int        `json:"version"`
}

func NewOfferCompleted() *OfferCompleted {
	return &OfferCompleted{
		Type:    OfferCompletedEvent,
		Version: 1,
	}
}

func (o *OfferCompleted) GetType() OfferEvent {
	return o.Type
}

func (o *OfferCompleted) GetVersion() int {
	return o.Version
}

var _ Event = &OfferCompleted{}
