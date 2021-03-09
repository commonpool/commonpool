package integration

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/client"
	"github.com/commonpool/backend/pkg/client/echo"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"
	"testing"
	"time"
)

type GroupTestSuite struct {
	suite.Suite
	*IntegrationTestBase
}

func TestGroupTestSuite(t *testing.T) {
	suite.Run(t, &GroupTestSuite{})
}

func (s *GroupTestSuite) SetupSuite() {
	s.IntegrationTestBase = &IntegrationTestBase{}
	s.IntegrationTestBase.Setup()
	time.Sleep(1 * time.Second)
}

func (s *GroupTestSuite) NewClient(user *models.UserSession) client.Client {
	return echo.NewEchoClient(s.server.Router, client.NewMockAuthentication(user))
}

func (s *GroupTestSuite) TestCreateGroup() {

	ctx := context.Background()

	var _, cli = s.testUserCli(s.T())

	var response = &handler.GetGroupResponse{}
	if !s.NoError(cli.CreateGroup(ctx, handler.NewCreateGroupRequest("sample", "description"), response)) {
		return
	}

	s.NotNil(response.Group)
	s.Equal(response.Group.Name, "sample")
	s.Equal(response.Group.Description, "description")
	s.NotEmpty(response.Group.GroupKey)
	s.NotEmpty(response.Group.CreatedAt)
}

func (s *GroupTestSuite) TestCreateGroupUnauthenticatedShouldFailWithUnauthorized() {

	err := s.NewClient(nil).CreateGroup(context.Background(), handler.NewCreateGroupRequest("sample", "description"), &handler.GetGroupResponse{})
	if !s.Error(err) {
		return
	}
}

func (s *GroupTestSuite) TestCreateGroupEmptyNameShouldFail() {

	ctx := context.Background()
	var _, cli = s.testUserCli(s.T())
	s.Error(cli.CreateGroup(ctx, handler.NewCreateGroupRequest("", "description"), &handler.GetGroupResponse{}))
	s.Error(cli.CreateGroup(ctx, handler.NewCreateGroupRequest("   ", "description"), &handler.GetGroupResponse{}))
}

func (s *GroupTestSuite) TestCreateGroupEmptyDescriptionShouldNotFail() {

	ctx := context.Background()
	var _, cli = s.testUserCli(s.T())
	s.NoError(cli.CreateGroup(ctx, handler.NewCreateGroupRequest("name", ""), &handler.GetGroupResponse{}))
}

func (s *GroupTestSuite) TestCreateGroupShouldCreateOwnerMembership() {

	ctx := context.TODO()

	owner, _ := s.testUserCli(s.T())
	group := &handler.GetGroupResponse{}
	if !s.NoError(s.testGroup2(s.T(), owner, group)) {
		return
	}

	time.Sleep(time.Second)

	membership, err := s.server.GetMembership.Get(ctx, keys.NewMembershipKey(group, owner))
	if !s.NoError(err) {
		return
	}
	s.Equal(owner.GetUserKey(), membership.UserKey)
	s.Equal(true, membership.IsOwner)
	s.Equal(true, membership.UserConfirmed)
	s.Equal(true, membership.GroupConfirmed)
	s.Equal(true, membership.IsAdmin)
	s.Equal(true, membership.IsMember)
}

