package integration

import (
	"context"
	"errors"
	"fmt"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/exceptions"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/web"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func CreateGroup(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *web.CreateGroupRequest) (*web.CreateGroupResponse, *http.Response, error) {
	c, recorder := NewRequest(ctx, userSession, http.MethodGet, "/api/v1/groups", request)
	err := a.CreateGroup(c)
	if err != nil {
		return nil, nil, err
	}
	response := &web.CreateGroupResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response), nil
}

func CreateOrAcceptInvitation(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *web.CreateOrAcceptInvitationRequest) (*web.CreateOrAcceptInvitationResponse, *http.Response, error) {
	c, recorder := NewRequest(ctx, userSession, http.MethodGet, "/api/v1/groups", request)
	err := a.CreateOrAcceptMembership(c)
	if err != nil {
		return nil, nil, err
	}
	response := &web.CreateOrAcceptInvitationResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response), nil
}

func DeclineOrCancelInvitation(t *testing.T, ctx context.Context, userSession *auth.UserSession, request *web.CancelOrDeclineInvitationRequest) (*web.CancelOrDeclineInvitationResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodGet, "/api/v1/groups", request)
	assert.NoError(t, a.CancelOrDeclineInvitation(c))
	response := &web.CancelOrDeclineInvitationResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response)
}

func GetUsersForInvitePicker(t *testing.T, ctx context.Context, groupId model.GroupKey, take int, skip int, userSession *auth.UserSession) (*web.GetUsersForGroupInvitePickerResponse, *http.Response) {
	groupIdStr := groupId.ID.String()
	c, recorder := NewRequest(ctx, userSession, http.MethodGet, fmt.Sprintf(`/api/v1/groups/%s/invite-member-picker?take=%d&skip=%d`, groupIdStr, take, skip), nil)
	c.SetParamNames("id")
	c.SetParamValues(groupIdStr)
	assert.NoError(t, a.GetUsersForGroupInvitePicker(c))
	response := &web.GetUsersForGroupInvitePickerResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response)
}

func GetLoggedInUserMemberships(t *testing.T, ctx context.Context, userSession *auth.UserSession) (*web.GetUserMembershipsResponse, *http.Response) {
	c, recorder := NewRequest(ctx, userSession, http.MethodGet, `/api/v1/my-memberships`, nil)
	assert.NoError(t, a.GetLoggedInUserMemberships(c))
	response := &web.GetUserMembershipsResponse{}
	t.Log(recorder.Body.String())
	return response, ReadResponse(t, recorder, response)

}

func TestCreateGroup(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()
	response, httpResponse, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusCreated, httpResponse.StatusCode)
	assert.NotNil(t, response.Group)
	assert.Equal(t, response.Group.Name, "sample")
	assert.Equal(t, response.Group.Description, "description")
	assert.NotEmpty(t, response.Group.ID)
	assert.NotEmpty(t, response.Group.CreatedAt)
}

func TestCreateGroupUnauthenticatedShouldFailWithUnauthorized(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	_, _, err := CreateGroup(t, ctx, nil, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err == nil {
		t.Fatal("err should not be nil")
	}
	assert.IsType(t, &exceptions.WebServiceException{}, err)
	assert.Equal(t, http.StatusUnauthorized, err.(*exceptions.WebServiceException).Status)
}

func TestCreateGroupEmptyName(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()
	_, httpResponse, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusBadRequest, httpResponse.StatusCode)
}

func TestCreateGroupEmptyDescriptionShouldNotFail(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()
	_, httpResponse, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "A Blibbers",
		Description: "",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusCreated, httpResponse.StatusCode)
}

