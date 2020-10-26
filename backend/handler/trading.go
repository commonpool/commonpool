package handler

import (
	"fmt"
	errs "github.com/commonpool/backend/errors"
	. "github.com/commonpool/backend/model"
	"github.com/commonpool/backend/trading"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func (h *Handler) SendOffer(c echo.Context) error {

	var err error

	authUserKey := h.authorization.GetAuthUserKey(c)

	req := web.SendOfferRequest{}
	if err = c.Bind(&req); err != nil {
		response := errs.ErrCreateResourceBadRequest(err)
		return NewErrResponse(c, &response)
	}

	if err = c.Validate(req); err != nil {
		resp := errs.ErrValidation(err.Error())
		return NewErrResponse(c, &resp)
	}

	offerKey := NewOfferKey(uuid.NewV4())
	offer := NewOffer(offerKey, authUserKey, nil)

	// keep track of resources so that one resource is not traded twice in the same transaction
	var usedResources = map[string]bool{}

	items := make([]OfferItem, len(req.Offer.Items))
	for i, reqItem := range req.Offer.Items {

		fromUser := NewUserKey(reqItem.From)
		toUser := NewUserKey(reqItem.To)

		// make sure that the offer item.from != item.to
		if fromUser == toUser {
			return NewErrResponse(c, fmt.Errorf("one cannot trade with oneself"))
		}

		// make sure the resource appears only once in an offer
		if reqItem.Type == ResourceItem {
			id := reqItem.ResourceId
			if _, ok := usedResources[*id]; ok {
				return NewErrResponse(c, fmt.Errorf("resources can only appear once in a transaction"))
			}
		}

		var item OfferItem
		itemKey := NewOfferItemKey(uuid.NewV4(), offerKey)

		if reqItem.Type == TimeItem {

			// make sure we're offering >0 hours
			seconds := reqItem.TimeInSeconds
			if *seconds <= 0 {
				return NewErrResponse(c, fmt.Errorf("time traded must be more than 0"))
			}
			item = NewTimeOfferItem(itemKey, fromUser, toUser, *seconds)

		} else if reqItem.Type == ResourceItem {

			// mark this resource as used
			rsId := reqItem.ResourceId
			usedResources[*rsId] = true

			// retrieve resource key
			resourceId, err := uuid.FromString(*rsId)
			if err != nil {
				return NewErrResponse(c, err)
			}
			resourceKey := NewResourceKey(resourceId)

			// make sure that the item.from = resource owner
			var resource Resource
			err = h.resourceStore.GetByKey(resourceKey, &resource)
			if err != nil {
				return NewErrResponse(c, err)
			}
			if resource.GetUserKey() != fromUser {
				res := errs.ErrTransactionResourceOwnerMismatch()
				return NewErrResponse(c, &res)
			}

			item = NewResourceOffer(itemKey, fromUser, toUser, resourceKey)
		}

		items[i] = item

	}

	err = h.tradingStore.SaveOffer(offer, items)
	if err != nil {
		return NewErrResponse(c, err)
	}

	o, err := h.tradingStore.GetOffer(offerKey)
	if err != nil {
		return NewErrResponse(c, err)
	}

	i, err := h.tradingStore.GetItems(offerKey)
	if err != nil {
		return NewErrResponse(c, err)
	}

	d, err := h.tradingStore.GetDecisions(offerKey)
	if err != nil {
		return NewErrResponse(c, err)
	}

	webOffer, err := h.mapToWebOffer(o, i, d)
	if err != nil {
		return NewErrResponse(c, err)
	}

	return c.JSON(http.StatusCreated, web.GetOfferResponse{Offer: *webOffer})

}

func (h *Handler) AcceptOffer(c echo.Context) error {

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

	var didAllOtherParticipantsAlreadyAccept = true
	var currentUserDecision *OfferDecision

	// retrieving the current user decision
	// also retrieve is this approval is the last one needed so that everyone approved the offer
	for _, decision := range decisions {
		decisionKey := decision.GetKey()
		decisionUserKey := decisionKey.GetUserKey()
		if decisionUserKey != authUserKey {
			if decision.Decision != AcceptedDecision {
				didAllOtherParticipantsAlreadyAccept = false
			}
			continue
		} else {
			currentUserDecision = &decision
		}
	}

	// in the case that a user is approving an offer he's not part of
	if currentUserDecision == nil {
		return fmt.Errorf("not involved in that offer")
	}

	// saving the decision
	err = h.tradingStore.SaveDecision(offerKey, authUserKey, AcceptedDecision)
	if err != nil {
		return NewErrResponse(c, err)
	}

	// complete offer if everyone approved
	var currentUserLastOneToDecide = currentUserDecision.Decision != AcceptedDecision && didAllOtherParticipantsAlreadyAccept
	if currentUserLastOneToDecide {
		err = h.tradingStore.CompleteOffer(offerKey, AcceptedOffer)
		if err != nil {
			return NewErrResponse(c, err)
		}
	}

	return h.getWebOffer(c, http.StatusAccepted, offerKey)

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

	authSession := h.authorization.GetAuthUserSession(c)

	qry := trading.GetOffersQuery{
		UserKeys: []UserKey{authSession.GetUserKey()},
	}
	result, err := h.tradingStore.GetOffers(qry)
	if err != nil {
		return NewErrResponse(c, err)
	}

	webOffers := make([]web.Offer, len(result.Items))
	for i, item := range result.Items {

		decisions, err := h.tradingStore.GetDecisions(item.Offer.GetKey())
		if err != nil {
			return NewErrResponse(c, err)
		}

		items, err := h.tradingStore.GetItems(item.Offer.GetKey())
		if err != nil {
			return NewErrResponse(c, err)
		}

		webOffer, err := h.mapToWebOffer(item.Offer, items, decisions)
		if err != nil {
			return NewErrResponse(c, err)
		}

		webOffers[i] = *webOffer
	}

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

func (h *Handler) mapToWebOffer(offer Offer, items []OfferItem, decisions []OfferDecision) (*web.Offer, error) {

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

	var responseItems = make([]web.OfferItem, len(items))

	for i, offerItem := range items {
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
