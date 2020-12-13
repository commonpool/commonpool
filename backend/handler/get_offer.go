package handler

import (
	"github.com/commonpool/backend/model"
	"github.com/labstack/echo/v4"
	"github.com/satori/go.uuid"
	"net/http"
)

func (h *Handler) HandleGetOffer(c echo.Context) error {

	var err error

	offerIdStr := c.Param("id")

	offerId, err := uuid.FromString(offerIdStr)
	if err != nil {
		return err
	}
	offerKey := model.NewOfferKey(offerId)

	offer, err := h.getWebOffer(offerKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, offer)

}
