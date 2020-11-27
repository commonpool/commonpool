package integration

import (
	"context"
	"errors"
	"fmt"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/web"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func CreateGroup(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *web.CreateGroupRequest) (*web.CreateGroupResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodGet, "/api/v1/groups", request)
	PanicIfError(a.CreateGroup(c))
	response := &web.CreateGroupResponse{}
	return response, ReadResponse(t, recorder, response)
}

func CreateOrAcceptInvitation(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *web.CreateOrAcceptInvitationRequest) (*web.CreateOrAcceptInvitationResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodGet, "/api/v1/groups", request)
	PanicIfError(a.CreateOrAcceptMembership(c))
	response := &web.CreateOrAcceptInvitationResponse{}
	return response, ReadResponse(t, recorder, response)
}

func DeclineOrCancelInvitation(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *web.CancelOrDeclineInvitationRequest) (*web.CancelOrDeclineInvitationResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodGet, "/api/v1/groups", request)
	PanicIfError(a.CancelOrDeclineInvitation(c))
	response := &web.CancelOrDeclineInvitationResponse{}
	return response, ReadResponse(t, recorder, response)
}

func GetUsersForInvitePicker(t *testing.T, ctx context.Context, groupId model.GroupKey, take int, skip int, userSession *auth.UserSession) (*web.GetUsersForGroupInvitePickerResponse, *http.Response) {
	groupIdStr := groupId.ID.String()
	c, recorder := NewRequest(ctx, userSession, http.MethodGet, fmt.Sprintf(`/api/v1/groups/%s/invite-member-picker?take=%d&skip=%d`, groupIdStr, take, skip), nil)
	c.SetParamNames("id")
	c.SetParamValues(groupIdStr)
	PanicIfError(a.GetUsersForGroupInvitePicker(c))
	response := &web.GetUsersForGroupInvitePickerResponse{}
	return response, ReadResponse(t, recorder, response)
}

func TestCreateGroup(t *testing.T) {
	ctx := context.Background()
	response, httpResponse := CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	assert.Equal(t, http.StatusCreated, httpResponse.StatusCode)
	assert.NotNil(t, response.Group)
	assert.Equal(t, response.Group.Name, "sample")
	assert.Equal(t, response.Group.Description, "description")
	assert.NotEmpty(t, response.Group.ID)
	assert.NotEmpty(t, response.Group.CreatedAt)
}

func TestCreateGroupUnauthenticatedShouldFailWithUnauthorized(t *testing.T) {
	ctx := context.Background()
	_, httpResponse := CreateGroup(t, ctx, nil, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	assert.Equal(t, http.StatusUnauthorized, httpResponse.StatusCode)
}

func TestCreateGroupEmptyName(t *testing.T) {
	ctx := context.Background()
	_, httpResponse := CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "",
		Description: "description",
	})
	assert.Equal(t, http.StatusBadRequest, httpResponse.StatusCode)
}

func TestCreateGroupEmptyDescriptionShouldNotFail(t *testing.T) {
	ctx := context.Background()
	_, httpResponse := CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "A Blibbers",
		Description: "",
	})
	assert.Equal(t, http.StatusCreated, httpResponse.StatusCode)
}

func TestCreateGroupShouldCreateOwnerMembership(t *testing.T) {
	ctx := context.Background()
	response, _ := CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	gk, _ := group.ParseGroupKey(response.Group.ID)
	grps, _ := GroupService.GetGroupsMemberships(ctx, group.NewGetMembershipsForGroupRequest(gk, nil))
	assert.Equal(t, 1, len(grps.Memberships.Items))
	assert.Equal(t, true, grps.Memberships.Items[0].IsOwner)
	assert.Equal(t, true, grps.Memberships.Items[0].UserConfirmed)
	assert.Equal(t, true, grps.Memberships.Items[0].GroupConfirmed)
	assert.Equal(t, false, grps.Memberships.Items[0].IsDeactivated)
	assert.Equal(t, true, grps.Memberships.Items[0].IsAdmin)
	assert.Equal(t, true, grps.Memberships.Items[0].IsMember)
	assert.Equal(t, User1.Subject, grps.Memberships.Items[0].UserID)
}

