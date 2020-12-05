package handler

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	errs "github.com/commonpool/backend/errors"
	. "github.com/commonpool/backend/model"
	resource2 "github.com/commonpool/backend/resource"
	route "github.com/commonpool/backend/router"
	"github.com/commonpool/backend/trading"
	"github.com/commonpool/backend/web"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"net/http"
)

type ValidError struct {
	Tag             string `json:"tag"`
	ActualTag       string `json:"actualTag"`
	Namespace       string `json:"namespace"`
	StructNamespace string `json:"structNamespace"`
	Field           string `json:"field"`
	StructField     string `json:"structField"`
	Param           string `json:"param"`
	Kind            string `json:"kind"`
	Type            string `json:"type"`
}

type ValidErrors struct {
	Errors  []ValidError      `json:"errors"`
	Message string            `json:"message"`
	Trans   map[string]string `json:"trans"`
}

func NewValidError(validerr validator.ValidationErrors) ValidErrors {

	var validErrors []ValidError

	for _, err := range validerr {
		validErrors = append(validErrors, ValidError{
			Tag:             err.Tag(),
			ActualTag:       err.ActualTag(),
			Namespace:       err.Namespace(),
			StructNamespace: err.StructNamespace(),
			Field:           err.Field(),
			StructField:     err.StructField(),
			Param:           err.Param(),
			Kind:            err.Kind().String(),
			Type:            err.Type().String(),
		})
	}

	translation := validerr.Translate(route.Trans)

	if validErrors == nil {
		validErrors = []ValidError{}
	}

	return ValidErrors{
		Message: validerr.Error(),
		Errors:  validErrors,
		Trans:   translation,
	}

}

type PersonOffer struct {
	OfferItem trading.OfferItem
	Resource  resource2.Resource
	Time      int64
	FromUser  auth.User
	ToUser    auth.User
}

func (h *Handler) HandleSendOffer(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "HandleSendOffer")

	var err error

	//  Getting Logged In User
	loggedInUser, err := auth.GetUserSession(ctx)
	if err != nil {
		return NewErrResponse(c, errs.ErrUnauthorized)
	}

	//  Unmarshaling Request Body
	req := web.SendOfferRequest{}
	if err = c.Bind(&req); err != nil {
		response := errs.ErrCreateResourceBadRequest(err)
		return NewErrResponse(c, &response)
	}

	//  Validating Request Body
	if err = c.Validate(req); err != nil {
		return c.JSON(400, NewValidError(err.(validator.ValidationErrors)))
	}

	reference := UserReference(loggedInUser)
	response, err := h.SendOffer(ctx, req, reference)
	if err != nil {
		return errs.ReturnException(c, err)
	}

	return c.JSON(http.StatusCreated, response)

}

func (h *Handler) SendOffer(ctx context.Context, req web.SendOfferRequest, fromUser UserReference) (*web.GetOfferResponse, error) {

	var tradingOfferItems []*trading.OfferItem
	for _, tradingOfferItem := range req.Offer.Items {

		itemKey := NewOfferItemKey(uuid.NewV4())
		fromUser := NewUserKey(tradingOfferItem.From)
		toUser := NewUserKey(tradingOfferItem.To)

		if tradingOfferItem.Type == trading.TimeItem {
			tradingOfferItems = append(tradingOfferItems, trading.NewTimeOfferItem(itemKey, fromUser, toUser, *tradingOfferItem.TimeInSeconds))
		} else if tradingOfferItem.Type == trading.ResourceItem {
			resourceKey, err := ParseResourceKey(*tradingOfferItem.ResourceId)
			if err != nil {
				return nil, err
			}
			tradingOfferItems = append(tradingOfferItems, trading.NewResourceOfferItem(itemKey, fromUser, toUser, resourceKey))
		}
	}

	offer, items, decisions, err := h.tradingService.SendOffer(ctx, trading.NewOfferItems(tradingOfferItems), req.Offer.Message)
	if err != nil {
		return nil, err
	}

	//  Mapping Offer to Web Response
	webOffer, err := h.mapToWebOffer(offer, items, decisions)
	if err != nil {
		return nil, err
	}

	response := web.GetOfferResponse{Offer: *webOffer}

	return &response, nil
}

