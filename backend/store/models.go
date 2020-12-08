package store

import (
	"fmt"
	"github.com/commonpool/backend/chat"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key"`
	ChannelID      string
	MessageType    chat.MessageType
	MessageSubType chat.MessageSubType
	SentById       string
	SentByUsername string
	SentAt         time.Time
	Text           string
	Blocks         string `gorm:"type:jsonb"`
	Attachments    string `gorm:"type:jsonb"`
	VisibleToUser  *string
}

type GraphGroup struct {
	ID          string    `mapstructure:"id"`
	CreatedAt   time.Time `mapstructure:"createdAt"`
	Description string    `mapstructure:"description"`
	CreatedBy   string    `mapstructure:"createdBy"`
	Name        string    `mapstructure:"name"`
}

type GraphMembership struct {
	GroupId        string `mapstructure:"groupId"`
	UserId         string `mapstructure:"userId"`
	IsMember       bool   `mapstructure:"isMember"`
	IsAdmin        bool   `mapstructure:"isAdmin"`
	IsOwner        bool   `mapstructure:"isOwner"`
	GroupConfirmed bool   `mapstructure:"groupConfirmed"`
	UserConfirmed  bool   `mapstructure:"userConfirmed"`
}

func MapGraphMembership(record neo4j.Record, key string) (*group.Membership, error) {

	graphMembership := GraphMembership{}
	field, ok := record.Get(key)
	if !ok {
		return nil, fmt.Errorf("could not get field " + key)
	}
	relationship, _ := field.(neo4j.Relationship)
	err := mapstructure.Decode(relationship.Props(), &graphMembership)
	if err != nil {
		return nil, err
	}

	groupKey, err := group.ParseGroupKey(graphMembership.GroupId)
	if err != nil {
		return nil, err
	}
	userKey := model.NewUserKey(graphMembership.UserId)

	membershipKey := model.NewMembershipKey(groupKey, userKey)

	return &group.Membership{
		Key:            membershipKey,
		IsMember:       graphMembership.IsMember,
		IsAdmin:        graphMembership.IsAdmin,
		IsOwner:        graphMembership.IsOwner,
		GroupConfirmed: graphMembership.GroupConfirmed,
		UserConfirmed:  graphMembership.UserConfirmed,
	}, nil

}

type GraphResource struct {
	ID               string           `mapstructure:"id"`
	CreatedAt        time.Time        `mapstructure:"createdAt"`
	UpdatedAt        time.Time        `mapstructure:"updatedAt"`
	DeletedAt        *time.Time       `mapstructure:"deletedAt"`
	Summary          string           `mapstructure:"summary"`
	Description      string           `mapstructure:"description"`
	CreatedBy        string           `mapstructure:"createdBy"`
	Type             resource.Type    `mapstructure:"type"`
	SubType          resource.SubType `mapstructure:"subType"`
	ValueInHoursFrom int              `mapstructure:"valueInHoursFrom"`
	ValueInHoursTo   int              `mapstructure:"valueInHoursTo"`
}
