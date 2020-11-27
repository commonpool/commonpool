package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/trading"
	"github.com/commonpool/backend/web"
	"github.com/go-playground/assert/v2"
	"net/http"
	"testing"
)

func SubmitOffer(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *web.SendOfferRequest) (*web.GetOfferResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/offers", request)
	PanicIfError(a.HandleSendOffer(c))
	response := &web.GetOfferResponse{}
	return response, ReadResponse(t, recorder, response)
}

func AcceptOffer(t *testing.T, ctx context.Context, userSession *auth.UserSession, offerKey model.OfferKey) *http.Response {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/accept", offerKey.ID.String()), nil)
	c.SetParamNames("id")
	c.SetParamValues(offerKey.ID.String())
	PanicIfError(a.HandleAcceptOffer(c))
	return recorder.Result()
}

func DeclineOffer(t *testing.T, ctx context.Context, userSession *auth.UserSession, offerKey model.OfferKey) *http.Response {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/decline", offerKey.ID.String()), nil)
	c.SetParamNames("id")
	c.SetParamValues(offerKey.ID.String())
	PanicIfError(a.DeclineOffer(c))
	return recorder.Result()
}


func TestUserCanSubmitOffer(t *testing.T) {

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, User1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary:          "Summary",
			Description:      "Description",
			Type:             resource.ResourceOffer,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith:       []web.InputResourceSharing{},
		},
	})

	var seconds int64 = 6000
	offerResp, httpOfferResp := SubmitOffer(t, ctx, User1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				{
					From:          User1.Subject,
					To:            User2.Subject,
					Type:          trading.ResourceItem,
					ResourceId:    &resp.Resource.Id,
					TimeInSeconds: nil,
				}, {
					From:          User2.Subject,
					To:            User1.Subject,
					Type:          trading.TimeItem,
					ResourceId:    nil,
					TimeInSeconds: &seconds,
				},
			},
			Message: "",
		},
	})
	assert.Equal(t, http.StatusCreated, httpOfferResp.StatusCode)
	assert.Equal(t, 2, len(offerResp.Offer.Items))

}

func TestCanAcceptOffer(t *testing.T) {

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, User1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary:          "Summary",
			Description:      "Description",
			Type:             resource.ResourceOffer,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith:       []web.InputResourceSharing{},
		},
	})

	var seconds int64 = 6000
	offerResp, _ := SubmitOffer(t, ctx, User1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				{
					From:          User1.Subject,
					To:            User2.Subject,
					Type:          trading.ResourceItem,
					ResourceId:    &resp.Resource.Id,
					TimeInSeconds: nil,
				}, {
					From:          User2.Subject,
					To:            User1.Subject,
					Type:          trading.TimeItem,
					ResourceId:    nil,
					TimeInSeconds: &seconds,
				},
			},
			Message: "Howdy :)",
		},
	})

	key, err := model.ParseOfferKey(offerResp.Offer.ID)
	PanicIfError(err)

	httpResp := AcceptOffer(t, ctx, User2, key)
	assert.Equal(t, http.StatusAccepted, httpResp.StatusCode)

	httpResp = AcceptOffer(t, ctx, User1, key)
	assert.Equal(t, http.StatusAccepted, httpResp.StatusCode)

}


func TestCanDeclineOffer(t *testing.T) {

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, User1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary:          "Summary",
			Description:      "Description",
			Type:             resource.ResourceOffer,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith:       []web.InputResourceSharing{},
		},
	})

	var seconds int64 = 6000
	offerResp, _ := SubmitOffer(t, ctx, User1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				{
					From:          User1.Subject,
					To:            User2.Subject,
					Type:          trading.ResourceItem,
					ResourceId:    &resp.Resource.Id,
					TimeInSeconds: nil,
				}, {
					From:          User2.Subject,
					To:            User1.Subject,
					Type:          trading.TimeItem,
					ResourceId:    nil,
					TimeInSeconds: &seconds,
				},
			},
			Message: "Howdy :)",
		},
	})

	key, err := model.ParseOfferKey(offerResp.Offer.ID)
	PanicIfError(err)

	httpResp := AcceptOffer(t, ctx, User2, key)
	assert.Equal(t, http.StatusAccepted, httpResp.StatusCode)

	httpResp = DeclineOffer(t, ctx, User1, key)
	assert.Equal(t, http.StatusAccepted, httpResp.StatusCode)

	of, err := TradingStore.GetOffer(key)
	PanicIfError(err)

	assert.Equal(t, trading.DeclinedOffer, of.Status)

	decisions, err := TradingStore.GetDecisions(key)
	PanicIfError(err)

	for _, decision := range decisions {
		if decision.UserID == User1KeyStr {
			assert.Equal(t, trading.DeclinedDecision, decision.Decision)
		} else if decision.UserID == User2KeyStr {
			assert.Equal(t, trading.AcceptedDecision, decision.Decision)
		} else {
			panic("")
		}
	}

}