func TestCreatingGroupShouldSubscribeOwnerToChanel(t *testing.T) {
	ctx := context.Background()

	amqpChan, err := AmqpClient.GetChannel()
	assert.NoError(t, err)
	defer amqpChan.Close()
	_, err = ChatService.CreateUserExchange(ctx, User1.GetUserKey())
	assert.NoError(t, err)
	err = amqpChan.QueueDeclare(ctx, "test", false, true, false, false, nil)
	assert.NoError(t, err)
	userKey := User1.GetUserKey()
	err = amqpChan.QueueBind(ctx, "test", "", userKey.GetExchangeName(), false, nil)
	assert.NoError(t, err)
	delivery, err := amqpChan.Consume(ctx, "test", "", false, false, false, false, nil)
	assert.NoError(t, err)

	_, _ = CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})

	select {
	case msg := <-delivery:
		fmt.Println("received message!")
		fmt.Println(string(msg.Body))
		return
	case <-time.After(1 * time.Second):
		t.FailNow()
	}

}

func TestOwnerShouldBeAbleToInviteUser(t *testing.T) {

	ctx := context.Background()
	grp, _ := CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	res, httpRes := CreateOrAcceptInvitation(t, ctx, User1, &web.CreateOrAcceptInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)
	assert.Equal(t, User2.Subject, res.Membership.UserID)
	assert.Equal(t, false, res.Membership.IsAdmin)
	assert.Equal(t, false, res.Membership.IsOwner)
	assert.Equal(t, false, res.Membership.IsDeactivated)
	assert.Equal(t, false, res.Membership.IsMember)
	assert.Equal(t, true, res.Membership.GroupConfirmed)
	assert.Equal(t, false, res.Membership.UserConfirmed)

}

