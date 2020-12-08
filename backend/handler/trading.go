package handler

import (
	"fmt"
	errs "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/group"
	. "github.com/commonpool/backend/model"
	"github.com/commonpool/backend/trading"
	"github.com/commonpool/backend/web"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (h *Handler) HandleSendOffer(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "HandleSendOffer")

	var err error

	req := web.SendOfferRequest{}
	if err = c.Bind(&req); err != nil {
		response := errs.ErrCreateResourceBadRequest(err)
		return NewErrResponse(c, &response)
	}

	if err = c.Validate(req); err != nil {
		return c.JSON(400, errs.NewValidError(err.(validator.ValidationErrors)))
	}

	var tradingOfferItems []trading.OfferItem2
	for _, tradingOfferItem := range req.Offer.Items {
		itemKey := NewOfferItemKey(uuid.NewV4())
		tradingOfferItem, err := mapNewOfferItem(tradingOfferItem, itemKey)
		if err != nil {
			return errs.ReturnException(c, err)
		}
		tradingOfferItems = append(tradingOfferItems, tradingOfferItem)
	}

	groupKey, err := group.ParseGroupKey(req.Offer.GroupID)
	if err != nil {
		return errs.ReturnException(c, err)
	}

	offer, offerItems, err := h.tradingService.SendOffer(ctx, groupKey, trading.NewOfferItems(tradingOfferItems), "")
	if err != nil {
		return errs.ReturnException(c, err)
	}

	webOffer, err := h.mapToWebOffer(offer, offerItems)
	if err != nil {
		return errs.ReturnException(c, err)
	}

	return c.JSON(http.StatusCreated, &web.GetOfferResponse{
		Offer: webOffer,
	})

}

func (h *Handler) HandleAcceptOffer(c echo.Context) error {

	ctx, l := GetEchoContext(c, "HandleAcceptOffer")

	offerKey, err := ParseOfferKey(c.Param("id"))
	if err != nil {
		l.Error("cannot parse offer key", zap.Error(err))
		return NewErrResponse(c, err)
	}

	_, err = h.tradingService.AcceptOffer(ctx, trading.NewAcceptOffer(offerKey))

	if err != nil {
		l.Error("could not accept offer", zap.Error(err))
		return err
	}

	return c.String(http.StatusOK, "")

}

func (h *Handler) HandleDeclineOffer(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "HandleDeclineOffer")

	offerKey, err := ParseOfferKey(c.Param("id"))
	if err != nil {
		return NewErrResponse(c, err)
	}

	err = h.tradingService.DeclineOffer(ctx, offerKey)
	if err != nil {
		return err
	}

	webOffer, err := h.getWebOffer(offerKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, webOffer)

}

func (h *Handler) HandleGetOffers(c echo.Context) error {

	authSession := h.authorization.GetAuthUserSession(c)
	loggedInUserKey := authSession.GetUserKey()

	result, err := h.tradingStore.GetOffersForUser(loggedInUserKey)
	if err != nil {
		return NewErrResponse(c, err)
	}

	var webOffers []web.Offer
	resultItems := result.Items
	for _, item := range resultItems {

		items, err := h.tradingStore.GetOfferItemsForOffer(item.Offer.GetKey())
		if err != nil {
			return NewErrResponse(c, err)
		}

		webOffer, err := h.mapToWebOffer(item.Offer, items)
		if err != nil {
			return NewErrResponse(c, err)
		}

		webOffers = append(webOffers, *webOffer)

	}

	return c.JSON(http.StatusOK, web.GetOffersResponse{
		Offers: webOffers,
	})

}

func (h *Handler) HandleGetOffer(c echo.Context) error {

	var err error

	offerIdStr := c.Param("id")
	offerId, err := uuid.FromString(offerIdStr)
	if err != nil {
		return NewErrResponse(c, err)
	}
	offerKey := NewOfferKey(offerId)

	offer, err := h.getWebOffer(offerKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, offer)

}

func (h *Handler) HandleConfirmServiceProvided(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "HandleConfirmServiceProvided")

	offerItemKey, err := ParseOfferItemKey(c.Param("id"))
	if err != nil {
		return err
	}

	err = h.tradingService.ConfirmServiceProvided(ctx, offerItemKey)
	if err != nil {
		return err
	}

	offerItem, err := h.tradingService.GetOfferItem(ctx, offerItemKey)
	if err != nil {
		return errs.ReturnException(c, err)
	}

	webResponse, err := mapWebOfferItem(offerItem)
	if err != nil {
		return errs.ReturnException(c, err)
	}

	return c.JSON(http.StatusOK, webResponse)

}

