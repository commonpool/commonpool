package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/trading"
	"github.com/commonpool/backend/web"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func SubmitOffer(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *web.SendOfferRequest) (*web.GetOfferResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/offers", request)
	assert.NoError(t, a.HandleSendOffer(c))
	response := &web.GetOfferResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response)
}

func AcceptOffer(t *testing.T, ctx context.Context, userSession *auth.UserSession, offerKey model.OfferKey) *http.Response {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/accept", offerKey.ID.String()), nil)
	c.SetParamNames("id")
	c.SetParamValues(offerKey.ID.String())
	assert.NoError(t, a.HandleAcceptOffer(c))
	return recorder.Result()
}

func ConfirmResourceTransfer(t *testing.T, ctx context.Context, userSession *auth.UserSession, offerItemKey model.OfferItemKey) *http.Response {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offer-items/%s/confirm/resource-transferred", offerItemKey.ID.String()), nil)
	c.SetParamNames("id")
	c.SetParamValues(offerItemKey.ID.String())
	assert.NoError(t, a.HandleConfirmResourceTransferred(c))
	return recorder.Result()
}

func DeclineOffer(t *testing.T, ctx context.Context, userSession *auth.UserSession, offerKey model.OfferKey) *http.Response {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/decline", offerKey.ID.String()), nil)
	c.SetParamNames("id")
	c.SetParamValues(offerKey.ID.String())
	assert.NoError(t, a.HandleDeclineOffer(c))
	return recorder.Result()
}

func GetTradingHistory(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *web.GetTradingHistoryRequest) (*web.GetTradingHistoryResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/transactions", request)
	assert.NoError(t, a.GetTradingHistory(c))
	response := &web.GetTradingHistoryResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response)
}

func TestUserCanSubmitOffer(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, user1)

	offerResp, httpOfferResp := SubmitOffer(t, ctx, user1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				*web.NewResourceTransferItem(web.NewUserTarget(user1.Subject), resp.Resource.Id),
				*web.NewCreditTransferItem(web.NewUserTarget(user2.Subject), web.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "",
		},
	})
	assert.Equal(t, http.StatusCreated, httpOfferResp.StatusCode)
	assert.Equal(t, 2, len(offerResp.Offer.Items))

}

func TestUsersCanAcceptOffer(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, user1)

	offerResp, _ := SubmitOffer(t, ctx, user1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				*web.NewResourceTransferItem(web.NewUserTarget(user2.Subject), resp.Resource.Id),
				*web.NewCreditTransferItem(web.NewUserTarget(user2.Subject), web.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
		},
	})

	key, err := model.ParseOfferKey(offerResp.Offer.ID)
	assert.NoError(t, err)

	httpResp := AcceptOffer(t, ctx, user2, key)
	assert.Equal(t, http.StatusOK, httpResp.StatusCode)

	offer, err := TradingStore.GetOffer(key)
	assert.NoError(t, err)
	assert.Equal(t, trading.PendingOffer, offer.Status)

	httpResp = AcceptOffer(t, ctx, user1, key)
	assert.Equal(t, http.StatusOK, httpResp.StatusCode)

	offer, err = TradingStore.GetOffer(key)
	assert.NoError(t, err)
	assert.Equal(t, trading.AcceptedOffer, offer.Status)

}

func TestCanDeclineOffer(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, user1)

	offerResp, _ := SubmitOffer(t, ctx, user1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				*web.NewResourceTransferItem(web.NewUserTarget(user1.Subject), resp.Resource.Id),
				*web.NewCreditTransferItem(web.NewUserTarget(user2.Subject), web.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
		},
	})

	key, err := model.ParseOfferKey(offerResp.Offer.ID)
	assert.NoError(t, err)

	httpResp := AcceptOffer(t, ctx, user2, key)
	assert.Equal(t, http.StatusOK, httpResp.StatusCode)

	offer, err := TradingStore.GetOffer(key)
	assert.NoError(t, err)
	assert.Equal(t, trading.PendingOffer, offer.Status)

	httpResp = DeclineOffer(t, ctx, user1, key)
	assert.Equal(t, http.StatusOK, httpResp.StatusCode)

	of, err := TradingStore.GetOffer(key)
	assert.NoError(t, err)
	assert.Equal(t, trading.DeclinedOffer, of.Status)

}

