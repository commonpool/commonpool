package integration

import (
	"context"
	"github.com/commonpool/backend/auth"
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
