package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"

	handler2 "github.com/commonpool/backend/pkg/resource/handler"
	"github.com/commonpool/backend/pkg/trading/handler"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func (s *IntegrationTestSuite) SubmitOffer(t *testing.T, ctx context.Context, userSession *models.UserSession, request *handler.SendOfferRequest) (*handler.OfferResponse, *http.Response) {
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/offers", request)
	s.server.Router.ServeHTTP(recorder, httpReq)
	response := &handler.OfferResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) AcceptOffer(t *testing.T, ctx context.Context, userSession *models.UserSession, offerKey keys.OfferKey) *http.Response {
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/accept", offerKey.ID.String()), nil)
	s.server.Router.ServeHTTP(recorder, httpReq)
	return recorder.Result()
}

func (s *IntegrationTestSuite) ConfirmResourceTransfer(t *testing.T, ctx context.Context, userSession *models.UserSession, offerItemKey keys.OfferItemKey) *http.Response {
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offer-items/%s/confirm/resource-transferred", offerItemKey.ID.String()), nil)
	s.server.Router.ServeHTTP(recorder, httpReq)
	return recorder.Result()
}

func (s *IntegrationTestSuite) DeclineOffer(t *testing.T, ctx context.Context, userSession *models.UserSession, offerKey keys.OfferKey) *http.Response {
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/decline", offerKey.ID.String()), nil)
	s.server.Router.ServeHTTP(recorder, httpReq)
	return recorder.Result()
}

// func GetTradingHistory(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *GetTradingHistoryRequest) (*GetTradingHistoryResponse, *http.Response) {
// 	c, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/transactions", request)
// 	assert.NoError(s.T(), a.GetTradingHistory(c))
// 	response := &web.GetTradingHistoryResponse{}
// 	t.Log(recorder.Body.String())
// 	return response, ReadResponse(s.T(), recorder, response)
// }

func (s *IntegrationTestSuite) TestUserCanSubmitOfferBetweenUsers() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	user2, delUser2 := s.testUser(s.T())
	defer delUser2()

	group, err := s.testGroup(s.T(), user1, user2)
	if !assert.NoError(s.T(), err) {
		return
	}

	ctx := context.Background()

	resp, _ := s.CreateResource(s.T(), ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})

	offerResp, httpResponse := s.SubmitOffer(s.T(), ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resp.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user1.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "",
			GroupID: group.ID,
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}
	assert.Equal(s.T(), 2, len(offerResp.Offer.Items))

}

func (s *IntegrationTestSuite) TestUserCanSubmitOfferBetweenUsersAndGroup() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	user2, delUser2 := s.testUser(s.T())
	defer delUser2()

	group, err := s.testGroup(s.T(), user1, user2)
	if !assert.NoError(s.T(), err) {
		return
	}

	ctx := context.Background()

	resource, httpResponse := s.CreateResource(s.T(), ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	offer, httpResponse := s.SubmitOffer(s.T(), ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resource.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewGroupTarget(group.ID), time.Hour*1),
			},
			Message: "",
			GroupID: group.ID,
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}
	assert.Equal(s.T(), 2, len(offer.Offer.Items))

}

