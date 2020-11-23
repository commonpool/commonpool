package handler

import (
	"fmt"
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
	OfferItem OfferItem
	Resource  Resource
	Time      int64
	FromUser  User
	ToUser    User
}

func (h *Handler) SendOffer(c echo.Context) error {
	c.Logger().Info("SendOffer")

	var err error

	//
	// > Getting Logged In User
	//
	c.Logger().Debug("SendOffer: getting logged in user key")
	loggedInUser := h.authorization.GetAuthUserSession(c)
	loggedInUserKey := loggedInUser.GetUserKey()

	//
	// > Unmarshaling Request Body
	//
	req := web.SendOfferRequest{}
	if err = c.Bind(&req); err != nil {
		c.Logger().Error(err, "SendOffer: could not unmarshal SendOfferRequest")
		response := errs.ErrCreateResourceBadRequest(err)
		return NewErrResponse(c, &response)
	}

	//
	// > Validating Request Body
	//
	//
	c.Logger().Debug("SendOffer: validating request payload")
	if err = c.Validate(req); err != nil {
		c.Logger().Error(err, "SendOffer: validating request payload... error!")
		return c.JSON(400, NewValidError(err.(validator.ValidationErrors)))
	}

	//
	// > Creating new Key for Offer
	//
	offerKey := NewOfferKey(uuid.NewV4())
	c.Logger().Debugf("SendOffer: new offer key is %s", offerKey.ID.String())

	// keep track of resources so that one resource is not traded twice in the same transaction
	var usedResources = NewResources([]Resource{})
	var participants = NewUsers([]User{})

	//
	// > Process each offer item
	//
	itemCount := len(req.Offer.Items)
	items := NewOfferItems([]OfferItem{})

	c.Logger().Debugf("SendOffer: processing %d offer items...", itemCount)
	for i, reqItem := range req.Offer.Items {
		c.Logger().Debugf("SendOffer: %d/%d processing offer item...", i+1, itemCount)

		var resourceKey *ResourceKey

		if reqItem.Type == ResourceItem {
			resourceKey, err = ParseResourceKey(*reqItem.ResourceId)
			if err != nil {
				return NewErrResponse(c, err)
			}
			//
			// > A resource can only appears once in an offer
			//
			c.Logger().Debugf("SendOffer: %d/%d checking if resource %s appear only once in offer...", i+1, itemCount, *reqItem.ResourceId)
			if usedResources.ContainsKey(*resourceKey) {
				err := fmt.Errorf("resources can only appear once in a transaction")
				c.Logger().Errorf("SendOffer: %d/%d %v", i+1, itemCount, err)
				return NewErrResponse(c, err)
			}
		}

		//
		// > Getting participants
		//
		fromUser := NewUserKey(reqItem.From)
		toUser := NewUserKey(reqItem.To)
		c.Logger().Debugf("SendOffer: %d/%d processing offer item from user %s to user %s", i+1, itemCount, fromUser.String(), toUser.String())

		//
		// > Making sure participants exist
		//
		c.Logger().Debugf("SendOffer: %d/%d making sure participants exist...", i+1, itemCount)
		fromUserFound := items.HasItemsForUser(fromUser)
		toUserFound := items.HasItemsForUser(toUser)

		//
		// > Make sure offer item.from != item.to
		//
		c.Logger().Debugf("SendOffer: %d/%d make sure offer item receiver != offer item sender", i+1, itemCount)
		if fromUser == toUser {
			err := NewErrResponse(c, fmt.Errorf("one cannot trade with oneself"))
			c.Logger().Errorf("SendOffer:%d/%d %v", i+1, itemCount, err)
			return err
		}

		if !fromUserFound || !toUserFound {
			c.Logger().Debugf("SendOffer: %d/%d some participants not found yet. Making sure they exist...", i+1, itemCount)

			var userKeys []UserKey
			if !fromUserFound {
				userKeys = append(userKeys, fromUser)
			}
			if !toUserFound {
				userKeys = append(userKeys, toUser)
			}

			c.Logger().Debugf("SendOffer: %d/%d expecting %d participant from query...", i+1, itemCount, len(userKeys))

			foundUsers, err := h.authStore.GetByKeys(userKeys)
			if err != nil {
				c.Logger().Errorf("SendOffer: %d/%d %v", i+1, itemCount, err)
				return c.JSON(http.StatusInternalServerError, err.Error())
			}

			userQueryCount := len(foundUsers.Items)
			c.Logger().Debugf("SendOffer: %d/%d got %d participant from query...", i+1, itemCount, userQueryCount)
			if userQueryCount != len(userKeys) {
				err := fmt.Errorf("participant not found")
				c.Logger().Errorf("SendOffer: %d/%d %v", i+1, itemCount, err)
				return c.JSON(http.StatusNotFound, err.Error())
			}
			participants.AppendAll(foundUsers)
		}

		var item OfferItem
		itemKey := NewOfferItemKey(uuid.NewV4(), offerKey)

		//
		// > Building the offer item
		//
		if reqItem.Type == TimeItem {

			//
			// > Resource Item is a TimeItem
			//
			c.Logger().Debugf("SendOffer: %d/%d offer item is a TimeItem. Validating...", i+1, itemCount)

			//
			// > Making sure the offer has Positive time value
			//
			seconds := reqItem.TimeInSeconds
			if *seconds <= 0 {
				err := fmt.Errorf("time traded must be more than 0")
				c.Logger().Errorf("SendOffer: %d/%d %v", i+1, itemCount, err)
				return NewErrResponse(c, err)
			}

			//
			// > Building the offer item
			//
			item = NewTimeOfferItem(itemKey, fromUser, toUser, *seconds)

		} else if reqItem.Type == ResourceItem {

			//
			// > Resource Item is a ResourceItem
			//
			c.Logger().Debugf("SendOffer: %d/%d offer item is a ResourceItem. Valiating...", i+1, itemCount)

			//
			// > Retrieving the resource key
			//
			c.Logger().Debugf("SendOffer: %d/%d getting ResourceItem resource key...", i+1, itemCount)
			resourceId, err := uuid.FromString(*reqItem.ResourceId)
			if err != nil {
				c.Logger().Debugf("SendOffer: %d/%d %v", i+1, itemCount, err)
				return NewErrResponse(c, err)
			}
			resourceKey := NewResourceKey(resourceId)

			//
			// > Retrieving the resource from the store
			//
			c.Logger().Debugf("SendOffer: %d/%d getting the ResourceItem Resource object...", i+1, itemCount)
			getResourceByKeyResponse := h.resourceStore.GetByKey(resource2.NewGetResourceByKeyQuery(resourceKey))
			if getResourceByKeyResponse.Error != nil {
				c.Logger().Errorf("SendOffer: %d/%d %v", i+1, itemCount, getResourceByKeyResponse.Error)
				return getResourceByKeyResponse.Error
			}
			resource := getResourceByKeyResponse.Resource
			usedResources.Append(*resource)

			//
			// > Making sure the Item.From = Resource.Owner
			// > Someone cannot make an offer where someone would trade a resource he's not the owner of
			//
			c.Logger().Debugf("SendOffer: %d/%d making sure the ResourceItem Resource Owner === item.from...", i+1, itemCount)
			if resource.GetUserKey() != fromUser {
				res := errs.ErrTransactionResourceOwnerMismatch()
				c.Logger().Errorf("SendOffer: %d/%d %v", i+1, itemCount, res)
				return NewErrResponse(c, &res)
			}

			//
			// > Building the Resource Offer item
			//
			item = NewResourceOfferItem(itemKey, fromUser, toUser, resourceKey)

		} else {

			//
			// > Resource Type is of unknown type
			//
			err := fmt.Errorf("unexpected resource type %d", reqItem.Type)
			c.Logger().Errorf("SendOffer: %d/%d %v", i+1, itemCount, err)
			return c.JSON(http.StatusBadRequest, err.Error())

		}

		//
		// > Appending the built item to the list of items
		//
		items = items.Append(item)
		c.Logger().Debugf("SendOffer: %d/%d processing offer item... done!", i+1, itemCount)

	}

	//
	// > Saving Offer
	//
	c.Logger().Debug("SendOffer: saving offer...")
	err = h.tradingStore.SaveOffer(NewOffer(offerKey, loggedInUserKey, nil), items)
	if err != nil {
		c.Logger().Errorf("SendOffer: %v", err)
		return NewErrResponse(c, err)
	}

	//
	// Sending message to concerned users
	//
	c.Logger().Debug("SendOffer: sending offer to chat...")
	createTopicRequest := chat.NewGetOrCreateConversationTopicRequest(participants.GetUserKeys())
	topicResponse := h.chatStore.GetOrCreateConversationTopic(&createTopicRequest)
	if topicResponse.Error != nil {
		return NewErrResponse(c, topicResponse.Error)
	}

	c.Logger().Debugf("SendOffer: sending offer to %d people", participants.GetUserCount())
	for _, userKey := range items.GetUserKeys().Items {

		offerItemsForUser := items.GetOfferItemsForUser(userKey)

		// Will contain the formatted message
		var blocks []Block

		// Adding a descriptive header
		blocks = append(blocks, *NewHeaderBlock(NewMarkdownObject(fmt.Sprintf("%s is proposing an exchange", loggedInUser.Username)), nil))

		// Looping through each item to construct the message
		c.Logger().Debugf("SendOffer: processing %d items for user %s", offerItemsForUser.ItemCount(), userKey.String())
		for _, offerItemForUser := range offerItemsForUser.Items {

			toUser, err := participants.GetUser(offerItemForUser.GetToUserKey())
			if err != nil {
				return NewErrResponse(c, err)
			}

			fromUser, err := participants.GetUser(offerItemForUser.GetFromUserKey())
			if err != nil {
				return NewErrResponse(c, err)
			}

			if offerItemForUser.GetFromUserKey() == userKey {
				c.Logger().Debug("SendOffer: user is giving that item")

				if offerItemForUser.IsTimeExchangeItem() {
					c.Logger().Debug("SendOffer: item is time exchange")

					message := fmt.Sprintf("**%s** would like **%s** of your time", toUser.Username, offerItemForUser.FormatOfferedTimeInSeconds())
					block := NewSectionBlock(NewMarkdownObject(message), nil, nil, nil)
					blocks = append(blocks, *block)

				} else if offerItemForUser.IsResourceExchangeItem() {
					c.Logger().Debug("SendOffer: item is resource exchange")

					resource, err := usedResources.GetResource(offerItemForUser.GetResourceKey())
					if err != nil {
						return NewErrResponse(c, err)
					}

					message := fmt.Sprintf("**%s** would like **%s**", fromUser.Username, resource.Summary)
					block := NewSectionBlock(NewMarkdownObject(message), nil, nil, nil)
					blocks = append(blocks, *block)
				}

			} else {
				c.Logger().Debug("SendOffer: user is receiving that item")

				if offerItemForUser.IsTimeExchangeItem() {
					c.Logger().Debug("SendOffer: item is time exchange")

					message := fmt.Sprintf("you would get **%s** of time bank credits from **%s**", offerItemForUser.FormatOfferedTimeInSeconds(), fromUser.Username)
					block := NewSectionBlock(NewMarkdownObject(message), nil, nil, nil)
					blocks = append(blocks, *block)

				} else if offerItemForUser.IsResourceExchangeItem() {
					c.Logger().Debug("SendOffer: item is resource exchange")

					resource, err := usedResources.GetResource(offerItemForUser.GetResourceKey())
					if err != nil {
						return NewErrResponse(c, err)
					}

					message := fmt.Sprintf("you would get **%s** from **%s**", resource.Summary, toUser.Username)
					block := NewSectionBlock(NewMarkdownObject(message), nil, nil, nil)
					blocks = append(blocks, *block)
				}
			}
		}
		primaryButtonStyle := Primary
		dangerButtonStyle := Danger
		acceptButton := NewButtonElement(NewPlainTextObject("Accept"), &primaryButtonStyle, nil, nil, nil, nil)
		declineButton := NewButtonElement(NewPlainTextObject("Decline"), &dangerButtonStyle, nil, nil, nil, nil)
		actionBlock := NewActionBlock([]BlockElement{
			acceptButton,
			declineButton,
		}, nil)
		blocks = append(blocks, *actionBlock)
		threadKey := NewThreadKey(topicResponse.TopicKey, userKey)
		sendMsgRequest := chat.NewSendMessageToThreadRequest(
			threadKey,
			loggedInUserKey,
			loggedInUser.Username,
			"New offer",
			blocks,
			[]Attachment{},
		)
		sendMsgResponse := h.chatStore.SendMessageToThread(&sendMsgRequest)
		if sendMsgResponse.Error != nil {
			return NewErrResponse(c, err)
		}
	}

	//
	// > Getting Offer
	//
	c.Logger().Debug("SendOffer: getting offer...")
	o, err := h.tradingStore.GetOffer(offerKey)
	if err != nil {
		c.Logger().Errorf("SendOffer: %v", err)
		return NewErrResponse(c, err)
	}

	//
	// > Getting Offer Items
	//
	c.Logger().Debug("SendOffer: getting offer items...")
	i, err := h.tradingStore.GetItems(offerKey)
	if err != nil {
		c.Logger().Errorf("SendOffer: %v", err)
		return NewErrResponse(c, err)
	}

	//
	// > Getting Offer Decisions
	//
	c.Logger().Debug("SendOffer: getting offer decisions...")
	d, err := h.tradingStore.GetDecisions(offerKey)
	if err != nil {
		c.Logger().Errorf("SendOffer: %v", err)
		return NewErrResponse(c, err)
	}

	//
	// > Mapping Offer to Web Response
	//
	c.Logger().Debug("SendOffer: mapping offer to web response...")
	webOffer, err := h.mapToWebOffer(o, i, d)
	if err != nil {
		c.Logger().Errorf("SendOffer: %v", err)
		return NewErrResponse(c, err)
	}

	return c.JSON(http.StatusCreated, web.GetOfferResponse{Offer: *webOffer})

}

