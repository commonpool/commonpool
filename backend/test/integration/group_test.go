package integration

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/client"
	"github.com/commonpool/backend/pkg/client/echo"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type GroupTestSuite struct {
	suite.Suite
	*IntegrationTestSuite
	groupOwner    *models.UserSession
	groupKey      keys.GroupKey
	group         *handler.GetGroupResponse
	delGroupOwner func()
}

func TestGroupTestSuite(t *testing.T) {
	suite.Run(t, new(GroupTestSuite))
}

func (s *GroupTestSuite) SetupSuite() {
	s.IntegrationTestSuite = &IntegrationTestSuite{}
	s.IntegrationTestSuite.SetupSuite()
	ctx := context.Background()
	cli := s.NewClient(s.groupOwner)
	var group handler.GetGroupResponse
	if !assert.NoError(s.T(), cli.CreateGroup(ctx, &handler.CreateGroupRequest{
		Name:        "sample",
		Description: "description",
	}, &group)) {
		return
	}
	s.group = &group
	s.groupKey = s.group.Group.GroupKey
	time.Sleep(1 * time.Second)
}

func (s *GroupTestSuite) NewClient(user *models.UserSession) client.Client {
	return echo.NewEchoClient(s.server.Router, client.NewMockAuthentication(user))
}

func (s *GroupTestSuite) TestCreateGroup() {
	ctx := context.Background()

	var _, cli = s.testUserCli(s.T())

	var response = &handler.GetGroupResponse{}
	if !assert.NoError(s.T(), cli.CreateGroup(ctx, handler.NewCreateGroupRequest("sample", "description"), response)) {
		return
	}

	assert.NotNil(s.T(), response.Group)
	assert.Equal(s.T(), response.Group.Name, "sample")
	assert.Equal(s.T(), response.Group.Description, "description")
	assert.NotEmpty(s.T(), response.Group.GroupKey)
	assert.NotEmpty(s.T(), response.Group.CreatedAt)
}

func (s *GroupTestSuite) TestCreateGroupUnauthenticatedShouldFailWithUnauthorized() {
	err := s.NewClient(nil).CreateGroup(context.Background(), handler.NewCreateGroupRequest("sample", "description"), &handler.GetGroupResponse{})
	if !assert.Error(s.T(), err) {
		return
	}
}

func (s *GroupTestSuite) TestCreateGroupEmptyNameShouldFail() {
	ctx := context.Background()
	var _, cli = s.testUserCli(s.T())
	assert.Error(s.T(), cli.CreateGroup(ctx, handler.NewCreateGroupRequest("", "description"), &handler.GetGroupResponse{}))
	assert.Error(s.T(), cli.CreateGroup(ctx, handler.NewCreateGroupRequest("   ", "description"), &handler.GetGroupResponse{}))
}

func (s *GroupTestSuite) TestCreateGroupEmptyDescriptionShouldNotFail() {
	ctx := context.Background()
	var _, cli = s.testUserCli(s.T())
	assert.NoError(s.T(), cli.CreateGroup(ctx, handler.NewCreateGroupRequest("name", ""), &handler.GetGroupResponse{}))
}

func (s *GroupTestSuite) TestCreateGroupShouldCreateOwnerMembership() {
	t := s.T()
	ctx := context.TODO()
	membership, err := s.server.GetMembership.Get(ctx, keys.NewMembershipKey(s.groupKey, s.groupOwner.GetUserKey()))
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, s.groupOwner.GetUserKey(), membership.UserKey)
	assert.Equal(t, true, membership.IsOwner)
	assert.Equal(t, true, membership.UserConfirmed)
	assert.Equal(t, true, membership.GroupConfirmed)
	assert.Equal(t, true, membership.IsAdmin)
	assert.Equal(t, true, membership.IsMember)
}

