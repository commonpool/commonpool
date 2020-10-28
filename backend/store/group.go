package store

import (
	"fmt"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"gorm.io/gorm"
)

type GroupStore struct {
	db *gorm.DB
}

var _ group.Store = &GroupStore{}

func (g *GroupStore) CreateGroup(request group.CreateGroupRequest) group.CreateGroupResponse {
	grp := model.Group{
		ID:          request.GroupKey.ID,
		CreatedBy:   request.CreatedBy.String(),
		Name:        request.Name,
		Description: request.Description,
	}
	err := g.db.Create(&grp).Error
	if err != nil {
		return group.CreateGroupResponse{
			Error: err,
		}
	}

	// Create membership for creator
	mbs := model.Membership{
		GroupID:        grp.ID,
		UserID:         request.CreatedBy.String(),
		IsMember:       true,
		IsAdmin:        true,
		IsOwner:        true,
		GroupConfirmed: true,
		UserConfirmed:  true,
		IsDeactivated:  false,
	}
	err = g.db.Create(&mbs).Error
	if err != nil {
		return group.CreateGroupResponse{
			Error: err,
		}
	}

	return group.CreateGroupResponse{
		Error: err,
	}
}

func (g *GroupStore) GetGroup(request group.GetGroupRequest) group.GetGroupResult {
	var grp model.Group
	err := g.db.Where("id = ?", request.Key.ID.String()).First(&grp).Error
	if err != nil {
		return group.GetGroupResult{Error: err}
	}
	return group.GetGroupResult{
		Group: grp,
	}
}

func (g *GroupStore) GrantPermission(request group.GrantPermissionRequest) group.GrantPermissionResult {
	err := g.updatePermission(request.MembershipKey, request.Permission, true)
	if err != nil {
		return group.GrantPermissionResult{
			Error: err,
		}
	}
	return group.GrantPermissionResult{}
}

func (g *GroupStore) RevokePermission(request group.RevokePermissionRequest) group.RevokePermissionResult {
	err := g.updatePermission(request.MembershipKey, request.Permission, false)
	if err != nil {
		return group.RevokePermissionResult{
			Error: err,
		}
	}
	return group.RevokePermissionResult{}
}