func (s *IntegrationTestSuite) TestUserCanSubmitOfferBetweenGroupAndMultipleUsers() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	user2, delUser2 := s.testUser(s.T())
	defer delUser2()

	user3, delUser2 := s.testUser(s.T())
	defer delUser2()

	group, err := s.testGroup(s.T(), user1, user2, user3)
	if !assert.NoError(s.T(), err) {
		return
	}

	ctx := context.Background()

	resp, httpResponse := s.CreateResource(s.T(), ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	offerResp, httpResponse := s.SubmitOffer(s.T(), ctx, user1, &handler.SendOfferRequest{
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
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	assert.Equal(s.T(), 3, len(offerResp.Offer.Items))

	if err := s.UsersAcceptOffer(s.T(), ctx, offerResp.Offer, []*models.UserSession{user1, user2}); !assert.NoError(s.T(), err) {
		return
	}

}

func (s *IntegrationTestSuite) TestUsersCanAcceptOfferBetweenUsers() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	user2, delUser2 := s.testUser(s.T())
	defer delUser2()

	group, err := s.testGroup(s.T(), user1, user2)
	if !assert.NoError(s.T(), err) {
		return
	}

	ctx := context.Background()

	resp, httpResponse := s.CreateResource(s.T(), ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	offerResp, httpResponse := s.SubmitOffer(s.T(), ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resp.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
			GroupID: group.ID,
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	key, err := keys.ParseOfferKey(offerResp.Offer.ID)
	assert.NoError(s.T(), err)

	httpResponse = s.AcceptOffer(s.T(), ctx, user2, key)
	if !AssertStatusAccepted(s.T(), httpResponse) {
		return
	}

	offer, err := s.server.TradingStore.GetOffer(key)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), domain.Pending, offer.Status)

	httpResponse = s.AcceptOffer(s.T(), ctx, user1, key)
	if !AssertStatusAccepted(s.T(), httpResponse) {
		return
	}

	offer, err = s.server.TradingStore.GetOffer(key)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), domain.Approved, offer.Status)

}

func (s *IntegrationTestSuite) TestUserCannotCreateOfferForResourceNotSharedWithGroup() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	user2, delUser2 := s.testUser(s.T())
	defer delUser2()

	ctx := context.Background()

	group1, err := s.testGroup(s.T(), user1, user2)
	if !assert.NoError(s.T(), err) {
		return
	}
	group2, err := s.testGroup(s.T(), user1, user2)
	if !assert.NoError(s.T(), err) {
		return
	}

	resource, httpResponse := s.CreateResource(s.T(), ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group2.ID,
				},
			},
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	_, httpResponse = s.SubmitOffer(s.T(), ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resource.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			GroupID: group1.ID,
			Message: "Howdy :)",
		},
	})
	if !AssertStatusForbidden(s.T(), httpResponse) {
		return
	}

}

func (s *IntegrationTestSuite) TestCannotCreateResourceTransferItemForResourceAlreadyOwned() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	user2, delUser2 := s.testUser(s.T())
	defer delUser2()

	ctx := context.Background()

	group, err := s.testGroup(s.T(), user1, user2)
	if !assert.NoError(s.T(), err) {
		return
	}

	resource, httpResponse := s.CreateResource(s.T(), ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	_, httpResponse = s.SubmitOffer(s.T(), ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resource.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			GroupID: group.ID,
			Message: "Howdy :)",
		},
	})
	if !AssertStatusForbidden(s.T(), httpResponse) {
		return
	}

}

func (s *IntegrationTestSuite) TestUsersCanDeclineOffer() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	user2, delUser2 := s.testUser(s.T())
	defer delUser2()

	group, err := s.testGroup(s.T(), user1, user2)
	if !assert.NoError(s.T(), err) {
		return
	}

	ctx := context.Background()

	resp, httpResponse := s.CreateResource(s.T(), ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	createOffer, httpResponse := s.SubmitOffer(s.T(), ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resp.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
			GroupID: group.ID,
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	offerKey := keys.MustParseOfferKey(createOffer.Offer.ID)

	acceptOffer := s.AcceptOffer(s.T(), ctx, user2, offerKey)
	AssertOK(s.T(), acceptOffer)

	offer, err := s.server.TradingStore.GetOffer(offerKey)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), domain.Pending, offer.Status)

	declineOffer := s.DeclineOffer(s.T(), ctx, user1, offerKey)
	AssertStatusNoContent(s.T(), declineOffer)

	offer, err = s.server.TradingStore.GetOffer(offerKey)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), domain.Declined, offer.Status)

}

