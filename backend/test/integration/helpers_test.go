package integration

import (
	"context"
	"github.com/commonpool/backend/auth"
	"github.com/commonpool/backend/web"
	"net/http"
	"strconv"
	"testing"
)

func testUser(t *testing.T) (*auth.UserSession, func()) {
	t.Helper()

	createUserLock.Lock()

	user := NewUser()
	upsertError := AuthStore.Upsert(user.GetUserKey(), user.Email, user.Username)

	var userXchangeErr error
	if upsertError != nil {
		_, userXchangeErr = ChatService.CreateUserExchange(context.TODO(), user.GetUserKey())
	}

	createUserLock.Unlock()

	if upsertError != nil {
		t.Fatalf("upsert error: %s", upsertError)
	}

	if userXchangeErr != nil {
		t.Fatalf("exchange error: %s", userXchangeErr)
	}

	return user, func(user *auth.UserSession) func() {
		return func() {
			createUserLock.Lock()
			_ = Db.Delete(auth.User{}, "id = ?", user.Subject).Error
			createUserLock.Unlock()
		}
	}(user)
}

var groupCounter = 0

func testGroup(t *testing.T, owner *auth.UserSession, members ...*auth.UserSession) *web.Group {
	ctx := context.Background()
	groupCounter++
	response, httpResponse := CreateGroup(t, ctx, owner, &web.CreateGroupRequest{
		Name:        "group-" + strconv.Itoa(groupCounter),
		Description: "group-" + strconv.Itoa(groupCounter),
	})
	if httpResponse.StatusCode != http.StatusCreated {
		t.Fatalf("could not create group")
	}

	for _, member := range members {
		CreateOrAcceptInvitation(t, ctx, owner, &web.CreateOrAcceptInvitationRequest{
			UserID:  member.Subject,
			GroupID: response.Group.ID,
		})
		CreateOrAcceptInvitation(t, ctx, member, &web.CreateOrAcceptInvitationRequest{
			UserID:  member.Subject,
			GroupID: response.Group.ID,
		})
	}

	return response.Group
}
