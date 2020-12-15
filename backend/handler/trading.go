package handler

import (
	"fmt"
	model2 "github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/exceptions"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	tradingmodel "github.com/commonpool/backend/pkg/trading/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"time"
)

func parseTargetFromQueryParams(c echo.Context, typeQueryParam string, valueQueryParam string) (*model2.Target, error) {
	typeParam := c.QueryParams().Get(typeQueryParam)
	if typeParam != "" {
		typeValue, err := model2.ParseOfferItemTargetType(typeParam)
		if err != nil {
			return nil, err
		}
		targetType := &typeValue
		targetIdStr := c.QueryParams().Get(valueQueryParam)
		if targetIdStr == "" {
			return nil, exceptions.ErrQueryParamRequired(valueQueryParam)
		}
		if targetType.IsGroup() {
			groupKey, err := groupmodel.ParseGroupKey(targetIdStr)
			if err != nil {
				return nil, err
			}
			return model2.NewGroupTarget(groupKey), nil
		} else if targetType.IsUser() {
			userKey := usermodel.NewUserKey(targetIdStr)
			return model2.NewUserTarget(userKey), nil
		}
	}
	return nil, nil
}

func mapWebOfferItem(offerItem tradingmodel.OfferItem, approvers *tradingmodel.OfferApprovers) (*web.OfferItem, error) {

	fromApprovers, hasFromApprovers := approvers.UsersAbleToGiveItem[offerItem.GetKey()]
	toApprovers, hasToApprovers := approvers.UsersAbleToReceiveItem[offerItem.GetKey()]

	if !hasFromApprovers {
		fromApprovers = usermodel.NewEmptyUserKeys()
	}
	if !hasToApprovers {
		toApprovers = usermodel.NewEmptyUserKeys()
	}

	if offerItem.IsCreditTransfer() {

		creditTransfer := offerItem.(*tradingmodel.CreditTransferItem)

		from, err := web.MapOfferItemTarget(creditTransfer.From)
		if err != nil {
			return nil, err
		}
		to, err := web.MapOfferItemTarget(creditTransfer.To)
		if err != nil {
			return nil, err
		}

		amount := int64(creditTransfer.Amount.Seconds())
		return &web.OfferItem{
			ID:                 creditTransfer.Key.String(),
			From:               from,
			To:                 to,
			Type:               tradingmodel.CreditTransfer,
			ReceivingApprovers: toApprovers.Strings(),
			GivingApprovers:    fromApprovers.Strings(),
			GiverApproved:      creditTransfer.GiverAccepted,
			ReceiverApproved:   creditTransfer.ReceiverAccepted,
			Amount:             &amount,
		}, nil

	} else if offerItem.IsBorrowingResource() {

		borrowResource := offerItem.(*tradingmodel.BorrowResourceItem)

		to, err := web.MapOfferItemTarget(borrowResource.To)
		if err != nil {
			return nil, err
		}

		resourceId := borrowResource.ResourceKey.String()
		duration := int64(borrowResource.Duration.Seconds())
		return &web.OfferItem{
			ID:                 borrowResource.Key.String(),
			To:                 to,
			ResourceId:         &resourceId,
			Duration:           &duration,
			Type:               tradingmodel.BorrowResource,
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

		resourceTransfer := offerItem.(*tradingmodel.ResourceTransferItem)

		to, err := web.MapOfferItemTarget(resourceTransfer.To)
		if err != nil {
			return nil, err
		}

		resourceId := resourceTransfer.ResourceKey.String()
		return &web.OfferItem{
			ID:                 resourceTransfer.Key.String(),
			To:                 to,
			ResourceId:         &resourceId,
			Type:               tradingmodel.ResourceTransfer,
			ReceivingApprovers: toApprovers.Strings(),
			GivingApprovers:    fromApprovers.Strings(),
			GiverApproved:      resourceTransfer.GiverAccepted,
			ReceiverApproved:   resourceTransfer.ReceiverAccepted,
			ItemGiven:          resourceTransfer.ItemGiven,
			ItemTaken:          resourceTransfer.ItemReceived,
		}, nil

	} else if offerItem.IsServiceProviding() {

		serviceProvision := offerItem.(*tradingmodel.ProvideServiceItem)

		to, err := web.MapOfferItemTarget(serviceProvision.To)
		if err != nil {
			return nil, err
		}

		resourceId := serviceProvision.ResourceKey.String()
		duration := int64(serviceProvision.Duration.Seconds())
		return &web.OfferItem{
			ID:                          serviceProvision.Key.String(),
			To:                          to,
			ResourceId:                  &resourceId,
			Duration:                    &duration,
			Type:                        tradingmodel.ProvideService,
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

func (h *Handler) getWebOffer(offerKey tradingmodel.OfferKey) (*web.GetOfferResponse, error) {

	offer, err := h.tradingStore.GetOffer(offerKey)
	if err != nil {
		return nil, err
	}

	items, err := h.tradingStore.GetOfferItemsForOffer(offerKey)
	if err != nil {
		return nil, err
	}

	approvers, err := h.tradingStore.FindApproversForOffer(offer.Key)
	if err != nil {
		return nil, err
	}

	webOffer, err := h.mapToWebOffer(offer, items, approvers)
	if err != nil {
		return nil, err
	}

	response := web.GetOfferResponse{
		Offer: webOffer,
	}

	return &response, nil
}

func (h *Handler) mapToWebOffer(offer *tradingmodel.Offer, items *tradingmodel.OfferItems, approvers *tradingmodel.OfferApprovers) (*web.Offer, error) {

	authorUsername, err := h.authStore.GetUsername(offer.GetAuthorKey())
	if err != nil {
		return nil, err
	}

	var responseItems []*web.OfferItem
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

	webOffer := web.Offer{
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

func mapNewOfferItem(tradingOfferItem web.SendOfferPayloadItem, itemKey tradingmodel.OfferItemKey) (tradingmodel.OfferItem, error) {

	itemType := tradingOfferItem.Type

	if itemType == tradingmodel.CreditTransfer {

		return mapCreateCreditTransferItem(tradingOfferItem, itemKey)

	} else if itemType == tradingmodel.ResourceTransfer {

		return mapCreateResourceTransferItem(tradingOfferItem, itemKey)

	} else if itemType == tradingmodel.ProvideService {

		return mapCreateProvideServiceItem(tradingOfferItem, itemKey)

	} else if itemType == tradingmodel.BorrowResource {

		return mapCreateBorrowItem(tradingOfferItem, itemKey)

	} else {

		return nil, fmt.Errorf("unexpected item type: %s", itemType)

	}
}

func mapCreateBorrowItem(tradingOfferItem web.SendOfferPayloadItem, itemKey tradingmodel.OfferItemKey) (*tradingmodel.BorrowResourceItem, error) {
	to, err := web.MapWebOfferItemTarget(tradingOfferItem.To)
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

	return &tradingmodel.BorrowResourceItem{
		OfferItemBase: tradingmodel.OfferItemBase{
			Type: tradingmodel.BorrowResource,
			Key:  itemKey,
			To:   to,
		},
		ResourceKey: resourceKey,
		Duration:    duration,
	}, nil
}

func mapCreateProvideServiceItem(tradingOfferItem web.SendOfferPayloadItem, itemKey tradingmodel.OfferItemKey) (*tradingmodel.ProvideServiceItem, error) {
	to, err := web.MapWebOfferItemTarget(tradingOfferItem.To)
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

	return &tradingmodel.ProvideServiceItem{
		OfferItemBase: tradingmodel.OfferItemBase{
			Type: tradingmodel.ProvideService,
			Key:  itemKey,
			To:   to,
		},
		ResourceKey: resourceKey,
		Duration:    duration,
	}, nil
}

func mapCreateResourceTransferItem(tradingOfferItem web.SendOfferPayloadItem, itemKey tradingmodel.OfferItemKey) (*tradingmodel.ResourceTransferItem, error) {

	to, err := web.MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, err
	}

	resourceKey, err := resourcemodel.ParseResourceKey(*tradingOfferItem.ResourceId)
	if err != nil {
		return nil, err
	}

	return &tradingmodel.ResourceTransferItem{
		OfferItemBase: tradingmodel.OfferItemBase{
			Type: tradingmodel.ResourceTransfer,
			Key:  itemKey,
			To:   to,
		},
		ResourceKey: resourceKey,
	}, nil
}

func mapCreateCreditTransferItem(tradingOfferItem web.SendOfferPayloadItem, itemKey tradingmodel.OfferItemKey) (*tradingmodel.CreditTransferItem, error) {
	to, err := web.MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, err
	}

	from, err := web.MapWebOfferItemTarget(*tradingOfferItem.From)
	if err != nil {
		return nil, err
	}

	amount, err := time.ParseDuration(*tradingOfferItem.Amount)
	if err != nil {
		return nil, err
	}

	return &tradingmodel.CreditTransferItem{
		OfferItemBase: tradingmodel.OfferItemBase{
			Type: tradingmodel.CreditTransfer,
			Key:  itemKey,
			To:   to,
		},
		From:   from,
		Amount: amount,
	}, nil
}