func TestCreateGroupShouldCreateOwnerMembership(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()
	createGroup, createGroupHttp, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, createGroupHttp.StatusCode)

	gk, _ := model.ParseGroupKey(createGroup.Group.ID)
	grps, _ := GroupService.GetGroupMemberships(ctx, group2.NewGetMembershipsForGroupRequest(gk, nil))
	assert.Equal(t, 1, len(grps.Memberships.Items))
	assert.Equal(t, true, grps.Memberships.Items[0].IsOwner)
	assert.Equal(t, true, grps.Memberships.Items[0].UserConfirmed)
	assert.Equal(t, true, grps.Memberships.Items[0].GroupConfirmed)
	assert.Equal(t, false, grps.Memberships.Items[0].IsDeactivated)
	assert.Equal(t, true, grps.Memberships.Items[0].IsAdmin)
	assert.Equal(t, true, grps.Memberships.Items[0].IsMember)
	assert.Equal(t, user1.Subject, grps.Memberships.Items[0].Key.UserKey.String())
}

func TestCreatingGroupShouldSubscribeOwnerToChanel(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()

	amqpChan, err := AmqpClient.GetChannel()
	assert.NoError(t, err)
	defer amqpChan.Close()
	_, err = ChatService.CreateUserExchange(ctx, user1.GetUserKey())
	assert.NoError(t, err)
	err = amqpChan.QueueDeclare(ctx, "test", false, true, false, false, nil)
	assert.NoError(t, err)
	userKey := user1.GetUserKey()
	err = amqpChan.QueueBind(ctx, "test", "", userKey.GetExchangeName(), false, nil)
	assert.NoError(t, err)
	delivery, err := amqpChan.Consume(ctx, "test", "", false, false, false, false, nil)
	assert.NoError(t, err)

	_, _, err = CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}

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
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()
	grp, _, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}

	res, httpRes, err := CreateOrAcceptInvitation(t, ctx, user1, &web.CreateOrAcceptInvitationRequest{
		UserID:  user2.Subject,
		GroupID: grp.Group.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, httpRes.StatusCode)
	assert.Equal(t, user2.Subject, res.Membership.UserID)
	assert.Equal(t, false, res.Membership.IsAdmin)
	assert.Equal(t, false, res.Membership.IsOwner)
	assert.Equal(t, false, res.Membership.IsDeactivated)
	assert.Equal(t, false, res.Membership.IsMember)
	assert.Equal(t, true, res.Membership.GroupConfirmed)
	assert.Equal(t, false, res.Membership.UserConfirmed)

}

func TestInviteeShouldBeAbleToAcceptInvitationFromOwner(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()
	createGroup, createGroupHttp, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusCreated, createGroupHttp.StatusCode)

	res, httpRes, err := CreateOrAcceptInvitation(t, ctx, user1, &web.CreateOrAcceptInvitationRequest{
		UserID:  user2.Subject,
		GroupID: createGroup.Group.ID,
	})

	if err != nil {
		t.Fatal(err)
	}

	res, httpRes, err = CreateOrAcceptInvitation(t, ctx, user2, &web.CreateOrAcceptInvitationRequest{
		UserID:  user2.Subject,
		GroupID: createGroup.Group.ID,
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, httpRes.StatusCode)
	assert.Equal(t, user2.Subject, res.Membership.UserID)
	assert.Equal(t, false, res.Membership.IsAdmin)
	assert.Equal(t, false, res.Membership.IsOwner)
	assert.Equal(t, false, res.Membership.IsDeactivated)
	assert.Equal(t, true, res.Membership.IsMember)
	assert.Equal(t, true, res.Membership.GroupConfirmed)
	assert.Equal(t, true, res.Membership.UserConfirmed)
}

