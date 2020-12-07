package handler

import (
	"encoding/json"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/trading"
	"github.com/commonpool/backend/web"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrading(t *testing.T) {

	mockLoggedInAs(user1)
	res1 := createResource(t, "summary1", "desc1", resource.Offer)

	mockLoggedInAs(user2)
	res2 := createResource(t, "summary2", "desc2", resource.Offer)

	mockLoggedInAs(user1)
	offer := sendOffer(t,
		aResourceOffer(user1.Subject, user2.Subject, res1.Resource.Id),
		aResourceOffer(user2.Subject, user1.Subject, res2.Resource.Id),
	)

	assert.Equal(t, 2, len(offer.Offer.Items))
	assert.Equal(t, user1.Subject, offer.Offer.Items[0].From)
	assert.Equal(t, user2.Subject, offer.Offer.Items[0].To)
	assert.Equal(t, res1.Resource.Id, offer.Offer.Items[0].ResourceId)
	assert.Equal(t, user2.Subject, offer.Offer.Items[1].From)
	assert.Equal(t, user1.Subject, offer.Offer.Items[1].To)
	assert.Equal(t, res2.Resource.Id, offer.Offer.Items[1].ResourceId)
	assert.Equal(t, 2, len(offer.Offer.Decisions))
	assert.Equal(t, user1.Subject, offer.Offer.Decisions[0].UserID)
	assert.Equal(t, trading.PendingDecision, offer.Offer.Decisions[0].Decision)
	assert.Equal(t, offer.Offer.ID, offer.Offer.Decisions[0].OfferID)
	assert.Equal(t, user2.Subject, offer.Offer.Decisions[1].UserID)
	assert.Equal(t, trading.PendingDecision, offer.Offer.Decisions[1].Decision)
	assert.Equal(t, offer.Offer.ID, offer.Offer.Decisions[1].OfferID)

	_ = getOffers(t)

}

func aResourceOffer(from string, to string, resource string) web.SendOfferPayloadItem {
	return web.SendOfferPayloadItem{
		From:          from,
		To:            to,
		Type:          trading.ResourceItem,
		ResourceId:    &resource,
		TimeInSeconds: nil,
	}
}
func aTimeOffer(from string, to string, time int64) web.SendOfferPayloadItem {
	return web.SendOfferPayloadItem{
		From:          from,
		To:            to,
		Type:          trading.ResourceItem,
		ResourceId:    nil,
		TimeInSeconds: &time,
	}
}

func getOffers(t *testing.T) web.GetOffersResponse {
	_, _, rec, c := newRequest(echo.GET, "/api/v1/offers", nil)
	err := h.HandleGetOffers(c)
	if err != nil {
		panic(err)
	}
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	response := web.GetOffersResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	return response
}

func newSendOfferRequest(js string) (*httptest.ResponseRecorder, echo.Context) {
	_, _, rec, c := newRequest(echo.POST, "/api/v1/offers", &js)
	return rec, c
}

func sendOffer(t *testing.T, items ...web.SendOfferPayloadItem) web.GetOfferResponse {

	request := web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: items,
		},
	}

	js, err := json.Marshal(request)
	if err != nil {
		panic(err)
	}

	rec, c := newSendOfferRequest(string(js))
	err = h.HandleSendOffer(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	res := web.GetOfferResponse{}
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	return res
}
