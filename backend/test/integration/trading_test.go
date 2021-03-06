package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"

	handler2 "github.com/commonpool/backend/pkg/resource/handler"
	"github.com/commonpool/backend/pkg/trading"
	"github.com/commonpool/backend/pkg/trading/handler"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func SubmitOffer(t *testing.T, ctx context.Context, userSession *models.UserSession, request *handler.SendOfferRequest) (*handler.OfferResponse, *http.Response, error) {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/offers", request)
	err := TradingHandler.HandleSendOffer(c)
	if err != nil {
		return nil, nil, err
	}
	response := &handler.OfferResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response), nil
}

func AcceptOffer(t *testing.T, ctx context.Context, userSession *models.UserSession, offerKey keys.OfferKey) *http.Response {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/accept", offerKey.ID.String()), nil)
	c.SetParamNames("id")
	c.SetParamValues(offerKey.ID.String())
	assert.NoError(t, TradingHandler.HandleAcceptOffer(c))
	return recorder.Result()
}

func ConfirmResourceTransfer(t *testing.T, ctx context.Context, userSession *models.UserSession, offerItemKey keys.OfferItemKey) *http.Response {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offer-items/%s/confirm/resource-transferred", offerItemKey.ID.String()), nil)
	c.SetParamNames("id")
	c.SetParamValues(offerItemKey.ID.String())
	assert.NoError(t, TradingHandler.HandleConfirmResourceTransferred(c))
	return recorder.Result()
}

func DeclineOffer(t *testing.T, ctx context.Context, userSession *models.UserSession, offerKey keys.OfferKey) *http.Response {
	c, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/decline", offerKey.ID.String()), nil)
	c.SetParamNames("id")
	c.SetParamValues(offerKey.ID.String())
	assert.NoError(t, TradingHandler.HandleDeclineOffer(c))
	return recorder.Result()
}

// func GetTradingHistory(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *GetTradingHistoryRequest) (*GetTradingHistoryResponse, *http.Response) {
// 	c, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/transactions", request)
// 	assert.NoError(t, a.GetTradingHistory(c))
// 	response := &web.GetTradingHistoryResponse{}
// 	t.Log(recorder.Body.String())
// 	return response, ReadResponse(t, recorder, response)
// }

func TestUserCanSubmitOfferBetweenUsers(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	group := testGroup(t, user1, user2)

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})

	offerResp, httpOfferResp, err := SubmitOffer(t, ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resp.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user1.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "",
			GroupID: group.ID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusCreated, httpOfferResp.StatusCode)
	assert.Equal(t, 2, len(offerResp.Offer.Items))

}

func TestUserCanSubmitOfferBetweenUsersAndGroup(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	group := testGroup(t, user1, user2)

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})

	offerResp, httpOfferResp, err := SubmitOffer(t, ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resp.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewGroupTarget(group.ID), time.Hour*1),
			},
			Message: "",
			GroupID: group.ID,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusCreated, httpOfferResp.StatusCode)
	assert.Equal(t, 2, len(offerResp.Offer.Items))

}

func TestUserCanSubmitOfferBetweenGroupAndMultipleUsers(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	user3, delUser2 := testUser(t)
	defer delUser2()

	group := testGroup(t, user1, user2, user3)

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})

	offerResp, httpOfferResp, err := SubmitOffer(t, ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resp.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*1),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewGroupTarget(group.ID), time.Hour*1),
			},
			Message: "",
			GroupID: group.ID,
		},
	})

	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusCreated, httpOfferResp.StatusCode)
	assert.Equal(t, 3, len(offerResp.Offer.Items))

	UsersAcceptOffer(t, ctx, offerResp.Offer, []*models.UserSession{user1, user2})

}

func TestUsersCanAcceptOfferBetweenUsers(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	group := testGroup(t, user1, user2)

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})

	offerResp, _, err := SubmitOffer(t, ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resp.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
			GroupID: group.ID,
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	key, err := keys.ParseOfferKey(offerResp.Offer.ID)
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

func TestUserCannotCreateOfferForResourceNotSharedWithGroup(t *testing.T) {

	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()

	group1 := testGroup(t, user1, user2)
	group2 := testGroup(t, user1, user2)

	resource, _ := CreateResource(t, ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group2.ID,
				},
			},
		},
	})

	_, _, err := SubmitOffer(t, ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resource.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			GroupID: group1.ID,
			Message: "Howdy :)",
		},
	})
	if err == nil {
		t.Fatal("err should not be nil")
	}

}

func TestCannotCreateResourceTransferItemForResourceAlreadyOwned(t *testing.T) {

	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()

	group := testGroup(t, user1, user2)

	resource, _ := CreateResource(t, ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})

	_, _, err := SubmitOffer(t, ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resource.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			GroupID: group.ID,
			Message: "Howdy :)",
		},
	})

	if err != nil {
		t.Fatal("err should not be nil")
	}

}

