package handler

import (
	fmt "fmt"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/pkg/trading/domain"
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

func mapWebOfferItem(offerItem domain.OfferItem, approvers trading.Approvers) (*OfferItem, error) {

	outboundApprovers := approvers.GetOutboundApprovers(offerItem.GetKey())
	inboundApprovers := approvers.GetInboundApprovers(offerItem.GetKey())

	if offerItem.IsCreditTransfer() {

		creditTransfer := offerItem.(*domain.CreditTransferItem)

		from, err := MapOfferItemTarget(creditTransfer.From)
		if err != nil {
			return nil, err
		}
		to, err := MapOfferItemTarget(creditTransfer.To)
		if err != nil {
			return nil, err
		}

		amount := int64(creditTransfer.Amount.Seconds())
		return &OfferItem{
			ID:                 creditTransfer.Key.String(),
			From:               from,
			To:                 to,
			Type:               domain.CreditTransfer,
			ReceivingApprovers: inboundApprovers.Strings(),
			GivingApprovers:    outboundApprovers.Strings(),
			GiverApproved:      creditTransfer.ApprovedOutbound,
			ReceiverApproved:   creditTransfer.ApprovedInbound,
			Amount:             &amount,
		}, nil

	} else if offerItem.IsBorrowingResource() {

		borrowResource := offerItem.(*domain.BorrowResourceItem)

		to, err := MapOfferItemTarget(borrowResource.To)
		if err != nil {
			return nil, err
		}

		resourceId := borrowResource.ResourceKey.String()
		duration := int64(borrowResource.Duration.Seconds())
		return &OfferItem{
			ID:                 borrowResource.Key.String(),
			To:                 to,
			ResourceId:         &resourceId,
			Duration:           &duration,
			Type:               domain.BorrowResource,
			ReceivingApprovers: inboundApprovers.Strings(),
			GivingApprovers:    outboundApprovers.Strings(),
			GiverApproved:      borrowResource.ApprovedOutbound,
			ReceiverApproved:   borrowResource.ApprovedInbound,
			ItemGiven:          borrowResource.ItemGiven,
			ItemTaken:          borrowResource.ItemTaken,
			ItemReceivedBack:   borrowResource.ItemReceivedBack,
			ItemReturnedBack:   borrowResource.ItemReturnedBack,
		}, nil

	} else if offerItem.IsResourceTransfer() {

		resourceTransfer := offerItem.(*domain.ResourceTransferItem)

		to, err := MapOfferItemTarget(resourceTransfer.To)
		if err != nil {
			return nil, err
		}

		resourceId := resourceTransfer.ResourceKey.String()
		return &OfferItem{
			ID:                 resourceTransfer.Key.String(),
			To:                 to,
			ResourceId:         &resourceId,
			Type:               domain.ResourceTransfer,
			ReceivingApprovers: inboundApprovers.Strings(),
			GivingApprovers:    outboundApprovers.Strings(),
			GiverApproved:      resourceTransfer.ApprovedOutbound,
			ReceiverApproved:   resourceTransfer.ApprovedInbound,
			ItemGiven:          resourceTransfer.ItemGiven,
			ItemTaken:          resourceTransfer.ItemReceived,
		}, nil

	} else if offerItem.IsServiceProviding() {

		serviceProvision := offerItem.(*domain.ProvideServiceItem)

		to, err := MapOfferItemTarget(serviceProvision.To)
		if err != nil {
			return nil, err
		}

		resourceId := serviceProvision.ResourceKey.String()
		duration := int64(serviceProvision.Duration.Seconds())
		return &OfferItem{
			ID:                          serviceProvision.Key.String(),
			To:                          to,
			ResourceId:                  &resourceId,
			Duration:                    &duration,
			Type:                        domain.ProvideService,
			ReceivingApprovers:          inboundApprovers.Strings(),
			GivingApprovers:             outboundApprovers.Strings(),
			GiverApproved:               serviceProvision.ApprovedOutbound,
			ReceiverApproved:            serviceProvision.ApprovedInbound,
			ServiceGivenConfirmation:    serviceProvision.ServiceGivenConfirmation,
			ServiceReceivedConfirmation: serviceProvision.ServiceReceivedConfirmation,
		}, nil
	} else {
		return nil, fmt.Errorf("unexpected offer item type")
	}

}

func (h *TradingHandler) mapToWebOffer(offer *trading.Offer, items *domain.OfferItems, approvers trading.Approvers) (*Offer, error) {

	authorUsername, err := h.userService.GetUsername(offer.GetAuthorKey())
	if err != nil {
		return nil, err
	}

	var responseItems []*OfferItem
	for _, offerItem := range items.Items {
		if offerItem.GetOfferKey() != offer.GetKey() {
			continue
		}
		webOfferItem, err := mapWebOfferItem(offerItem, approvers)
		if err != nil {
			return nil, err
		}
		responseItems = append(responseItems, webOfferItem)
	}

	webOffer := Offer{
		ID:             offer.Key.String(),
		CreatedAt:      offer.CreatedAt,
		CompletedAt:    offer.CompletedAt,
		Status:         offer.Status,
		Items:          responseItems,
		AuthorID:       offer.CreatedByKey.String(),
		AuthorUsername: authorUsername,
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
