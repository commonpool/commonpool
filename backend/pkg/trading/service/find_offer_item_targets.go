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
	from *domain2.Target,
	to *domain2.Target) (*domain2.Targets, error) {

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

	var targets []*domain2.Target

	groupTarget := domain2.NewGroupTarget(groupKey)
	if to == nil || !to.Equals(groupTarget) {
		targets = append(targets, groupTarget)
	}

	for _, membership := range membershipsForGroup.Memberships {
		userKey := keys.NewUserKey(membership.UserKey)
		userTarget := domain2.NewUserTarget(userKey)
		if to == nil || !to.Equals(userTarget) {
			targets = append(targets, userTarget)
		}
	}

	return domain2.NewTargets(targets), nil
}