func (s *IntegrationTestSuite) TestSendingOfferShouldCreateChatChannelBetweenUsers() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	user2, delUser2 := s.testUser(s.T())
	defer delUser2()

	group, err := s.testGroup(s.T(), user1, user2)
	if !assert.NoError(s.T(), err) {
		return
	}

	ctx := context.Background()

	resp, httpResponse := s.CreateResource(s.T(), ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupID: group.ID,
				},
			},
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	_, httpResponse = s.SubmitOffer(s.T(), ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []handler.SendOfferPayloadItem{
				*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resp.Resource.Id),
				*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
			},
			Message: "Howdy :)",
			GroupID: group.ID,
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	channelKey, err := s.server.ChatService.GetConversationChannelKey(ctx, keys.NewUserKeys([]keys.UserKey{user1.GetUserKey(), user2.GetUserKey()}))
	assert.NoError(s.T(), err)

	subscriptions, err := s.server.ChatStore.GetSubscriptionsForChannel(ctx, channelKey)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 2, len(subscriptions))

	_, err = s.server.ChatStore.GetChannel(ctx, channelKey)
	assert.NoError(s.T(), err)

}

func (s *IntegrationTestSuite) TestSendingOfferBetweenMultiplePeopleShouldCreateChatChannelBetweenUsers() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	user2, delUser2 := s.testUser(s.T())
	defer delUser2()

	user3, delUser2 := s.testUser(s.T())
	defer delUser2()

	group, err := s.testGroup(s.T(), user1, user2, user3)
	if !assert.NoError(s.T(), err) {
		return
	}

	ctx := context.Background()

	res1, httpResponse := s.CreateResource(s.T(), ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{GroupID: group.ID},
			},
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}
	res2, httpResponse := s.CreateResource(s.T(), ctx, user2, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{GroupID: group.ID},
			},
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	_, httpResponse = s.SubmitOffer(s.T(), ctx, user1, &handler.SendOfferRequest{
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
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	channelKey, err := s.server.ChatService.GetConversationChannelKey(ctx, keys.NewUserKeys([]keys.UserKey{user1.GetUserKey(), user2.GetUserKey(), user3.GetUserKey()}))
	assert.NoError(s.T(), err)

	subscriptions, err := s.server.ChatStore.GetSubscriptionsForChannel(ctx, channelKey)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 3, len(subscriptions))

	_, err = s.server.ChatStore.GetChannel(ctx, channelKey)
	assert.NoError(s.T(), err)

}

func (s *IntegrationTestSuite) TestCanGetTradingHistory() {
	return

	//
	//
	// user1, delUser1 := s.testUser(s.T())
	// defer delUser1()
	//
	// user2, delUser2 := s.testUser(s.T())
	// defer delUser2()
	//
	// ctx := context.Background()
	//
	// resource1, _ := CreateResource(s.T(), ctx, user1)
	// resource2, _ := CreateResource(s.T(), ctx, user2)
	//
	// offer1, offer1Http, _ := SubmitOffer(s.T(), ctx, user1, &handler.SendOfferRequest{
	// 	Offer: handler.SendOfferPayload{
	// 		Items: []handler.SendOfferPayloadItem{
	// 			*handler.NewResourceTransferItem(handler.NewUserTarget(user1.Subject), resource1.Resource.Id),
	// 			*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
	// 		},
	// 		Message: "Howdy :)",
	// 	},
	// })
	// AssertStatusCreated(s.T(), offer1Http)
	//
	// assert.NoError(s.T(), UsersAcceptOffer(s.T(), ctx, offer1.Offer, []*models.UserSession{user1, user2}))
	// assert.NoError(s.T(), UsersConfirmResourceTransferred(s.T(), ctx, offer1.Offer, []*models.UserSession{user1, user2}))
	//
	// offer2, offer2Http, _ := SubmitOffer(s.T(), ctx, user1, &handler.SendOfferRequest{
	// 	Offer: handler.SendOfferPayload{
	// 		Items: []handler.SendOfferPayloadItem{
	// 			*handler.NewResourceTransferItem(handler.NewUserTarget(user2.Subject), resource2.Resource.Id),
	// 			*handler.NewCreditTransferItem(handler.NewUserTarget(user2.Subject), handler.NewUserTarget(user1.Subject), time.Hour*2),
	// 		},
	// 		Message: "Howdy :)",
	// 	},
	// })
	// AssertStatusCreated(s.T(), offer2Http)
	//
	// assert.NoError(s.T(), UsersAcceptOffer(s.T(), ctx, offer2.Offer, []*models.UserSession{user1, user2}))
	// assert.NoError(s.T(), UsersConfirmResourceTransferred(s.T(), ctx, offer2.Offer, []*models.UserSession{user1, user2}))

}

