package integration

import (
	"context"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/client"
	"github.com/commonpool/backend/pkg/client/echo"
	"github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/group/readmodels"
	"github.com/commonpool/backend/pkg/keys"
	res "github.com/commonpool/backend/pkg/resource/handler"
	"github.com/commonpool/backend/test"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func testUser(t *testing.T) (*models.UserSession, func()) {
	t.Helper()

	createUserLock.Lock()
	defer createUserLock.Unlock()

	u := NewUser()
	upsertError := srv.User.Store.Upsert(u.GetUserKey(), u.Email, u.Username)

	var userXchangeErr error
	if upsertError != nil {
		_, userXchangeErr = srv.ChatService.CreateUserExchange(context.TODO(), u.GetUserKey())
	}

	if upsertError != nil {
		t.Fatalf("upsert error: %s", upsertError)
	}

	if userXchangeErr != nil {
		t.Fatalf("exchange error: %s", userXchangeErr)
	}

	if err := getUserClient(u).GetLoggedInUserMemberships(context.TODO(), &handler.GetMembershipsResponse{}); err != nil {
		t.Fatal(err)
	}

	return u, func(user *models.UserSession) func() {
		return func() {
			session := srv.GraphDriver.GetSession()
			defer session.Close()
			result, err := session.Run(`MATCH (u:User{id:$id}) detach delete u`, map[string]interface{}{
				"id": user.Subject,
			})
			if err != nil {
				t.Fatal(err)
			}
			if result.Err() != nil {
				t.Fatal(result.Err())
			}
		}
	}(u)
}

func testUserCli(t *testing.T) (*models.UserSession, client.Client) {
	user, _ := testUser(t)
	return user, getUserClient(user)
}

func getUserClient(user *models.UserSession) client.Client {
	return echo.NewEchoClient(srv.Router, client.NewMockAuthentication(user))
}

func testGroup2(t *testing.T, owner *models.UserSession, output *handler.GetGroupResponse, members ...*models.UserSession) error {
	ctx := context.Background()
	groupCounter++
	ownerCli := getUserClient(owner)
	if err := ownerCli.CreateGroup(ctx, handler.NewCreateGroupRequest("group-"+strconv.Itoa(groupCounter), "group-"+strconv.Itoa(groupCounter)), output); !assert.NoError(t, err) {
		return err
	}
	for _, member := range members {
		if err := getUserClient(member).JoinGroup(ctx, keys.NewMembershipKey(output.Group.GroupKey, member.GetUserKey())); err != nil {
			return err
		}
		if err := ownerCli.JoinGroup(ctx, keys.NewMembershipKey(output.Group.GroupKey, member.GetUserKey())); err != nil {
			return err
		}
	}
	return nil
}

func testResource(ctx context.Context, cli client.Client, out *res.GetResourceResponse, groups ...keys.GroupKeyGetter) error {
	return cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo(), groups...).AsRequest(), out)
}

func testGroup(t *testing.T, owner *models.UserSession, members ...*models.UserSession) (*readmodels.GroupReadModel, error) {
	var response handler.GetGroupResponse
	if err := testGroup2(t, owner, &response, members...); err != nil {
		return nil, err
	}
	return response.Group, nil
}
