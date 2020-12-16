package handler

import (
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/handler"
	"github.com/labstack/echo/v4"
	"net/http"
)

type GetOffersResponse struct {
	Offers []Offer `json:"offers"`
}

func (h *TradingHandler) HandleGetOffers(c echo.Context) error {

	ctx, _ := handler.GetEchoContext(c, "HandleGetOffers")

	loggedInUser, err := auth.GetLoggedInUser(ctx)
	if err != nil {
		return err
	}
	loggedInUserKey := loggedInUser.GetUserKey()

	offersForUser, err := h.tradingService.GetOffersForUser(loggedInUserKey)
	if err != nil {
		return err
	}

	approversForOffers, err := h.tradingService.FindApproversForOffers(offersForUser.GetOfferKeys())
	if err != nil {
		return err
	}

	var webOffers []Offer
	for _, offerForUser := range offersForUser.Items {

		approversForOffer, err := approversForOffers.GetApproversForOffer(offerForUser.Offer.Key)
		if err != nil {
			return err
		}

		webOffer, err := h.mapToWebOffer(offerForUser.Offer, offerForUser.OfferItems, approversForOffer)
		if err != nil {
			return err
		}

		webOffers = append(webOffers, *webOffer)

	}

	return c.JSON(http.StatusOK, GetOffersResponse{
		Offers: webOffers,
	})

}