func (h *Handler) HandleAcceptOffer(c echo.Context) error {

	ctx, l := GetEchoContext(c, "HandleAcceptOffer")

	//  Parsing Offer key from query params
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

func (h *Handler) DeclineOffer(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "DeclineOffer")

	// retrieve current user
	authSession, err := auth.GetUserSession(ctx)
	if err != nil {
		return errs.ErrUnauthorized
	}

	authUserKey := authSession.GetUserKey()

	// retrieve offer key
	offerKey, err := ParseOfferKey(c.Param("id"))
	if err != nil {
		return NewErrResponse(c, err)
	}

	// retrieve offer
	offer, err := h.tradingStore.GetOffer(offerKey)
	if err != nil {
		return NewErrResponse(c, err)
	}

	// can only approve pending offers
	if offer.Status != trading.PendingOffer {
		return NewErrResponse(c, fmt.Errorf("offer is not pending approval"))
	}

	// get offer decisions
	decisions, err := h.tradingStore.GetDecisions(offerKey)
	if err != nil {
		return NewErrResponse(c, err)
	}

	var currentUserDecision *trading.OfferDecision

	// retrieving the current user decision
	// also retrieve is this approval is the last one needed so that everyone approved the offer
	for _, decision := range decisions.Items {
		decisionKey := decision.GetKey()
		decisionUserKey := decisionKey.GetUserKey()
		if decisionUserKey == authUserKey {
			currentUserDecision = decision
			break
		}
	}

	// in the case that a user is approving an offer he's not part of
	if currentUserDecision == nil {
		return fmt.Errorf("not involved in that offer")
	}

	// saving the decision
	err = h.tradingStore.SaveDecision(offerKey, authUserKey, trading.DeclinedDecision)
	if err != nil {
		return NewErrResponse(c, err)
	}

	// complete offer if everyone approved
	err = h.tradingStore.SaveOfferStatus(offerKey, trading.DeclinedOffer)
	if err != nil {
		return NewErrResponse(c, err)
	}

	webOffer, err := h.GetWebOffer(offerKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, webOffer)

}

func (h *Handler) GetOffers(c echo.Context) error {

	//  Retrieving currently logged in user
	authSession := h.authorization.GetAuthUserSession(c)
	loggedInUserKey := authSession.GetUserKey()

	//  Retrieving logged in user offers
	qry := trading.GetOffersQuery{
		UserKeys: []UserKey{loggedInUserKey},
	}
	result, err := h.tradingStore.GetOffers(qry)
	if err != nil {
		return NewErrResponse(c, err)
	}

	//  Mapping offers to web response
	var webOffers []web.Offer
	resultItems := result.Items
	for _, item := range resultItems {

		//  Retrieving decisions for offer
		decisions, err := h.tradingStore.GetDecisions(item.Offer.GetKey())
		if err != nil {
			return NewErrResponse(c, err)
		}

		//  Filtering item if user already declined this offer
		loggedInUserDeclined := false
		for _, decision := range decisions.Items {
			if decision.GetUserKey() == loggedInUserKey && decision.Decision == trading.DeclinedDecision {
				loggedInUserDeclined = true
			}
		}
		if loggedInUserDeclined {
			continue
		}

		//  Getting offer items
		items, err := h.tradingStore.GetItems(item.Offer.GetKey())
		if err != nil {
			return NewErrResponse(c, err)
		}

		//  Mapping offer to web response
		webOffer, err := h.mapToWebOffer(&item.Offer, items, decisions)
		if err != nil {
			return NewErrResponse(c, err)
		}

		webOffers = append(webOffers, *webOffer)

	}

	//  Mapping to web response
	return c.JSON(http.StatusOK, web.GetOffersResponse{
		Offers: webOffers,
	})

}

func (h *Handler) GetOffer(c echo.Context) error {

	var err error

	offerIdStr := c.Param("id")
	offerId, err := uuid.FromString(offerIdStr)
	if err != nil {
		return NewErrResponse(c, err)
	}
	offerKey := NewOfferKey(offerId)

	offer, err := h.GetWebOffer(offerKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, offer)

}