func (h *Handler) HandleConfirmResourceTransferred(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "HandleConfirmResourceTransferred")

	offerItemKey, err := ParseOfferItemKey(c.Param("id"))
	if err != nil {
		return err
	}

	err = h.tradingService.ConfirmResourceTransferred(ctx, offerItemKey)
	if err != nil {
		return err
	}

	offerItem, err := h.tradingService.GetOfferItem(ctx, offerItemKey)
	if err != nil {
		return errs.ReturnException(c, err)
	}

	webResponse, err := mapWebOfferItem(offerItem)
	if err != nil {
		return errs.ReturnException(c, err)
	}

	return c.JSON(http.StatusOK, webResponse)

}

func (h *Handler) HandleConfirmResourceBorrowed(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "HandleConfirmResourceBorrowed")

	offerItemKey, err := ParseOfferItemKey(c.Param("id"))
	if err != nil {
		return err
	}

	err = h.tradingService.ConfirmResourceBorrowed(ctx, offerItemKey)
	if err != nil {
		return err
	}

	offerItem, err := h.tradingService.GetOfferItem(ctx, offerItemKey)
	if err != nil {
		return errs.ReturnException(c, err)
	}

	webResponse, err := mapWebOfferItem(offerItem)
	if err != nil {
		return errs.ReturnException(c, err)
	}

	return c.JSON(http.StatusOK, webResponse)

}

func (h *Handler) HandleConfirmBorrowedResourceReturned(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "HandleConfirmBorrowedResourceReturned")

	offerItemKey, err := ParseOfferItemKey(c.Param("id"))
	if err != nil {
		return err
	}

	err = h.tradingService.ConfirmBorrowedResourceReturned(ctx, offerItemKey)
	if err != nil {
		return err
	}

	offerItem, err := h.tradingService.GetOfferItem(ctx, offerItemKey)
	if err != nil {
		return errs.ReturnException(c, err)
	}

	webResponse, err := mapWebOfferItem(offerItem)
	if err != nil {
		return errs.ReturnException(c, err)
	}

	return c.JSON(http.StatusOK, webResponse)

}

func (h *Handler) GetTradingHistory(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "GetTradingHistory")

	req := web.GetTradingHistoryRequest{}
	if err := c.Bind(&req); err != nil {
		return err
	}

	var userKeys []UserKey
	for _, userId := range req.UserIDs {
		userKey := NewUserKey(userId)
		userKeys = append(userKeys, userKey)
	}

	tradingHistory, err := h.tradingService.GetTradingHistory(ctx, NewUserKeys(userKeys))
	if err != nil {
		return err
	}

	tradingUserKeys := NewUserKeys([]UserKey{})
	for _, entry := range tradingHistory {
		tradingUserKeys = tradingUserKeys.Append(entry.ToUserID)
		tradingUserKeys = tradingUserKeys.Append(entry.FromUserID)
	}

	users, err := h.authStore.GetByKeys(ctx, tradingUserKeys.Items)
	if err != nil {
		return err
	}

	var responseEntries []web.TradingHistoryEntry
	for _, entry := range tradingHistory {
		var resourceId *string
		if entry.ResourceID != nil {
			resourceIdStr := entry.ResourceID.String()
			resourceId = &resourceIdStr
		}
		fromUser, err := users.GetUser(entry.FromUserID)
		if err != nil {
			return err
		}
		toUser, err := users.GetUser(entry.ToUserID)
		if err != nil {
			return err
		}
		webEntry := web.TradingHistoryEntry{
			Timestamp:         entry.Timestamp.String(),
			FromUserID:        entry.FromUserID.String(),
			FromUsername:      fromUser.Username,
			ToUserID:          entry.ToUserID.String(),
			ToUsername:        toUser.Username,
			ResourceID:        resourceId,
			TimeAmountSeconds: entry.TimeAmountSeconds,
		}
		responseEntries = append(responseEntries, webEntry)
	}

	return c.JSON(http.StatusOK, web.GetTradingHistoryResponse{
		Entries: responseEntries,
	})
}