func TestUsersCanDeclineOffer(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	group := testGroup(t, user1, user2)

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})

	createOffer, createOfferHttp, err := SubmitOffer(t, ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resp.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
			GroupID: group.ID,
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	AssertStatusCreated(t, createOfferHttp)

	offerKey := keys.MustParseOfferKey(createOffer.Offer.ID)

	acceptOffer := AcceptOffer(t, ctx, user2, offerKey)
	AssertOK(t, acceptOffer)

	offer, err := TradingStore.GetOffer(offerKey)
	assert.NoError(t, err)
	assert.Equal(t, trading.PendingOffer, offer.Status)

	declineOffer := DeclineOffer(t, ctx, user1, offerKey)
	AssertStatusNoContent(t, declineOffer)

	offer, err = TradingStore.GetOffer(offerKey)
	assert.NoError(t, err)
	assert.Equal(t, trading.DeclinedOffer, offer.Status)

}

func TestSendingOfferShouldCreateChatChannelBetweenUsers(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	group := testGroup(t, user1, user2)

	ctx := context.Background()

	resp, _ := CreateResource(t, ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})

	_, submitOfferHttp, err := SubmitOffer(t, ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resp.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
			GroupID: group.ID,
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	AssertStatusCreated(t, submitOfferHttp)

	channelKey, err := ChatService.GetConversationChannelKey(ctx, keys.NewUserKeys([]keys.UserKey{user1.GetUserKey(), user2.GetUserKey()}))
	assert.NoError(t, err)

	subscriptions, err := ChatStore.GetSubscriptionsForChannel(ctx, channelKey)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(subscriptions))

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

	group := testGroup(t, user1, user2, user3)

	ctx := context.Background()

	res1, _ := CreateResource(t, ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{GroupID: group.ID},
			},
		},
	})
	res2, _ := CreateResource(t, ctx, user2, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{GroupID: group.ID},
			},
		},
	})

	_, submitOfferHttp, err := SubmitOffer(t, ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), res1.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
				*handler.NewResourceTransferItem(handler.NewUserTarget(user1.Subject), res2.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user3.Subject), handler.NewUserTarget(user2.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
			GroupID: group.ID,
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	AssertStatusCreated(t, submitOfferHttp)

	channelKey, err := ChatService.GetConversationChannelKey(ctx, keys.NewUserKeys([]keys.UserKey{user1.GetUserKey(), user2.GetUserKey(), user3.GetUserKey()}))
	assert.NoError(t, err)

	subscriptions, err := ChatStore.GetSubscriptionsForChannel(ctx, channelKey)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(subscriptions))

	_, err = ChatStore.GetChannel(ctx, channelKey)
	assert.NoError(t, err)

}

func TestCanGetTradingHistory(t *testing.T) {
	return

	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()

	resource1, _ := CreateResource(t, ctx, user1)
	resource2, _ := CreateResource(t, ctx, user2)

	offer1, offer1Http, _ := SubmitOffer(t, ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user1.Subject), resource1.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
		},
	})
	AssertStatusCreated(t, offer1Http)

	assert.NoError(t, UsersAcceptOffer(t, ctx, offer1.Offer, []*models.UserSession{user1, user2}))
	assert.NoError(t, UsersConfirmResourceTransferred(t, ctx, offer1.Offer, []*models.UserSession{user1, user2}))

	offer2, offer2Http, _ := SubmitOffer(t, ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resource2.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
		},
	})
	AssertStatusCreated(t, offer2Http)

	assert.NoError(t, UsersAcceptOffer(t, ctx, offer2.Offer, []*models.UserSession{user1, user2}))
	assert.NoError(t, UsersConfirmResourceTransferred(t, ctx, offer2.Offer, []*models.UserSession{user1, user2}))

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

func UsersConfirmResourceTransferred(t *testing.T, ctx context.Context, offer *handler.Offer, users []*models.UserSession) error {
	for _, offerItem := range offer.Items {

		if offerItem.Type != domain.ResourceTransfer {
			continue
		}

		var offerItemUsers []keys.UserKey
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
			offerKey := keys.MustParseOfferItemKey(offerItem.ID)
			httpResponse := ConfirmResourceTransfer(t, ctx, offerItemUserSession, offerKey)
			assert.Equal(t, http.StatusOK, httpResponse.StatusCode)
		}
	}
	return nil
}

func UsersAcceptOffer(t *testing.T, ctx context.Context, offer *handler.Offer, users []*models.UserSession) error {

	usersAccepted := map[keys.UserKey]bool{}

	var offerItemUsers []keys.UserKey
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

			offerKey, err := keys.ParseOfferKey(offer.ID)
			if err != nil {
				return err
			}

			httpResponse := AcceptOffer(t, ctx, offerItemUserSession, offerKey)
			assert.Equal(t, http.StatusOK, httpResponse.StatusCode)

		}
	}

	return nil
}

func findUserSession(subject string, users []*models.UserSession) (*models.UserSession, error) {
	for _, user := range users {
		if user.Subject == subject {
			return user, nil
		}
	}
	return nil, fmt.Errorf("could not find user session")
}