func TestInviteeShouldBeAbleToDeclineInvitationFromOwner(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()
	createGroup, createGroupHttp, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusCreated, createGroupHttp.StatusCode)

	_, acceptInvitationHttp, err := CreateOrAcceptInvitation(t, ctx, user1, &web.CreateOrAcceptInvitationRequest{
		UserID:  user2.Subject,
		GroupID: createGroup.Group.ID,
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, acceptInvitationHttp.StatusCode)

	_, declineInvitationHttp := DeclineOrCancelInvitation(t, ctx, user2, &web.CancelOrDeclineInvitationRequest{
		UserID:  user2.Subject,
		GroupID: createGroup.Group.ID,
	})
	assert.Equal(t, http.StatusAccepted, declineInvitationHttp.StatusCode)

	grpKey, _ := model.ParseGroupKey(createGroup.Group.ID)
	_, err = GroupStore.GetMembership(ctx, model.NewMembershipKey(grpKey, user2.GetUserKey()))
	assert.True(t, errors.Is(err, exceptions.ErrMembershipNotFound))
}

func TestOwnerShouldBeAbleToDeclineInvitationFromOwner(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()
	createGroup, createGroupHttp, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusCreated, createGroupHttp.StatusCode)
	_, _, err = CreateOrAcceptInvitation(t, ctx, user2, &web.CreateOrAcceptInvitationRequest{
		UserID:  user2.Subject,
		GroupID: createGroup.Group.ID,
	})

	if err != nil {
		t.Fatal(err)
	}

	_, httpRes := DeclineOrCancelInvitation(t, ctx, user1, &web.CancelOrDeclineInvitationRequest{
		UserID:  user2.Subject,
		GroupID: createGroup.Group.ID,
	})
	assert.Equal(t, http.StatusAccepted, httpRes.StatusCode)
	grpKey, _ := model.ParseGroupKey(createGroup.Group.ID)
	_, err = GroupStore.GetMembership(ctx, model.NewMembershipKey(grpKey, user2.GetUserKey()))
	assert.True(t, errors.Is(err, exceptions.ErrMembershipNotFound))
}

func TestRandomUserShouldNotBeAbleToAcceptInvitation(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	user3, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()
	grp, _, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = CreateOrAcceptInvitation(t, ctx, user1, &web.CreateOrAcceptInvitationRequest{
		UserID:  user2.Subject,
		GroupID: grp.Group.ID,
	})

	if err != nil {
		t.Fatal(err)
	}

	_, _, err = CreateOrAcceptInvitation(t, ctx, user3, &web.CreateOrAcceptInvitationRequest{
		UserID:  user2.Subject,
		GroupID: grp.Group.ID,
	})

	if err == nil {
		t.Fatal("err should not be nil")
	}

	assert.IsType(t, &exceptions.WebServiceException{}, err)
	assert.Equal(t, http.StatusForbidden, err.(*exceptions.WebServiceException).Status)
}

func TestPersonShouldBeAbleToRequestBeingInvitedInGroup(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()
	grp, _, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	res, httpRes, err := CreateOrAcceptInvitation(t, ctx, user2, &web.CreateOrAcceptInvitationRequest{
		UserID:  user2.Subject,
		GroupID: grp.Group.ID,
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, httpRes.StatusCode)
	assert.Equal(t, user2.Subject, res.Membership.UserID)
	assert.Equal(t, false, res.Membership.IsAdmin)
	assert.Equal(t, false, res.Membership.IsOwner)
	assert.Equal(t, false, res.Membership.IsDeactivated)
	assert.Equal(t, false, res.Membership.IsMember)
	assert.Equal(t, false, res.Membership.GroupConfirmed)
	assert.Equal(t, true, res.Membership.UserConfirmed)
}

func TestOwnerShouldBeAbleToAcceptInvitationRequest(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()
	grp, _, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	res, httpRes, err := CreateOrAcceptInvitation(t, ctx, user2, &web.CreateOrAcceptInvitationRequest{
		UserID:  user2.Subject,
		GroupID: grp.Group.ID,
	})

	if err != nil {
		t.Fatal(err)
	}

	res, httpRes, err = CreateOrAcceptInvitation(t, ctx, user1, &web.CreateOrAcceptInvitationRequest{
		UserID:  user2.Subject,
		GroupID: grp.Group.ID,
	})

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, httpRes.StatusCode)
	assert.Equal(t, user2.Subject, res.Membership.UserID)
	assert.Equal(t, false, res.Membership.IsAdmin)
	assert.Equal(t, false, res.Membership.IsOwner)
	assert.Equal(t, false, res.Membership.IsDeactivated)
	assert.Equal(t, true, res.Membership.IsMember)
	assert.Equal(t, true, res.Membership.GroupConfirmed)
	assert.Equal(t, true, res.Membership.UserConfirmed)
}

