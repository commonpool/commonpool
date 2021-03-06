package handler

import (
	"github.com/commonpool/backend/pkg/group/readmodels"
	"time"
)

type Membership struct {
	UserID         string    `json:"userId"`
	GroupID        string    `json:"groupId"`
	IsAdmin        bool      `json:"isAdmin"`
	IsMember       bool      `json:"isMember"`
	IsOwner        bool      `json:"isOwner"`
	GroupConfirmed bool      `json:"groupConfirmed"`
	UserConfirmed  bool      `json:"userConfirmed"`
	CreatedAt      time.Time `json:"createdAt"`
	IsDeactivated  bool      `json:"isDeactivated"`
	GroupName      string    `json:"groupName"`
	UserName       string    `json:"userName"`
}

func NewMembership(membership *readmodels.MembershipReadModel) Membership {
	return Membership{
		UserID:         membership.UserKey,
		GroupID:        membership.GroupKey,
		IsAdmin:        membership.IsAdmin,
		IsMember:       membership.IsMember,
		IsOwner:        membership.IsOwner,
		GroupConfirmed: membership.GroupConfirmed,
		UserConfirmed:  membership.UserConfirmed,
		// TODO: CreatedAt:      membership.CreatedAt,
		GroupName: membership.GroupName,
		UserName:  membership.UserName,
	}
}
