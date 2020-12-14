// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/user"
	"sync"
)

// Ensure, that UserStore does implement user.Store.
// If this is not the case, regenerate this file with moq.
var _ user.Store = &UserStore{}

// UserStore is a mock implementation of user.Store.
//
//     func TestSomethingThatUsesStore(t *testing.T) {
//
//         // make and configure a mocked user.Store
//         mockedStore := &UserStore{
//             FindFunc: func(query user.Query) ([]*user.User, error) {
// 	               panic("mock out the Find method")
//             },
//             GetByKeyFunc: func(key model.UserKey) (*user.User, error) {
// 	               panic("mock out the GetByKey method")
//             },
//             GetByKeysFunc: func(ctx context.Context, keys []model.UserKey) (*user.Users, error) {
// 	               panic("mock out the GetByKeys method")
//             },
//             GetUsernameFunc: func(key model.UserKey) (string, error) {
// 	               panic("mock out the GetUsername method")
//             },
//             UpsertFunc: func(key model.UserKey, email string, username string) error {
// 	               panic("mock out the Upsert method")
//             },
//         }
//
//         // use mockedStore in code that requires user.Store
//         // and then make assertions.
//
//     }
type UserStore struct {
	// FindFunc mocks the Find method.
	FindFunc func(query user.Query) ([]*user.User, error)

	// GetByKeyFunc mocks the GetByKey method.
	GetByKeyFunc func(key model.UserKey) (*user.User, error)

	// GetByKeysFunc mocks the GetByKeys method.
	GetByKeysFunc func(ctx context.Context, keys []model.UserKey) (*user.Users, error)

	// GetUsernameFunc mocks the GetUsername method.
	GetUsernameFunc func(key model.UserKey) (string, error)

	// UpsertFunc mocks the Upsert method.
	UpsertFunc func(key model.UserKey, email string, username string) error

	// calls tracks calls to the methods.
	calls struct {
		// Find holds details about calls to the Find method.
		Find []struct {
			// Query is the query argument value.
			Query user.Query
		}
		// GetByKey holds details about calls to the GetByKey method.
		GetByKey []struct {
			// Key is the key argument value.
			Key model.UserKey
		}
		// GetByKeys holds details about calls to the GetByKeys method.
		GetByKeys []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Keys is the keys argument value.
			Keys []model.UserKey
		}
		// GetUsername holds details about calls to the GetUsername method.
		GetUsername []struct {
			// Key is the key argument value.
			Key model.UserKey
		}
		// Upsert holds details about calls to the Upsert method.
		Upsert []struct {
			// Key is the key argument value.
			Key model.UserKey
			// Email is the email argument value.
			Email string
			// Username is the username argument value.
			Username string
		}
	}
	lockFind        sync.RWMutex
	lockGetByKey    sync.RWMutex
	lockGetByKeys   sync.RWMutex
	lockGetUsername sync.RWMutex
	lockUpsert      sync.RWMutex
}

// Find calls FindFunc.
func (mock *UserStore) Find(query user.Query) ([]*user.User, error) {
	if mock.FindFunc == nil {
		panic("UserStore.FindFunc: method is nil but Store.Find was just called")
	}
	callInfo := struct {
		Query user.Query
	}{
		Query: query,
	}
	mock.lockFind.Lock()
	mock.calls.Find = append(mock.calls.Find, callInfo)
	mock.lockFind.Unlock()
	return mock.FindFunc(query)
}

// FindCalls gets all the calls that were made to Find.
// Check the length with:
//     len(mockedStore.FindCalls())
func (mock *UserStore) FindCalls() []struct {
	Query user.Query
} {
	var calls []struct {
		Query user.Query
	}
	mock.lockFind.RLock()
	calls = mock.calls.Find
	mock.lockFind.RUnlock()
	return calls
}

