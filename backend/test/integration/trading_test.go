package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	readmodels "github.com/commonpool/backend/pkg/trading/readmodels"

	handler2 "github.com/commonpool/backend/pkg/resource/handler"
	"github.com/commonpool/backend/pkg/trading/handler"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func (s *IntegrationTestSuite) SubmitOffer(ctx context.Context, userSession *models.UserSession, request *handler.SendOfferRequest) (*handler.OfferResponse, *http.Response) {
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/offers", request)
	s.server.Router.ServeHTTP(recorder, httpReq)
	response := &handler.OfferResponse{}
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) AcceptOffer(ctx context.Context, userSession *models.UserSession, offerKey keys.OfferKey) *http.Response {
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/accept", offerKey.ID.String()), nil)
	s.server.Router.ServeHTTP(recorder, httpReq)
	return recorder.Result()
}

func (s *IntegrationTestSuite) ConfirmResourceTransfer(ctx context.Context, userSession *models.UserSession, offerItemKey keys.OfferItemKey) *http.Response {
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offer-items/%s/confirm/resource-transferred", offerItemKey.ID.String()), nil)
	s.server.Router.ServeHTTP(recorder, httpReq)
	return recorder.Result()
}

func (s *IntegrationTestSuite) DeclineOffer(ctx context.Context, userSession *models.UserSession, offerKey keys.OfferKey) *http.Response {
	httpReq, recorder := NewRequest(ctx, userSession, http.MethodPost, fmt.Sprintf("/api/v1/offers/%s/decline", offerKey.ID.String()), nil)
	s.server.Router.ServeHTTP(recorder, httpReq)
	return recorder.Result()
}

func (s *IntegrationTestSuite) AssertOfferStatus(ctx context.Context, offerKey keys.OfferKey, status domain.OfferStatus) bool {
	offer, err := s.server.GetOffer.Get(ctx, offerKey)
	if !assert.NoError(s.T(), err) {
		return false
	}
	return assert.Equal(s.T(), status, offer.Status)
}

