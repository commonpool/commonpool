package store

import (
	"fmt"
	"github.com/commonpool/backend/model"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
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

func mapMembership(record neo4j.Record, key string) (*group2.Membership, error) {

	graphMembership := Membership{}
	field, ok := record.Get(key)
	if !ok {
		return nil, fmt.Errorf("could not get field " + key)
	}
	relationship, _ := field.(neo4j.Relationship)
	err := mapstructure.Decode(relationship.Props(), &graphMembership)
	if err != nil {
		return nil, err
	}

	groupKey, err := model.ParseGroupKey(graphMembership.GroupId)
	if err != nil {
		return nil, err
	}
	userKey := model.NewUserKey(graphMembership.UserId)

	membershipKey := model.NewMembershipKey(groupKey, userKey)

	return &group2.Membership{
		Key:            membershipKey,
		IsMember:       graphMembership.IsMember,
		IsAdmin:        graphMembership.IsAdmin,
		IsOwner:        graphMembership.IsOwner,
		GroupConfirmed: graphMembership.GroupConfirmed,
		UserConfirmed:  graphMembership.UserConfirmed,
	}, nil

}