func (h *Handler) AcceptOffer(c echo.Context) error {
	c.Logger().Info("AcceptOffer")

	//
	// > Retrieving the logged in user
	//
	c.Logger().Debug("AcceptOffer: getting authenticated user...")
	authSession := h.authorization.GetAuthUserSession(c)
	authUserKey := authSession.GetUserKey()
	c.Logger().Debug("AcceptOffer: getting authenticated user... done!")

	//
	// > Parsing Offer key from query params
	//
	c.Logger().Debug("AcceptOffer: parsing Offer key...")
	offerKey, err := ParseOfferKey(c.Param("id"))
	if err != nil {
		c.Logger().Error(err, "AcceptOffer: parsing Offer key... error!")
		return NewErrResponse(c, err)
	}
	c.Logger().Debug("AcceptOffer: parsing Offer key... done!")

	//
	// > Retrieving offer
	//
	c.Logger().Debug("AcceptOffer: retrieving offer...")
	offer, err := h.tradingStore.GetOffer(offerKey)
	if err != nil {
		c.Logger().Error(err, "AcceptOffer: retrieving offer... error!")
		return NewErrResponse(c, err)
	}
	c.Logger().Debug("AcceptOffer: retrieving offer... done!")

	//
	// > Ensure offer is still pending approval
	//
	c.Logger().Debug("AcceptOffer: ensure offer is still pending approval...")
	if offer.Status != PendingOffer {
		err := fmt.Errorf("offer is not pending approval")
		c.Logger().Warn(err, "AcceptOffer: ensure offer is still pending approval... error!")
		return NewErrResponse(c, err)
	}
	c.Logger().Debug("AcceptOffer: ensure offer is still pending approval... done!")

	//
	// > Retrieve Offer decisions
	//
	c.Logger().Debug("AcceptOffer: retrieving offer decisions...")
	decisions, err := h.tradingStore.GetDecisions(offerKey)
	if err != nil {
		c.Logger().Error(err, "AcceptOffer: retrieving offer decisions... error!")
		return NewErrResponse(c, err)
	}
	c.Logger().Debug("AcceptOffer: retrieving offer decisions... done!")

	var didAllOtherParticipantsAlreadyAccept = true
	var currentUserDecision *OfferDecision

	//
	// > Retrieving current user decision, and check if everyone else accepted the offer already
	//
	c.Logger().Debug("AcceptOffer: checking if everyone else already accepted the offer...")
	for _, decision := range decisions {
		decisionKey := decision.GetKey()
		decisionUserKey := decisionKey.GetUserKey()
		if decisionUserKey != authUserKey {
			c.Logger().Debugf("AcceptOffer: decision for user %s is %d", decisionUserKey.String(), decision.Decision)
			if decision.Decision != AcceptedDecision {
				c.Logger().Debug("AcceptOffer: decision is not accepted")
				didAllOtherParticipantsAlreadyAccept = false
			}
		} else {
			currentUserDecision = &decision
		}
	}
	c.Logger().Debug("AcceptOffer: checking if everyone else already accepted the offer... done!")

	// in the case that a user is approving an offer he's not part of
	if currentUserDecision == nil {
		err := fmt.Errorf("could not find current user decision")
		c.Logger().Errorf("AcceptOffer: cannot accept this offer: %v", err)
		return err
	}

	//
	// > Persisting the decision
	//
	c.Logger().Debug("AcceptOffer: saving decision...")
	err = h.tradingStore.SaveDecision(offerKey, authUserKey, AcceptedDecision)
	if err != nil {
		c.Logger().Errorf("AcceptOffer: saving decision... error! %v", err)
		return NewErrResponse(c, err)
	}
	c.Logger().Debug("AcceptOffer: saving decision... done!")

	//
	// > Complete offer if everyone accepted already
	//
	var currentUserLastOneToDecide = didAllOtherParticipantsAlreadyAccept
	if currentUserLastOneToDecide {
		c.Logger().Debug("AcceptOffer: everyone already accepted. completing offer...")
		err = h.tradingStore.CompleteOffer(offerKey, AcceptedOffer)
		if err != nil {
			c.Logger().Errorf("AcceptOffer: everyone already accepted. completing offer... error! : %v", err)
			return NewErrResponse(c, err)
		}
		c.Logger().Debug("AcceptOffer: everyone already accepted. completing offer... done!")
	}

	//
	// > Convert to web response
	//
	c.Logger().Debug("AcceptOffer: converting to response body...")
	err = h.getWebOffer(c, http.StatusAccepted, offerKey)
	if err != nil {
		c.Logger().Errorf("AcceptOffer: converting to response body... error! %v", err)
		return err
	}
	c.Logger().Debug("AcceptOffer: converting to response body... done!")

	return nil

}

