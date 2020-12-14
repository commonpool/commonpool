package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) HandleGetOffers(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleGetOffers")

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	result, err := h.tradingStore.GetOffersForUser(loggedInUserKey)
	if err != nil {
		return err
	}

	approvers, err := h.tradingStore.FindApproversForOffers(result.GetOfferKeys())
	if err != nil {
		return err
	}

	var webOffers []web.Offer
	resultItems := result.Items
	for _, item := range resultItems {

		approversForOffer, err := approvers.GetApproversForOffer(item.Offer.Key)
		if err != nil {
			return err
		}

		webOffer, err := h.mapToWebOffer(item.Offer, item.OfferItems, approversForOffer)
		if err != nil {
			return err
		}

		webOffers = append(webOffers, *webOffer)

	}

	return c.JSON(http.StatusOK, web.GetOffersResponse{
		Offers: webOffers,
	})

}
