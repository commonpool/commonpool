// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mock

import (
	"context"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/keys"
	"sync"
)

// Ensure, that GroupService does implement group.Service.
// If this is not the case, regenerate this file with moq.
var _ group.Service = &GroupService{}

// GroupService is a mock implementation of group.Service.
//
//     func TestSomethingThatUsesService(t *testing.T) {
//
//         // make and configure a mocked group.Service
//         mockedService := &GroupService{
//             CancelOrDeclineInvitationFunc: func(ctx context.Context, request *group.CancelOrDeclineInvitationRequest) error {
// 	               panic("mock out the CancelOrDeclineInvitation method")
//             },
//             CreateGroupFunc: func(ctx context.Context, request *group.CreateGroupRequest) (*group.CreateGroupResponse, error) {
// 	               panic("mock out the CreateGroup method")
//             },
//             CreateOrAcceptInvitationFunc: func(ctx context.Context, request *group.CreateOrAcceptInvitationRequest) (*group.CreateOrAcceptInvitationResponse, error) {
// 	               panic("mock out the CreateOrAcceptInvitation method")
//             },
//             GetGroupFunc: func(ctx context.Context, request *group.GetGroupRequest) (*group.GetGroupResult, error) {
// 	               panic("mock out the GetGroup method")
//             },
//             GetGroupMembershipsFunc: func(ctx context.Context, request *group.GetMembershipsForGroupRequest) (*group.GetMembershipsForGroupResponse, error) {
// 	               panic("mock out the GetGroupMemberships method")
//             },
//             GetGroupsFunc: func(ctx context.Context, request *group.GetGroupsRequest) (*group.GetGroupsResult, error) {
// 	               panic("mock out the GetGroups method")
//             },
//             GetGroupsByKeysFunc: func(ctx context.Context, groupKeys *keys.GroupKeys) (*group.Groups, error) {
// 	               panic("mock out the GetGroupsByKeys method")
//             },
//             GetMembershipFunc: func(ctx context.Context, request *group.GetMembershipRequest) (*group.GetMembershipResponse, error) {
// 	               panic("mock out the GetMembership method")
//             },
//             GetUserMembershipsFunc: func(ctx context.Context, request *group.GetMembershipsForUserRequest) (*group.GetMembershipsForUserResponse, error) {
// 	               panic("mock out the GetUserMemberships method")
//             },
//         }
//
//         // use mockedService in code that requires group.Service
//         // and then make assertions.
//
//     }
type GroupService struct {
	// CancelOrDeclineInvitationFunc mocks the CancelOrDeclineInvitation method.
	CancelOrDeclineInvitationFunc func(ctx context.Context, request *group.CancelOrDeclineInvitationRequest) error

	// CreateGroupFunc mocks the CreateGroup method.
	CreateGroupFunc func(ctx context.Context, request *group.CreateGroupRequest) (*group.CreateGroupResponse, error)

	// CreateOrAcceptInvitationFunc mocks the CreateOrAcceptInvitation method.
	CreateOrAcceptInvitationFunc func(ctx context.Context, request *group.CreateOrAcceptInvitationRequest) (*group.CreateOrAcceptInvitationResponse, error)

	// GetGroupFunc mocks the GetGroup method.
	GetGroupFunc func(ctx context.Context, request *group.GetGroupRequest) (*group.GetGroupResult, error)

	// GetGroupMembershipsFunc mocks the GetGroupMemberships method.
	GetGroupMembershipsFunc func(ctx context.Context, request *group.GetMembershipsForGroupRequest) (*group.GetMembershipsForGroupResponse, error)

	// GetGroupsFunc mocks the GetGroups method.
	GetGroupsFunc func(ctx context.Context, request *group.GetGroupsRequest) (*group.GetGroupsResult, error)

	// GetGroupsByKeysFunc mocks the GetGroupsByKeys method.
	GetGroupsByKeysFunc func(ctx context.Context, groupKeys *keys.GroupKeys) (*group.Groups, error)

	// GetMembershipFunc mocks the GetMembership method.
	GetMembershipFunc func(ctx context.Context, request *group.GetMembershipRequest) (*group.GetMembershipResponse, error)

	// GetUserMembershipsFunc mocks the GetUserMemberships method.
	GetUserMembershipsFunc func(ctx context.Context, request *group.GetMembershipsForUserRequest) (*group.GetMembershipsForUserResponse, error)

	// calls tracks calls to the methods.
	calls struct {
		// CancelOrDeclineInvitation holds details about calls to the CancelOrDeclineInvitation method.
		CancelOrDeclineInvitation []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *group.CancelOrDeclineInvitationRequest
		}
		// CreateGroup holds details about calls to the CreateGroup method.
		CreateGroup []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *group.CreateGroupRequest
		}
		// CreateOrAcceptInvitation holds details about calls to the CreateOrAcceptInvitation method.
		CreateOrAcceptInvitation []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *group.CreateOrAcceptInvitationRequest
		}
		// GetGroup holds details about calls to the GetGroup method.
		GetGroup []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *group.GetGroupRequest
		}
		// GetGroupMemberships holds details about calls to the GetGroupMemberships method.
		GetGroupMemberships []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *group.GetMembershipsForGroupRequest
		}
		// GetGroups holds details about calls to the GetGroups method.
		GetGroups []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *group.GetGroupsRequest
		}
		// GetGroupsByKeys holds details about calls to the GetGroupsByKeys method.
		GetGroupsByKeys []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// GroupKeys is the groupKeys argument value.
			GroupKeys *keys.GroupKeys
		}
		// GetMembership holds details about calls to the GetMembership method.
		GetMembership []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *group.GetMembershipRequest
		}
		// GetUserMemberships holds details about calls to the GetUserMemberships method.
		GetUserMemberships []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Request is the request argument value.
			Request *group.GetMembershipsForUserRequest
		}
	}
	lockCancelOrDeclineInvitation sync.RWMutex
	lockCreateGroup               sync.RWMutex
	lockCreateOrAcceptInvitation  sync.RWMutex
	lockGetGroup                  sync.RWMutex
	lockGetGroupMemberships       sync.RWMutex
	lockGetGroups                 sync.RWMutex
	lockGetGroupsByKeys           sync.RWMutex
	lockGetMembership             sync.RWMutex
	lockGetUserMemberships        sync.RWMutex
}