func (s *GroupTestSuite) TestOwnerShouldBeAbleToInviteUser() {

	ctx := context.TODO()

	owner, ownerCli := s.testUserCli(s.T())
	group := &handler.GetGroupResponse{}
	if !s.NoError(s.testGroup2(s.T(), owner, group)) {
		return
	}

	user, _ := s.testUser(s.T())
	membershipKey := keys.NewMembershipKey(group, user)

	if !s.NoError(ownerCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	var membership handler.GetMembershipResponse
	if !s.NoError(ownerCli.GetMembership(ctx, membershipKey, &membership)) {
		return
	}

	s.Equal(user.GetUserKey(), membership.Membership.UserKey)
	s.Equal(false, membership.Membership.IsAdmin)
	s.Equal(false, membership.Membership.IsOwner)
	s.Equal(false, membership.Membership.IsMember)
	s.Equal(true, membership.Membership.GroupConfirmed)
	s.Equal(false, membership.Membership.UserConfirmed)
}

func (s *GroupTestSuite) TestInviteeShouldBeAbleToAcceptInvitationFromOwner() {

	ctx := context.Background()

	owner, ownerCli := s.testUserCli(s.T())
	group := &handler.GetGroupResponse{}
	if !s.NoError(s.testGroup2(s.T(), owner, group)) {
		return
	}

	user, _ := s.testUser(s.T())
	userCli := s.NewClient(user)
	membershipKey := keys.NewMembershipKey(group, user)

	if !s.NoError(ownerCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	if !s.NoError(userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	time.Sleep(time.Second)

	var membership handler.GetMembershipResponse
	if !s.NoError(userCli.GetMembership(ctx, membershipKey, &membership)) {
		return
	}

	s.Equal(user.GetUserKey(), membership.Membership.UserKey)
	s.Equal(false, membership.Membership.IsAdmin)
	s.Equal(false, membership.Membership.IsOwner)
	s.Equal(true, membership.Membership.IsMember)
	s.Equal(true, membership.Membership.GroupConfirmed)
	s.Equal(true, membership.Membership.UserConfirmed)

}

func (s *GroupTestSuite) TestInviteeShouldBeAbleToDeclineInvitationFromOwner() {

	ctx := context.Background()

	owner, ownerCli := s.testUserCli(s.T())
	group := &handler.GetGroupResponse{}
	if !s.NoError(s.testGroup2(s.T(), owner, group)) {
		return
	}

	user, userCli := s.testUserCli(s.T())
	membershipKey := keys.NewMembershipKey(group, user)

	// owner invites user
	if !s.NoError(ownerCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	time.Sleep(time.Second)

	if !s.NoError(userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})) {
		return
	}

	// user declines
	if !s.NoError(userCli.LeaveGroup(ctx, membershipKey)) {
		return
	}

	time.Sleep(time.Second)

	err := userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})
	s.Error(err)
	s.ErrorIs(err, exceptions.ErrMembershipNotFound)

}

func (s *GroupTestSuite) TestUserShouldBeAbleToDeclineInvitationFromOwner() {

	ctx := context.Background()

	owner, ownerCli := s.testUserCli(s.T())
	group := &handler.GetGroupResponse{}
	if !s.NoError(s.testGroup2(s.T(), owner, group)) {
		return
	}

	user, userCli := s.testUserCli(s.T())
	membershipKey := keys.NewMembershipKey(group, user)

	// user requests invitation
	if !s.NoError(userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	if !s.NoError(userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})) {
		return
	}

	// owner declines
	if !s.NoError(ownerCli.LeaveGroup(ctx, membershipKey)) {
		return
	}

	time.Sleep(time.Second)

	err := userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})
	s.Error(err)
	s.ErrorIs(err, exceptions.ErrMembershipNotFound)

}

