package handler

import (
	fmt "fmt"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/pkg/trading/domain"
	readmodels "github.com/commonpool/backend/pkg/trading/readmodels"
	"github.com/labstack/echo/v4"
	"time"
)

func parseTargetFromQueryParams(c echo.Context, typeQueryParam string, valueQueryParam string) (*domain.Target, error) {
	typeParam := c.QueryParams().Get(typeQueryParam)
	if typeParam != "" {
		typeValue, err := domain.ParseOfferItemTargetType(typeParam)
		if err != nil {
			return nil, err
		}
		targetType := &typeValue
		targetIdStr := c.QueryParams().Get(valueQueryParam)
		if targetIdStr == "" {
			return nil, exceptions.ErrQueryParamRequired(valueQueryParam)
		}
		if targetType.IsGroup() {
			groupKey, err := keys.ParseGroupKey(targetIdStr)
			if err != nil {
				return nil, err
			}
			return domain.NewGroupTarget(groupKey), nil
		} else if targetType.IsUser() {
			userKey := keys.NewUserKey(targetIdStr)
			return domain.NewUserTarget(userKey), nil
		}
	}
	return nil, nil
}

func mapWebOfferItem(offerItem *readmodels.OfferItemReadModel2, approvers trading.Approvers) (*OfferItem, error) {

	outboundApprovers := approvers.GetOutboundApprovers(offerItem.OfferItemKey)
	inboundApprovers := approvers.GetInboundApprovers(offerItem.OfferItemKey)

	if offerItem.Type == domain.ProvideService {

		from, err := MapOfferItemTarget(offerItem.FromType, offerItem.FromID)
		if err != nil {
			return nil, err
		}
		to, err := MapOfferItemTarget(offerItem.ToType, offerItem.ToID)
		if err != nil {
			return nil, err
		}

		amount := int64(offerItem.Amount.Seconds())
		return &OfferItem{
			ID:                 offerItem.OfferItemKey,
			From:               from,
			To:                 to,
			Type:               domain.CreditTransfer,
			ReceivingApprovers: inboundApprovers.Strings(),
			GivingApprovers:    outboundApprovers.Strings(),
			GiverApproved:      offerItem.ApprovedOutbound,
			ReceiverApproved:   offerItem.ApprovedInbound,
			Amount:             &amount,
		}, nil

	} else if offerItem.Type == string(domain.BorrowResource) {

		to, err := MapOfferItemTarget(offerItem.ToType, offerItem.ToID)
		if err != nil {
			return nil, err
		}
		duration := int64(offerItem.Duration.Seconds())
		return &OfferItem{
			ID:                 offerItem.OfferItemKey,
			To:                 to,
			ResourceId:         &offerItem.ResourceID,
			Duration:           &duration,
			Type:               domain.BorrowResource,
			ReceivingApprovers: inboundApprovers.Strings(),
			GivingApprovers:    outboundApprovers.Strings(),
			GiverApproved:      offerItem.ApprovedOutbound,
			ReceiverApproved:   offerItem.ApprovedInbound,
			ItemGiven:          offerItem.ResourceLent,
			ItemTaken:          offerItem.ResourceBorrowed,
			ItemReceivedBack:   offerItem.LentItemReceived,
			ItemReturnedBack:   offerItem.BorrowedItemReturned,
		}, nil

	} else if offerItem.Type == string(domain.ResourceTransfer) {

		to, err := MapOfferItemTarget(offerItem.ToType, offerItem.ToID)
		if err != nil {
			return nil, err
		}

		return &OfferItem{
			ID:                 offerItem.OfferItemKey,
			To:                 to,
			ResourceId:         &offerItem.ResourceID,
			Type:               domain.ResourceTransfer,
			ReceivingApprovers: inboundApprovers.Strings(),
			GivingApprovers:    outboundApprovers.Strings(),
			GiverApproved:      offerItem.ApprovedOutbound,
			ReceiverApproved:   offerItem.ApprovedInbound,
			ItemGiven:          offerItem.ResourceGiven,
			ItemTaken:          offerItem.ResourceTaken,
		}, nil

	} else if offerItem.Type == string(domain.ProvideService) {

		to, err := MapOfferItemTarget(offerItem.ToType, offerItem.ToID)
		if err != nil {
			return nil, err
		}

		duration := int64(offerItem.Duration.Seconds())
		return &OfferItem{
			ID:                          offerItem.OfferItemKey,
			To:                          to,
			ResourceId:                  &offerItem.ResourceID,
			Duration:                    &duration,
			Type:                        domain.ProvideService,
			ReceivingApprovers:          inboundApprovers.Strings(),
			GivingApprovers:             outboundApprovers.Strings(),
			GiverApproved:               offerItem.ApprovedOutbound,
			ReceiverApproved:            offerItem.ApprovedInbound,
			ServiceGivenConfirmation:    offerItem.ServiceGiven,
			ServiceReceivedConfirmation: offerItem.ServiceReceived,
		}, nil
	} else {
		return nil, fmt.Errorf("unexpected offer item type")
	}

}

