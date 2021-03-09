package integration

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/client"
	"github.com/commonpool/backend/pkg/client/echo"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/server"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"os"
	"testing"
	"time"
)

var srv *server.Server

func TestMain(m *testing.M) {
	var err error
	srv, err = server.NewServer()
	cleanDb()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}

func NewClient(user *models.UserSession) client.Client {
	return echo.NewEchoClient(srv.Router, client.NewMockAuthentication(user))
}

func TestCreateGroup(t *testing.T) {

	ctx := context.Background()

	var _, cli = testUserCli(t)

	var response = &handler.GetGroupResponse{}
	if !assert.NoError(t, cli.CreateGroup(ctx, handler.NewCreateGroupRequest("sample", "description"), response)) {
		return
	}

	assert.NotNil(t, response.Group)
	assert.Equal(t, response.Group.Name, "sample")
	assert.Equal(t, response.Group.Description, "description")
	assert.NotEmpty(t, response.Group.GroupKey)
	assert.NotEmpty(t, response.Group.CreatedAt)
}

func TestCreateGroupUnauthenticatedShouldFailWithUnauthorized(t *testing.T) {

	err := NewClient(nil).CreateGroup(context.Background(), handler.NewCreateGroupRequest("sample", "description"), &handler.GetGroupResponse{})
	if !assert.Error(t, err) {
		return
	}
}

func TestCreateGroupEmptyNameShouldFail(t *testing.T) {

	ctx := context.Background()
	var _, cli = testUserCli(t)
	assert.Error(t, cli.CreateGroup(ctx, handler.NewCreateGroupRequest("", "description"), &handler.GetGroupResponse{}))
	assert.Error(t, cli.CreateGroup(ctx, handler.NewCreateGroupRequest("   ", "description"), &handler.GetGroupResponse{}))
}

func TestCreateGroupEmptyDescriptionShouldNotFail(t *testing.T) {

	ctx := context.Background()
	var _, cli = testUserCli(t)
	assert.NoError(t, cli.CreateGroup(ctx, handler.NewCreateGroupRequest("name", ""), &handler.GetGroupResponse{}))
}

func TestCreateGroupShouldCreateOwnerMembership(t *testing.T) {

	ctx := context.TODO()

	owner, _ := testUserCli(t)
	group := &handler.GetGroupResponse{}
	if !assert.NoError(t, testGroup2(t, owner, group)) {
		return
	}

	time.Sleep(time.Second)

	membership, err := srv.GetMembership.Get(ctx, keys.NewMembershipKey(group, owner))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, owner.GetUserKey(), membership.UserKey)
	assert.Equal(t, true, membership.IsOwner)
	assert.Equal(t, true, membership.UserConfirmed)
	assert.Equal(t, true, membership.GroupConfirmed)
	assert.Equal(t, true, membership.IsAdmin)
	assert.Equal(t, true, membership.IsMember)
}

