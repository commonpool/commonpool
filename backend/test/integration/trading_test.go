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

func ConfirmItemGivenOrReceived(t *testing.T, ctx context.Context, userSession *auth.UserSession, offerItemKey model.OfferItemKey) *http.Response {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offer-items/%s/confirmation", offerItemKey.ID.String()), nil)
	c.SetParamNames("id")
	c.SetParamValues(offerItemKey.ID.String())
	PanicIfError(a.ConfirmItemReceivedOrGiven(c))
	return recorder.Result()
}

func DeclineOffer(t *testing.T, ctx context.Context, userSession *auth.UserSession, offerKey model.OfferKey) *http.Response {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/decline", offerKey.ID.String()), nil)
	c.SetParamNames("id")
	c.SetParamValues(offerKey.ID.String())
	PanicIfError(a.DeclineOffer(c))
	return recorder.Result()
}

func GetTradingHistory(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *web.GetTradingHistoryRequest) (*web.GetTradingHistoryResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/transactions", request)
	PanicIfError(a.GetTradingHistory(c))
	response := &web.GetTradingHistoryResponse{}
	return response, ReadResponse(t, recorder, response)
}

func TestUserCanSubmitOffer(t *testing.T) {
	teardown()
	setup()

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
	teardown()
	setup()

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
	assert.Equal(t, http.StatusOK, httpResp.StatusCode)

	httpResp = AcceptOffer(t, ctx, User1, key)
	assert.Equal(t, http.StatusOK, httpResp.StatusCode)

}

func TestCanDeclineOffer(t *testing.T) {
	teardown()
	setup()

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
	assert.Equal(t, http.StatusOK, httpResp.StatusCode)

	httpResp = DeclineOffer(t, ctx, User1, key)
	assert.Equal(t, http.StatusOK, httpResp.StatusCode)

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

func TestSendingOfferShouldCreateChatChannelBetweenUsers(t *testing.T) {
	teardown()
	setup()

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
	_, _ = SubmitOffer(t, ctx, User1, &web.SendOfferRequest{
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

	channelKey, err := ChatService.GetConversationChannelKey(ctx, model.NewUserKeys([]model.UserKey{User1.GetUserKey(), User2.GetUserKey()}))
	PanicIfError(err)

	subs, err := ChatStore.GetSubscriptionsForChannel(ctx, channelKey)
	PanicIfError(err)

	assert.Equal(t, 2, len(subs))

	_, err = ChatStore.GetChannel(ctx, channelKey)
	PanicIfError(err)

}

func TestSendingOfferBetweenMultiplePeopleShouldCreateChatChannelBetweenUsers(t *testing.T) {
	teardown()
	setup()

	ctx := context.Background()

	resp1, _ := CreateResource(t, ctx, User1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary:          "Summary",
			Description:      "Description",
			Type:             resource.ResourceOffer,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith:       []web.InputResourceSharing{},
		},
	})

	resp2, _ := CreateResource(t, ctx, User2, &web.CreateResourceRequest{
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
	_, _ = SubmitOffer(t, ctx, User1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				{
					From:          User1.Subject,
					To:            User2.Subject,
					Type:          trading.ResourceItem,
					ResourceId:    &resp1.Resource.Id,
					TimeInSeconds: nil,
				}, {
					From:          User2.Subject,
					To:            User1.Subject,
					Type:          trading.TimeItem,
					ResourceId:    nil,
					TimeInSeconds: &seconds,
				}, {
					From:          User2.Subject,
					To:            User3.Subject,
					Type:          trading.ResourceItem,
					ResourceId:    &resp2.Resource.Id,
					TimeInSeconds: nil,
				}, {
					From:          User3.Subject,
					To:            User2.Subject,
					Type:          trading.TimeItem,
					ResourceId:    nil,
					TimeInSeconds: &seconds,
				},
			},
			Message: "Howdy :)",
		},
	})

	channelKey, err := ChatService.GetConversationChannelKey(ctx, model.NewUserKeys([]model.UserKey{User1.GetUserKey(), User2.GetUserKey(), User3.GetUserKey()}))
	PanicIfError(err)

	subs, err := ChatStore.GetSubscriptionsForChannel(ctx, channelKey)
	PanicIfError(err)

	assert.Equal(t, 3, len(subs))

	_, err = ChatStore.GetChannel(ctx, channelKey)
	PanicIfError(err)

}