func mapWebOfferItem(offerItem trading.OfferItem2) (*web.OfferItem, error) {

	if offerItem.IsCreditTransfer() {

		creditTransfer := offerItem.(*trading.CreditTransferItem)

		from, err := web.MapOfferItemTarget(creditTransfer.From)
		if err != nil {
			return nil, err
		}
		to, err := web.MapOfferItemTarget(creditTransfer.To)
		if err != nil {
			return nil, err
		}

		return &web.OfferItem{
			ID:   creditTransfer.Key.String(),
			From: from,
			To:   to,
			Type: trading.CreditTransfer,
		}, nil

	} else if offerItem.IsBorrowingResource() {

		borrowResource := offerItem.(*trading.BorrowResourceItem)

		to, err := web.MapOfferItemTarget(borrowResource.To)
		if err != nil {
			return nil, err
		}

		resourceId := borrowResource.ResourceKey.String()
		duration := int64(borrowResource.Duration.Seconds())
		return &web.OfferItem{
			ID:         borrowResource.Key.String(),
			To:         to,
			ResourceId: &resourceId,
			Duration:   &duration,
			Type:       trading.BorrowResource,
		}, nil

	} else if offerItem.IsResourceTransfer() {

		resourceTransfer := offerItem.(*trading.ResourceTransferItem)

		to, err := web.MapOfferItemTarget(resourceTransfer.To)
		if err != nil {
			return nil, err
		}

		resourceId := resourceTransfer.ResourceKey.String()
		return &web.OfferItem{
			ID:         resourceTransfer.Key.String(),
			To:         to,
			ResourceId: &resourceId,
			Type:       trading.ResourceTransfer,
		}, nil

	} else if offerItem.IsServiceProviding() {

		serviceProvision := offerItem.(*trading.ProvideServiceItem)

		to, err := web.MapOfferItemTarget(serviceProvision.To)
		if err != nil {
			return nil, err
		}

		resourceId := serviceProvision.ResourceKey.String()
		duration := int64(serviceProvision.Duration.Seconds())
		return &web.OfferItem{
			ID:         serviceProvision.Key.String(),
			To:         to,
			ResourceId: &resourceId,
			Duration:   &duration,
			Type:       trading.ProvideService,
		}, nil
	} else {
		return nil, fmt.Errorf("unexpected offer item type")
	}

}

func (h *Handler) getWebOffer(offerKey OfferKey) (*web.GetOfferResponse, error) {

	offer, err := h.tradingStore.GetOffer(offerKey)
	if err != nil {
		return nil, err
	}

	items, err := h.tradingStore.GetOfferItemsForOffer(offerKey)
	if err != nil {
		return nil, err
	}

	webOffer, err := h.mapToWebOffer(offer, items)
	if err != nil {
		return nil, err
	}

	response := web.GetOfferResponse{
		Offer: webOffer,
	}

	return &response, nil
}

func (h *Handler) mapToWebOffer(offer *trading.Offer, items *trading.OfferItems) (*web.Offer, error) {

	authorUsername, err := h.authStore.GetUsername(offer.GetAuthorKey())
	if err != nil {
		return nil, err
	}

	webOffer := web.Offer{
		ID:             offer.Key.String(),
		CreatedAt:      offer.CreatedAt,
		CompletedAt:    offer.CompletedAt,
		Status:         offer.Status,
		Items:          nil,
		AuthorID:       offer.CreatedByKey.String(),
		AuthorUsername: authorUsername,
	}

	var responseItems []*web.OfferItem

	for _, offerItem := range items.Items {
		webOfferItem, err := mapWebOfferItem(offerItem)
		if err != nil {
			return nil, err
		}
		responseItems = append(responseItems, webOfferItem)
	}

	webOffer.Items = responseItems

	return &webOffer, nil
}

func mapNewOfferItem(tradingOfferItem web.SendOfferPayloadItem, itemKey OfferItemKey) (trading.OfferItem2, error) {

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

func mapCreateBorrowItem(tradingOfferItem web.SendOfferPayloadItem, itemKey OfferItemKey) (*trading.BorrowResourceItem, error) {
	to, err := web.MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, err
	}

	resourceKey, err := ParseResourceKey(*tradingOfferItem.ResourceId)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(int64(time.Second) * *tradingOfferItem.Duration)

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

func mapCreateProvideServiceItem(tradingOfferItem web.SendOfferPayloadItem, itemKey OfferItemKey) (*trading.ProvideServiceItem, error) {
	to, err := web.MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, err
	}

	resourceKey, err := ParseResourceKey(*tradingOfferItem.ResourceId)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(int64(time.Second) * *tradingOfferItem.Duration)

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

func mapCreateResourceTransferItem(tradingOfferItem web.SendOfferPayloadItem, itemKey OfferItemKey) (*trading.ResourceTransferItem, error) {

	to, err := web.MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, err
	}

	resourceKey, err := ParseResourceKey(*tradingOfferItem.ResourceId)
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

func mapCreateCreditTransferItem(tradingOfferItem web.SendOfferPayloadItem, itemKey OfferItemKey) (*trading.CreditTransferItem, error) {
	to, err := web.MapWebOfferItemTarget(tradingOfferItem.To)
	if err != nil {
		return nil, err
	}

	from, err := web.MapWebOfferItemTarget(*tradingOfferItem.From)
	if err != nil {
		return nil, err
	}

	amount := time.Duration(int64(time.Second) * *tradingOfferItem.Amount)

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