func (h *TradingHandler) mapToWebOffer(offer *readmodels.OfferReadModel, items []*readmodels.OfferItemReadModel, approvers trading.Approvers) (*Offer, error) {

	var responseItems []*OfferItem
	for _, offerItem := range offer.OfferItems {
		webOfferItem, err := mapWebOfferItem(offerItem, approvers)
		if err != nil {
			return nil, err
		}
		responseItems = append(responseItems, webOfferItem)
	}

	webOffer := Offer{
		ID: offer.ID,
		// TODO: CreatedAt:      offer.CreatedAt,
		CompletedAt:    offer.CompletedAt,
		Status:         offer.Status,
		Items:          responseItems,
		AuthorID:       offer.SubmittedByID,
		AuthorUsername: offer.SubmittedByID,
	}

	return &webOffer, nil
}

func mapNewOfferItem(tradingOfferItem SendOfferPayloadItem, itemKey keys.OfferItemKey) (domain.OfferItem, error) {

	itemType := tradingOfferItem.Type

	if itemType == domain.CreditTransfer {

		res, err := mapCreateCreditTransferItem(tradingOfferItem, itemKey)
		if err != nil {
			return nil, fmt.Errorf("error mapping tradingOfferItem as trading.CreditTransferItem: %v", err)
		}
		return res, nil

	} else if itemType == domain.ResourceTransfer {

		res, err := mapCreateResourceTransferItem(tradingOfferItem, itemKey)
		if err != nil {
			return nil, fmt.Errorf("error mapping tradingOfferItem as trading.ResourceTransferItem: %v", err)
		}
		return res, nil

	} else if itemType == domain.ProvideService {

		res, err := mapCreateProvideServiceItem(tradingOfferItem, itemKey)
		if err != nil {
			return nil, fmt.Errorf("error mapping tradingOfferItem as trading.ProvideServiceItem: %v", err)
		}
		return res, nil

	} else if itemType == domain.BorrowResource {

		res, err := mapCreateBorrowItem(tradingOfferItem, itemKey)
		if err != nil {
			return nil, fmt.Errorf("error mapping tradingOfferItem as trading.BorrowItem: %v", err)
		}
		return res, nil

	} else {

		return nil, fmt.Errorf("unexpected item type: %s", itemType)

	}
}

func mapCreateBorrowItem(tradingOfferItem SendOfferPayloadItem, itemKey keys.OfferItemKey) (*domain.BorrowResourceItem, error) {
	to, err := MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, fmt.Errorf("error mapping 'tradingOfferItem.To' as trading.Target: %v", err)
	}

	resourceKey, err := keys.ParseResourceKey(*tradingOfferItem.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'tradingOfferItem.ResourceId' as keys.ResourceKey: %v", err)
	}

	duration, err := time.ParseDuration(*tradingOfferItem.Duration)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'tradingOfferItem.Duration' as time.Duration: %v", err)
	}

	return &domain.BorrowResourceItem{
		OfferItemBase: domain.OfferItemBase{
			Type: domain.BorrowResource,
			Key:  itemKey,
			To:   to,
		},
		ResourceKey: resourceKey,
		Duration:    duration,
	}, nil
}

func mapCreateProvideServiceItem(tradingOfferItem SendOfferPayloadItem, itemKey keys.OfferItemKey) (*domain.ProvideServiceItem, error) {
	to, err := MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, fmt.Errorf("error mapping 'tradingOfferItem.To' as trading.Target: %v", err)
	}

	resourceKey, err := keys.ParseResourceKey(*tradingOfferItem.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'tradingOfferItem.ResourceId' as keys.ResourceKey: %v", err)
	}

	duration, err := time.ParseDuration(*tradingOfferItem.Duration)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'tradingOfferItem.Duration' as time.Duration: %v", err)
	}

	return &domain.ProvideServiceItem{
		OfferItemBase: domain.OfferItemBase{
			Type: domain.ProvideService,
			Key:  itemKey,
			To:   to,
		},
		ResourceKey: resourceKey,
		Duration:    duration,
	}, nil
}

func mapCreateResourceTransferItem(tradingOfferItem SendOfferPayloadItem, itemKey keys.OfferItemKey) (*domain.ResourceTransferItem, error) {

	to, err := MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, fmt.Errorf("error mapping 'tradingOfferItem.To' to trading.Target: %v", err)
	}

	resourceKey, err := keys.ParseResourceKey(*tradingOfferItem.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'tradingOfferItem.ResourceId' as keys.ResourceKey: %v", err)
	}

	return &domain.ResourceTransferItem{
		OfferItemBase: domain.OfferItemBase{
			Type: domain.ResourceTransfer,
			Key:  itemKey,
			To:   to,
		},
		ResourceKey: resourceKey,
	}, nil
}

func mapCreateCreditTransferItem(tradingOfferItem SendOfferPayloadItem, itemKey keys.OfferItemKey) (*domain.CreditTransferItem, error) {
	to, err := MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, fmt.Errorf("error mapping 'tradingOfferItem.To' as trading.Target: %v", err)
	}

	from, err := MapWebOfferItemTarget(*tradingOfferItem.From)
	if err != nil {
		return nil, fmt.Errorf("error mapping 'tradingOfferItem.From' as trading.Target: %v", err)
	}

	amount, err := time.ParseDuration(*tradingOfferItem.Amount)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'tradingOfferItem.Amount' as time.Duration: %v", err)
	}

	return &domain.CreditTransferItem{
		OfferItemBase: domain.OfferItemBase{
			Type: domain.CreditTransfer,
			Key:  itemKey,
			To:   to,
		},
		From:   from,
		Amount: amount,
	}, nil
}
