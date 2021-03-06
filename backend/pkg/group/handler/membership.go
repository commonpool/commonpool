package handler

import (
	"github.com/commonpool/backend/pkg/auth/models"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/domain"
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

func NewMembership(membership *domain.Membership, groupNames group.Names, names models.UserNames) Membership {
	return Membership{
		UserID:         membership.Key.UserKey.String(),
		GroupID:        membership.Key.GroupKey.String(),
		IsAdmin:        membership.IsAdmin(),
		IsMember:       membership.IsMember(),
		IsOwner:        membership.IsOwner(),
		GroupConfirmed: membership.HasGroupConfirmed(),
		UserConfirmed:  membership.HasUserConfirmed(),
		CreatedAt:      membership.CreatedAt,
		GroupName:      groupNames[membership.GetGroupKey()],
		UserName:       names[membership.GetUserKey()],
	}
}
