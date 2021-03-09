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

func (s *IntegrationTestSuite) CreateGroup(t *testing.T, ctx context.Context, userSession *models.UserSession, request *handler.CreateGroupRequest) (*handler.GetGroupResponse, *http.Response) {
	req, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/groups", request)
	s.server.Router.ServeHTTP(recorder, req)
	response := &handler.GetGroupResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(s.T(), recorder, response)
}

func (s *IntegrationTestSuite) CreateOrAcceptInvitation(ctx context.Context, userSession *models.UserSession, request *handler.CreateOrAcceptInvitationRequest) *http.Response {
	req, recorder := NewRequest(ctx, userSession, http.MethodPost, "/api/v1/memberships", request)
	s.server.Router.ServeHTTP(recorder, req)
	return ReadResponse(s.T(), recorder, nil)
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

	groupOwner, delGroupOwner := s.testUser(s.T())
	defer delGroupOwner()

	ctx := context.Background()
	group, httpResponse := s.CreateGroup(s.T(), ctx, groupOwner, &handler.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if !AssertStatusCreated(s.T(), httpResponse) {
		return
	}
	groupKey := group.Group.GroupKey

	time.Sleep(1 * time.Second)

	s.T().Run("CreateGroup", func(t *testing.T) {

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

		membership, err := s.server.GetMembership.Get(ctx, keys.NewMembershipKey(groupKey, groupOwner.GetUserKey()))
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, groupOwner.Subject, membership.UserKey)
		assert.Equal(t, true, membership.IsOwner)
		assert.Equal(t, true, membership.UserConfirmed)
		assert.Equal(t, true, membership.GroupConfirmed)
		assert.Equal(t, true, membership.IsAdmin)
		assert.Equal(t, true, membership.IsMember)
	})

	s.T().Run("CreatingGroupShouldSubscribeOwnerToChanel", func(t *testing.T) {

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

		user, delUser := s.testUser(s.T())
		defer delUser()

		httpResponse = s.CreateOrAcceptInvitation(ctx, groupOwner, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user.GetUserKey(),
			GroupKey: groupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		membership, httpResponse := s.GetMembership(t, ctx, groupOwner, user.GetUserKey(), groupKey)
		if !AssertOK(t, httpResponse) {
			return
		}

		assert.Equal(t, user.Subject, membership.Membership.UserKey)
		assert.Equal(t, false, membership.Membership.IsAdmin)
		assert.Equal(t, false, membership.Membership.IsOwner)
		assert.Equal(t, false, membership.Membership.IsMember)
		assert.Equal(t, true, membership.Membership.GroupConfirmed)
		assert.Equal(t, false, membership.Membership.UserConfirmed)

	})

	s.T().Run("InviteeShouldBeAbleToAcceptInvitationFromOwner", func(t *testing.T) {

		ctx := context.Background()

		user, delUser := s.testUser(s.T())
		defer delUser()

		httpResponse = s.CreateOrAcceptInvitation(ctx, groupOwner, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user.GetUserKey(),
			GroupKey: groupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		httpResponse = s.CreateOrAcceptInvitation(ctx, user, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user.GetUserKey(),
			GroupKey: groupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		var membership *handler.GetMembershipResponse
		err := retry.Do(func() error {
			membership, httpResponse = s.GetMembership(t, ctx, user, user.GetUserKey(), groupKey)
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

		assert.Equal(t, user.Subject, membership.Membership.UserKey)
		assert.Equal(t, false, membership.Membership.IsAdmin)
		assert.Equal(t, false, membership.Membership.IsOwner)
		assert.Equal(t, true, membership.Membership.IsMember)
		assert.Equal(t, true, membership.Membership.GroupConfirmed)
		assert.Equal(t, true, membership.Membership.UserConfirmed)
	})

	s.T().Run("InviteeShouldBeAbleToDeclineInvitationFromOwner", func(t *testing.T) {

		ctx := context.Background()

		user, delUser := s.testUser(s.T())
		defer delUser()

		httpResponse = s.CreateOrAcceptInvitation(ctx, groupOwner, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user.GetUserKey(),
			GroupKey: groupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		httpResponse = s.DeclineOrCancelInvitation(t, ctx, user, &handler.CancelOrDeclineInvitationRequest{
			UserKey:  user.GetUserKey(),
			GroupKey: groupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		time.Sleep(1 * time.Second)

		_, err := s.server.GetMembership.Get(ctx, keys.NewMembershipKey(groupKey, user.GetUserKey()))
		assert.True(t, errors.Is(err, exceptions.ErrMembershipNotFound))

	})

	s.T().Run("OwnerShouldBeAbleToDeclineInvitationFromOwner", func(t *testing.T) {

		ctx := context.Background()

		user, delUser := s.testUser(s.T())
		defer delUser()

		httpResponse = s.CreateOrAcceptInvitation(ctx, user, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user.GetUserKey(),
			GroupKey: groupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		httpResponse = s.DeclineOrCancelInvitation(t, ctx, groupOwner, &handler.CancelOrDeclineInvitationRequest{
			UserKey:  user.GetUserKey(),
			GroupKey: groupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		time.Sleep(1 * time.Second)

		_, err := s.server.GetMembership.Get(ctx, keys.NewMembershipKey(groupKey, user.GetUserKey()))
		assert.True(t, errors.Is(err, exceptions.ErrMembershipNotFound))

	})

	s.T().Run("RandomUserShouldNotBeAbleToAcceptInvitation", func(t *testing.T) {

		ctx := context.Background()

		user, delUser := s.testUser(s.T())
		defer delUser()

		randomUser, delrandomUser := s.testUser(s.T())
		defer delrandomUser()

		httpResponse = s.CreateOrAcceptInvitation(ctx, groupOwner, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user.GetUserKey(),
			GroupKey: groupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		httpResponse = s.CreateOrAcceptInvitation(ctx, randomUser, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user.GetUserKey(),
			GroupKey: groupKey,
		})
		if !AssertStatusForbidden(t, httpResponse) {
			return
		}

	})

	s.T().Run("PersonShouldBeAbleToRequestBeingInvitedInGroup", func(t *testing.T) {

		ctx := context.Background()

		user, delUser := s.testUser(s.T())
		defer delUser()

		httpResponse = s.CreateOrAcceptInvitation(ctx, user, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user.GetUserKey(),
			GroupKey: groupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		membership, httpResponse := s.GetMembership(t, ctx, user, user.GetUserKey(), groupKey)
		if !AssertOK(t, httpResponse) {
			return
		}

		assert.Equal(t, user.Subject, membership.Membership.UserKey)
		assert.Equal(t, false, membership.Membership.IsAdmin)
		assert.Equal(t, false, membership.Membership.IsOwner)
		assert.Equal(t, false, membership.Membership.IsMember)
		assert.Equal(t, false, membership.Membership.GroupConfirmed)
		assert.Equal(t, true, membership.Membership.UserConfirmed)

	})

	s.T().Run("OwnerShouldBeAbleToAcceptInvitationRequest", func(t *testing.T) {

		ctx := context.Background()

		user, delUser := s.testUser(s.T())
		defer delUser()

		httpResponse = s.CreateOrAcceptInvitation(ctx, groupOwner, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user.GetUserKey(),
			GroupKey: groupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		httpResponse = s.CreateOrAcceptInvitation(ctx, user, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user.GetUserKey(),
			GroupKey: groupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		time.Sleep(1 * time.Second)

		membershipResponse, httpResponse := s.GetMembership(t, ctx, groupOwner, user.GetUserKey(), groupKey)
		if !AssertOK(t, httpResponse) {
			return
		}

		membership := membershipResponse.Membership
		assert.Equal(t, user.Subject, membership.UserKey)
		assert.Equal(t, false, membership.IsAdmin)
		assert.Equal(t, false, membership.IsOwner)
		assert.Equal(t, true, membership.IsMember)
		assert.Equal(t, true, membership.GroupConfirmed)
		assert.Equal(t, true, membership.UserConfirmed)
	})

	s.T().Run("GroupShouldReceiveMessageWhenUserJoined", func(t *testing.T) {

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

		httpResponse = s.CreateOrAcceptInvitation(ctx, user2, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  user2.GetUserKey(),
			GroupKey: grp.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			return
		}

		httpResponse = s.CreateOrAcceptInvitation(ctx, user1, &handler.CreateOrAcceptInvitationRequest{
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
		grp, httpResponse := s.CreateGroup(t, ctx, user3, &handler.CreateGroupRequest{
			Name:        "sample",
			Description: "description",
		})
		if !AssertStatusCreated(t, httpResponse) {
			return
		}

		resp, httpResponse := s.GetUsersForInvitePicker(t, ctx, grp.Group.GroupKey, 100, 0, user3)

		if !AssertOK(t, httpResponse) {
			return
		}

		seen := map[string]bool{}
		for _, user := range resp.Users {
			_, ok := seen[user.Username]
			if !assert.False(t, ok) {
				return
			}
			seen[user.Username] = true
		}
	})

	s.T().Run("GetLoggedInUserMembershipsWithoutGroup", func(t *testing.T) {

		ctx := context.Background()

		user, delUser := s.testUser(t)
		defer delUser()

		getMemberships, getMembershipsHttp := s.GetLoggedInUserMemberships(t, ctx, user)
		assert.Equal(t, http.StatusOK, getMembershipsHttp.StatusCode)
		assert.Equal(t, 0, len(getMemberships.Memberships))
	})

	s.T().Run("TestGetLoggedInUserMembershipsWithGroup", func(t *testing.T) {

		memberships, httpResponse := s.GetLoggedInUserMemberships(t, ctx, groupOwner)
		if !AssertOK(t, httpResponse) {
			return
		}

		if !assert.Len(t, memberships.Memberships, 1) {
			return
		}

		membership := memberships.Memberships[0]
		assert.Equal(t, groupKey, membership.GroupKey)
		assert.Equal(t, group.Group.Name, membership.GroupName)
		assert.Equal(t, group.Group.GroupKey, membership.GroupKey)
		assert.Equal(t, groupOwner.Subject, membership.UserKey)
		assert.Equal(t, groupOwner.Username, membership.UserName)
		assert.Equal(t, true, membership.UserConfirmed)
		assert.Equal(t, true, membership.GroupConfirmed)
		assert.Equal(t, true, membership.IsMember)
		assert.Equal(t, true, membership.IsOwner)

	})

}
