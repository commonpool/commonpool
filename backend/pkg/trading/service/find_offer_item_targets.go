package service

import (
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
	domain2 "github.com/commonpool/backend/pkg/trading/domain"
	"golang.org/x/net/context"
)

func (t TradingService) FindTargetsForOfferItem(
	ctx context.Context,
	groupKey keys.GroupKey,
	itemType domain2.OfferItemType,
	from *keys.Target,
	to *keys.Target) (*keys.Targets, error) {

	membershipStatus := domain.ApprovedMembershipStatus
	membershipsForGroup, err := t.groupService.GetGroupMemberships(ctx, &group2.GetMembershipsForGroupRequest{
		GroupKey:         groupKey,
		MembershipStatus: &membershipStatus,
	})
	if err != nil {
		return nil, err
	}

	_, err = t.groupService.GetGroup(ctx, groupKey)
	if err != nil {
		return nil, err
	}

	var targets []*keys.Target

	groupTarget := keys.NewGroupTarget(groupKey)
	if to == nil || !to.Equals(*groupTarget) {
		targets = append(targets, groupTarget)
	}

	for _, membership := range membershipsForGroup.Memberships {
		userKey := keys.NewUserKey(membership.UserKey)
		userTarget := keys.NewUserTarget(userKey)
		if to == nil || !to.Equals(*userTarget) {
			targets = append(targets, userTarget)
		}
	}

	return keys.NewTargets(targets), nil
}