func TestGroupShouldReceiveMessageWhenUserJoined(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()

	x1 := ListenOnUserExchange(t, ctx, user1.GetUserKey())
	defer x1.Close()
	x2 := ListenOnUserExchange(t, ctx, user2.GetUserKey())
	defer x2.Close()

	grp, _, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = CreateOrAcceptInvitation(t, ctx, user2, &web.CreateOrAcceptInvitationRequest{
		UserID:  user2.Subject,
		GroupID: grp.Group.ID,
	})

	if err != nil {
		t.Fatal(err)
	}

	_, _, err = CreateOrAcceptInvitation(t, ctx, user1, &web.CreateOrAcceptInvitationRequest{
		UserID:  user2.Subject,
		GroupID: grp.Group.ID,
	})

	if err != nil {
		t.Fatal(err)
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

}

func TestGetUsersForInvitePickerShouldNotReturnDuplicates(t *testing.T) {
	cleanDb()

	user1, delUser1 := testUser(t)
	defer delUser1()

	user2, delUser2 := testUser(t)
	defer delUser2()

	user3, delUser2 := testUser(t)
	defer delUser2()

	ctx := context.Background()

	_, _, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = CreateGroup(t, ctx, user2, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = CreateGroup(t, ctx, user2, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = CreateGroup(t, ctx, user2, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = CreateGroup(t, ctx, user3, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = CreateGroup(t, ctx, user3, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}
	grp, _, err := CreateGroup(t, ctx, user3, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}

	grpKey, _ := model.ParseGroupKey(grp.Group.ID)
	resp, httpResp := GetUsersForInvitePicker(t, ctx, grpKey, 100, 0, user3)

	assert.Equal(t, http.StatusOK, httpResp.StatusCode)
	assert.Equal(t, 2, len(resp.Users))

}

func TestGetLoggedInUserMembershipsWithoutGroup(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()
	getMemberships, getMembershipsHttp := GetLoggedInUserMemberships(t, ctx, user1)
	assert.Equal(t, http.StatusOK, getMembershipsHttp.StatusCode)
	assert.Equal(t, 0, len(getMemberships.Memberships))
}

func TestGetLoggedInUserMembershipsWithGroup(t *testing.T) {
	t.Parallel()

	user1, delUser1 := testUser(t)
	defer delUser1()

	ctx := context.Background()
	createGroup, _, err := CreateGroup(t, ctx, user1, &web.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	})
	if err != nil {
		t.Fatal(err)
	}

	getMemberships, getMembershipsHttp := GetLoggedInUserMemberships(t, ctx, user1)

	assert.Equal(t, http.StatusOK, getMembershipsHttp.StatusCode)
	assert.Equal(t, 1, len(getMemberships.Memberships))
	assert.Equal(t, createGroup.Group.ID, getMemberships.Memberships[0].GroupID)
	assert.Equal(t, "sample", getMemberships.Memberships[0].GroupName)
	assert.Equal(t, user1.Subject, getMemberships.Memberships[0].UserID)
	assert.Equal(t, user1.Username, getMemberships.Memberships[0].UserName)
	assert.Equal(t, true, getMemberships.Memberships[0].UserConfirmed)
	assert.Equal(t, true, getMemberships.Memberships[0].GroupConfirmed)
	assert.Equal(t, true, getMemberships.Memberships[0].IsMember)
	assert.Equal(t, true, getMemberships.Memberships[0].IsOwner)
}