//
// func TestGetExpandedTradingHistory() {
//
//
//
// 	user1, delUser1 := s.testUser(s.T())
// 	defer delUser1()
//
// 	user2, delUser2 := s.testUser(s.T())
// 	defer delUser2()
//
// 	ctx := context.Background()
//
// 	resource1, resource1Http := CreateResource(s.T(), ctx, user1)
// 	assert.Equal(s.T(), http.StatusCreated, resource1Http.StatusCode)
//
// 	resource2, resource2Http := CreateResource(s.T(), ctx, user2)
// 	assert.Equal(s.T(), http.StatusCreated, resource2Http.StatusCode)
//
// 	SubmitConfirmAcceptOffer(s.T(), ctx, user1, []*auth.UserSession{user1, user2, user3}, &web.SendOfferRequest{
// 		Offer: web.SendOfferPayload{
// 			Items: []web.SendOfferPayloadItem{
// 				*web.NewResourceTransferItem(user1.Subject, user3.Subject, resource1.Resource.Id),
// 				*web.NewCreditTransferItem(user3.Subject, user1.Subject, 6000),
// 			},
// 			Message: "Howdy :)",
// 		},
// 	})
// 	SubmitConfirmAcceptOffer(s.T(), ctx, user2, []*auth.UserSession{user1, user2, user3}, &web.SendOfferRequest{
// 		Offer: web.SendOfferPayload{
// 			Items: []web.SendOfferPayloadItem{
// 				*web.NewResourceTransferItem(user2.Subject, user3.Subject, resource2.Resource.Id),
// 				*web.NewCreditTransferItem(user3.Subject, user1.Subject, 6000),
// 			},
// 			Message: "Howdy :)",
// 		},
// 	})
// 	SubmitConfirmAcceptOffer(s.T(), ctx, user2, []*auth.UserSession{user1, user2, user3}, &web.SendOfferRequest{
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

func (s *IntegrationTestSuite) UsersConfirmResourceTransferred(t *testing.T, ctx context.Context, offer *handler.Offer, users []*models.UserSession) error {
	for _, offerItem := range offer.Items {

		if offerItem.Type != domain.ResourceTransfer {
			continue
		}

		var offerItemUsers []keys.UserKey
		for _, user := range users {
			offerItemUsers = append(offerItemUsers, user.GetUserKey())
		}

		for _, offerItemUser := range offerItemUsers {
			offerItemUserSession, err := s.findUserSession(offerItemUser.String(), users)
			if err != nil {
				panic(err)
			}
			if offerItemUserSession == nil {
				continue
			}
			offerKey := keys.MustParseOfferItemKey(offerItem.ID)
			httpResponse := s.ConfirmResourceTransfer(s.T(), ctx, offerItemUserSession, offerKey)
			assert.Equal(s.T(), http.StatusOK, httpResponse.StatusCode)
		}
	}
	return nil
}

func (s *IntegrationTestSuite) UsersAcceptOffer(t *testing.T, ctx context.Context, offer *handler.Offer, users []*models.UserSession) error {

	usersAccepted := map[keys.UserKey]bool{}

	var offerItemUsers []keys.UserKey
	for _, user := range users {
		offerItemUsers = append(offerItemUsers, user.GetUserKey())
	}

	for _, offerItemUser := range offerItemUsers {
		if alreadyAccepted, ok := usersAccepted[offerItemUser]; !alreadyAccepted || !ok {
			usersAccepted[offerItemUser] = true

			offerItemUserSession, err := s.findUserSession(offerItemUser.String(), users)
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

			httpResponse := s.AcceptOffer(s.T(), ctx, offerItemUserSession, offerKey)
			assert.Equal(s.T(), http.StatusOK, httpResponse.StatusCode)

		}
	}

	return nil
}

func (s *IntegrationTestSuite) findUserSession(subject string, users []*models.UserSession) (*models.UserSession, error) {
	for _, user := range users {
		if user.Subject == subject {
			return user, nil
		}
	}
	return nil, fmt.Errorf("could not find user session")
}
