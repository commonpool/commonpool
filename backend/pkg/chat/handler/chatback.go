package handler

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/auth/authenticator/oidc"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/labstack/echo/v4"
	"net/http"
)

type InteractionPayloadType string

const (
	BlockActions InteractionPayloadType = "block_actions"
)

type InteractionCallback struct {
	Payload InteractionCallbackPayload `json:"payload"`
	Token   string                     `json:"token"`
}

type InteractionPayloadUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type InteractionCallbackPayload struct {
	Type        InteractionPayloadType             `json:"type"`
	User        InteractionPayloadUser             `json:"user"`
	TriggerId   string                             `json:"triggerId"`
	ResponseURL string                             `json:"responseUrl"`
	Message     *Message                           `json:"message"`
	Actions     []Action                           `json:"actions"`
	State       map[string]map[string]ElementState `json:"state"`
}

func (h *Handler) Chatback(c echo.Context) error {

	// Unmarshal request
	req := InteractionCallback{}
	if err := c.Bind(&req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	c.Logger().Info(req.Payload.Actions[0].ActionID)

	if req.Payload.Actions[0].ActionID == "accept_offer" {
		return h.HandleChatbackOfferAccepted(c, req)
	} else if req.Payload.Actions[0].ActionID == "confirm_service_provided" {
		return h.HandleChatbackConfirmServiceProvided(c, req)
	} else if req.Payload.Actions[0].ActionID == "confirm_resource_transferred" {
		return h.HandleChatbackConfirmResourceTransferred(c, req)
	} else if req.Payload.Actions[0].ActionID == "confirm_resource_borrowed" {
		return h.HandleChatbackConfirmResourceBorrowed(c, req)
	} else if req.Payload.Actions[0].ActionID == "confirm_resource_borrowed_returned" {
		return h.HandleChatbackConfirmResourceBorrowedReturned(c, req)
	}

	return nil

}

func (h *Handler) HandleChatbackConfirmServiceProvided(c echo.Context, req InteractionCallback) error {

	ctx, _ := handler.GetEchoContext(c, "HandleChatbackConfirmServiceProvided")
	offerKey, offerItemKey, err := h.getOfferKeyForChatbackAction(ctx, req)
	if err != nil {
		return err
	}

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	err = h.confirmServiceGiven.Execute(ctx, domain.NewConfirmServiceGiven(ctx, offerKey, offerItemKey, loggedInUser.GetUserKey()))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "OK")

}

func (h *Handler) HandleChatbackConfirmResourceTransferred(c echo.Context, req InteractionCallback) error {

	ctx, _ := handler.GetEchoContext(c, "HandleChatbackConfirmResourceTransferred")
	offerKey, offerItemKey, err := h.getOfferKeyForChatbackAction(ctx, req)
	if err != nil {
		return err
	}

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	err = h.confirmResourceGiven.Execute(ctx, domain.NewConfirmResourceGiven(ctx, offerKey, offerItemKey, loggedInUser.GetUserKey()))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "OK")

}

func (h *Handler) HandleChatbackConfirmResourceBorrowed(c echo.Context, req InteractionCallback) error {

	ctx, _ := handler.GetEchoContext(c, "HandleChatbackConfirmResourceBorrowed")
	offerKey, offerItemKey, err := h.getOfferKeyForChatbackAction(ctx, req)
	if err != nil {
		return err
	}

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	err = h.confirmBorrowed.Execute(ctx, domain.NewConfirmResourceBorrowed(ctx, offerKey, offerItemKey, loggedInUser.GetUserKey()))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "OK")

}

func (h *Handler) HandleChatbackConfirmResourceBorrowedReturned(c echo.Context, req InteractionCallback) error {

	ctx, _ := handler.GetEchoContext(c, "HandleChatbackConfirmResourceBorrowedReturned")
	offerKey, offerItemKey, err := h.getOfferKeyForChatbackAction(ctx, req)
	if err != nil {
		return err
	}

	loggedInUser, err := oidc.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}

	err = h.confirmReturned.Execute(ctx, domain.NewConfirmResourceReturned(ctx, offerKey, offerItemKey, loggedInUser.GetUserKey()))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "OK")
}

func (h *Handler) HandleChatbackOfferAccepted(c echo.Context, req InteractionCallback) error {

	ctx, _ := handler.GetEchoContext(c, "HandleChatbackOfferAccepted")
	offerKey, err := h.parseChatbackOfferKey(req)
	if err != nil {
		return err
	}

	err = h.acceptOffer.Execute(ctx, domain.NewAcceptOffer(ctx, offerKey))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "OK")
}

func (h *Handler) parseChatbackOfferKey(req InteractionCallback) (keys.OfferKey, error) {
	offerId := req.Payload.Actions[0].Value
	if offerId == nil {
		return keys.OfferKey{}, fmt.Errorf("offerId value is required")
	}
	return keys.ParseOfferKey(*offerId)
}

func (h *Handler) parseChatbackOfferItemKey(req InteractionCallback) (keys.OfferItemKey, error) {
	offerId := req.Payload.Actions[0].Value
	if offerId == nil {
		return keys.OfferItemKey{}, fmt.Errorf("offerItemId value is required")
	}
	return keys.ParseOfferItemKey(*offerId)
}

func (h *Handler) getOfferKeyForChatbackAction(ctx context.Context, req InteractionCallback) (keys.OfferKey, keys.OfferItemKey, error) {
	offerItemKey, err := h.parseChatbackOfferItemKey(req)
	if err != nil {
		return keys.OfferKey{}, keys.OfferItemKey{}, err
	}
	offerKey, err := h.getOfferKeyForOfferItemKey.Get(ctx, offerItemKey)
	if err != nil {
		return keys.OfferKey{}, keys.OfferItemKey{}, err
	}
	return offerKey, offerItemKey, nil
}
