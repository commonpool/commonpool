// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource"
	"sync"
)

// Ensure, that ResourceStore does implement resource.Store.
// If this is not the case, regenerate this file with moq.
var _ resource.Store = &ResourceStore{}

// ResourceStore is a mock implementation of resource.Store.
//
//     func TestSomethingThatUsesStore(t *testing.T) {
//
//         // make and configure a mocked resource.Store
//         mockedStore := &ResourceStore{
//             CreateFunc: func(ctx context.Context, createResourceQuery *resource.CreateResourceQuery) error {
// 	               panic("mock out the Create method")
//             },
//             DeleteFunc: func(ctx context.Context, resourceKey model.ResourceKey) error {
// 	               panic("mock out the Delete method")
//             },
//             GetByKeyFunc: func(ctx context.Context, getResourceByKeyQuery *resource.GetResourceByKeyQuery) (*resource.GetResourceByKeyResponse, error) {
// 	               panic("mock out the GetByKey method")
//             },
//             GetByKeysFunc: func(ctx context.Context, resourceKeys *model.ResourceKeys) (*resource.GetResourceByKeysResponse, error) {
// 	               panic("mock out the GetByKeys method")
//             },
//             SearchFunc: func(ctx context.Context, searchResourcesQuery *resource.SearchResourcesQuery) (*resource.SearchResourcesResponse, error) {
// 	               panic("mock out the Search method")
//             },
//             UpdateFunc: func(ctx context.Context, updateResourceQuery *resource.UpdateResourceQuery) error {
// 	               panic("mock out the Update method")
//             },
//         }
//
//         // use mockedStore in code that requires resource.Store
//         // and then make assertions.
//
//     }
type ResourceStore struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(ctx context.Context, createResourceQuery *resource.CreateResourceQuery) error

	// DeleteFunc mocks the Delete method.
	DeleteFunc func(ctx context.Context, resourceKey keys.ResourceKey) error

	// GetByKeyFunc mocks the GetByKey method.
	GetByKeyFunc func(ctx context.Context, getResourceByKeyQuery *resource.GetResourceByKeyQuery) (*resource.GetResourceByKeyResponse, error)

	// GetByKeysFunc mocks the GetByKeys method.
	GetByKeysFunc func(ctx context.Context, resourceKeys *resource.ResourceKeys) (*resource.GetResourceByKeysResponse, error)

	// SearchFunc mocks the Search method.
	SearchFunc func(ctx context.Context, searchResourcesQuery *resource.SearchResourcesQuery) (*resource.SearchResourcesResponse, error)

	// UpdateFunc mocks the Update method.
	UpdateFunc func(ctx context.Context, updateResourceQuery *resource.UpdateResourceQuery) error

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// CreateResourceQuery is the createResourceQuery argument value.
			CreateResourceQuery *resource.CreateResourceQuery
		}
		// Delete holds details about calls to the Delete method.
		Delete []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ResourceKey is the resourceKey argument value.
			ResourceKey keys.ResourceKey
		}
		// GetByKey holds details about calls to the GetByKey method.
		GetByKey []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// GetResourceByKeyQuery is the getResourceByKeyQuery argument value.
			GetResourceByKeyQuery *resource.GetResourceByKeyQuery
		}
		// GetByKeys holds details about calls to the GetByKeys method.
		GetByKeys []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ResourceKeys is the resourceKeys argument value.
			ResourceKeys *resource.ResourceKeys
		}
		// Search holds details about calls to the Search method.
		Search []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// SearchResourcesQuery is the searchResourcesQuery argument value.
			SearchResourcesQuery *resource.SearchResourcesQuery
		}
		// Update holds details about calls to the Update method.
		Update []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UpdateResourceQuery is the updateResourceQuery argument value.
			UpdateResourceQuery *resource.UpdateResourceQuery
		}
	}
	lockCreate    sync.RWMutex
	lockDelete    sync.RWMutex
	lockGetByKey  sync.RWMutex
	lockGetByKeys sync.RWMutex
	lockSearch    sync.RWMutex
	lockUpdate    sync.RWMutex
}