func (s *IntegrationTestSuite) AssertOfferPending(ctx context.Context, offerKey keys.OfferKey) bool {
	return s.AssertOfferStatus(ctx, offerKey, domain.Pending)
}
func (s *IntegrationTestSuite) AssertOfferAccepted(ctx context.Context, offerKey keys.OfferKey) bool {
	return s.AssertOfferStatus(ctx, offerKey, domain.Approved)
}
func (s *IntegrationTestSuite) AssertOfferCompleted(ctx context.Context, offerKey keys.OfferKey) bool {
	return s.AssertOfferStatus(ctx, offerKey, domain.Completed)
}
func (s *IntegrationTestSuite) AssertOfferDeclined(ctx context.Context, offerKey keys.OfferKey) bool {
	return s.AssertOfferStatus(ctx, offerKey, domain.Declined)
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

	resp, _ := s.CreateResource(ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: handler2.NewInputResourceSharings().WithGroups(group.GroupKey),
		},
	})

	time.Sleep(1 * time.Second)

	offerResp, httpResponse := s.SubmitOffer(ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []domain.SubmitOfferItemBase{
				domain.NewResourceTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), resp.Resource.ResourceKey),
				domain.NewCreditTransferItemInputBase(keys.NewUserTarget(user1.GetUserKey()), keys.NewUserTarget(user1.GetUserKey()), time.Hour*2),
			},
			Message:  "",
			GroupKey: group.GroupKey,
		},
	})

	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	assert.Equal(s.T(), 2, len(offerResp.Offer.OfferItems))

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

	resource, httpResponse := s.CreateResource(ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: handler2.NewInputResourceSharings().WithGroups(group.GroupKey),
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	offer, httpResponse := s.SubmitOffer(ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []domain.SubmitOfferItemBase{
				domain.NewResourceTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), resource.Resource.ResourceKey),
				domain.NewCreditTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), keys.NewGroupTarget(group.GroupKey), time.Hour*1),
			},
			Message:  "",
			GroupKey: group.GroupKey,
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}
	assert.Equal(s.T(), 2, len(offer.Offer.OfferItems))

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

	resp, httpResponse := s.CreateResource(ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: handler2.NewInputResourceSharings().WithGroups(group.GroupKey),
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	offerResp, httpResponse := s.SubmitOffer(ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []domain.SubmitOfferItemBase{
				domain.NewResourceTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), resp.Resource.ResourceKey),
				domain.NewCreditTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), keys.NewUserTarget(user1.GetUserKey()), time.Hour*1),
				domain.NewCreditTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), keys.NewGroupTarget(group.GroupKey), time.Hour*1),
			},
			Message:  "",
			GroupKey: group.GroupKey,
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	if !assert.Equal(s.T(), 3, len(offerResp.Offer.OfferItems)) {
		return
	}

	time.Sleep(1 * time.Second)

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

	resp, httpResponse := s.CreateResource(ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: handler2.NewInputResourceSharings().WithGroups(group.GroupKey),
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	offerResp, httpResponse := s.SubmitOffer(ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []domain.SubmitOfferItemBase{
				domain.NewResourceTransferItemInputBase(user2.GetUserKey().Target(), resp.Resource.ResourceKey),
				domain.NewCreditTransferItemInputBase(user2.GetUserKey().Target(), user1.GetUserKey().Target(), time.Hour*2),
			},
			Message:  "Howdy :)",
			GroupKey: group.GroupKey,
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	httpResponse = s.AcceptOffer(ctx, user2, offerResp.Offer.OfferKey)
	if !AssertStatusAccepted(s.T(), httpResponse) {
		return
	}

	time.Sleep(1 * time.Second)

	if !s.AssertOfferPending(ctx, offerResp.Offer.OfferKey) {
		return
	}

	httpResponse = s.AcceptOffer(ctx, user1, offerResp.Offer.OfferKey)
	if !AssertStatusAccepted(s.T(), httpResponse) {
		return
	}

	time.Sleep(1 * time.Second)

	if !s.AssertOfferAccepted(ctx, offerResp.Offer.OfferKey) {
		return
	}

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

	resource, httpResponse := s.CreateResource(ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: handler2.NewInputResourceSharings().WithGroups(group2.GroupKey),
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	_, httpResponse = s.SubmitOffer(ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []domain.SubmitOfferItemBase{
				domain.NewResourceTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), resource.Resource.ResourceKey),
				domain.NewCreditTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), keys.NewUserTarget(user1.GetUserKey()), time.Hour*2),
			},
			GroupKey: group1.GroupKey,
			Message:  "Howdy :)",
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

	resource, httpResponse := s.CreateResource(ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: handler2.NewInputResourceSharings().WithGroups(group.GroupKey),
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	_, httpResponse = s.SubmitOffer(ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []domain.SubmitOfferItemBase{
				domain.NewResourceTransferItemInputBase(keys.NewUserTarget(user1.GetUserKey()), resource.Resource.ResourceKey),
			},
			GroupKey: group.GroupKey,
			Message:  "Howdy :)",
		},
	})
	if !AssertStatusBadRequest(s.T(), httpResponse) {
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

	resp, httpResponse := s.CreateResource(ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: handler2.NewInputResourceSharings().WithGroups(group.GroupKey),
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	createOffer, httpResponse := s.SubmitOffer(ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []domain.SubmitOfferItemBase{
				domain.NewResourceTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), resp.Resource.ResourceKey),
				domain.NewCreditTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), keys.NewUserTarget(user1.GetUserKey()), time.Hour*2),
			},
			Message:  "Howdy :)",
			GroupKey: group.GroupKey,
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	httpResponse = s.AcceptOffer(ctx, user2, createOffer.Offer.OfferKey)
	if !AssertStatusAccepted(s.T(), httpResponse) {
		return
	}

	if !s.AssertOfferPending(ctx, createOffer.Offer.OfferKey) {
		return
	}

	declineOffer := s.DeclineOffer(ctx, user1, createOffer.Offer.OfferKey)
	AssertStatusAccepted(s.T(), declineOffer)

	time.Sleep(600 * time.Millisecond)

	if !s.AssertOfferDeclined(ctx, createOffer.Offer.OfferKey) {
		return
	}

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

	resp, httpResponse := s.CreateResource(ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: handler2.NewInputResourceSharings().WithGroups(group.GroupKey),
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	_, httpResponse = s.SubmitOffer(ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []domain.SubmitOfferItemBase{
				domain.NewResourceTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), resp.Resource.ResourceKey),
				domain.NewCreditTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), keys.NewUserTarget(user1.GetUserKey()), time.Hour*2),
			},
			Message:  "Howdy :)",
			GroupKey: group.GroupKey,
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

	res1, httpResponse := s.CreateResource(ctx, user1, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: handler2.NewInputResourceSharings().WithGroups(group.GroupKey),
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}
	res2, httpResponse := s.CreateResource(ctx, user2, &handler2.CreateResourceRequest{
		Resource: handler2.CreateResourcePayload{
			SharedWith: []handler2.InputResourceSharing{
				{
					GroupKey: group.GroupKey,
				},
			},
		},
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}

	_, httpResponse = s.SubmitOffer(ctx, user1, &handler.SendOfferRequest{
		Offer: handler.SendOfferPayload{
			Items: []domain.SubmitOfferItemBase{
				domain.NewResourceTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), res1.Resource.ResourceKey),
				domain.NewCreditTransferItemInputBase(keys.NewUserTarget(user2.GetUserKey()), keys.NewUserTarget(user1.GetUserKey()), time.Hour*2),
				domain.NewResourceTransferItemInputBase(keys.NewUserTarget(user1.GetUserKey()), res2.Resource.ResourceKey),
				domain.NewCreditTransferItemInputBase(keys.NewUserTarget(user3.GetUserKey()), keys.NewUserTarget(user2.GetUserKey()), time.Hour*2),
			},
			Message:  "Howdy :)",
			GroupKey: group.GroupKey,
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
	// 		Items: []domain.SubmitOfferItemBase{
	// 			*domain.NewResourceTransferItemInputBase(domain.NewUserTarget(user1.GetUserKey()), resource1.Resource.Id),
	// 			*domain.NewCreditTransferItemInputBase(domain.NewUserTarget(user2.GetUserKey()), domain.NewUserTarget(user1.GetUserKey()), time.Hour*2),
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
	// 		Items: []domain.SubmitOfferItemBase{
	// 			*domain.NewResourceTransferItemInputBase(domain.NewUserTarget(user2.GetUserKey()), resource2.Resource.Id),
	// 			*domain.NewCreditTransferItemInputBase(domain.NewUserTarget(user2.GetUserKey()), domain.NewUserTarget(user1.GetUserKey()), time.Hour*2),
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
// 				*web.NewResourceTransferItemInputBase(user1.Subject, user3.Subject, resource1.Resource.Id),
// 				*web.NewCreditTransferItem(user3.Subject, user1.Subject, 6000),
// 			},
// 			Message: "Howdy :)",
// 		},
// 	})
// 	SubmitConfirmAcceptOffer(s.T(), ctx, user2, []*auth.UserSession{user1, user2, user3}, &web.SendOfferRequest{
// 		Offer: web.SendOfferPayload{
// 			Items: []web.SendOfferPayloadItem{
// 				*web.NewResourceTransferItemInputBase(user2.Subject, user3.Subject, resource2.Resource.Id),
// 				*web.NewCreditTransferItem(user3.Subject, user1.Subject, 6000),
// 			},
// 			Message: "Howdy :)",
// 		},
// 	})
// 	SubmitConfirmAcceptOffer(s.T(), ctx, user2, []*auth.UserSession{user1, user2, user3}, &web.SendOfferRequest{
// 		Offer: web.SendOfferPayload{
// 			Items: []web.SendOfferPayloadItem{
// 				*web.NewResourceTransferItemInputBase(user3.Subject, user1.Subject, resource1.Resource.Id),
// 				*web.NewCreditTransferItem(user3.Subject, user1.Subject, 6000),
// 			},
// 			Message: "Howdy :)",
// 		},
// 	})
//
// }

func (s *IntegrationTestSuite) UsersConfirmResourceTransferred(t *testing.T, ctx context.Context, offer *readmodels.OfferReadModel, users []*models.UserSession) error {

	for _, offerItem := range offer.OfferItems {
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
			httpResponse := s.ConfirmResourceTransfer(ctx, offerItemUserSession, offerItem.OfferItemKey)
			assert.Equal(s.T(), http.StatusOK, httpResponse.StatusCode)
		}
	}
	return nil
}

func (s *IntegrationTestSuite) UsersAcceptOffer(t *testing.T, ctx context.Context, offer *readmodels.OfferReadModel, users []*models.UserSession) error {

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

			httpResponse := s.AcceptOffer(ctx, offerItemUserSession, offer.OfferKey)
			assert.Equal(s.T(), http.StatusAccepted, httpResponse.StatusCode)

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
