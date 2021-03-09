package integration

import (
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/exceptions"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func (s *IntegrationTestSuite) CreateGroup(t *testing.T, ctx context.Context, userSession *models.UserSession, request *handler.CreateGroupRequest) (*handler.CreateGroupResponse, *http.Response) {
	req, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/groups", request)
	s.server.Router.ServeHTTP(recorder, req)
	response := &handler.CreateGroupResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) CreateOrAcceptInvitation(t *testing.T, ctx context.Context, userSession *models.UserSession, request *handler.CreateOrAcceptInvitationRequest) (*handler.CreateOrAcceptInvitationResponse, *http.Response) {
	req, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/memberships", request)
	s.server.Router.ServeHTTP(recorder, req)
	response := &handler.CreateOrAcceptInvitationResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) DeclineOrCancelInvitation(t *testing.T, ctx context.Context, userSession *models.UserSession, request *handler.CancelOrDeclineInvitationRequest) *http.Response {
	req, recorder := NewRequest(ctx, userSession, http.MethodDelete, "/api/v1/memberships", request)
	s.server.Router.ServeHTTP(recorder, req)
	response := group2.CancelOrDeclineInvitationRequest{}
	t.Log(recorder.Body.String())
	return ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) GetUsersForInvitePicker(t *testing.T, ctx context.Context, groupId keys.GroupKey, take int, skip int, userSession *models.UserSession) (*handler.GetUsersForGroupInvitePickerResponse, *http.Response) {
	groupIdStr := groupId.ID.String()
	req, recorder := NewRequest(ctx, userSession, http.MethodGet, fmt.Sprintf(`/api/v1/groups/%s/invite-member-picker?take=%d&skip=%d`, groupIdStr, take, skip), nil)
	s.server.Router.ServeHTTP(recorder, req)
	response := &handler.GetUsersForGroupInvitePickerResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) GetLoggedInUserMemberships(t *testing.T, ctx context.Context, userSession *models.UserSession) (*handler.GetMembershipsResponse, *http.Response) {
	req, recorder := NewRequest(ctx, userSession, http.MethodGet, `/api/v1/memberships`, nil)
	s.server.Router.ServeHTTP(recorder, req)
	response := &handler.GetMembershipsResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) GetMembership(t *testing.T, ctx context.Context, userSession *models.UserSession, userKey keys.UserKey, groupKey keys.GroupKey) (*handler.GetMembershipResponse, *http.Response) {
	req, recorder := NewRequest(ctx, userSession, http.MethodGet, fmt.Sprintf(`/api/v1/groups/%s/memberships/%s`, groupKey.String(), userKey.String()), nil)
	s.server.Router.ServeHTTP(recorder, req)
	response := &handler.GetMembershipResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) TestGroups() {

	user1, delUser1 := s.testUser(s.T())
	defer delUser1()

	user2, delUser2 := s.testUser(s.T())
	defer delUser2()

	user3, delUser2 := s.testUser(s.T())
	defer delUser2()

	s.T().Run("CreateGroup", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		response, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}
		assert.NotNil(t, response.Group)
		assert.Equal(t, response.Group.Name, "sample")
		assert.Equal(t, response.Group.Description, "description")
		assert.NotEmpty(t, response.Group.GroupKey)
		assert.NotEmpty(t, response.Group.CreatedAt)

	})

	s.T().Run("CreateGroupUnauthenticatedShouldFailWithUnauthorized", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, httpResponse := s.CreateGroup(t, ctx, nil, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusUnauthorized(t, httpResponse) {
			return
		}
	})

	s.T().Run("CreateGroupEmptyNameShouldFail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "",
			Description: "description",
		})

		if !AssertStatusBadRequest(t, httpResponse) {
			return
		}

	})

	s.T().Run("CreateGroupEmptyDescriptionShouldNotFail", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		_, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "A Blibbers",
			Description: "",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}

	})

	s.T().Run("CreateGroupShouldCreateOwnerMembership", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		createGroup, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})

		if !AssertStatusCreated(t, httpResponse) {
			return
		}

		groupKey := createGroup.Group.GroupKey
		grps, _ := s.server.Group.Service.GetGroupMemberships(ctx, group2.NewGetMembershipsForGroupRequest(groupKey, nil))
		if !assert.Len(t, grps.Memberships, 1) {
			return
		}
		assert.Equal(t, true, grps.Memberships[0].IsOwner)
		assert.Equal(t, true, grps.Memberships[0].UserConfirmed)
		assert.Equal(t, true, grps.Memberships[0].GroupConfirmed)
		assert.Equal(t, true, grps.Memberships[0].IsAdmin)
		assert.Equal(t, true, grps.Memberships[0].IsMember)
		assert.Equal(t, user1.Subject, grps.Memberships[0].UserKey)

	})

	s.T().Run("CreatingGroupShouldSubscribeOwnerToChanel", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		amqpChan, err := s.server.AmqpClient.GetChannel()
		assert.NoError(t, err)
		defer amqpChan.Close()
		_, err = s.server.ChatService.CreateUserExchange(ctx, user1.GetUserKey())
		assert.NoError(t, err)
		err = amqpChan.QueueDeclare(ctx, "test", false, true, false, false, nil)
		assert.NoError(t, err)
		userKey := user1.GetUserKey()
		err = amqpChan.QueueBind(ctx, "test", "", userKey.GetExchangeName(), false, nil)
		assert.NoError(t, err)
		delivery, err := amqpChan.Consume(ctx, "test", "", false, false, false, false, nil)
		assert.NoError(t, err)

		_, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}

		select {
		case msg := <-delivery:
			fmt.Println("received message!")
			fmt.Println(string(msg.Body))
			return
		case <-time.After(1 * time.Second):
			t.FailNow()
		}

	})

	s.T().Run("OwnerShouldBeAbleToInviteUser", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		grp, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}

		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, user1, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: grp.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		membership, httpResponse := s.GetMembership(t, ctx, user1, user2.GetUserKey(), grp.Group.GroupKey)
		if !AssertOK(t, httpResponse) {
			return
		}

		assert.Equal(t, user2.Subject, membership.Membership.UserID)
		assert.Equal(t, false, membership.Membership.IsAdmin)
		assert.Equal(t, false, membership.Membership.IsOwner)
		assert.Equal(t, false, membership.Membership.IsDeactivated)
		assert.Equal(t, false, membership.Membership.IsMember)
		assert.Equal(t, true, membership.Membership.GroupConfirmed)
		assert.Equal(t, false, membership.Membership.UserConfirmed)

	})

	s.T().Run("InviteeShouldBeAbleToAcceptInvitationFromOwner", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		createGroup, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}

		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, user1, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: createGroup.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, user2, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: createGroup.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		var membership *handler.GetMembershipResponse
		err := retry.Do(func() error {
			membership, httpResponse = s.GetMembership(t, ctx, user2, user2.GetUserKey(), createGroup.Group.GroupKey)
			if httpResponse.StatusCode != http.StatusOK {
				return fmt.Errorf("invalid status code")
			}
			if !membership.Membership.UserConfirmed {
				return fmt.Errorf("retrying")
			}
			return nil
		}, retry.Attempts(10), retry.MaxDelay(20*time.Millisecond))

		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, user2.Subject, membership.Membership.UserID)
		assert.Equal(t, false, membership.Membership.IsAdmin)
		assert.Equal(t, false, membership.Membership.IsOwner)
		assert.Equal(t, false, membership.Membership.IsDeactivated)
		assert.Equal(t, true, membership.Membership.IsMember)
		assert.Equal(t, true, membership.Membership.GroupConfirmed)
		assert.Equal(t, true, membership.Membership.UserConfirmed)
	})

	s.T().Run("InviteeShouldBeAbleToDeclineInvitationFromOwner", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		createGroup, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}
		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, user1, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: createGroup.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		httpResponse = s.DeclineOrCancelInvitation(t, ctx, user2, &handler.CancelOrDeclineInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: createGroup.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		_, err := s.server.Group.Store.GetMembership(ctx, keys.NewMembershipKey(createGroup.Group.GroupKey, user2.GetUserKey()))
		assert.True(t, errors.Is(err, exceptions.ErrMembershipNotFound))

	})

	s.T().Run("OwnerShouldBeAbleToDeclineInvitationFromOwner", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		createGroup, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}

		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, user2, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: createGroup.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		httpResponse = s.DeclineOrCancelInvitation(t, ctx, user1, &handler.CancelOrDeclineInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: createGroup.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		_, err := s.server.Group.Store.GetMembership(ctx, keys.NewMembershipKey(createGroup.Group.GroupKey, user2.GetUserKey()))
		assert.True(t, errors.Is(err, exceptions.ErrMembershipNotFound))

	})

	s.T().Run("RandomUserShouldNotBeAbleToAcceptInvitation", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		grp, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}

		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, user1, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: grp.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, user3, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: grp.Group.GroupKey,
		})
		if !AssertStatusForbidden(t, httpResponse) {
			return
		}

	})

	s.T().Run("PersonShouldBeAbleToRequestBeingInvitedInGroup", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		grp, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}

		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, user2, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: grp.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		membership, httpResponse := s.GetMembership(t, ctx, user2, user2.GetUserKey(), grp.Group.GroupKey)
		if !AssertOK(t, httpResponse) {
			return
		}

		assert.Equal(t, user2.Subject, membership.Membership.UserID)
		assert.Equal(t, false, membership.Membership.IsAdmin)
		assert.Equal(t, false, membership.Membership.IsOwner)
		assert.Equal(t, false, membership.Membership.IsDeactivated)
		assert.Equal(t, false, membership.Membership.IsMember)
		assert.Equal(t, false, membership.Membership.GroupConfirmed)
		assert.Equal(t, true, membership.Membership.UserConfirmed)

	})

	s.T().Run("OwnerShouldBeAbleToAcceptInvitationRequest", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		grp, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}

		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, user2, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: grp.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, user1, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: grp.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		membership, httpResponse := s.GetMembership(t, ctx, user2, user2.GetUserKey(), grp.Group.GroupKey)
		if !AssertOK(t, httpResponse) {
			return
		}

		assert.Equal(t, user2.Subject, membership.Membership.UserID)
		assert.Equal(t, false, membership.Membership.IsAdmin)
		assert.Equal(t, false, membership.Membership.IsOwner)
		assert.Equal(t, false, membership.Membership.IsDeactivated)
		assert.Equal(t, true, membership.Membership.IsMember)
		assert.Equal(t, true, membership.Membership.GroupConfirmed)
		assert.Equal(t, true, membership.Membership.UserConfirmed)
	})

	s.T().Run("GroupShouldReceiveMessageWhenUserJoined", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		x1 := s.ListenOnUserExchange(t, ctx, user1.GetUserKey())
		defer x1.Close()
		x2 := s.ListenOnUserExchange(t, ctx, user2.GetUserKey())
		defer x2.Close()

		grp, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}

		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, user2, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: grp.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, user1, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: grp.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		x1done := make(chan bool)
		x2done := make(chan bool)
		fail := make(chan bool)

		go func() {
			i1 := 0
			for delivery := range x1.Delivery {
				i1++
				fmt.Println("received message on queue 1")
				fmt.Println("Message type: " + delivery.Type)
				fmt.Println("Message body: " + string(delivery.Body))
				if i1 == 2 {
					x1done <- true
				}
			}
		}()

		go func() {
			i1 := 0
			for delivery := range x2.Delivery {
				i1++
				fmt.Println("received message on queue 2")
				fmt.Println("Message type: " + delivery.Type)
				fmt.Println("Message body: " + string(delivery.Body))
				if i1 == 1 {
					x2done <- true
				}
			}
		}()

		go func() {
			time.Sleep(time.Second * 10)
			fail <- true
		}()

		i1done := false
		i2done := false
		for {
			select {
			case _ = <-x1done:
				i1done = true
				if i2done {
					return
				}
				break
			case _ = <-x2done:
				i2done = true
				if i1done {
					return
				}
				break
			case _ = <-fail:
				t.FailNow()
			}

		}

	})

	s.T().Run("GetUsersForInvitePickerShouldNotReturnDuplicates", func(t *testing.T) {
		s.cleanDb()
		ctx := context.Background()

		_, httpResponse := s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}
		_, httpResponse = s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}
		_, httpResponse = s.CreateGroup(t, ctx, user1, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}
		_, httpResponse = s.CreateGroup(t, ctx, user2, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}
		_, httpResponse = s.CreateGroup(t, ctx, user2, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}
		_, httpResponse = s.CreateGroup(t, ctx, user2, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}
		_, httpResponse = s.CreateGroup(t, ctx, user3, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}
		_, httpResponse = s.CreateGroup(t, ctx, user3, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}
		grp, httpResponse := s.CreateGroup(t, ctx, user3, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}

		resp, httpResp := s.GetUsersForInvitePicker(t, ctx, grp.Group.GroupKey, 100, 0, user3)

		assert.Equal(t, http.StatusOK, httpResp.StatusCode)
		assert.Equal(t, 2, len(resp.Users))

	})

	s.T().Run("GetLoggedInUserMembershipsWithoutGroup", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		getMemberships, getMembershipsHttp := s.GetLoggedInUserMemberships(t, ctx, user1)
		assert.Equal(t, http.StatusOK, getMembershipsHttp.StatusCode)
		assert.Equal(t, 0, len(getMemberships.Memberships))
	})

	s.T().Run("TestGetLoggedInUserMembershipsWithGroup", func(t *testing.T) {
		t.Parallel()
		user, delUser := s.testUser(t)
		defer delUser()

		ctx := context.Background()
		createGroup, httpResponse := s.CreateGroup(t, ctx, user, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}

		getMemberships, httpResponse := s.GetLoggedInUserMemberships(t, ctx, user)
		if !AssertOK(t, httpResponse) {
			return
		}

		if !assert.Len(t, getMemberships.Memberships, 1) {
			return
		}
		assert.Equal(t, createGroup.Group.GroupKey, getMemberships.Memberships[0].GroupID)
		assert.Equal(t, "sample", getMemberships.Memberships[0].GroupName)
		assert.Equal(t, user.Subject, getMemberships.Memberships[0].UserID)
		assert.Equal(t, user.Username, getMemberships.Memberships[0].UserName)
		assert.Equal(t, true, getMemberships.Memberships[0].UserConfirmed)
		assert.Equal(t, true, getMemberships.Memberships[0].GroupConfirmed)
		assert.Equal(t, true, getMemberships.Memberships[0].IsMember)
		assert.Equal(t, true, getMemberships.Memberships[0].IsOwner)

	})

}
