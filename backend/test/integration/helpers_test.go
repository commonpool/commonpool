package integration

import (
	"context"
	"github.com/commonpool/backend/pkg/auth"
	"github.com/commonpool/backend/pkg/group/handler"
	"net/http"
	"strconv"
	"testing"
)

func testUser(t *testing.T) (*auth.UserSession, func()) {
	t.Helper()

	createUserLock.Lock()

	u := NewUser()
	upsertError := AuthStore.Upsert(u.GetUserKey(), u.Email, u.Username)

	var userXchangeErr error
	if upsertError != nil {
		_, userXchangeErr = ChatService.CreateUserExchange(context.TODO(), u.GetUserKey())
	}

	createUserLock.Unlock()

	if upsertError != nil {
		t.Fatalf("upsert error: %s", upsertError)
	}

	if userXchangeErr != nil {
		t.Fatalf("exchange error: %s", userXchangeErr)
	}

	return u, func(user *auth.UserSession) func() {
		return func() {
			session := Driver.GetSession()
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

var groupCounter = 0

func testGroup(t *testing.T, owner *auth.UserSession, members ...*auth.UserSession) *handler.Group {
	ctx := context.Background()
	groupCounter++
	response, httpResponse, err := CreateGroup(t, ctx, owner, &handler.CreateGroupRequest{
		Name:        "group-" + strconv.Itoa(groupCounter),
		Description: "group-" + strconv.Itoa(groupCounter),
	})
	if err != nil {
		t.Fatal(err)
	}
	if httpResponse.StatusCode != http.StatusCreated {
		t.Fatalf("could not create group")
	}

	for _, member := range members {
		CreateOrAcceptInvitation(t, ctx, owner, &handler.CreateOrAcceptInvitationRequest{
			UserID:  member.Subject,
			GroupID: response.Group.ID,
		})
		CreateOrAcceptInvitation(t, ctx, member, &handler.CreateOrAcceptInvitationRequest{
			UserID:  member.Subject,
			GroupID: response.Group.ID,
		})
	}

	return response.Group
}
