package store

import (
	"fmt"
	"github.com/commonpool/backend/pkg/group/domain"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Membership struct {
	GroupId        string `mapstructure:"groupId"`
	UserId         string `mapstructure:"userId"`
	IsMember       bool   `mapstructure:"isMember"`
	IsAdmin        bool   `mapstructure:"isAdmin"`
	IsOwner        bool   `mapstructure:"isOwner"`
	GroupConfirmed bool   `mapstructure:"groupConfirmed"`
	UserConfirmed  bool   `mapstructure:"userConfirmed"`
}

func mapMembership(record *neo4j.Record, key string) (*domain.Membership, error) {

	graphMembership := Membership{}
	field, ok := record.Get(key)
	if !ok {
		return nil, fmt.Errorf("could not get field " + key)
	}
	relationship, _ := field.(neo4j.Relationship)
	err := mapstructure.Decode(relationship.Props, &graphMembership)
	if err != nil {
		return nil, err
	}

	groupKey, err := keys.ParseGroupKey(graphMembership.GroupId)
	if err != nil {
		return nil, err
	}
	userKey := keys.NewUserKey(graphMembership.UserId)

	membershipKey := keys.NewMembershipKey(groupKey, userKey)

	var permission domain.PermissionLevel
	if graphMembership.IsOwner {
		permission = domain.Owner
	} else if graphMembership.IsAdmin {
		permission = domain.Admin
	} else if graphMembership.IsMember {
		permission = domain.Member
	}

	var status domain.MembershipStatus
	if graphMembership.UserConfirmed && graphMembership.GroupConfirmed {
		status = domain.ApprovedMembershipStatus
	} else if graphMembership.UserConfirmed {
		status = domain.PendingGroupMembershipStatus
	} else if graphMembership.GroupConfirmed {
		status = domain.PendingUserMembershipStatus
	}

	return &domain.Membership{
		Key:             membershipKey,
		PermissionLevel: permission,
		Status:          status,
	}, nil

}