func TestSendingOfferShouldCreateChatChannelBetweenUsers(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, user1)

	_, _ = SubmitOffer(t, ctx, user1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				*web.NewResourceTransferItem(web.NewUserTarget(user1.Subject), resp.Resource.Id),
				*web.NewCreditTransferItem(web.NewUserTarget(user2.Subject), web.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
		},
	})

	channelKey, err := ChatService.GetConversationChannelKey(ctx, model.NewUserKeys([]model.UserKey{user1.GetUserKey(), user2.GetUserKey()}))
	assert.NoError(t, err)

	subs, err := ChatStore.GetSubscriptionsForChannel(ctx, channelKey)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(subs))

	_, err = ChatStore.GetChannel(ctx, channelKey)
	assert.NoError(t, err)

}

func TestSendingOfferBetweenMultiplePeopleShouldCreateChatChannelBetweenUsers(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	user3, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()

	resp1, _ := CreateResource(t, ctx, user1)
	resp2, _ := CreateResource(t, ctx, user2)

	_, _ = SubmitOffer(t, ctx, user1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				*web.NewResourceTransferItem(web.NewUserTarget(user1.Subject), resp1.Resource.Id),
				*web.NewCreditTransferItem(web.NewUserTarget(user2.Subject), web.NewUserTarget(user1.Subject), time.Hour*2),
				*web.NewResourceTransferItem(web.NewUserTarget(user2.Subject), resp2.Resource.Id),
				*web.NewCreditTransferItem(web.NewUserTarget(user3.Subject), web.NewUserTarget(user2.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
		},
	})

	channelKey, err := ChatService.GetConversationChannelKey(ctx, model.NewUserKeys([]model.UserKey{user1.GetUserKey(), user2.GetUserKey(), user3.GetUserKey()}))
	assert.NoError(t, err)

	subs, err := ChatStore.GetSubscriptionsForChannel(ctx, channelKey)
	assert.NoError(t, err)

	assert.Equal(t, 3, len(subs))

	_, err = ChatStore.GetChannel(ctx, channelKey)
	assert.NoError(t, err)

}

