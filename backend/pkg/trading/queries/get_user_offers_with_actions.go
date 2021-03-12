package queries

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	readmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type GetUserOffersWithActions struct {
	db                   *gorm.DB
	getOffersPermissions *GetOffersPermissions
	getOfferKeysForUser  *GetOfferKeysForUser
}

func NewGetUserOffersWithActions(db *gorm.DB, getOffersPermissions *GetOffersPermissions, getOfferKeysForUser *GetOfferKeysForUser) *GetUserOffersWithActions {
	return &GetUserOffersWithActions{
		db:                   db,
		getOffersPermissions: getOffersPermissions,
		getOfferKeysForUser:  getOfferKeysForUser,
	}
}

func (q *GetUserOffersWithActions) Get(ctx context.Context, userKey keys.UserKey) ([]*readmodels.OfferReadModelWithActions, error) {

	offerKeys, err := q.getOfferKeysForUser.Get(ctx, userKey)
	if err != nil {
		return nil, err
	}

	var offerPermissions map[keys.OfferKey]*PermissionTuples
	var offers []*readmodels.OfferReadModel

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var err error
		offerPermissions, err = q.getOffersPermissions.Get(ctx, offerKeys)
		return err
	})

	g.Go(func() error {
		var err error
		offers, err = getOffers(ctx, offerKeys, q.db)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	var result []*readmodels.OfferReadModelWithActions

	for _, offer := range offers {
		permissions := offerPermissions[offer.OfferKey]
		var actions []readmodels.OfferActionReadModel

		var canApprove = false

		if offer.Status == domain.Pending {
			for _, item := range offer.OfferItems {
				if !item.ApprovedOutbound {
					if permissions.Can(userKey, item, domain.Outbound) {
						canApprove = true
					}
				}
				if !item.ApprovedInbound {
					if permissions.Can(userKey, item, domain.Inbound) {
						canApprove = true
					}
				}
			}
		}

		if offer.Status == domain.Pending {
			actions = append(actions, readmodels.OfferActionReadModel{
				Name:         "Approve",
				Enabled:      canApprove,
				Completed:    offer.Status != domain.Pending,
				ActionURL:    fmt.Sprintf("/api/v1/offers/%s/actions/approve", offer.OfferKey.String()),
				OfferItemKey: nil,
				Style:        "success",
			})

			actions = append(actions, readmodels.OfferActionReadModel{
				Name:         "Decline",
				Enabled:      offer.Status == domain.Pending && permissions.CanAny(userKey),
				Completed:    offer.Status != domain.Pending,
				ActionURL:    fmt.Sprintf("/api/v1/offers/%s/actions/decline", offer.OfferKey.String()),
				OfferItemKey: nil,
				Style:        "danger",
			})
		}
		if offer.Status == domain.Approved {
			for _, offerItem := range offer.OfferItems {
				switch offerItem.Type {
				case domain.ResourceTransfer:
					actions = append(actions, readmodels.OfferActionReadModel{
						Name: "Resource was given to " + offerItem.To.GetName(),
						Enabled: offer.Status == domain.Approved &&
							((!offerItem.ResourceGiven && permissions.Can(userKey, offerItem, domain.Outbound)) ||
								(!offerItem.ResourceTaken && permissions.Can(userKey, offerItem, domain.Inbound))),
						Completed:    offerItem.ResourceGiven && offerItem.ResourceTaken,
						ActionURL:    fmt.Sprintf("/api/v1/offers/%s/offer-items/%s/actions/resource-given", offer.OfferKey.String(), offerItem.OfferItemKey.String()),
						Style:        "success",
						OfferItemKey: &offerItem.OfferItemKey,
					})
				case domain.BorrowResource:
					actions = append(actions, readmodels.OfferActionReadModel{
						Name: "Resource was borrowed",
						Enabled: offer.Status == domain.Approved &&
							((!offerItem.ResourceLent && permissions.Can(userKey, offerItem, domain.Outbound)) ||
								(!offerItem.ResourceBorrowed && permissions.Can(userKey, offerItem, domain.Inbound))),
						Completed:    offerItem.ResourceBorrowed && offerItem.ResourceLent,
						ActionURL:    fmt.Sprintf("/api/v1/offers/%s/offer-items/%s/actions/resource-borrowed", offer.OfferKey.String(), offerItem.OfferItemKey.String()),
						Style:        "success",
						OfferItemKey: &offerItem.OfferItemKey,
					})
					actions = append(actions, readmodels.OfferActionReadModel{
						Name: "Resource was returned",
						Enabled: offer.Status == domain.Approved &&
							((offerItem.ResourceLent && !offerItem.LentItemReceived && permissions.Can(userKey, offerItem, domain.Outbound)) ||
								(offerItem.ResourceBorrowed && !offerItem.BorrowedItemReturned && permissions.Can(userKey, offerItem, domain.Inbound))),
						Completed:    offerItem.LentItemReceived && offerItem.BorrowedItemReturned,
						ActionURL:    fmt.Sprintf("/api/v1/offers/%s/offer-items/%s/actions/resource-returned", offer.OfferKey.String(), offerItem.OfferItemKey.String()),
						Style:        "success",
						OfferItemKey: &offerItem.OfferItemKey,
					})
				case domain.ProvideService:
					actions = append(actions, readmodels.OfferActionReadModel{
						Name: "Service was given to " + offerItem.To.GetName(),
						Enabled: offer.Status == domain.Approved &&
							((!offerItem.ServiceGiven && permissions.Can(userKey, offerItem, domain.Outbound)) ||
								(!offerItem.ServiceReceived && permissions.Can(userKey, offerItem, domain.Inbound))),
						Completed:    offerItem.ServiceGiven && offerItem.ServiceReceived,
						ActionURL:    fmt.Sprintf("/api/v1/offers/%s/offer-items/%s/actions/service-given", offer.OfferKey.String(), offerItem.OfferItemKey.String()),
						Style:        "success",
						OfferItemKey: &offerItem.OfferItemKey,
					})
				}
			}
		}
		result = append(result, &readmodels.OfferReadModelWithActions{
			OfferReadModel: offer,
			Actions:        actions,
		})
	}

	if result == nil {
		return []*readmodels.OfferReadModelWithActions{}, nil
	}

	return result, nil
}
