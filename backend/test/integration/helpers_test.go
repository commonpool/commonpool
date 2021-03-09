package integration

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/group/handler"
	"io/ioutil"
	"strconv"
	"testing"
)

func (s *IntegrationTestSuite) testUser(t *testing.T) (*models.UserSession, func()) {
	t.Helper()

	s.createUserLock.Lock()

	u := s.NewUser()
	upsertError := s.server.User.Store.Upsert(u.GetUserKey(), u.Email, u.Username)

	var userXchangeErr error
	if upsertError != nil {
		_, userXchangeErr = s.server.ChatService.CreateUserExchange(context.TODO(), u.GetUserKey())
	}

	s.createUserLock.Unlock()

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

func (s *IntegrationTestSuite) testGroup(t *testing.T, owner *models.UserSession, members ...*models.UserSession) (*handler.Group, error) {

	ctx := context.Background()
	s.groupCounter++

	response, httpResponse := s.CreateGroup(t, ctx, owner, &handler.CreateGroupRequest{
		Name:        "group-" + strconv.Itoa(s.groupCounter),
		Description: "group-" + strconv.Itoa(s.groupCounter),
	})
	if !AssertStatusCreated(t, httpResponse) {
		bytes, bytesErr := ioutil.ReadAll(httpResponse.Body)
		if bytesErr != nil {
			return nil, bytesErr
		}
		return nil, fmt.Errorf(string(bytes))
	}

	for _, member := range members {
		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, owner, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  member.GetUserKey(),
			GroupKey: response.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			bytes, bytesErr := ioutil.ReadAll(httpResponse.Body)
			if bytesErr != nil {
				return nil, bytesErr
			}
			return nil, fmt.Errorf(string(bytes))
		}
		_, httpResponse = s.CreateOrAcceptInvitation(t, ctx, member, &handler.CreateOrAcceptInvitationRequest{
			UserKey:  member.GetUserKey(),
			GroupKey: response.Group.GroupKey,
		})
		if !AssertStatusAccepted(t, httpResponse) {
			bytes, bytesErr := ioutil.ReadAll(httpResponse.Body)
			if bytesErr != nil {
				return nil, bytesErr
			}
			return nil, fmt.Errorf(string(bytes))
		}
	}

	return response.Group, nil
}
