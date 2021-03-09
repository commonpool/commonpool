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
	"golang.org/x/sync/errgroup"
	"strconv"
	"testing"
)

func (s *IntegrationTestBase) testUser(t *testing.T) (*models.UserSession, func()) {
	t.Helper()

	createUserLock.Lock()

	u := s.NewUser()
	upsertError := s.server.User.Store.Upsert(u.GetUserKey(), u.Email, u.Username)

	var userXchangeErr error
	if upsertError != nil {
		_, userXchangeErr = s.server.ChatService.CreateUserExchange(context.TODO(), u.GetUserKey())
	}

	createUserLock.Unlock()

	if upsertError != nil {
		t.Fatalf("upsert error: %s", upsertError)
	}

	if userXchangeErr != nil {
		t.Fatalf("exchange error: %s", userXchangeErr)
	}

	return u, func(user *models.UserSession) func() {
		return func() {
			session := s.server.GraphDriver.GetSession()
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

func (s *IntegrationTestBase) testUserCli(t *testing.T) (*models.UserSession, client.Client) {
	user, _ := s.testUser(t)
	return user, s.getUserClient(user)
}

func (s *IntegrationTestBase) getUserClient(user *models.UserSession) client.Client {
	return echo.NewEchoClient(s.server.Router, client.NewMockAuthentication(user))
}

func (s *IntegrationTestBase) testGroup2(t *testing.T, owner *models.UserSession, output *handler.GetGroupResponse, members ...*models.UserSession) error {
	ctx := context.Background()
	groupCounter++
	ownerCli := s.getUserClient(owner)
	if err := ownerCli.CreateGroup(ctx, handler.NewCreateGroupRequest("group-"+strconv.Itoa(groupCounter), "group-"+strconv.Itoa(groupCounter)), output); !assert.NoError(t, err) {
		return err
	}
	g, ctx := errgroup.WithContext(ctx)
	for _, member := range members {
		g.Go(func() error {
			return s.getUserClient(member).JoinGroup(ctx, keys.NewMembershipKey(output.Group.GroupKey, member.GetUserKey()))
		})
		g.Go(func() error {
			return ownerCli.JoinGroup(ctx, keys.NewMembershipKey(output.Group.GroupKey, member.GetUserKey()))
		})
	}
	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *IntegrationTestBase) testResource(ctx context.Context, cli client.Client, out *res.GetResourceResponse, groups ...keys.GroupKeyGetter) error {
	return cli.CreateResource(ctx, res.NewCreateResourcePayload(test.AResourceInfo(), groups...).AsRequest(), out)
}

func (s *IntegrationTestBase) testGroup(t *testing.T, owner *models.UserSession, members ...*models.UserSession) (*readmodels.GroupReadModel, error) {
	var response handler.GetGroupResponse
	if err := s.testGroup2(t, owner, &response, members...); err != nil {
		return nil, err
	}
	return response.Group, nil
}