func (s *GroupTestSuite) TestRandomUserShouldNotBeAbleToAcceptOrDeclineInvitation() {

	ctx := context.Background()

	owner, _ := s.testUserCli(s.T())
	group := &handler.GetGroupResponse{}
	if !s.NoError(s.testGroup2(s.T(), owner, group)) {
		return
	}

	user, userCli := s.testUserCli(s.T())
	membershipKey := keys.NewMembershipKey(group, user)

	_, randomUserCli := s.testUserCli(s.T())

	// user requests invitation
	if !s.NoError(userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	// check membership exists
	if !s.NoError(userCli.GetMembership(ctx, membershipKey, &handler.GetMembershipResponse{})) {
		return
	}

	// random user tries to decline
	err := randomUserCli.LeaveGroup(ctx, membershipKey)
	s.Error(err)
	s.ErrorIs(err, exceptions.ErrForbidden)

	// random user tries to approve
	err = randomUserCli.JoinGroup(ctx, membershipKey)
	s.Error(err)
	s.ErrorIs(err, exceptions.ErrForbidden)

}

func (s *GroupTestSuite) TestPersonShouldBeAbleToRequestBeingInvitedInGroup() {

	ctx := context.Background()

	owner, _ := s.testUserCli(s.T())
	group := &handler.GetGroupResponse{}
	if !s.NoError(s.testGroup2(s.T(), owner, group)) {
		return
	}

	user, userCli := s.testUserCli(s.T())
	membershipKey := keys.NewMembershipKey(group, user)

	// user requests invitation
	if !s.NoError(userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	// get membership
	membership := &handler.GetMembershipResponse{}
	if !s.NoError(userCli.GetMembership(ctx, membershipKey, membership)) {
		return
	}

	s.Equal(user.GetUserKey(), membership.Membership.UserKey)
	s.Equal(false, membership.Membership.IsAdmin)
	s.Equal(false, membership.Membership.IsOwner)
	s.Equal(false, membership.Membership.IsMember)
	s.Equal(false, membership.Membership.GroupConfirmed)
	s.Equal(true, membership.Membership.UserConfirmed)

}

func (s *GroupTestSuite) TestOwnerShouldBeAbleToAcceptInvitationRequest() {

	ctx := context.Background()

	owner, ownerCli := s.testUserCli(s.T())
	group := &handler.GetGroupResponse{}
	if !s.NoError(s.testGroup2(s.T(), owner, group)) {
		return
	}

	user, userCli := s.testUserCli(s.T())
	membershipKey := keys.NewMembershipKey(group, user)

	// user requests invitation
	if !s.NoError(userCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	// owner accepts invitation
	if !s.NoError(ownerCli.JoinGroup(ctx, membershipKey)) {
		return
	}

	time.Sleep(time.Second)

	// get membership
	membership := &handler.GetMembershipResponse{}
	if !s.NoError(userCli.GetMembership(ctx, membershipKey, membership)) {
		return
	}

	s.Equal(user.GetUserKey(), membership.Membership.UserKey)
	s.Equal(false, membership.Membership.IsAdmin)
	s.Equal(false, membership.Membership.IsOwner)
	s.Equal(true, membership.Membership.IsMember)
	s.Equal(true, membership.Membership.GroupConfirmed)
	s.Equal(true, membership.Membership.UserConfirmed)
}

func (s *GroupTestSuite) TestGetLoggedInUserMembershipsWithoutGroup() {

	ctx := context.Background()
	_, userCli := s.testUserCli(s.T())
	var memberships handler.GetMembershipsResponse
	if !s.NoError(userCli.GetLoggedInUserMemberships(ctx, &memberships)) {
		return
	}
	s.Len(memberships.Memberships, 0)
}

func (s *GroupTestSuite) TestTestGetLoggedInUserMembershipsWithGroup() {

	ctx := context.Background()

	owner, ownerCli := s.testUserCli(s.T())
	group := &handler.GetGroupResponse{}
	if !s.NoError(s.testGroup2(s.T(), owner, group)) {
		return
	}

	time.Sleep(time.Second)

	var memberships handler.GetMembershipsResponse
	if !s.NoError(ownerCli.GetLoggedInUserMemberships(ctx, &memberships)) {
		return
	}
	if !s.Len(memberships.Memberships, 1) {
		return
	}
	membership := memberships.Memberships[0]
	s.Equal(group.GetGroupKey(), membership.GroupKey)
	s.Equal(group.Group.Name, membership.GroupName)
	s.Equal(group.Group.GroupKey, membership.GroupKey)
	s.Equal(owner.GetUserKey(), membership.UserKey)
	s.Equal(owner.Username, membership.UserName)
	s.Equal(true, membership.UserConfirmed)
	s.Equal(true, membership.GroupConfirmed)
	s.Equal(true, membership.IsMember)
	s.Equal(true, membership.IsOwner)
}

func (s *GroupTestSuite) TestGroupShouldReceiveMessageWhenUserJoined() {

	ctx := context.Background()

	_, user1Cli := s.testUserCli(s.T())
	user2, user2Cli := s.testUserCli(s.T())

	ws, err := user1Cli.GetWebsocketClient()

	if !s.NoError(err) {
		return
	}
	defer ws.Close()

	groupResponse := &handler.GetGroupResponse{}
	if !s.NoError(user1Cli.CreateGroup(ctx, handler.NewCreateGroupRequest("sample", "description"), groupResponse)) {
		return
	}
	if !s.NoError(user2Cli.JoinGroup(ctx, keys.NewMembershipKey(groupResponse.Group.GroupKey, user2.GetUserKey()))) {
		return
	}
	if !s.NoError(user1Cli.JoinGroup(ctx, keys.NewMembershipKey(groupResponse.Group.GroupKey, user2.GetUserKey()))) {
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

	ctx := context.Background()

	var userClis []client.Client
	var groups []*handler.GetGroupResponse
	var userCount = 3
	var groupCount = 2

	g1, ctx := errgroup.WithContext(ctx)
	for i := 0; i < userCount; i++ {
		g1.Go(func() error {
			user, userCli := s.testUserCli(s.T())
			userClis = append(userClis, userCli)
			g2, ctx := errgroup.WithContext(ctx)
			for j := 0; j < groupCount; j++ {
				func(userCli client.Client, user *models.UserSession) {
					g2.Go(func() error {
						s.T().Log(user.Username)
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

	if !s.NoError(g1.Wait()) {
		return
	}

	for i := 0; i < userCount*groupCount; i++ {
		group := groups[i]
		var response handler.GetUsersForGroupInvitePickerResponse
		if !s.NoError(userClis[0].GetMemberInvitationPicker(ctx, group.Group.GroupKey, "", 0, 10, &response)) {
			return
		}
		seen := map[string]bool{}
		for _, user := range response.Users {
			_, ok := seen[user.Username]
			if !s.False(ok, "found duplicate user %s in results", user.Username) {
				return
			}
			seen[user.Username] = true
		}
	}
}