func (h *Handler) GetWebOffer(offerKey OfferKey) (*web.GetOfferResponse, error) {

	offer, err := h.tradingStore.GetOffer(offerKey)
	if err != nil {
		return nil, err
	}

	items, err := h.tradingStore.GetItems(offerKey)
	if err != nil {
		return nil, err
	}

	decisions, err := h.tradingStore.GetDecisions(offerKey)
	if err != nil {
		return nil, err
	}

	webOffer, err := h.mapToWebOffer(offer, items, decisions)
	if err != nil {
		return nil, err
	}

	response := web.GetOfferResponse{
		Offer: *webOffer,
	}

	return &response, nil
}

func (h *Handler) ConfirmItemReceivedOrGiven(c echo.Context) error {

	ctx, _ := GetEchoContext(c, "ConfirmItemReceivedOrGiven")

	offerItemKey, err := ParseOfferItemKey(c.Param("id"))
	if err != nil {
		return err
	}
	err = h.tradingService.ConfirmItemReceivedOrGiven(ctx, offerItemKey)
	if err != nil {
		return err
	}

	offerItem, err := h.tradingService.GetOfferItem(ctx, offerItemKey)
	if err != nil {
		return errs.ReturnException(c, err)
	}

	webResponse := h.mapOfferItem(offerItem)

	return c.JSON(http.StatusOK, webResponse)

}

func (h *Handler) mapOfferItem(offerItem *trading.OfferItem) web.OfferItem {
	resourceId := offerItem.ResourceID
	var resourceIdResult *string = nil
	if resourceId != nil {
		resourceIdStr := resourceId.String()
		resourceIdResult = &resourceIdStr
	}
	webResponse := web.OfferItem{
		ID:            offerItem.ID.String(),
		FromUserID:    offerItem.FromUserID,
		ToUserID:      offerItem.ToUserID,
		Type:          offerItem.ItemType,
		ResourceId:    resourceIdResult,
		TimeInSeconds: offerItem.OfferedTimeInSeconds,
	}
	return webResponse
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

func (h *Handler) mapToWebOffer(offer *trading.Offer, items *trading.OfferItems, decisions *trading.OfferDecisions) (*web.Offer, error) {

	authorUsername, err := h.authStore.GetUsername(offer.GetAuthorKey())
	if err != nil {
		return nil, err
	}

	webOffer := web.Offer{
		ID:             offer.ID.String(),
		CreatedAt:      offer.CreatedAt,
		CompletedAt:    offer.CompletedAt,
		Status:         offer.Status,
		Items:          nil,
		AuthorID:       offer.AuthorID,
		AuthorUsername: authorUsername,
	}

	var responseItems = make([]web.OfferItem, items.ItemCount())

	for i, offerItem := range items.Items {
		webItem := web.OfferItem{
			ID:         offerItem.ID.String(),
			FromUserID: offerItem.FromUserID,
			ToUserID:   offerItem.ToUserID,
			Type:       offerItem.ItemType,
		}
		if offerItem.ItemType == trading.ResourceItem {
			resourceUid := offerItem.ResourceID
			var resourceIdResult *string = nil
			if resourceUid != nil {
				resourceIdStr := resourceUid.String()
				resourceIdResult = &resourceIdStr
			}
			webItem.ResourceId = resourceIdResult
		} else if offerItem.ItemType == trading.TimeItem {
			webItem.TimeInSeconds = offerItem.OfferedTimeInSeconds
		}
		responseItems[i] = webItem
	}
	webOffer.Items = responseItems

	var responseDecisions = make([]web.OfferDecision, len(decisions.Items))
	for i, decision := range decisions.Items {
		webDecision := web.OfferDecision{
			OfferID:  decision.OfferID.String(),
			UserID:   decision.UserID,
			Decision: decision.Decision,
		}
		responseDecisions[i] = webDecision
	}
	webOffer.Decisions = responseDecisions

	return &webOffer, nil
}
