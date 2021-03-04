package domain

import (
	"github.com/commonpool/backend/pkg/keys"
)

type OfferDeclined struct {
	Type       OfferEvent   `json:"type"`
	DeclinedBy keys.UserKey `json:"declined_by"`
	Version    int          `json:"version"`
}

func NewOfferDeclined(declinedBy keys.UserKey) *OfferDeclined {
	return &OfferDeclined{
		Type:       OfferDeclinedEvent,
		DeclinedBy: declinedBy,
		Version:    1,
	}
}

func (o *OfferDeclined) GetType() OfferEvent {
	return o.Type
}

func (o *OfferDeclined) GetVersion() int {
	return o.Version
}

var _ Event = &OfferDeclined{}