func TestOwnerShouldBeAbleToInviteUser(t *testing.T) {

	ctx := context.TODO()

	owner, ownerCli := testUserCli(t)
	group := &handler.GetGroupResponse{}
	if !assert.NoError(t, testGroup2(t, owner, group)) {
		return
	}

	user, _ := testUser(t)
	membershipKey := keys.NewMembershipKey(group, user)

	if !assert.NoError(t, ownerCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	var membership handler.GetMembershipResponse
	if !assert.NoError(t, ownerCli.GetMembership(ctx, membershipKey, &membership)) {
		return
	}

	assert.Equal(t, user.GetUserKey(), membership.Membership.UserKey)
	assert.Equal(t, false, membership.Membership.IsAdmin)
	assert.Equal(t, false, membership.Membership.IsOwner)
	assert.Equal(t, false, membership.Membership.IsMember)
	assert.Equal(t, true, membership.Membership.GroupConfirmed)
	assert.Equal(t, false, membership.Membership.UserConfirmed)
}

func TestInviteeShouldBeAbleToAcceptInvitationFromOwner(t *testing.T) {

	ctx := context.Background()

	owner, ownerCli := testUserCli(t)
	group := &handler.GetGroupResponse{}
	if !assert.NoError(t, testGroup2(t, owner, group)) {
		return
	}

	user, _ := testUser(t)
	userCli := NewClient(user)
	membershipKey := keys.NewMembershipKey(group, user)

	if !assert.NoError(t, ownerCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	if !assert.NoError(t, userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	time.Sleep(time.Second)

	var membership handler.GetMembershipResponse
	if !assert.NoError(t, userCli.GetMembership(ctx, membershipKey, &membership)) {
		return
	}

	assert.Equal(t, user.GetUserKey(), membership.Membership.UserKey)
	assert.Equal(t, false, membership.Membership.IsAdmin)
	assert.Equal(t, false, membership.Membership.IsOwner)
	assert.Equal(t, true, membership.Membership.IsMember)
	assert.Equal(t, true, membership.Membership.GroupConfirmed)
	assert.Equal(t, true, membership.Membership.UserConfirmed)

}

func TestInviteeShouldBeAbleToDeclineInvitationFromOwner(t *testing.T) {

	ctx := context.Background()

	owner, ownerCli := testUserCli(t)
	group := &handler.GetGroupResponse{}
	if !assert.NoError(t, testGroup2(t, owner, group)) {
		return
	}

	user, userCli := testUserCli(t)
	membershipKey := keys.NewMembershipKey(group, user)

	// owner invites user
	if !assert.NoError(t, ownerCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	time.Sleep(time.Second)

	if !assert.NoError(t, userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})) {
		return
	}

	// user declines
	if !assert.NoError(t, userCli.LeaveGroup(ctx, membershipKey)) {
		return
	}

	time.Sleep(time.Second)

	err := userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})
	assert.Error(t, err)
	assert.ErrorIs(t, err, exceptions.ErrMembershipNotFound)

}

func TestUserShouldBeAbleToDeclineInvitationFromOwner(t *testing.T) {

	ctx := context.Background()

	owner, ownerCli := testUserCli(t)
	group := &handler.GetGroupResponse{}
	if !assert.NoError(t, testGroup2(t, owner, group)) {
		return
	}

	user, userCli := testUserCli(t)
	membershipKey := keys.NewMembershipKey(group, user)

	// user requests invitation
	if !assert.NoError(t, userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	if !assert.NoError(t, userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})) {
		return
	}

	// owner declines
	if !assert.NoError(t, ownerCli.LeaveGroup(ctx, membershipKey)) {
		return
	}

	time.Sleep(time.Second)

	err := userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})
	assert.Error(t, err)
	assert.ErrorIs(t, err, exceptions.ErrMembershipNotFound)

}

func TestRandomUserShouldNotBeAbleToAcceptOrDeclineInvitation(t *testing.T) {

	ctx := context.Background()

	owner, _ := testUserCli(t)
	group := &handler.GetGroupResponse{}
	if !assert.NoError(t, testGroup2(t, owner, group)) {
		return
	}

	user, userCli := testUserCli(t)
	membershipKey := keys.NewMembershipKey(group, user)

	_, randomUserCli := testUserCli(t)

	// user requests invitation
	if !assert.NoError(t, userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	// check membership exists
	if !assert.NoError(t, userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})) {
		return
	}

	// random user tries to decline
	err := randomUserCli.LeaveGroup(ctx, membershipKey)
	assert.Error(t, err)
	assert.ErrorIs(t, err, exceptions.ErrForbidden)

	// random user tries to approve
	err = randomUserCli.JoinGroup(ctx, membershipKey)
	assert.Error(t, err)
	assert.ErrorIs(t, err, exceptions.ErrForbidden)

}

