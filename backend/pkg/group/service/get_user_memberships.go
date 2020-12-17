package service

import (
	"context"
	group2 "github.com/commonpool/backend/pkg/group"
)

func (g GroupService) GetUserMemberships(ctx context.Context, request *group2.GetMembershipsForUserRequest) (*group2.GetMembershipsForUserResponse, error) {

	memberships, err := g.groupStore.GetMembershipsForUser(ctx, request.UserKey, request.MembershipStatus)
	if err != nil {
		return nil, err
	}

	return &group2.GetMembershipsForUserResponse{
		Memberships: memberships,
	}, nil
}
