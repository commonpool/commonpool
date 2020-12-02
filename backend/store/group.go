package store

import (
	"context"
	"errors"
	"fmt"
	"github.com/commonpool/backend/amqp"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

type GroupStore struct {
	db *gorm.DB
	mq amqp.Client
}

var _ group.Store = &GroupStore{}

func NewGroupStore(db *gorm.DB, mq amqp.Client) *GroupStore {
	return &GroupStore{
		db: db,
		mq: mq,
	}
}

func (g *GroupStore) GetGroups(take int, skip int) ([]group.Group, int64, error) {

	var groups []group.Group

	var count int64
	qry := *g.db.Model(group.Group{})

	err := qry.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = qry.Limit(take).
		Offset(skip).
		Find(&groups).
		Error

	if err != nil {
		return nil, 0, err
	}

	return groups, count, nil

}

func (g *GroupStore) CreateGroup(ctx context.Context, groupKey model.GroupKey, createdBy model.UserKey, name string, description string) (*group.Group, error) {

	ctx, _ = GetCtx(ctx, "GroupStore", "CreateGroup")

	grp := group.Group{
		ID:          groupKey.ID,
		CreatedBy:   createdBy.String(),
		Name:        name,
		Description: description,
	}
	err := g.db.Create(&grp).Error
	if err != nil {
		return nil, err
	}

	return &grp, nil

}

func (g *GroupStore) GetGroup(ctx context.Context, groupKey model.GroupKey) (*group.Group, error) {
	var grp group.Group

	ctx, l := GetCtx(ctx, "GroupStore", "GetGroup")

	l = l.With(zap.Object("group", groupKey))
	l.Debug("getting group")

	err := g.db.Where("id = ?", groupKey.ID.String()).First(&grp).Error

	if err != nil {
		l.Error("could not get group", zap.Error(err))
		return nil, err
	}

	return &grp, nil
}

func (g *GroupStore) GetGroupsByKeys(ctx context.Context, groupKeys []model.GroupKey) (*group.Groups, error) {

	ctx, l := GetCtx(ctx, "GroupStore", "GetGroupsByKeys")

	var result []group.Group

	err := utils.Partition(len(groupKeys), 999, func(i1 int, i2 int) error {

		l.Debug(fmt.Sprintf("getting groups from index %d - %d", i1, i2))

		qryParts := make([]string, i2-i1)
		qryParams := make([]interface{}, i2-i1)

		for i, item := range groupKeys[i1:i2] {
			qryParts[i] = "?"
			qryParams[i] = item.ID.String()
		}

		qry := "id in (" + strings.Join(qryParts, ",") + ")"

		var partition []group.Group
		qryResult := g.db.Model(group.Group{}).Where(qry, qryParams...).Find(&partition)

		if qryResult.Error != nil {
			l.Error("could not get groups partition", zap.Error(qryResult.Error))
			return qryResult.Error
		}

		for _, item := range partition {
			result = append(result, item)
		}

		return nil

	})

	if err != nil {
		l.Error("could not get groups", zap.Error(err))
		return nil, err
	}

	return group.NewGroups(result), nil
}

func (g *GroupStore) CreateMembership(ctx context.Context, membershipKey model.MembershipKey, isMember bool, isAdmin bool, isOwner bool, isDeactivated bool, groupConfirmed bool, userConfirmed bool) (*group.Membership, error) {

	ctx, l := GetCtx(ctx, "GroupStore", "CreateMembership")

	l.Debug("creating membership")

	membership := group.Membership{
		GroupID:        membershipKey.GroupKey.ID,
		UserID:         membershipKey.UserKey.String(),
		IsMember:       isMember,
		IsAdmin:        isAdmin,
		IsOwner:        isOwner,
		IsDeactivated:  isDeactivated,
		GroupConfirmed: groupConfirmed,
		UserConfirmed:  userConfirmed,
	}
	err := g.db.Create(&membership).Error

	if err != nil {
		l.Error("could not create membership", zap.Error(err))
		return nil, err
	}

	return &membership, nil

}

func (g *GroupStore) MarkInvitationAsAccepted(ctx context.Context, membershipKey model.MembershipKey, decisionFrom group.MembershipParty) error {

	ctx, l := GetCtx(ctx, "GroupStore", "MarkInvitationAsAccepted")

	l = l.With(zap.Object("membership", membershipKey))
	l.Debug("marking invitation as accepted")

	err := g.updateInvitationAcceptance(ctx, membershipKey, true, decisionFrom)
	if err != nil {
		l.Error("could not update invitation acceptance status", zap.Error(err))
		return err
	}

	return nil
}

func (g *GroupStore) updateInvitationAcceptance(ctx context.Context, membershipKey model.MembershipKey, isAccepted bool, decisionFrom group.MembershipParty) error {

	ctx, l := GetCtx(ctx, "GroupStore", "updateInvitationAcceptance")

	l.Debug("updating invitation acceptance")

	var statusQry string
	var updates = map[string]interface{}{}

	if decisionFrom == group.PartyUser {
		statusQry = "user_confirmed = ?"
		updates["user_confirmed"] = isAccepted
	} else if decisionFrom == group.PartyGroup {
		statusQry = "group_confirmed = ?"
		updates["group_confirmed"] = isAccepted
	} else {
		err := fmt.Errorf("unexpected value for request.From")
		l.Error(err.Error())
		return err
	}

	qry := g.db.
		Model(group.Membership{}).
		Where("group_id = ? AND user_id = ? AND "+statusQry, membershipKey.GroupKey.ID.String(), membershipKey.UserKey.String(), !isAccepted).
		Updates(updates)

	err := qry.Error
	if err != nil {
		l.Error("could not update membership", zap.Error(err))
		return err
	}

	if qry.RowsAffected == 0 {
		err := fmt.Errorf("membership not found")
		l.Error(err.Error())
		return err
	}

	var membership = group.Membership{}
	err = g.db.Model(group.Membership{}).
		Where("group_id = ? AND user_id = ?", membershipKey.GroupKey.ID.String(), membershipKey.UserKey.String()).
		First(&membership).
		Error

	if membership.GroupConfirmed && membership.UserConfirmed {
		err = g.db.Model(&membership).Update("is_member", true).Error
		if err != nil {
			return err
		}
	}

	if err != nil {
		l.Error("could not get membership", zap.Error(err))
		return err
	}

	return nil
}

func (g *GroupStore) GetMembershipsForUser(ctx context.Context, userKey model.UserKey, membershipStatus *group.MembershipStatus) (*group.Memberships, error) {

	ctx, l := GetCtx(ctx, "GroupStore", "GetMembershipsForUser")

	l = l.With(zap.Object("user", userKey))

	var memberships []group.Membership
	chain := g.db.Where("user_id = ?", userKey.String())
	chain = g.filterMembershipStatus(chain, membershipStatus)
	err := chain.Find(&memberships).Error

	if err != nil {
		l.Error("could not get memberships for user", zap.Error(err))
		return nil, err
	}

	return group.NewMemberships(memberships), nil

}

func (g *GroupStore) filterMembershipStatus(chain *gorm.DB, membershipStatus *group.MembershipStatus) *gorm.DB {
	if membershipStatus != nil {
		if *membershipStatus == group.ApprovedMembershipStatus {
			chain = chain.Where("group_confirmed = true AND user_confirmed = true")
		} else if *membershipStatus == group.PendingGroupMembershipStatus {
			chain = chain.Where("group_confirmed = false AND user_confirmed = true")
		} else if *membershipStatus == group.PendingUserMembershipStatus {
			chain = chain.Where("group_confirmed = true AND user_confirmed = false")
		} else if *membershipStatus == group.PendingStatus {
			chain = chain.Where("group_confirmed = false OR user_confirmed = false")
		}
	}
	return chain
}

func (g *GroupStore) GetMembershipsForGroup(ctx context.Context, groupKey model.GroupKey, membershipStatus *group.MembershipStatus) (*group.Memberships, error) {

	ctx, l := GetCtx(ctx, "GroupStore", "GetMembershipsForGroup")

	l = l.With(zap.Object("group", groupKey))

	var memberships []group.Membership
	chain := g.db.Where("group_id = ?", groupKey.ID.String())
	chain = g.filterMembershipStatus(chain, membershipStatus)
	err := chain.Find(&memberships).Error

	if err != nil {
		l.Error("could not get memberships for group", zap.Error(err))
		return nil, err
	}

	return group.NewMemberships(memberships), nil

}

func (g *GroupStore) GetMembership(ctx context.Context, membershipKey model.MembershipKey) (*group.Membership, error) {

	ctx, l := GetCtx(ctx, "GroupStore", "GetMembership")
	l = l.With(zap.Object("membership", membershipKey))

	l.Debug("getting membership")

	var membership group.Membership

	err := g.db.First(&membership, "user_id = ? AND group_id = ?",
		membershipKey.UserKey.String(),
		membershipKey.GroupKey.ID.String()).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, group.ErrMembershipNotFound
	}

	if err != nil {
		l.Error("could not get membership", zap.Error(err))
		return nil, err
	}

	return &membership, nil
}

func (g *GroupStore) DeleteMembership(ctx context.Context, membershipKey model.MembershipKey) error {

	ctx, l := GetCtx(ctx, "GroupStore", "GetMembership")

	l = l.With(zap.Object("membership", membershipKey))

	l.Debug("deleting membership")

	req := g.db.Delete(&group.Membership{}, "user_id = ? AND group_id = ?",
		membershipKey.UserKey.String(),
		membershipKey.GroupKey.ID.String())

	err := req.Error
	if err != nil {
		l.Error("could not delete membership", zap.Error(err))
		return err
	}
	if req.RowsAffected == 0 {
		err := fmt.Errorf("membership not found")
		l.Error(err.Error())
		return err
	}

	return nil

}
