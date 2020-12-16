package handler

import (
	fmt "fmt"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/group"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	"github.com/commonpool/backend/pkg/trading"
	usermodel "github.com/commonpool/backend/pkg/user/usermodel"
	"github.com/labstack/echo/v4"
	"time"
)

func parseTargetFromQueryParams(c echo.Context, typeQueryParam string, valueQueryParam string) (*resourcemodel.Target, error) {
	typeParam := c.QueryParams().Get(typeQueryParam)
	if typeParam != "" {
		typeValue, err := resourcemodel.ParseOfferItemTargetType(typeParam)
		if err != nil {
			return nil, err
		}
		targetType := &typeValue
		targetIdStr := c.QueryParams().Get(valueQueryParam)
		if targetIdStr == "" {
			return nil, exceptions.ErrQueryParamRequired(valueQueryParam)
		}
		if targetType.IsGroup() {
			groupKey, err := group.ParseGroupKey(targetIdStr)
			if err != nil {
				return nil, err
			}
			return resourcemodel.NewGroupTarget(groupKey), nil
		} else if targetType.IsUser() {
			userKey := usermodel.NewUserKey(targetIdStr)
			return resourcemodel.NewUserTarget(userKey), nil
		}
	}
	return nil, nil
}

func mapWebOfferItem(offerItem trading.OfferItem, approvers *trading.OfferApprovers) (*OfferItem, error) {

	fromApprovers, hasFromApprovers := approvers.UsersAbleToGiveItem[offerItem.GetKey()]
	toApprovers, hasToApprovers := approvers.UsersAbleToReceiveItem[offerItem.GetKey()]

	if !hasFromApprovers {
		fromApprovers = usermodel.NewEmptyUserKeys()
	}
	if !hasToApprovers {
		toApprovers = usermodel.NewEmptyUserKeys()
	}

	if offerItem.IsCreditTransfer() {

		creditTransfer := offerItem.(*trading.CreditTransferItem)

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
			Type:               trading.CreditTransfer,
			ReceivingApprovers: toApprovers.Strings(),
			GivingApprovers:    fromApprovers.Strings(),
			GiverApproved:      creditTransfer.GiverAccepted,
			ReceiverApproved:   creditTransfer.ReceiverAccepted,
			Amount:             &amount,
		}, nil

	} else if offerItem.IsBorrowingResource() {

		borrowResource := offerItem.(*trading.BorrowResourceItem)

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
			Type:               trading.BorrowResource,
			ReceivingApprovers: toApprovers.Strings(),
			GivingApprovers:    fromApprovers.Strings(),
			GiverApproved:      borrowResource.GiverAccepted,
			ReceiverApproved:   borrowResource.ReceiverAccepted,
			ItemGiven:          borrowResource.ItemGiven,
			ItemTaken:          borrowResource.ItemTaken,
			ItemReceivedBack:   borrowResource.ItemReceivedBack,
			ItemReturnedBack:   borrowResource.ItemReturnedBack,
		}, nil

	} else if offerItem.IsResourceTransfer() {

		resourceTransfer := offerItem.(*trading.ResourceTransferItem)

		to, err := MapOfferItemTarget(resourceTransfer.To)
		if err != nil {
			return nil, err
		}

		resourceId := resourceTransfer.ResourceKey.String()
		return &OfferItem{
			ID:                 resourceTransfer.Key.String(),
			To:                 to,
			ResourceId:         &resourceId,
			Type:               trading.ResourceTransfer,
			ReceivingApprovers: toApprovers.Strings(),
			GivingApprovers:    fromApprovers.Strings(),
			GiverApproved:      resourceTransfer.GiverAccepted,
			ReceiverApproved:   resourceTransfer.ReceiverAccepted,
			ItemGiven:          resourceTransfer.ItemGiven,
			ItemTaken:          resourceTransfer.ItemReceived,
		}, nil

	} else if offerItem.IsServiceProviding() {

		serviceProvision := offerItem.(*trading.ProvideServiceItem)

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
			Type:                        trading.ProvideService,
			ReceivingApprovers:          toApprovers.Strings(),
			GivingApprovers:             fromApprovers.Strings(),
			GiverApproved:               serviceProvision.GiverAccepted,
			ReceiverApproved:            serviceProvision.ReceiverAccepted,
			ServiceGivenConfirmation:    serviceProvision.ServiceGivenConfirmation,
			ServiceReceivedConfirmation: serviceProvision.ServiceReceivedConfirmation,
		}, nil
	} else {
		return nil, fmt.Errorf("unexpected offer item type")
	}

}