func TestCanGetTradingHistory(t *testing.T) {
	teardown()
	setup()

	ctx := context.Background()

	Db.Delete(trading.Offer{}, "1 = 1")
	Db.Delete(trading.OfferItem{}, "1 = 1")
	Db.Delete(trading.OfferDecision{}, "1 = 1")

	resource1, _ := CreateResource(t, ctx, User1, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary:          "Summary",
			Description:      "Description",
			Type:             resource.ResourceOffer,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith:       []web.InputResourceSharing{},
		},
	})

	resource2, _ := CreateResource(t, ctx, User2, &web.CreateResourceRequest{
		Resource: web.CreateResourcePayload{
			Summary:          "Summary",
			Description:      "Description",
			Type:             resource.ResourceOffer,
			ValueInHoursFrom: 1,
			ValueInHoursTo:   3,
			SharedWith:       []web.InputResourceSharing{},
		},
	})

	offer1, _ := SubmitOffer(t, ctx, User1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				AResourceOfferItem(User1.Subject, User2.Subject, resource1.Resource.Id),
				ATimeOfferItem(User2.Subject, User1.Subject, 6000),
			},
			Message: "Howdy :)",
		},
	})

	PanicIfError(UsersAcceptOffer(t, ctx, offer1.Offer, []*auth.UserSession{User1, User2}))
	PanicIfError(UsersConfirmItems(t, ctx, offer1.Offer, []*auth.UserSession{User1, User2}))

	offer2, _ := SubmitOffer(t, ctx, User1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				AResourceOfferItem(User2.Subject, User1.Subject, resource2.Resource.Id),
				ATimeOfferItem(User1.Subject, User2.Subject, 6000),
			},
			Message: "Howdy :)",
		},
	})

	PanicIfError(UsersAcceptOffer(t, ctx, offer2.Offer, []*auth.UserSession{User1, User2}))
	PanicIfError(UsersConfirmItems(t, ctx, offer2.Offer, []*auth.UserSession{User1, User2}))

	tradingHistory, httpRes := GetTradingHistory(t, ctx, User1, &web.GetTradingHistoryRequest{
		UserIDs: []string{User1Key.String(), User2Key.String()},
	})
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)

	assert.Equal(t, 4, len(tradingHistory.Entries))

}

func AResourceOfferItem(from string, to string, resourceId string) web.SendOfferPayloadItem {
	return web.SendOfferPayloadItem{
		From:          from,
		To:            to,
		Type:          trading.ResourceItem,
		ResourceId:    &resourceId,
		TimeInSeconds: nil,
	}
}

func ATimeOfferItem(from string, to string, seconds int64) web.SendOfferPayloadItem {
	return web.SendOfferPayloadItem{
		From:          from,
		To:            to,
		Type:          trading.TimeItem,
		ResourceId:    nil,
		TimeInSeconds: &seconds,
	}
}

func UsersConfirmItems(t *testing.T, ctx context.Context, offer web.Offer, users []*auth.UserSession) error {
	for _, offerItem := range offer.Items {

		if offerItem.Type == trading.TimeItem {
			continue
		}

		offerItemUsers := []model.UserKey{model.NewUserKey(offerItem.FromUserID), model.NewUserKey(offerItem.ToUserID)}
		for _, offerItemUser := range offerItemUsers {
			offerItemUserSession := findUserSession(offerItemUser.String(), users)
			if offerItemUserSession == nil {
				continue
			}
			offerKey := model.MustParseOfferItemKey(offerItem.ID)
			httpResponse := ConfirmItemGivenOrReceived(t, ctx, offerItemUserSession, offerKey)
			assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
		}
	}
	return nil
}

func UsersAcceptOffer(t *testing.T, ctx context.Context, offer web.Offer, users []*auth.UserSession) error {
	usersAccepted := map[model.UserKey]bool{}
	for _, offerItem := range offer.Items {
		offerItemUsers := []model.UserKey{model.NewUserKey(offerItem.FromUserID), model.NewUserKey(offerItem.ToUserID)}
		for _, offerItemUser := range offerItemUsers {
			if alreadyAccepted, ok := usersAccepted[offerItemUser]; !alreadyAccepted || !ok {
				usersAccepted[offerItemUser] = true

				offerItemUserSession := findUserSession(offerItemUser.String(), users)
				if offerItemUserSession == nil {
					continue
				}

				offerKey, err := model.ParseOfferKey(offer.ID)
				if err != nil {
					return err
				}

				httpResponse := AcceptOffer(t, ctx, offerItemUserSession, offerKey)
				assert.Equal(t, http.StatusOK, httpResponse.StatusCode)

			}
		}
	}
	return nil
}

func findUserSession(subject string, users []*auth.UserSession) *auth.UserSession {
	var offerItemUserSession *auth.UserSession
	for _, user := range users {
		if user.Subject == subject {
			offerItemUserSession = user
			continue
		}
	}
	return offerItemUserSession
}

func getOfferItemsReceivedByUser(offerItems []web.OfferItem, userKey model.UserKey) []web.OfferItem {
	var result []web.OfferItem
	for _, offerItem := range offerItems {
		if offerItem.ToUserID == userKey.String() {
			result = append(result, offerItem)
		}
	}
	return result
}