func TestCanGetTradingHistory(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()

	resource1, _ := CreateResource(t, ctx, user1)
	resource2, _ := CreateResource(t, ctx, user2)

	offer1, _ := SubmitOffer(t, ctx, user1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				*web.NewResourceTransferItem(web.NewUserTarget(user1.Subject), resource1.Resource.Id),
				*web.NewCreditTransferItem(web.NewUserTarget(user2.Subject), web.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
		},
	})

	assert.NoError(t, UsersAcceptOffer(t, ctx, offer1.Offer, []*auth.UserSession{user1, user2}))
	assert.NoError(t, UsersConfirmResourceTransferred(t, ctx, offer1.Offer, []*auth.UserSession{user1, user2}))

	offer2, _ := SubmitOffer(t, ctx, user1, &web.SendOfferRequest{
		Offer: web.SendOfferPayload{
			Items: []web.SendOfferPayloadItem{
				*web.NewResourceTransferItem(web.NewUserTarget(user2.Subject), resource2.Resource.Id),
				*web.NewCreditTransferItem(web.NewUserTarget(user2.Subject), web.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
		},
	})

	assert.NoError(t, UsersAcceptOffer(t, ctx, offer2.Offer, []*auth.UserSession{user1, user2}))
	assert.NoError(t, UsersConfirmResourceTransferred(t, ctx, offer2.Offer, []*auth.UserSession{user1, user2}))

	tradingHistory, httpRes := GetTradingHistory(t, ctx, user1, &web.GetTradingHistoryRequest{
		UserIDs: []string{user1.Subject, user2.Subject},
	})

	assert.Equal(t, http.StatusOK, httpRes.StatusCode)
	assert.Equal(t, 4, len(tradingHistory.Entries))

}

//
// func TestGetExpandedTradingHistory(t *testing.T) {
//
// 	t.Parallel()
//
// 	user1, delUser1 := testUser(t)
// 	defer delUser1()
//
// 	user2, delUser2 := testUser(t)
// 	defer delUser2()
//
// 	ctx := context.Background()
//
// 	resource1, resource1Http := CreateResource(t, ctx, user1)
// 	assert.Equal(t, http.StatusCreated, resource1Http.StatusCode)
//
// 	resource2, resource2Http := CreateResource(t, ctx, user2)
// 	assert.Equal(t, http.StatusCreated, resource2Http.StatusCode)
//
// 	SubmitConfirmAcceptOffer(t, ctx, user1, []*auth.UserSession{user1, user2, user3}, &web.SendOfferRequest{
// 		Offer: web.SendOfferPayload{
// 			Items: []web.SendOfferPayloadItem{
// 				*web.NewResourceTransferItem(user1.Subject, user3.Subject, resource1.Resource.Id),
// 				*web.NewCreditTransferItem(user3.Subject, user1.Subject, 6000),
// 			},
// 			Message: "Howdy :)",
// 		},
// 	})
// 	SubmitConfirmAcceptOffer(t, ctx, user2, []*auth.UserSession{user1, user2, user3}, &web.SendOfferRequest{
// 		Offer: web.SendOfferPayload{
// 			Items: []web.SendOfferPayloadItem{
// 				*web.NewResourceTransferItem(user2.Subject, user3.Subject, resource2.Resource.Id),
// 				*web.NewCreditTransferItem(user3.Subject, user1.Subject, 6000),
// 			},
// 			Message: "Howdy :)",
// 		},
// 	})
// 	SubmitConfirmAcceptOffer(t, ctx, user2, []*auth.UserSession{user1, user2, user3}, &web.SendOfferRequest{
// 		Offer: web.SendOfferPayload{
// 			Items: []web.SendOfferPayloadItem{
// 				*web.NewResourceTransferItem(user3.Subject, user1.Subject, resource1.Resource.Id),
// 				*web.NewCreditTransferItem(user3.Subject, user1.Subject, 6000),
// 			},
// 			Message: "Howdy :)",
// 		},
// 	})
//
// }

func SubmitConfirmAcceptOffer(t *testing.T, ctx context.Context, user *auth.UserSession, allUsers []*auth.UserSession, offer *web.SendOfferRequest) {
	createdOffer, createdOfferHttp := SubmitOffer(t, ctx, user, offer)
	assert.Equal(t, http.StatusCreated, createdOfferHttp.StatusCode)
	assert.NoError(t, UsersAcceptOffer(t, ctx, createdOffer.Offer, allUsers))
	assert.NoError(t, UsersConfirmResourceTransferred(t, ctx, createdOffer.Offer, allUsers))
}

func UsersConfirmResourceTransferred(t *testing.T, ctx context.Context, offer web.Offer, users []*auth.UserSession) error {
	for _, offerItem := range offer.Items {

		if offerItem.Type != trading.ResourceTransfer {
			continue
		}

		var offerItemUsers []model.UserKey
		for _, user := range users {
			offerItemUsers = append(offerItemUsers, user.GetUserKey())
		}

		for _, offerItemUser := range offerItemUsers {
			offerItemUserSession, err := findUserSession(offerItemUser.String(), users)
			if err != nil {
				panic(err)
			}
			if offerItemUserSession == nil {
				continue
			}
			offerKey := model.MustParseOfferItemKey(offerItem.ID)
			httpResponse := ConfirmResourceTransfer(t, ctx, offerItemUserSession, offerKey)
			assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
		}
	}
	return nil
}

func UsersAcceptOffer(t *testing.T, ctx context.Context, offer web.Offer, users []*auth.UserSession) error {

	usersAccepted := map[model.UserKey]bool{}

	var offerItemUsers []model.UserKey
	for _, user := range users {
		offerItemUsers = append(offerItemUsers, user.GetUserKey())
	}

	for _, offerItemUser := range offerItemUsers {
		if alreadyAccepted, ok := usersAccepted[offerItemUser]; !alreadyAccepted || !ok {
			usersAccepted[offerItemUser] = true

			offerItemUserSession, err := findUserSession(offerItemUser.String(), users)
			if err != nil {
				panic(err)
			}

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

	return nil
}

func findUserSession(subject string, users []*auth.UserSession) (*auth.UserSession, error) {
	for _, user := range users {
		if user.Subject == subject {
			return user, nil
		}
	}
	return nil, fmt.Errorf("could not find user session")
}