func (s *GroupTestSuite) TestOwnerShouldBeAbleToInviteUser() {

	ctx := context.TODO()
	user, _ := s.testUser(s.T())
	ownerCli := s.NewClient(s.groupOwner)
	membershipKey := keys.NewMembershipKey(s.groupKey, user.GetUserKey())

	if !assert.NoError(s.T(), ownerCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	var membership handler.GetMembershipResponse
	if !assert.NoError(s.T(), ownerCli.GetMembership(ctx, membershipKey, &membership)) {
		return
	}

	assert.Equal(s.T(), user.GetUserKey(), membership.Membership.UserKey)
	assert.Equal(s.T(), false, membership.Membership.IsAdmin)
	assert.Equal(s.T(), false, membership.Membership.IsOwner)
	assert.Equal(s.T(), false, membership.Membership.IsMember)
	assert.Equal(s.T(), true, membership.Membership.GroupConfirmed)
	assert.Equal(s.T(), false, membership.Membership.UserConfirmed)
}

func (s *GroupTestSuite) TestInviteeShouldBeAbleToAcceptInvitationFromOwner() {
	ctx := context.Background()

	user, _ := s.testUser(s.T())
	membershipKey := keys.NewMembershipKey(s.groupKey, user.GetUserKey())

	ownerCli := s.NewClient(s.groupOwner)
	userCli := s.NewClient(user)

	if !assert.NoError(s.T(), ownerCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	if !assert.NoError(s.T(), userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	var membership handler.GetMembershipResponse
	if !assert.NoError(s.T(), userCli.GetMembership(ctx, membershipKey, &membership)) {
		return
	}

	assert.Equal(s.T(), user.GetUserKey(), membership.Membership.UserKey)
	assert.Equal(s.T(), false, membership.Membership.IsAdmin)
	assert.Equal(s.T(), false, membership.Membership.IsOwner)
	assert.Equal(s.T(), true, membership.Membership.IsMember)
	assert.Equal(s.T(), true, membership.Membership.GroupConfirmed)
	assert.Equal(s.T(), true, membership.Membership.UserConfirmed)

}

func (s *GroupTestSuite) TestInviteeShouldBeAbleToDeclineInvitationFromOwner() {
	ctx := context.Background()

	user, _ := s.testUser(s.T())
	membershipKey := keys.NewMembershipKey(s.groupKey, user.GetUserKey())

	ownerCli := s.NewClient(s.groupOwner)
	userCli := s.NewClient(user)

	// owner invites user
	if !assert.NoError(s.T(), ownerCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	if !assert.NoError(s.T(), userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})) {
		return
	}

	// user declines
	if !assert.NoError(s.T(), userCli.LeaveGroup(ctx, membershipKey)) {
		return
	}

	err := userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, exceptions.ErrMembershipNotFound)

}

func (s *GroupTestSuite) TestOwnerShouldBeAbleToDeclineInvitationFromOwner() {
	ctx := context.Background()

	user, _ := s.testUser(s.T())
	membershipKey := keys.NewMembershipKey(s.groupKey, user.GetUserKey())

	ownerCli := s.NewClient(s.groupOwner)
	userCli := s.NewClient(user)

	// user requests invitation
	if !assert.NoError(s.T(), userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	if !assert.NoError(s.T(), userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})) {
		return
	}

	// owner declines
	if !assert.NoError(s.T(), ownerCli.LeaveGroup(ctx, membershipKey)) {
		return
	}

	err := userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, exceptions.ErrMembershipNotFound)

}

func (s *GroupTestSuite) TestRandomUserShouldNotBeAbleToAcceptOrDeclineInvitation() {
	ctx := context.Background()

	user, _ := s.testUser(s.T())
	membershipKey := keys.NewMembershipKey(s.groupKey, user.GetUserKey())
	userCli := s.NewClient(user)

	randomUser, _ := s.testUser(s.T())
	randomUserCli := s.NewClient(randomUser)

	// user requests invitation
	if !assert.NoError(s.T(), userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	// check membership exists
	if !assert.NoError(s.T(), userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})) {
		return
	}

	// random user tries to decline
	err := randomUserCli.LeaveGroup(ctx, membershipKey)
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, exceptions.ErrForbidden)

	// random user tries to approve
	err = randomUserCli.JoinGroup(ctx, membershipKey)
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, exceptions.ErrForbidden)

}

