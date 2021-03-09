package commandhandlers

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	queries2 "github.com/commonpool/backend/pkg/resource/queries"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/pkg/trading/queries"
)

type SubmitOfferHandler struct {
	offerRepo          domain.OfferRepository
	getPermission      *queries.GetOfferPermissions
	getResourcesByKeys *queries2.GetResourcesByKeysWithSharings
}

func NewSubmitOfferHandler(
	offerRepo domain.OfferRepository,
	getPermission *queries.GetOfferPermissions,
	getResourcesByKeys *queries2.GetResourcesByKeysWithSharings) *SubmitOfferHandler {
	return &SubmitOfferHandler{
		offerRepo:          offerRepo,
		getPermission:      getPermission,
		getResourcesByKeys: getResourcesByKeys,
	}
}

func (c *SubmitOfferHandler) Execute(ctx context.Context, command domain.SubmitOffer) error {

	offerKey, err := keys.ParseOfferKey(command.AggregateID)
	if err != nil {
		return err
	}

	return doWithOffer(ctx, offerKey, c.offerRepo, c.submitOffer(ctx, command))
}

func (c *SubmitOfferHandler) submitOffer(ctx context.Context, command domain.SubmitOffer) func(offer *domain.Offer) error {
	return func(offer *domain.Offer) error {

		loggedInUser, err := oidc.GetLoggedInUser(ctx)
		if err != nil {
			return err
		}

		err = offer.Submit(loggedInUser.GetUserKey(), command.Payload.GroupKey, command.Payload.OfferItems)
		if err != nil {
			return err
		}

		resourcesByKeys, err := c.getResourcesByKeys.Get(ctx, offer.GetResourceKeys())
		if err != nil {
			return err
		}

		var resourceMap = map[keys.ResourceKey]*readmodel.ResourceWithSharingsReadModel{}
		for _, resource := range resourcesByKeys {
			resourceMap[resource.ResourceKey] = resource
		}

		for _, offerItem := range offer.GetOfferItems().Items {
			roi, ok := offerItem.(domain.ResourceOfferItem)
			if !ok {
				continue
			}
			resource, ok := resourceMap[roi.GetResourceKey()]
			if !ok {
				return exceptions.ErrResourceNotFound
			}

			var found = false
			for _, sharing := range resource.Sharings {
				if sharing.GroupKey == offer.GetGroupKey() {
					found = true
					break
				}
			}
			if !found {
				return exceptions.ErrBadRequest("Resource is not shared with the group")
			}

			if resource.Owner.Equals(roi.GetTo()) {
				return exceptions.ErrBadRequestf("OfferItem Resource destination is the same as the resource owner")
			}
		}

		return nil
	}
}