func (h *TradingHandler) getWebOffer(offerKey trading.OfferKey) (*OfferResponse, error) {

	offer, err := h.tradingService.GetOffer(offerKey)
	if err != nil {
		return nil, err
	}

	items, err := h.tradingService.GetOfferItemsForOffer(offerKey)
	if err != nil {
		return nil, err
	}

	approvers, err := h.tradingService.FindApproversForOffer(offer.Key)
	if err != nil {
		return nil, err
	}

	webOffer, err := h.mapToWebOffer(offer, items, approvers)
	if err != nil {
		return nil, err
	}

	response := OfferResponse{
		Offer: webOffer,
	}

	return &response, nil
}

func (h *TradingHandler) mapToWebOffer(offer *trading.Offer, items *trading.OfferItems, approvers *trading.OfferApprovers) (*Offer, error) {

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

func mapNewOfferItem(tradingOfferItem SendOfferPayloadItem, itemKey trading.OfferItemKey) (trading.OfferItem, error) {

	itemType := tradingOfferItem.Type

	if itemType == trading.CreditTransfer {

		return mapCreateCreditTransferItem(tradingOfferItem, itemKey)

	} else if itemType == trading.ResourceTransfer {

		return mapCreateResourceTransferItem(tradingOfferItem, itemKey)

	} else if itemType == trading.ProvideService {

		return mapCreateProvideServiceItem(tradingOfferItem, itemKey)

	} else if itemType == trading.BorrowResource {

		return mapCreateBorrowItem(tradingOfferItem, itemKey)

	} else {

		return nil, fmt.Errorf("unexpected item type: %s", itemType)

	}
}

func mapCreateBorrowItem(tradingOfferItem SendOfferPayloadItem, itemKey trading.OfferItemKey) (*trading.BorrowResourceItem, error) {
	to, err := MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, err
	}

	resourceKey, err := resourcemodel.ParseResourceKey(*tradingOfferItem.ResourceId)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(*tradingOfferItem.Duration)
	if err != nil {
		return nil, err
	}

	return &trading.BorrowResourceItem{
		OfferItemBase: trading.OfferItemBase{
			Type: trading.BorrowResource,
			Key:  itemKey,
			To:   to,
		},
		ResourceKey: resourceKey,
		Duration:    duration,
	}, nil
}

func mapCreateProvideServiceItem(tradingOfferItem SendOfferPayloadItem, itemKey trading.OfferItemKey) (*trading.ProvideServiceItem, error) {
	to, err := MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, err
	}

	resourceKey, err := resourcemodel.ParseResourceKey(*tradingOfferItem.ResourceId)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(*tradingOfferItem.Duration)
	if err != nil {
		return nil, err
	}

	return &trading.ProvideServiceItem{
		OfferItemBase: trading.OfferItemBase{
			Type: trading.ProvideService,
			Key:  itemKey,
			To:   to,
		},
		ResourceKey: resourceKey,
		Duration:    duration,
	}, nil
}

func mapCreateResourceTransferItem(tradingOfferItem SendOfferPayloadItem, itemKey trading.OfferItemKey) (*trading.ResourceTransferItem, error) {

	to, err := MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, err
	}

	resourceKey, err := resourcemodel.ParseResourceKey(*tradingOfferItem.ResourceId)
	if err != nil {
		return nil, err
	}

	return &trading.ResourceTransferItem{
		OfferItemBase: trading.OfferItemBase{
			Type: trading.ResourceTransfer,
			Key:  itemKey,
			To:   to,
		},
		ResourceKey: resourceKey,
	}, nil
}

func mapCreateCreditTransferItem(tradingOfferItem SendOfferPayloadItem, itemKey trading.OfferItemKey) (*trading.CreditTransferItem, error) {
	to, err := MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, err
	}

	from, err := MapWebOfferItemTarget(*tradingOfferItem.From)
	if err != nil {
		return nil, err
	}

	amount, err := time.ParseDuration(*tradingOfferItem.Amount)
	if err != nil {
		return nil, err
	}

	return &trading.CreditTransferItem{
		OfferItemBase: trading.OfferItemBase{
			Type: trading.CreditTransfer,
			Key:  itemKey,
			To:   to,
		},
		From:   from,
		Amount: amount,
	}, nil
}