// GetByKey calls GetByKeyFunc.
func (mock *UserStore) GetByKey(key model.UserKey) (*user.User, error) {
	if mock.GetByKeyFunc == nil {
		panic("UserStore.GetByKeyFunc: method is nil but Store.GetByKey was just called")
	}
	callInfo := struct {
		Key model.UserKey
	}{
		Key: key,
	}
	mock.lockGetByKey.Lock()
	mock.calls.GetByKey = append(mock.calls.GetByKey, callInfo)
	mock.lockGetByKey.Unlock()
	return mock.GetByKeyFunc(key)
}

// GetByKeyCalls gets all the calls that were made to GetByKey.
// Check the length with:
//     len(mockedStore.GetByKeyCalls())
func (mock *UserStore) GetByKeyCalls() []struct {
	Key model.UserKey
} {
	var calls []struct {
		Key model.UserKey
	}
	mock.lockGetByKey.RLock()
	calls = mock.calls.GetByKey
	mock.lockGetByKey.RUnlock()
	return calls
}

// GetByKeys calls GetByKeysFunc.
func (mock *UserStore) GetByKeys(ctx context.Context, keys []model.UserKey) (*user.Users, error) {
	if mock.GetByKeysFunc == nil {
		panic("UserStore.GetByKeysFunc: method is nil but Store.GetByKeys was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Keys []model.UserKey
	}{
		Ctx:  ctx,
		Keys: keys,
	}
	mock.lockGetByKeys.Lock()
	mock.calls.GetByKeys = append(mock.calls.GetByKeys, callInfo)
	mock.lockGetByKeys.Unlock()
	return mock.GetByKeysFunc(ctx, keys)
}

// GetByKeysCalls gets all the calls that were made to GetByKeys.
// Check the length with:
//     len(mockedStore.GetByKeysCalls())
func (mock *UserStore) GetByKeysCalls() []struct {
	Ctx  context.Context
	Keys []model.UserKey
} {
	var calls []struct {
		Ctx  context.Context
		Keys []model.UserKey
	}
	mock.lockGetByKeys.RLock()
	calls = mock.calls.GetByKeys
	mock.lockGetByKeys.RUnlock()
	return calls
}

// GetUsername calls GetUsernameFunc.
func (mock *UserStore) GetUsername(key model.UserKey) (string, error) {
	if mock.GetUsernameFunc == nil {
		panic("UserStore.GetUsernameFunc: method is nil but Store.GetUsername was just called")
	}
	callInfo := struct {
		Key model.UserKey
	}{
		Key: key,
	}
	mock.lockGetUsername.Lock()
	mock.calls.GetUsername = append(mock.calls.GetUsername, callInfo)
	mock.lockGetUsername.Unlock()
	return mock.GetUsernameFunc(key)
}

// GetUsernameCalls gets all the calls that were made to GetUsername.
// Check the length with:
//     len(mockedStore.GetUsernameCalls())
func (mock *UserStore) GetUsernameCalls() []struct {
	Key model.UserKey
} {
	var calls []struct {
		Key model.UserKey
	}
	mock.lockGetUsername.RLock()
	calls = mock.calls.GetUsername
	mock.lockGetUsername.RUnlock()
	return calls
}

// Upsert calls UpsertFunc.
func (mock *UserStore) Upsert(key model.UserKey, email string, username string) error {
	if mock.UpsertFunc == nil {
		panic("UserStore.UpsertFunc: method is nil but Store.Upsert was just called")
	}
	callInfo := struct {
		Key      model.UserKey
		Email    string
		Username string
	}{
		Key:      key,
		Email:    email,
		Username: username,
	}
	mock.lockUpsert.Lock()
	mock.calls.Upsert = append(mock.calls.Upsert, callInfo)
	mock.lockUpsert.Unlock()
	return mock.UpsertFunc(key, email, username)
}

// UpsertCalls gets all the calls that were made to Upsert.
// Check the length with:
//     len(mockedStore.UpsertCalls())
func (mock *UserStore) UpsertCalls() []struct {
	Key      model.UserKey
	Email    string
	Username string
} {
	var calls []struct {
		Key      model.UserKey
		Email    string
		Username string
	}
	mock.lockUpsert.RLock()
	calls = mock.calls.Upsert
	mock.lockUpsert.RUnlock()
	return calls
}
