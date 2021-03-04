package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

type OfferRepository interface {
	Load(ctx context.Context, offerKey keys.OfferKey) (*Offer, error)
	Save(ctx context.Context, offer *Offer) error
}
