package handler

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/chat"
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

	validErrors := []ValidError{}

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

	ctx := c.Request().Context()

	var err error

	//  Getting Logged In User
	loggedInUser := h.authorization.GetAuthUserSession(c)

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

	reference := UserReference(&loggedInUser)
	response, err := h.SendOffer(ctx, req, reference)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, response)

}

func (h *Handler) SendOffer(ctx context.Context, req web.SendOfferRequest, fromUser UserReference) (*web.GetOfferResponse, error) {

	//  Creating new Key for Offer
	offerKey := NewOfferKey(uuid.NewV4())

	// keep track of resources so that one resource is not traded twice in the same transaction
	var usedResources = resource2.NewResources([]resource2.Resource{})
	var participants = auth.NewUsers([]auth.User{})

	//  Process each offer item
	items := trading.NewOfferItems([]trading.OfferItem{})

	for _, reqItem := range req.Offer.Items {

		if reqItem.Type == trading.ResourceItem {
			resourceKey, err := ParseResourceKey(*reqItem.ResourceId)
			if err != nil {
				return nil, err
			}
			//  A resource can only appears once in an offer
			if usedResources.ContainsKey(*resourceKey) {
				err := fmt.Errorf("resources can only appear once in a transaction")
				return nil, err
			}
		}

		//  Getting participants
		fromUser := NewUserKey(reqItem.From)
		toUser := NewUserKey(reqItem.To)

		//  Making sure participants exist
		fromUserFound := items.HasItemsForUser(fromUser)
		toUserFound := items.HasItemsForUser(toUser)

		//  Make sure offer item.from != item.to
		if fromUser == toUser {
			return nil, fmt.Errorf("one cannot trade with oneself")
		}

		if !fromUserFound || !toUserFound {

			var userKeys []UserKey
			if !fromUserFound {
				userKeys = append(userKeys, fromUser)
			}
			if !toUserFound {
				userKeys = append(userKeys, toUser)
			}

			foundUsers, err := h.authStore.GetByKeys(ctx, userKeys)
			if err != nil {
				return nil, err
			}

			userQueryCount := len(foundUsers.Items)
			if userQueryCount != len(userKeys) {
				err := fmt.Errorf("participant not found")
				return nil, err
			}

			participants = participants.AppendAll(foundUsers)
		}

		var item trading.OfferItem
		itemKey := NewOfferItemKey(uuid.NewV4())

		//  Building the offer item
		if reqItem.Type == trading.TimeItem {

			//  Resource Item is a TimeItem

			//  Making sure the offer has Positive time value
			seconds := reqItem.TimeInSeconds
			if *seconds <= 0 {
				err := fmt.Errorf("time traded must be more than 0")
				return nil, err
			}

			//  Building the offer item
			item = trading.NewTimeOfferItem(offerKey, itemKey, fromUser, toUser, *seconds)

		} else if reqItem.Type == trading.ResourceItem {

			//  Resource Item is a ResourceItem
			//  Retrieving the resource key
			resourceId, err := uuid.FromString(*reqItem.ResourceId)
			if err != nil {
				return nil, err
			}
			resourceKey := NewResourceKey(resourceId)

			//  Retrieving the resource from the store
			getResourceByKeyResponse := h.resourceStore.GetByKey(ctx, resource2.NewGetResourceByKeyQuery(resourceKey))
			if getResourceByKeyResponse.Error != nil {
				return nil, err
			}
			resource := getResourceByKeyResponse.Resource
			usedResources = usedResources.Append(*resource)

			//  Making sure the Item.From = Resource.Owner > Someone cannot make an offer where someone would trade a resource he's not the owner of
			//
			if resource.GetOwnerKey() != fromUser {
				res := errs.ErrTransactionResourceOwnerMismatch()
				return nil, &res
			}

			//  Building the Resource Offer item
			item = trading.NewResourceOfferItem(offerKey, itemKey, fromUser, toUser, resourceKey)

		} else {

			//  Resource Type is of unknown type
			err := fmt.Errorf("unexpected resource type %d", reqItem.Type)
			return nil, err

		}

		//  Appending the built item to the list of items
		items = items.Append(item)

	}

	//  Saving Offer
	err := h.tradingStore.SaveOffer(trading.NewOffer(offerKey, fromUser.GetUserKey(), req.Offer.Message, nil), items)
	if err != nil {
		return nil, err
	}

	//
	// Sending message to concerned users
	//

	for i, item := range participants.GetUserKeys().Items {
		fmt.Println(fmt.Sprintf("Participant %d : %s", i, item.String()))
	}

	for _, itemUserKey := range items.GetUserKeys().Items {

		userOfferItems := items.GetOfferItemsForUser(itemUserKey)

		// Will contain the formatted message
		var blocks []chat.Block

		// Adding a descriptive header
		blocks = append(blocks, *chat.NewHeaderBlock(
			chat.NewMarkdownObject(
				fmt.Sprintf("%s is proposing an exchange", h.chatService.GetUserLink(fromUser.GetUserKey())),
			), nil))

		// Looping through each item to construct the message
		for _, userOfferItem := range userOfferItems.Items {

			toUser, err := participants.GetUser(userOfferItem.GetToUserKey())
			if err != nil {
				return nil, err
			}

			fromUser, err := participants.GetUser(userOfferItem.GetFromUserKey())
			if err != nil {
				return nil, err
			}

			if userOfferItem.GetFromUserKey() == itemUserKey {

				if userOfferItem.IsTimeExchangeItem() {
					message := fmt.Sprintf("%s would like %s of your time",
						h.chatService.GetUserLink(toUser.GetUserKey()),
						userOfferItem.FormatOfferedTimeInSeconds())
					block := chat.NewSectionBlock(chat.NewMarkdownObject(message), nil, nil, nil)
					blocks = append(blocks, *block)
				} else if userOfferItem.IsResourceExchangeItem() {
					message := fmt.Sprintf("%s would like %s",
						h.chatService.GetUserLink(toUser.GetUserKey()),
						h.chatService.GetResourceLink(userOfferItem.GetResourceKey()))
					block := chat.NewSectionBlock(chat.NewMarkdownObject(message), nil, nil, nil)
					blocks = append(blocks, *block)
				}

			} else {

				if userOfferItem.IsTimeExchangeItem() {
					message := fmt.Sprintf("you would get %s of time bank credits from %s",
						userOfferItem.FormatOfferedTimeInSeconds(),
						h.chatService.GetUserLink(fromUser.GetUserKey()))
					block := chat.NewSectionBlock(chat.NewMarkdownObject(message), nil, nil, nil)
					blocks = append(blocks, *block)
				} else if userOfferItem.IsResourceExchangeItem() {
					resourceKey := userOfferItem.GetResourceKey()
					message := fmt.Sprintf("you would get %s from %s",
						h.chatService.GetResourceLink(resourceKey),
						h.chatService.GetUserLink(fromUser.GetUserKey()))
					block := chat.NewSectionBlock(chat.NewMarkdownObject(message), nil, nil, nil)
					blocks = append(blocks, *block)
				}
			}
		}

		primaryButtonStyle := chat.Primary
		dangerButtonStyle := chat.Danger
		acceptOfferActionId := "accept_offer"
		declineOfferActionId := "decline_offer"
		offerId := offerKey.ID.String()
		acceptButton := chat.NewButtonElement(chat.NewPlainTextObject("Accept"), &primaryButtonStyle, &acceptOfferActionId, nil, &offerId, nil)
		declineButton := chat.NewButtonElement(chat.NewPlainTextObject("Decline"), &dangerButtonStyle, &declineOfferActionId, nil, &offerId, nil)
		actionBlock := chat.NewActionBlock([]chat.BlockElement{
			*acceptButton,
			*declineButton,
		}, nil)
		blocks = append(blocks, *actionBlock)

		linkBlock := chat.NewSectionBlock(
			chat.NewMarkdownObject(
				fmt.Sprintf("[View offer details](/users/%s/transactions/%s)", itemUserKey.String(), offerKey.ID)),
			nil,
			nil,
			nil)
		blocks = append(blocks, *linkBlock)

		sendMsgRequest := chat.NewSendConversationMessage(
			fromUser.GetUserKey(),
			fromUser.GetUsername(),
			participants.GetUserKeys(),
			"New offer",
			blocks,
			[]chat.Attachment{},
			&itemUserKey,
		)
		_, err := h.chatService.SendConversationMessage(ctx, sendMsgRequest)
		if err != nil {
			return nil, err
		}
	}

	if req.Offer.Message != "" {
		sendMsgRequest := chat.NewSendConversationMessage(
			fromUser.GetUserKey(),
			fromUser.GetUsername(),
			participants.GetUserKeys(),
			req.Offer.Message,
			[]chat.Block{},
			[]chat.Attachment{},
			nil,
		)
		_, err := h.chatService.SendConversationMessage(ctx, sendMsgRequest)
		if err != nil {
			return nil, err
		}
	}

	//  Getting Offer
	o, err := h.tradingStore.GetOffer(offerKey)
	if err != nil {
		return nil, err
	}

	//  Getting Offer Items
	i, err := h.tradingStore.GetItems(offerKey)
	if err != nil {
		return nil, err
	}

	//  Getting Offer Decisions
	d, err := h.tradingStore.GetDecisions(offerKey)
	if err != nil {
		return nil, err
	}

	//  Mapping Offer to Web Response
	webOffer, err := h.mapToWebOffer(o, i, d)
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

	// retrieve current user
	authSession := h.authorization.GetAuthUserSession(c)
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
	for _, decision := range decisions {
		decisionKey := decision.GetKey()
		decisionUserKey := decisionKey.GetUserKey()
		if decisionUserKey == authUserKey {
			currentUserDecision = &decision
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
		for _, decision := range decisions {
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
		webOffer, err := h.mapToWebOffer(item.Offer, items, decisions)
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

func (h *Handler) mapToWebOffer(offer trading.Offer, items *trading.OfferItems, decisions []trading.OfferDecision) (*web.Offer, error) {

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

	var responseDecisions = make([]web.OfferDecision, len(decisions))
	for i, decision := range decisions {
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