func (h *Handler) DeclineOffer(c echo.Context) error {
	c.Logger().Info("DeclineOffer")

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
	if offer.Status != PendingOffer {
		return NewErrResponse(c, fmt.Errorf("offer is not pending approval"))
	}

	// get offer decisions
	decisions, err := h.tradingStore.GetDecisions(offerKey)
	if err != nil {
		return NewErrResponse(c, err)
	}

	var currentUserDecision *OfferDecision

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
	err = h.tradingStore.SaveDecision(offerKey, authUserKey, DeclinedDecision)
	if err != nil {
		return NewErrResponse(c, err)
	}

	// complete offer if everyone approved
	err = h.tradingStore.CompleteOffer(offerKey, DeclinedOffer)
	if err != nil {
		return NewErrResponse(c, err)
	}

	return h.getWebOffer(c, http.StatusAccepted, offerKey)

}

func (h *Handler) GetOffers(c echo.Context) error {
	c.Logger().Info("GetOffers")

	//
	// > Retrieving currently logged in user
	//
	c.Logger().Debug("GetOffers: getting logged in user...")
	authSession := h.authorization.GetAuthUserSession(c)
	loggedInUserKey := authSession.GetUserKey()
	c.Logger().Debug("GetOffers: getting logged in user... done!")

	//
	// > Retrieving logged in user offers
	//
	c.Logger().Debug("GetOffers: retrieving logged in user offers...")
	qry := trading.GetOffersQuery{
		UserKeys: []UserKey{loggedInUserKey},
	}
	result, err := h.tradingStore.GetOffers(qry)
	if err != nil {
		c.Logger().Debugf("GetOffers: retrieving logged in user offers... error! %v", err)
		return NewErrResponse(c, err)
	}
	c.Logger().Debug("GetOffers: retrieving logged in user offers... done!")

	//
	// > Mapping offers to web response
	//
	var webOffers []web.Offer
	resultItems := result.Items
	resultCount := len(resultItems)
	c.Logger().Debugf("GetOffers: processing %d offers...", resultCount)
	for i, item := range resultItems {

		c.Logger().Debugf("GetOffers: %d/%d processing offer...", i+1, resultCount)

		//
		// > Retrieving decisions for offer
		//
		c.Logger().Debugf("GetOffers: %d/%d getting offer decisions...", i+1, resultCount)
		decisions, err := h.tradingStore.GetDecisions(item.Offer.GetKey())
		if err != nil {
			c.Logger().Errorf("GetOffers: %d/%d getting offer decisions... error! %v", i+1, resultCount, err)
			return NewErrResponse(c, err)
		}
		c.Logger().Debugf("GetOffers: %d/%d getting offer decisions... done!", i+1, resultCount)

		//
		// > Filtering item if user already declined this offer
		//
		c.Logger().Debugf("GetOffers: %d/%d filtering if user already declined the offer...", i+1, resultCount)
		loggedInUserDeclined := false
		for _, decision := range decisions {
			if decision.GetUserKey() == loggedInUserKey && decision.Decision == DeclinedDecision {
				loggedInUserDeclined = true
			}
		}
		if loggedInUserDeclined {
			c.Logger().Debugf("GetOffers: %d/%d user already declined the offer ... skipping", i+1, resultCount)
			continue
		}
		c.Logger().Debugf("GetOffers: %d/%d user did not decline the offer. continuing...", i+1, resultCount)

		//
		// > Getting offer items
		//
		c.Logger().Debugf("GetOffers: %d/%d getting offer items...", i+1, resultCount)
		items, err := h.tradingStore.GetItems(item.Offer.GetKey())
		if err != nil {
			c.Logger().Debugf("GetOffers: %d/%d getting offer items... error! %v", i+1, resultCount, err)
			return NewErrResponse(c, err)
		}
		c.Logger().Debugf("GetOffers: %d/%d getting offer items... done!", i+1, resultCount)

		//
		// > Mapping offer to web response
		//
		c.Logger().Debugf("GetOffers: %d/%d mapping offer to web response...", i+1, resultCount)
		webOffer, err := h.mapToWebOffer(item.Offer, items, decisions)
		if err != nil {
			c.Logger().Errorf("GetOffers: %d/%d mapping offer to web response... error! %v", i+1, resultCount, err)
			return NewErrResponse(c, err)
		}
		c.Logger().Debugf("GetOffers: %d/%d mapping offer to web response... done!", i+1, resultCount)

		webOffers = append(webOffers, *webOffer)

		c.Logger().Debugf("GetOffers: %d/%d processing offer... done!", i+1, resultCount)
	}

	//
	// > Mapping to web response
	//
	return c.JSON(http.StatusOK, web.GetOffersResponse{
		Offers: webOffers,
	})

}