// Create calls CreateFunc.
func (mock *ResourceStore) Create(ctx context.Context, createResourceQuery *resource.CreateResourceQuery) error {
	if mock.CreateFunc == nil {
		panic("ResourceStore.CreateFunc: method is nil but Store.Create was just called")
	}
	callInfo := struct {
		Ctx                 context.Context
		CreateResourceQuery *resource.CreateResourceQuery
	}{
		Ctx:                 ctx,
		CreateResourceQuery: createResourceQuery,
	}
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(ctx, createResourceQuery)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//     len(mockedStore.CreateCalls())
func (mock *ResourceStore) CreateCalls() []struct {
	Ctx                 context.Context
	CreateResourceQuery *resource.CreateResourceQuery
} {
	var calls []struct {
		Ctx                 context.Context
		CreateResourceQuery *resource.CreateResourceQuery
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
	return calls
}

// Delete calls DeleteFunc.
func (mock *ResourceStore) Delete(ctx context.Context, resourceKey keys.ResourceKey) error {
	if mock.DeleteFunc == nil {
		panic("ResourceStore.DeleteFunc: method is nil but Store.Delete was just called")
	}
	callInfo := struct {
		Ctx         context.Context
		ResourceKey keys.ResourceKey
	}{
		Ctx:         ctx,
		ResourceKey: resourceKey,
	}
	mock.lockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	mock.lockDelete.Unlock()
	return mock.DeleteFunc(ctx, resourceKey)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//     len(mockedStore.DeleteCalls())
func (mock *ResourceStore) DeleteCalls() []struct {
	Ctx         context.Context
	ResourceKey keys.ResourceKey
} {
	var calls []struct {
		Ctx         context.Context
		ResourceKey keys.ResourceKey
	}
	mock.lockDelete.RLock()
	calls = mock.calls.Delete
	mock.lockDelete.RUnlock()
	return calls
}

// GetByKey calls GetByKeyFunc.
func (mock *ResourceStore) GetByKey(ctx context.Context, getResourceByKeyQuery *resource.GetResourceByKeyQuery) (*resource.GetResourceByKeyResponse, error) {
	if mock.GetByKeyFunc == nil {
		panic("ResourceStore.GetByKeyFunc: method is nil but Store.GetByKey was just called")
	}
	callInfo := struct {
		Ctx                   context.Context
		GetResourceByKeyQuery *resource.GetResourceByKeyQuery
	}{
		Ctx:                   ctx,
		GetResourceByKeyQuery: getResourceByKeyQuery,
	}
	mock.lockGetByKey.Lock()
	mock.calls.GetByKey = append(mock.calls.GetByKey, callInfo)
	mock.lockGetByKey.Unlock()
	return mock.GetByKeyFunc(ctx, getResourceByKeyQuery)
}

// GetByKeyCalls gets all the calls that were made to GetByKey.
// Check the length with:
//     len(mockedStore.GetByKeyCalls())
func (mock *ResourceStore) GetByKeyCalls() []struct {
	Ctx                   context.Context
	GetResourceByKeyQuery *resource.GetResourceByKeyQuery
} {
	var calls []struct {
		Ctx                   context.Context
		GetResourceByKeyQuery *resource.GetResourceByKeyQuery
	}
	mock.lockGetByKey.RLock()
	calls = mock.calls.GetByKey
	mock.lockGetByKey.RUnlock()
	return calls
}

// GetByKeys calls GetByKeysFunc.
func (mock *ResourceStore) GetByKeys(ctx context.Context, resourceKeys *resource.ResourceKeys) (*resource.GetResourceByKeysResponse, error) {
	if mock.GetByKeysFunc == nil {
		panic("ResourceStore.GetByKeysFunc: method is nil but Store.GetByKeys was just called")
	}
	callInfo := struct {
		Ctx          context.Context
		ResourceKeys *resource.ResourceKeys
	}{
		Ctx:          ctx,
		ResourceKeys: resourceKeys,
	}
	mock.lockGetByKeys.Lock()
	mock.calls.GetByKeys = append(mock.calls.GetByKeys, callInfo)
	mock.lockGetByKeys.Unlock()
	return mock.GetByKeysFunc(ctx, resourceKeys)
}

// GetByKeysCalls gets all the calls that were made to GetByKeys.
// Check the length with:
//     len(mockedStore.GetByKeysCalls())
func (mock *ResourceStore) GetByKeysCalls() []struct {
	Ctx          context.Context
	ResourceKeys *resource.ResourceKeys
} {
	var calls []struct {
		Ctx          context.Context
		ResourceKeys *resource.ResourceKeys
	}
	mock.lockGetByKeys.RLock()
	calls = mock.calls.GetByKeys
	mock.lockGetByKeys.RUnlock()
	return calls
}

// Search calls SearchFunc.
func (mock *ResourceStore) Search(ctx context.Context, searchResourcesQuery *resource.SearchResourcesQuery) (*resource.SearchResourcesResponse, error) {
	if mock.SearchFunc == nil {
		panic("ResourceStore.SearchFunc: method is nil but Store.Search was just called")
	}
	callInfo := struct {
		Ctx                  context.Context
		SearchResourcesQuery *resource.SearchResourcesQuery
	}{
		Ctx:                  ctx,
		SearchResourcesQuery: searchResourcesQuery,
	}
	mock.lockSearch.Lock()
	mock.calls.Search = append(mock.calls.Search, callInfo)
	mock.lockSearch.Unlock()
	return mock.SearchFunc(ctx, searchResourcesQuery)
}

// SearchCalls gets all the calls that were made to Search.
// Check the length with:
//     len(mockedStore.SearchCalls())
func (mock *ResourceStore) SearchCalls() []struct {
	Ctx                  context.Context
	SearchResourcesQuery *resource.SearchResourcesQuery
} {
	var calls []struct {
		Ctx                  context.Context
		SearchResourcesQuery *resource.SearchResourcesQuery
	}
	mock.lockSearch.RLock()
	calls = mock.calls.Search
	mock.lockSearch.RUnlock()
	return calls
}

// Update calls UpdateFunc.
func (mock *ResourceStore) Update(ctx context.Context, updateResourceQuery *resource.UpdateResourceQuery) error {
	if mock.UpdateFunc == nil {
		panic("ResourceStore.UpdateFunc: method is nil but Store.Update was just called")
	}
	callInfo := struct {
		Ctx                 context.Context
		UpdateResourceQuery *resource.UpdateResourceQuery
	}{
		Ctx:                 ctx,
		UpdateResourceQuery: updateResourceQuery,
	}
	mock.lockUpdate.Lock()
	mock.calls.Update = append(mock.calls.Update, callInfo)
	mock.lockUpdate.Unlock()
	return mock.UpdateFunc(ctx, updateResourceQuery)
}

// UpdateCalls gets all the calls that were made to Update.
// Check the length with:
//     len(mockedStore.UpdateCalls())
func (mock *ResourceStore) UpdateCalls() []struct {
	Ctx                 context.Context
	UpdateResourceQuery *resource.UpdateResourceQuery
} {
	var calls []struct {
		Ctx                 context.Context
		UpdateResourceQuery *resource.UpdateResourceQuery
	}
	mock.lockUpdate.RLock()
	calls = mock.calls.Update
	mock.lockUpdate.RUnlock()
	return calls
}
