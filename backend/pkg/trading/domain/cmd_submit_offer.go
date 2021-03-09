package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/commands"
	"github.com/commonpool/backend/pkg/keys"
)

type SubmitOfferPayload struct {
	GroupKey   keys.GroupKey    `json:"group_key"`
	OfferItems SubmitOfferItems `json:"offer_items"`
}

type SubmitOffer struct {
	commands.CommandEnvelope
	Payload SubmitOfferPayload `json:"payload"`
}

func NewPostOffer(
	ctx context.Context,
	offerKey keys.OfferKey,
	groupKey keys.GroupKey,
	offerItems SubmitOfferItems) SubmitOffer {
	return SubmitOffer{
		commands.NewCommandEnvelope(ctx, SubmitOfferCommand, "offer", offerKey.String()),
		SubmitOfferPayload{
			groupKey,
			offerItems,
		},
	}
}