func (h *Handler) GetOffer(c echo.Context) error {
	c.Logger().Info("GetOffer")

	var err error

	offerIdStr := c.Param("id")
	offerId, err := uuid.FromString(offerIdStr)
	if err != nil {
		return NewErrResponse(c, err)
	}
	offerKey := NewOfferKey(offerId)

	return h.getWebOffer(c, http.StatusOK, offerKey)

}

func (h *Handler) getWebOffer(c echo.Context, statusCode int, offerKey OfferKey) error {
	offer, err := h.tradingStore.GetOffer(offerKey)
	if err != nil {
		return NewErrResponse(c, err)
	}

	items, err := h.tradingStore.GetItems(offerKey)
	if err != nil {
		return NewErrResponse(c, err)
	}

	decisions, err := h.tradingStore.GetDecisions(offerKey)
	if err != nil {
		return NewErrResponse(c, err)
	}

	webOffer, err := h.mapToWebOffer(offer, items, decisions)
	if err != nil {
		return NewErrResponse(c, err)
	}

	response := web.GetOfferResponse{
		Offer: *webOffer,
	}
	return c.JSON(statusCode, response)
}

func (h *Handler) mapToWebOffer(offer Offer, items *OfferItems, decisions []OfferDecision) (*web.Offer, error) {

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
		if offerItem.ItemType == ResourceItem {
			webItem.ResourceId = offerItem.ResourceID.String()
		} else if offerItem.ItemType == TimeItem {
			webItem.TimeInSeconds = *offerItem.OfferedTimeInSeconds
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