func (s *GroupTestSuite) TestPersonShouldBeAbleToRequestBeingInvitedInGroup() {
	ctx := context.Background()
	user, userCli := s.testUserCli(s.T())
	membershipKey := keys.NewMembershipKey(s.groupKey, user.GetUserKey())

	// user requests invitation
	if !assert.NoError(s.T(), userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	// get membership
	membership := &handler.GetMembershipResponse{}
	if !assert.NoError(s.T(), userCli.GetMembership(ctx, membershipKey, membership)) {
		return
	}

	assert.Equal(s.T(), user.GetUserKey(), membership.Membership.UserKey)
	assert.Equal(s.T(), false, membership.Membership.IsAdmin)
	assert.Equal(s.T(), false, membership.Membership.IsOwner)
	assert.Equal(s.T(), false, membership.Membership.IsMember)
	assert.Equal(s.T(), false, membership.Membership.GroupConfirmed)
	assert.Equal(s.T(), true, membership.Membership.UserConfirmed)

}

func (s *GroupTestSuite) TestOwnerShouldBeAbleToAcceptInvitationRequest() {
	ctx := context.Background()

	ownerCli := s.NewClient(s.groupOwner)
	user, userCli := s.testUserCli(s.T())
	membershipKey := keys.NewMembershipKey(s.groupKey, user.GetUserKey())

	// user requests invitation
	if !assert.NoError(s.T(), userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	// owner accepts invitation
	if !assert.NoError(s.T(), ownerCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	// get membership
	membership := &handler.GetMembershipResponse{}
	if !assert.NoError(s.T(), userCli.GetMembership(ctx, membershipKey, membership)) {
		return
	}

	assert.Equal(s.T(), user.GetUserKey(), membership.Membership.UserKey)
	assert.Equal(s.T(), false, membership.Membership.IsAdmin)
	assert.Equal(s.T(), false, membership.Membership.IsOwner)
	assert.Equal(s.T(), true, membership.Membership.IsMember)
	assert.Equal(s.T(), true, membership.Membership.GroupConfirmed)
	assert.Equal(s.T(), true, membership.Membership.UserConfirmed)
}

func (s *GroupTestSuite) TestGetLoggedInUserMembershipsWithoutGroup() {
	ctx := context.Background()
	_, userCli := s.testUserCli(s.T())
	var memberships handler.GetMembershipsResponse
	if !assert.NoError(s.T(), userCli.GetLoggedInUserMemberships(ctx, &memberships)) {
		return
	}
	assert.Len(s.T(), memberships.Memberships, 0)
}

func (s *GroupTestSuite) TestTestGetLoggedInUserMembershipsWithGroup() {

	ctx := context.Background()
	ownerCli := s.NewClient(s.groupOwner)

	var memberships handler.GetMembershipsResponse
	if !assert.NoError(s.T(), ownerCli.GetLoggedInUserMemberships(ctx, &memberships)) {
		return
	}

	if !assert.Len(s.T(), memberships.Memberships, 1) {
		return
	}

	membership := memberships.Memberships[0]
	assert.Equal(s.T(), s.groupKey, membership.GroupKey)
	assert.Equal(s.T(), s.group.Group.Name, membership.GroupName)
	assert.Equal(s.T(), s.group.Group.GroupKey, membership.GroupKey)
	assert.Equal(s.T(), s.groupOwner.GetUserKey(), membership.UserKey)
	assert.Equal(s.T(), s.groupOwner.Username, membership.UserName)
	assert.Equal(s.T(), true, membership.UserConfirmed)
	assert.Equal(s.T(), true, membership.GroupConfirmed)
	assert.Equal(s.T(), true, membership.IsMember)
	assert.Equal(s.T(), true, membership.IsOwner)
}

func (s *GroupTestSuite) TestGroupShouldReceiveMessageWhenUserJoined() {
	ctx := context.Background()

	_, user1Cli := s.testUserCli(s.T())
	user2, user2Cli := s.testUserCli(s.T())

	ws, err := user1Cli.GetWebsocketClient()
	if !assert.NoError(s.T(), err) {
		return
	}
	defer ws.Close()

	groupResponse := &handler.GetGroupResponse{}
	if !assert.NoError(s.T(), user1Cli.CreateGroup(ctx, handler.NewCreateGroupRequest("sample", "description"), groupResponse)) {
		return
	}
	if !assert.NoError(s.T(), user2Cli.JoinGroup(ctx, keys.NewMembershipKey(groupResponse.Group.GroupKey, user2.GetUserKey()))) {
		return
	}
	if !assert.NoError(s.T(), user1Cli.JoinGroup(ctx, keys.NewMembershipKey(groupResponse.Group.GroupKey, user2.GetUserKey()))) {
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

func (s *GroupTestSuite) TestGetUsersForInvitePickerShouldNotReturnDuplicates() {
	t := s.T()
	ctx := context.Background()

	var userClis []client.Client
	var groups []*handler.GetGroupResponse
	var userCount = 3
	var groupCount = 2
	for i := 0; i < userCount; i++ {
		_, userCli := s.testUserCli(s.T())
		userClis = append(userClis, userCli)
		for j := 0; j < groupCount; j++ {
			var response handler.GetGroupResponse
			if !assert.NoError(s.T(), userCli.CreateGroup(ctx, handler.NewCreateGroupRequest("sample", "description"), &response)) {
				return
			}
			groups = append(groups, &response)
		}
	}

	for i := 0; i < userCount*groupCount; i++ {
		group := groups[i]
		var response handler.GetUsersForGroupInvitePickerResponse
		if !assert.NoError(s.T(), userClis[0].GetMemberInvitationPicker(ctx, group.Group.GroupKey, "", 0, 10, &response)) {
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