func TestInviteeShouldBeAbleToAcceptInvitationFromOwner(t *testing.T) {
	ctx := context.Background()
	grp, _ := CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	res, httpRes := CreateOrAcceptInvitation(t, ctx, User1, &web.CreateOrAcceptInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	res, httpRes = CreateOrAcceptInvitation(t, ctx, User2, &web.CreateOrAcceptInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)
	assert.Equal(t, User2.Subject, res.Membership.UserID)
	assert.Equal(t, false, res.Membership.IsAdmin)
	assert.Equal(t, false, res.Membership.IsOwner)
	assert.Equal(t, false, res.Membership.IsDeactivated)
	assert.Equal(t, true, res.Membership.IsMember)
	assert.Equal(t, true, res.Membership.GroupConfirmed)
	assert.Equal(t, true, res.Membership.UserConfirmed)
}

func TestInviteeShouldBeAbleToDeclineInvitationFromOwner(t *testing.T) {
	ctx := context.Background()
	grp, _ := CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	_, _ = CreateOrAcceptInvitation(t, ctx, User1, &web.CreateOrAcceptInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	_, httpRes := DeclineOrCancelInvitation(t, ctx, User2, &web.CancelOrDeclineInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	assert.Equal(t, http.StatusAccepted, httpRes.StatusCode)
	grpKey, _ := group.ParseGroupKey(grp.Group.ID)
	_, err := GroupStore.GetMembership(ctx, model.NewMembershipKey(grpKey, User2.GetUserKey()))
	assert.True(t, errors.Is(err, group.ErrMembershipNotFound))
}

func TestOwnerShouldBeAbleToDeclineInvitationFromOwner(t *testing.T) {
	ctx := context.Background()
	grp, _ := CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	_, _ = CreateOrAcceptInvitation(t, ctx, User2, &web.CreateOrAcceptInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	_, httpRes := DeclineOrCancelInvitation(t, ctx, User1, &web.CancelOrDeclineInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	assert.Equal(t, http.StatusAccepted, httpRes.StatusCode)
	grpKey, _ := group.ParseGroupKey(grp.Group.ID)
	_, err := GroupStore.GetMembership(ctx, model.NewMembershipKey(grpKey, User2.GetUserKey()))
	assert.True(t, errors.Is(err, group.ErrMembershipNotFound))
}

func TestRandomUserShouldNotBeAbleToAcceptInvitation(t *testing.T) {
	ctx := context.Background()
	grp, _ := CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	_, _ = CreateOrAcceptInvitation(t, ctx, User1, &web.CreateOrAcceptInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	_, httpRes := CreateOrAcceptInvitation(t, ctx, User3, &web.CreateOrAcceptInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	assert.Equal(t, http.StatusForbidden, httpRes.StatusCode)
}

func TestPersonShouldBeAbleToRequestBeingInvitedInGroup(t *testing.T) {
	ctx := context.Background()
	grp, _ := CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	res, httpRes := CreateOrAcceptInvitation(t, ctx, User2, &web.CreateOrAcceptInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)
	assert.Equal(t, User2.Subject, res.Membership.UserID)
	assert.Equal(t, false, res.Membership.IsAdmin)
	assert.Equal(t, false, res.Membership.IsOwner)
	assert.Equal(t, false, res.Membership.IsDeactivated)
	assert.Equal(t, false, res.Membership.IsMember)
	assert.Equal(t, false, res.Membership.GroupConfirmed)
	assert.Equal(t, true, res.Membership.UserConfirmed)
}

func TestOwnerShouldBeAbleToAcceptInvitationRequest(t *testing.T) {
	ctx := context.Background()
	grp, _ := CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	res, httpRes := CreateOrAcceptInvitation(t, ctx, User2, &web.CreateOrAcceptInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	res, httpRes = CreateOrAcceptInvitation(t, ctx, User1, &web.CreateOrAcceptInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)
	assert.Equal(t, User2.Subject, res.Membership.UserID)
	assert.Equal(t, false, res.Membership.IsAdmin)
	assert.Equal(t, false, res.Membership.IsOwner)
	assert.Equal(t, false, res.Membership.IsDeactivated)
	assert.Equal(t, true, res.Membership.IsMember)
	assert.Equal(t, true, res.Membership.GroupConfirmed)
	assert.Equal(t, true, res.Membership.UserConfirmed)
}

func TestGroupShouldReceiveMessageWhenUserJoined(t *testing.T) {
	ctx := context.Background()

	x1 := ListenOnUserExchange(t, ctx, User1.GetUserKey())
	defer x1.Close()
	x2 := ListenOnUserExchange(t, ctx, User2.GetUserKey())
	defer x2.Close()

	grp, _ := CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	_, _ = CreateOrAcceptInvitation(t, ctx, User2, &web.CreateOrAcceptInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})
	_, _ = CreateOrAcceptInvitation(t, ctx, User1, &web.CreateOrAcceptInvitationRequest{
		UserID:  User2.Subject,
		GroupID: grp.Group.ID,
	})

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
		time.Sleep(time.Second * 5)
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

}

func TestGetUsersForInvitePickerShouldNotReturnDuplicates(t *testing.T) {
	ctx := context.Background()

	_, _ = CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	_, _ = CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	_, _ = CreateGroup(t, ctx, User1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	_, _ = CreateGroup(t, ctx, User2, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	_, _ = CreateGroup(t, ctx, User2, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	_, _ = CreateGroup(t, ctx, User2, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	_, _ = CreateGroup(t, ctx, User3, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	_, _ = CreateGroup(t, ctx, User3, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	grp, _ := CreateGroup(t, ctx, User3, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})

	grpKey, _ := group.ParseGroupKey(grp.Group.ID)
	resp, httpResp := GetUsersForInvitePicker(t, ctx, grpKey, 100, 0, User3)

	assert.Equal(t, http.StatusOK, httpResp.StatusCode)
	assert.Equal(t, 2, len(resp.Users))
	assert.Equal(t, "user1", resp.Users[0].Username)
	assert.Equal(t, "user2", resp.Users[1].Username)

}