func (g *GroupStore) updatePermission(membershipKey model.MembershipKey, permission model.PermissionType, hasPermission bool) error {
	updates := map[string]interface{}{}
	if permission == model.MemberPermission {
		updates["is_member"] = hasPermission
	} else if permission == model.AdminPermission {
		updates["is_admin"] = hasPermission
	}
	req := g.db.
		Model(model.Membership{}).
		Where("user_id = ? AND group_id = ?", membershipKey.UserKey.String(), membershipKey.GroupKey.ID.String()).
		Updates(updates)
	err := req.Error
	if err != nil {
		return err
	}
	if req.RowsAffected == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func (g *GroupStore) Invite(request group.InviteRequest) group.InviteResponse {
	membership := model.Membership{
		GroupID:        request.MembershipKey.GroupKey.ID,
		UserID:         request.MembershipKey.UserKey.String(),
		IsMember:       true,
		IsAdmin:        false,
		IsOwner:        false,
		IsDeactivated:  false,
		GroupConfirmed: request.InvitedBy == group.GroupParty,
		UserConfirmed:  request.InvitedBy == group.UserParty,
	}
	err := g.db.Create(&membership).Error
	return group.InviteResponse{
		Error: err,
	}
}

func (g *GroupStore) Exclude(request group.ExcludeRequest) group.ExcludeResponse {
	req := g.db.Delete(model.Membership{}, "group_id = ? AND user_id = ?", request.MembershipKey.GroupKey.ID.String(), request.MembershipKey.UserKey.String())
	err := req.Error
	if err != nil {
		return group.ExcludeResponse{
			Error: err,
		}
	}
	if req.RowsAffected == 0 {
		return group.ExcludeResponse{
			Error: fmt.Errorf("not found"),
		}
	}
	return group.ExcludeResponse{}
}

func (g *GroupStore) GetGroupPermissionsForUser(request group.GetMembershipPermissionsRequest) group.GetMembershipPermissionsResponse {
	var membership model.Membership
	err := g.db.First(&membership, "group_id = ? AND user_id = ?", request.MembershipKey.GroupKey.ID.String(), request.MembershipKey.UserKey.String()).Error
	if err != nil {
		return group.GetMembershipPermissionsResponse{
			Error: err,
		}
	}
	return group.GetMembershipPermissionsResponse{
		MembershipPermissions: group.MembershipPermissions{
			MembershipKey: request.MembershipKey,
			IsMember:      membership.IsMember,
			IsAdmin:       membership.IsAdmin,
		},
	}
}

func (g *GroupStore) MarkInvitationAsAccepted(request group.MarkInvitationAsAcceptedRequest) group.MarkInvitationAsAcceptedResponse {
	err := g.updateInvitationAcceptance(request.MembershipKey, true, request.From)
	if err != nil {
		return group.MarkInvitationAsAcceptedResponse{
			Error: err,
		}
	}
	return group.MarkInvitationAsAcceptedResponse{}
}

func (g *GroupStore) MarkInvitationAsDeclined(request group.MarkInvitationAsDeclinedRequest) group.MarkInvitationAsDeclinedResponse {
	err := g.updateInvitationAcceptance(request.MembershipKey, false, request.From)
	if err != nil {
		return group.MarkInvitationAsDeclinedResponse{
			Error: err,
		}
	}
	return group.MarkInvitationAsDeclinedResponse{}
}

func (g *GroupStore) updateInvitationAcceptance(membershipKey model.MembershipKey, isAccepted bool, decisionFrom group.MembershipParty) error {

	var statusQry string
	var updates = map[string]interface{}{}

	if decisionFrom == group.UserParty {
		statusQry = "user_confirmed = ?"
		updates["user_confirmed"] = isAccepted
	} else if decisionFrom == group.GroupParty {
		statusQry = "group_confirmed = ?"
		updates["group_confirmed"] = isAccepted
	} else {
		return fmt.Errorf("unexpected value for request.From")
	}

	qry := g.db.
		Model(model.Membership{}).
		Where("group_id = ? AND user_id = ? AND "+statusQry, membershipKey.GroupKey.ID.String(), membershipKey.UserKey.String(), !isAccepted).
		Updates(updates)

	err := qry.Error
	if err != nil {
		return err
	}

	if qry.RowsAffected == 0 {
		return fmt.Errorf("membership not found")
	}

	return nil
}

func (g *GroupStore) GetMembershipsForUser(request group.GetMembershipsForUserRequest) group.GetMembershipsForUserResponse {
	var memberships []model.Membership
	chain := g.db.Where("user_id = ?", request.UserKey.String())
	chain = g.filterMembershipStatus(chain, request.MembershipStatus)

	err := chain.Find(&memberships).Error
	return group.GetMembershipsForUserResponse{
		Error:       err,
		Memberships: memberships,
	}
}

func (g *GroupStore) filterMembershipStatus(chain *gorm.DB, membershipStatus *model.MembershipStatus) *gorm.DB {
	if membershipStatus != nil {
		if *membershipStatus == model.ApprovedMembershipStatus {
			chain = chain.Where("group_confirmed = true AND user_confirmed = true")
		} else if *membershipStatus == model.PendingGroupMembershipStatus {
			chain = chain.Where("group_confirmed = false AND user_confirmed = true")
		} else if *membershipStatus == model.PendingUserMembershipStatus {
			chain = chain.Where("group_confirmed = true AND user_confirmed = false")
		} else if *membershipStatus == model.PendingStatus {
			chain = chain.Where("group_confirmed = false OR user_confirmed = false")
		}
	}
	return chain
}

func (g *GroupStore) GetMembershipsForGroup(request group.GetMembershipsForGroupRequest) group.GetMembershipsForGroupResponse {
	var memberships []model.Membership
	chain := g.db.Where("group_id = ?", request.GroupKey.ID.String())
	chain = g.filterMembershipStatus(chain, request.MembershipStatus)
	err := chain.Find(&memberships).Error
	return group.GetMembershipsForGroupResponse{
		Error:       err,
		Memberships: memberships,
	}
}

func (g *GroupStore) GetMembership(request group.GetMembershipRequest) group.GetMembershipResponse {
	var membership model.Membership
	err := g.db.First(&membership, "user_id = ? AND group_id = ?",
		request.MembershipKey.UserKey.String(),
		request.MembershipKey.GroupKey.ID.String()).
		Error

	return group.GetMembershipResponse{
		Error:      err,
		Membership: membership,
	}
}

func (g *GroupStore) DeleteMembership(request group.DeleteMembershipRequest) group.DeleteMembershipResponse {
	req := g.db.Delete(&model.Membership{}, "user_id = ? AND group_id = ?", request.MembershipKey.UserKey.String(), request.MembershipKey.GroupKey.ID.String())
	err := req.Error
	if err != nil {
		return group.DeleteMembershipResponse{
			Error: err,
		}
	}
	if req.RowsAffected == 0 {
		return group.DeleteMembershipResponse{
			Error: fmt.Errorf("cannot delete invitation: not found"),
		}
	}
	return group.DeleteMembershipResponse{}
}

func NewGroupStore(db *gorm.DB) *GroupStore {
	return &GroupStore{
		db: db,
	}
}
