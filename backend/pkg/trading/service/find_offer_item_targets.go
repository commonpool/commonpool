package service

import (
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading"
	"golang.org/x/net/context"
)

func (t TradingService) FindTargetsForOfferItem(
	ctx context.Context,
	groupKey keys.GroupKey,
	itemType trading.OfferItemType,
	from *trading.Target,
	to *trading.Target) (*trading.Targets, error) {

	membershipStatus := group2.ApprovedMembershipStatus
	membershipsForGroup, err := t.groupService.GetGroupMemberships(ctx, &group2.GetMembershipsForGroupRequest{
		GroupKey:         groupKey,
		MembershipStatus: &membershipStatus,
	})
	if err != nil {
		return nil, err
	}

	group, err := t.groupService.GetGroup(ctx, &group2.GetGroupRequest{
		Key: groupKey,
	})
	if err != nil {
		return nil, err
	}

	var targets []*trading.Target

	groupTarget := trading.NewGroupTarget(group.Group.Key)

	if to == nil || !to.Equals(groupTarget) {
		targets = append(targets, groupTarget)
	}

	for _, membership := range membershipsForGroup.Memberships.Items {
		userTarget := trading.NewUserTarget(membership.GetUserKey())
		if to == nil || !to.Equals(userTarget) {
			targets = append(targets, userTarget)
		}
	}

	return trading.NewTargets(targets), nil
}