// CancelOrDeclineInvitation calls CancelOrDeclineInvitationFunc.
func (mock *GroupService) CancelOrDeclineInvitation(ctx context.Context, request *group.CancelOrDeclineInvitationRequest) error {
	if mock.CancelOrDeclineInvitationFunc == nil {
		panic("GroupService.CancelOrDeclineInvitationFunc: method is nil but Service.CancelOrDeclineInvitation was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *group.CancelOrDeclineInvitationRequest
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockCancelOrDeclineInvitation.Lock()
	mock.calls.CancelOrDeclineInvitation = append(mock.calls.CancelOrDeclineInvitation, callInfo)
	mock.lockCancelOrDeclineInvitation.Unlock()
	return mock.CancelOrDeclineInvitationFunc(ctx, request)
}

// CancelOrDeclineInvitationCalls gets all the calls that were made to CancelOrDeclineInvitation.
// Check the length with:
//     len(mockedService.CancelOrDeclineInvitationCalls())
func (mock *GroupService) CancelOrDeclineInvitationCalls() []struct {
	Ctx     context.Context
	Request *group.CancelOrDeclineInvitationRequest
} {
	var calls []struct {
		Ctx     context.Context
		Request *group.CancelOrDeclineInvitationRequest
	}
	mock.lockCancelOrDeclineInvitation.RLock()
	calls = mock.calls.CancelOrDeclineInvitation
	mock.lockCancelOrDeclineInvitation.RUnlock()
	return calls
}

// CreateGroup calls CreateGroupFunc.
func (mock *GroupService) CreateGroup(ctx context.Context, request *group.CreateGroupRequest) (*group.CreateGroupResponse, error) {
	if mock.CreateGroupFunc == nil {
		panic("GroupService.CreateGroupFunc: method is nil but Service.CreateGroup was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *group.CreateGroupRequest
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockCreateGroup.Lock()
	mock.calls.CreateGroup = append(mock.calls.CreateGroup, callInfo)
	mock.lockCreateGroup.Unlock()
	return mock.CreateGroupFunc(ctx, request)
}

// CreateGroupCalls gets all the calls that were made to CreateGroup.
// Check the length with:
//     len(mockedService.CreateGroupCalls())
func (mock *GroupService) CreateGroupCalls() []struct {
	Ctx     context.Context
	Request *group.CreateGroupRequest
} {
	var calls []struct {
		Ctx     context.Context
		Request *group.CreateGroupRequest
	}
	mock.lockCreateGroup.RLock()
	calls = mock.calls.CreateGroup
	mock.lockCreateGroup.RUnlock()
	return calls
}

// CreateOrAcceptInvitation calls CreateOrAcceptInvitationFunc.
func (mock *GroupService) CreateOrAcceptInvitation(ctx context.Context, request *group.CreateOrAcceptInvitationRequest) (*group.CreateOrAcceptInvitationResponse, error) {
	if mock.CreateOrAcceptInvitationFunc == nil {
		panic("GroupService.CreateOrAcceptInvitationFunc: method is nil but Service.CreateOrAcceptInvitation was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *group.CreateOrAcceptInvitationRequest
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockCreateOrAcceptInvitation.Lock()
	mock.calls.CreateOrAcceptInvitation = append(mock.calls.CreateOrAcceptInvitation, callInfo)
	mock.lockCreateOrAcceptInvitation.Unlock()
	return mock.CreateOrAcceptInvitationFunc(ctx, request)
}

// CreateOrAcceptInvitationCalls gets all the calls that were made to CreateOrAcceptInvitation.
// Check the length with:
//     len(mockedService.CreateOrAcceptInvitationCalls())
func (mock *GroupService) CreateOrAcceptInvitationCalls() []struct {
	Ctx     context.Context
	Request *group.CreateOrAcceptInvitationRequest
} {
	var calls []struct {
		Ctx     context.Context
		Request *group.CreateOrAcceptInvitationRequest
	}
	mock.lockCreateOrAcceptInvitation.RLock()
	calls = mock.calls.CreateOrAcceptInvitation
	mock.lockCreateOrAcceptInvitation.RUnlock()
	return calls
}

// GetGroup calls GetGroupFunc.
func (mock *GroupService) GetGroup(ctx context.Context, request *group.GetGroupRequest) (*group.GetGroupResult, error) {
	if mock.GetGroupFunc == nil {
		panic("GroupService.GetGroupFunc: method is nil but Service.GetGroup was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *group.GetGroupRequest
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockGetGroup.Lock()
	mock.calls.GetGroup = append(mock.calls.GetGroup, callInfo)
	mock.lockGetGroup.Unlock()
	return mock.GetGroupFunc(ctx, request)
}

// GetGroupCalls gets all the calls that were made to GetGroup.
// Check the length with:
//     len(mockedService.GetGroupCalls())
func (mock *GroupService) GetGroupCalls() []struct {
	Ctx     context.Context
	Request *group.GetGroupRequest
} {
	var calls []struct {
		Ctx     context.Context
		Request *group.GetGroupRequest
	}
	mock.lockGetGroup.RLock()
	calls = mock.calls.GetGroup
	mock.lockGetGroup.RUnlock()
	return calls
}

// GetGroupMemberships calls GetGroupMembershipsFunc.
func (mock *GroupService) GetGroupMemberships(ctx context.Context, request *group.GetMembershipsForGroupRequest) (*group.GetMembershipsForGroupResponse, error) {
	if mock.GetGroupMembershipsFunc == nil {
		panic("GroupService.GetGroupMembershipsFunc: method is nil but Service.GetGroupMemberships was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *group.GetMembershipsForGroupRequest
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockGetGroupMemberships.Lock()
	mock.calls.GetGroupMemberships = append(mock.calls.GetGroupMemberships, callInfo)
	mock.lockGetGroupMemberships.Unlock()
	return mock.GetGroupMembershipsFunc(ctx, request)
}

// GetGroupMembershipsCalls gets all the calls that were made to GetGroupMemberships.
// Check the length with:
//     len(mockedService.GetGroupMembershipsCalls())
func (mock *GroupService) GetGroupMembershipsCalls() []struct {
	Ctx     context.Context
	Request *group.GetMembershipsForGroupRequest
} {
	var calls []struct {
		Ctx     context.Context
		Request *group.GetMembershipsForGroupRequest
	}
	mock.lockGetGroupMemberships.RLock()
	calls = mock.calls.GetGroupMemberships
	mock.lockGetGroupMemberships.RUnlock()
	return calls
}

// GetGroups calls GetGroupsFunc.
func (mock *GroupService) GetGroups(ctx context.Context, request *group.GetGroupsRequest) (*group.GetGroupsResult, error) {
	if mock.GetGroupsFunc == nil {
		panic("GroupService.GetGroupsFunc: method is nil but Service.GetGroups was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *group.GetGroupsRequest
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockGetGroups.Lock()
	mock.calls.GetGroups = append(mock.calls.GetGroups, callInfo)
	mock.lockGetGroups.Unlock()
	return mock.GetGroupsFunc(ctx, request)
}

// GetGroupsCalls gets all the calls that were made to GetGroups.
// Check the length with:
//     len(mockedService.GetGroupsCalls())
func (mock *GroupService) GetGroupsCalls() []struct {
	Ctx     context.Context
	Request *group.GetGroupsRequest
} {
	var calls []struct {
		Ctx     context.Context
		Request *group.GetGroupsRequest
	}
	mock.lockGetGroups.RLock()
	calls = mock.calls.GetGroups
	mock.lockGetGroups.RUnlock()
	return calls
}

// GetGroupsByKeys calls GetGroupsByKeysFunc.
func (mock *GroupService) GetGroupsByKeys(ctx context.Context, groupKeys *keys.GroupKeys) (*group.Groups, error) {
	if mock.GetGroupsByKeysFunc == nil {
		panic("GroupService.GetGroupsByKeysFunc: method is nil but Service.GetGroupsByKeys was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		GroupKeys *keys.GroupKeys
	}{
		Ctx:       ctx,
		GroupKeys: groupKeys,
	}
	mock.lockGetGroupsByKeys.Lock()
	mock.calls.GetGroupsByKeys = append(mock.calls.GetGroupsByKeys, callInfo)
	mock.lockGetGroupsByKeys.Unlock()
	return mock.GetGroupsByKeysFunc(ctx, groupKeys)
}

// GetGroupsByKeysCalls gets all the calls that were made to GetGroupsByKeys.
// Check the length with:
//     len(mockedService.GetGroupsByKeysCalls())
func (mock *GroupService) GetGroupsByKeysCalls() []struct {
	Ctx       context.Context
	GroupKeys *keys.GroupKeys
} {
	var calls []struct {
		Ctx       context.Context
		GroupKeys *keys.GroupKeys
	}
	mock.lockGetGroupsByKeys.RLock()
	calls = mock.calls.GetGroupsByKeys
	mock.lockGetGroupsByKeys.RUnlock()
	return calls
}

// GetMembership calls GetMembershipFunc.
func (mock *GroupService) GetMembership(ctx context.Context, request *group.GetMembershipRequest) (*group.GetMembershipResponse, error) {
	if mock.GetMembershipFunc == nil {
		panic("GroupService.GetMembershipFunc: method is nil but Service.GetMembership was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *group.GetMembershipRequest
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockGetMembership.Lock()
	mock.calls.GetMembership = append(mock.calls.GetMembership, callInfo)
	mock.lockGetMembership.Unlock()
	return mock.GetMembershipFunc(ctx, request)
}

// GetMembershipCalls gets all the calls that were made to GetMembership.
// Check the length with:
//     len(mockedService.GetMembershipCalls())
func (mock *GroupService) GetMembershipCalls() []struct {
	Ctx     context.Context
	Request *group.GetMembershipRequest
} {
	var calls []struct {
		Ctx     context.Context
		Request *group.GetMembershipRequest
	}
	mock.lockGetMembership.RLock()
	calls = mock.calls.GetMembership
	mock.lockGetMembership.RUnlock()
	return calls
}

// GetUserMemberships calls GetUserMembershipsFunc.
func (mock *GroupService) GetUserMemberships(ctx context.Context, request *group.GetMembershipsForUserRequest) (*group.GetMembershipsForUserResponse, error) {
	if mock.GetUserMembershipsFunc == nil {
		panic("GroupService.GetUserMembershipsFunc: method is nil but Service.GetUserMemberships was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		Request *group.GetMembershipsForUserRequest
	}{
		Ctx:     ctx,
		Request: request,
	}
	mock.lockGetUserMemberships.Lock()
	mock.calls.GetUserMemberships = append(mock.calls.GetUserMemberships, callInfo)
	mock.lockGetUserMemberships.Unlock()
	return mock.GetUserMembershipsFunc(ctx, request)
}

// GetUserMembershipsCalls gets all the calls that were made to GetUserMemberships.
// Check the length with:
//     len(mockedService.GetUserMembershipsCalls())
func (mock *GroupService) GetUserMembershipsCalls() []struct {
	Ctx     context.Context
	Request *group.GetMembershipsForUserRequest
} {
	var calls []struct {
		Ctx     context.Context
		Request *group.GetMembershipsForUserRequest
	}
	mock.lockGetUserMemberships.RLock()
	calls = mock.calls.GetUserMemberships
	mock.lockGetUserMemberships.RUnlock()
	return calls
}