func TestPersonShouldBeAbleToRequestBeingInvitedInGroup(t *testing.T) {

	ctx := context.Background()

	owner, _ := testUserCli(t)
	group := &handler.GetGroupResponse{}
	if !assert.NoError(t, testGroup2(t, owner, group)) {
		return
	}

	user, userCli := testUserCli(t)
	membershipKey := keys.NewMembershipKey(group, user)

	// user requests invitation
	if !assert.NoError(t, userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	// get membership
	membership := &handler.GetMembershipResponse{}
	if !assert.NoError(t, userCli.GetMembership(ctx, membershipKey, membership)) {
		return
	}

	assert.Equal(t, user.GetUserKey(), membership.Membership.UserKey)
	assert.Equal(t, false, membership.Membership.IsAdmin)
	assert.Equal(t, false, membership.Membership.IsOwner)
	assert.Equal(t, false, membership.Membership.IsMember)
	assert.Equal(t, false, membership.Membership.GroupConfirmed)
	assert.Equal(t, true, membership.Membership.UserConfirmed)

}

func TestOwnerShouldBeAbleToAcceptInvitationRequest(t *testing.T) {

	ctx := context.Background()

	owner, ownerCli := testUserCli(t)
	group := &handler.GetGroupResponse{}
	if !assert.NoError(t, testGroup2(t, owner, group)) {
		return
	}

	user, userCli := testUserCli(t)
	membershipKey := keys.NewMembershipKey(group, user)

	// user requests invitation
	if !assert.NoError(t, userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	// owner accepts invitation
	if !assert.NoError(t, ownerCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	time.Sleep(time.Second)

	// get membership
	membership := &handler.GetMembershipResponse{}
	if !assert.NoError(t, userCli.GetMembership(ctx, membershipKey, membership)) {
		return
	}

	assert.Equal(t, user.GetUserKey(), membership.Membership.UserKey)
	assert.Equal(t, false, membership.Membership.IsAdmin)
	assert.Equal(t, false, membership.Membership.IsOwner)
	assert.Equal(t, true, membership.Membership.IsMember)
	assert.Equal(t, true, membership.Membership.GroupConfirmed)
	assert.Equal(t, true, membership.Membership.UserConfirmed)
}

func TestGetLoggedInUserMembershipsWithoutGroup(t *testing.T) {

	ctx := context.Background()
	_, userCli := testUserCli(t)
	var memberships handler.GetMembershipsResponse
	if !assert.NoError(t, userCli.GetLoggedInUserMemberships(ctx, &memberships)) {
		return
	}
	assert.Len(t, memberships.Memberships, 0)
}

func TestTestGetLoggedInUserMembershipsWithGroup(t *testing.T) {

	ctx := context.Background()

	owner, ownerCli := testUserCli(t)
	group := &handler.GetGroupResponse{}
	if !assert.NoError(t, testGroup2(t, owner, group)) {
		return
	}

	time.Sleep(time.Second)

	var memberships handler.GetMembershipsResponse
	if !assert.NoError(t, ownerCli.GetLoggedInUserMemberships(ctx, &memberships)) {
		return
	}
	if !assert.Len(t, memberships.Memberships, 1) {
		return
	}
	membership := memberships.Memberships[0]
	assert.Equal(t, group.GetGroupKey(), membership.GroupKey)
	assert.Equal(t, group.Group.Name, membership.GroupName)
	assert.Equal(t, group.Group.GroupKey, membership.GroupKey)
	assert.Equal(t, owner.GetUserKey(), membership.UserKey)
	assert.Equal(t, owner.Username, membership.UserName)
	assert.Equal(t, true, membership.UserConfirmed)
	assert.Equal(t, true, membership.GroupConfirmed)
	assert.Equal(t, true, membership.IsMember)
	assert.Equal(t, true, membership.IsOwner)
}

func TestGroupShouldReceiveMessageWhenUserJoined(t *testing.T) {

	ctx := context.Background()

	_, user1Cli := testUserCli(t)
	user2, user2Cli := testUserCli(t)

	ws, err := user1Cli.GetWebsocketClient()

	if !assert.NoError(t, err) {
		return
	}
	defer ws.Close()

	groupResponse := &handler.GetGroupResponse{}
	if !assert.NoError(t, user1Cli.CreateGroup(ctx, handler.NewCreateGroupRequest("sample", "description"), groupResponse)) {
		return
	}
	if !assert.NoError(t, user2Cli.JoinGroup(ctx, keys.NewMembershipKey(groupResponse.Group.GroupKey, user2.GetUserKey()))) {
		return
	}
	if !assert.NoError(t, user1Cli.JoinGroup(ctx, keys.NewMembershipKey(groupResponse.Group.GroupKey, user2.GetUserKey()))) {
		return
	}

	// TODO :
	// x1done := make(chan bool)
	// x2done := make(chan bool)
	// fail := make(chan bool)
	//
	// go func() {
	// 	i1 := 0
	// 	for delivery := range x1.Delivery {
	// 		i1++
	// 		fmt.Println("received message on queue 1")
	// 		fmt.Println("Message type: " + delivery.Type)
	// 		fmt.Println("Message body: " + string(delivery.Body))
	// 		if i1 == 2 {
	// 			x1done <- true
	// 		}
	// 	}
	// }()
	//
	// go func() {
	// 	i1 := 0
	// 	for delivery := range x2.Delivery {
	// 		i1++
	// 		fmt.Println("received message on queue 2")
	// 		fmt.Println("Message type: " + delivery.Type)
	// 		fmt.Println("Message body: " + string(delivery.Body))
	// 		if i1 == 1 {
	// 			x2done <- true
	// 		}
	// 	}
	// }()
	//
	// go func() {
	// 	time.Sleep(time.Second * 10)
	// 	fail <- true
	// }()
	//
	// i1done := false
	// i2done := false
	// for {
	// 	select {
	// 	case _ = <-x1done:
	// 		i1done = true
	// 		if i2done {
	// 			return
	// 		}
	// 		break
	// 	case _ = <-x2done:
	// 		i2done = true
	// 		if i1done {
	// 			return
	// 		}
	// 		break
	// 	case _ = <-fail:
	// 		t.FailNow()
	// 	}
	//
	// }

}

func TestGetUsersForInvitePickerShouldNotReturnDuplicates(t *testing.T) {

	ctx := context.Background()

	var userClis []client.Client
	var groups []*handler.GetGroupResponse
	var userCount = 3
	var groupCount = 2

	g1, ctx := errgroup.WithContext(ctx)
	for i := 0; i < userCount; i++ {
		g1.Go(func() error {
			user, userCli := testUserCli(t)
			userClis = append(userClis, userCli)
			g2, ctx := errgroup.WithContext(ctx)
			for j := 0; j < groupCount; j++ {
				func(userCli client.Client, user *models.UserSession) {
					g2.Go(func() error {
						t.Log(user.Username)
						var response handler.GetGroupResponse
						if err := userCli.CreateGroup(ctx, handler.NewCreateGroupRequest("sample", "description"), &response); err != nil {
							return err
						}
						groups = append(groups, &response)
						return nil
					})
				}(userCli, user)
			}
			return g2.Wait()
		})
	}

	if !assert.NoError(t, g1.Wait()) {
		return
	}

	for i := 0; i < userCount*groupCount; i++ {
		group := groups[i]
		var response handler.GetUsersForGroupInvitePickerResponse
		if !assert.NoError(t, userClis[0].GetMemberInvitationPicker(ctx, group.Group.GroupKey, "", 0, 10, &response)) {
			return
		}
		seen := map[string]bool{}
		for _, user := range response.Users {
			_, ok := seen[user.Username]
			if !assert.False(t, ok, "found duplicate user %s in results", user.Username) {
				return
			}
			seen[user.Username] = true
		}
	}
}
