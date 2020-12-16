package service

import (
	group2 "github.com/commonpool/backend/pkg/group"
	model3 "github.com/commonpool/backend/pkg/resource/model"
	model2 "github.com/commonpool/backend/pkg/trading/model"
	"golang.org/x/net/context"
)

func (t TradingService) FindTargetsForOfferItem(
	ctx context.Context,
	groupKey group2.GroupKey,
	itemType model2.OfferItemType,
	from *model3.Target,
	to *model3.Target) (*model3.Targets, error) {

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

	var targets []*model3.Target

	groupTarget := model3.NewGroupTarget(group.Group.Key)

	if to == nil || !to.Equals(groupTarget) {
		targets = append(targets, groupTarget)
	}

	for _, membership := range membershipsForGroup.Memberships.Items {
		userTarget := model3.NewUserTarget(membership.GetUserKey())
		if to == nil || !to.Equals(userTarget) {
			targets = append(targets, userTarget)
		}
	}

	return model3.NewTargets(targets), nil
}
