package store

import (
	"context"
	"github.com/commonpool/backend/graph"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/exceptions"
	graph2 "github.com/commonpool/backend/pkg/graph"
	group2 "github.com/commonpool/backend/pkg/group"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"gorm.io/gorm"
	"strings"
	"time"
)

type GroupStore struct {
	graphDriver graph.Driver
}

var _ group2.Store = &GroupStore{}

func NewGroupStore(graphDriver graph.Driver) *GroupStore {
	return &GroupStore{
		graphDriver: graphDriver,
	}
}

func (g *GroupStore) GetGroups(take int, skip int) (*group2.Groups, int64, error) {

	session, err := g.graphDriver.GetSession()
	if err != nil {
		return nil, 0, err
	}
	defer session.Close()

	result, err := session.Run(`
		MATCH (g:Group)
		ORDER BY g.name
		SKIP $skip 
		LIMIT $take
		RETURN g`,
		map[string]interface{}{
			"take": take,
			"skip": skip,
		})

	if err != nil {
		return nil, 0, err
	}
	if result.Err() != nil {
		return nil, 0, result.Err()
	}

	var groups []*group2.Group

	for result.Next() {
		mappedGroup, err := mapRecordToGroup(result.Record(), "g")
		if err != nil {
			return nil, 0, err
		}
		groups = append(groups, mappedGroup)
	}

	countResult, err := session.Run("match (g:Group) RETURN count(g) as count", map[string]interface{}{})

	if err != nil {
		return nil, 0, err
	}
	if countResult.Err() != nil {
		return nil, 0, countResult.Err()
	}

	countResult.Next()
	countField, _ := countResult.Record().Get("count")
	count := countField.(*int64)

	return group2.NewGroups(groups), *count, nil

}

func mapRecordToGroup(record neo4j.Record, key string) (*group2.Group, error) {
	field, _ := record.Get(key)
	node := field.(neo4j.Node)
	return MapGroupNode(node)
}

func IsGroupNode(node neo4j.Node) bool {
	return graph2.NodeHasLabel(node, "Group")
}

func MapGroupNode(node neo4j.Node) (*group2.Group, error) {
	graphGroup := Group{}
	err := mapstructure.Decode(node.Props(), &graphGroup)
	if err != nil {
		return nil, err
	}
	mappedGroup, err := mapGraphGroupToGroup(&graphGroup)
	if err != nil {
		return nil, err
	}
	return mappedGroup, nil
}

func mapGraphGroupToGroup(graphGroup *Group) (*group2.Group, error) {

	groupKey, err := model.ParseGroupKey(graphGroup.ID)
	if err != nil {
		return nil, err
	}
	userKey := model.NewUserKey(graphGroup.CreatedBy)

	return &group2.Group{
		Key:         groupKey,
		CreatedBy:   userKey,
		CreatedAt:   graphGroup.CreatedAt,
		Name:        graphGroup.Name,
		Description: graphGroup.Description,
	}, nil

}

func (g *GroupStore) CreateGroupAndMembership(
	ctx context.Context,
	groupKey model.GroupKey,
	createdBy model.UserKey,
	name string,
	description string) (*group2.Group, *group2.Membership, error) {

	session, err := g.graphDriver.GetSession()
	if err != nil {
		return nil, nil, err
	}
	defer session.Close()

	now := time.Now().UTC().UnixNano() / 1e6
	result, err := session.Run(`
		MATCH (u:User {id:$userId})
		CREATE (g:Group {
			id:$id,
			createdBy:$userId,
			createdAt:datetime({epochMillis:$createdAt}),
			name:$name,
			description:$description
		})-[:CreatedBy]->(u),
		(g)<-[m:IsMemberOf{
			groupId:$id,	
			userId:$userId,
			isMember:true,
			isAdmin:true,
			isOwner:true,
			groupConfirmed:true,
			userConfirmed:true
		}]-(u)
		RETURN g, m`,
		map[string]interface{}{
			"userId":      createdBy.String(),
			"id":          groupKey.String(),
			"createdAt":   now,
			"name":        name,
			"description": description,
		})

	if err != nil {
		return nil, nil, err
	}
	if result.Err() != nil {
		return nil, nil, result.Err()
	}

	result.Next()

	record := result.Record()

	grp, err := mapRecordToGroup(record, "g")
	if err != nil {
		return nil, nil, err
	}

	membership, err := mapMembership(record, "m")

	return grp, membership, nil
}

func (g *GroupStore) GetGroup(ctx context.Context, groupKey model.GroupKey) (*group2.Group, error) {

	session, err := g.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`
		MATCH (g:Group {id:$id})
		RETURN g`,
		map[string]interface{}{
			"id": groupKey.String(),
		})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	if !result.Next() {
		return nil, exceptions.ErrGroupNotFound
	}

	grp, err := mapRecordToGroup(result.Record(), "g")
	if err != nil {
		return nil, err
	}

	return grp, nil

}

func (g *GroupStore) GetGroupsByKeys(ctx context.Context, groupKeys *model.GroupKeys) (*group2.Groups, error) {

	session, err := g.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`
		MATCH (g:Group)
		WHERE g.id IN $ids
		RETURN g`,
		map[string]interface{}{
			"ids": groupKeys.Strings(),
		})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	var groups []*group2.Group
	for result.Next() {
		grp, err := mapRecordToGroup(result.Record(), "g")
		if err != nil {
			return nil, err
		}
		groups = append(groups, grp)
	}

	return group2.NewGroups(groups), nil
}

func (g *GroupStore) CreateMembership(ctx context.Context, membershipKey model.MembershipKey, isMember bool, isAdmin bool, isOwner bool, isDeactivated bool, groupConfirmed bool, userConfirmed bool) (*group2.Membership, error) {

	session, err := g.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`
		MATCH (g:Group {id:$groupId}),(u:User {id:$userId})
		CREATE (g)<-[m:IsMemberOf {
			userId:$userId,
			groupId:$groupId,
			isMember:$isMember,
			isAdmin:$isAdmin,
			isOwner:$isOwner,
			groupConfirmed:$groupConfirmed,
			userConfirmed:$userConfirmed
		}]-(u)
		RETURN m`,
		map[string]interface{}{
			"groupId":        membershipKey.GroupKey.String(),
			"userId":         membershipKey.UserKey.String(),
			"isMember":       isMember,
			"isAdmin":        isAdmin,
			"isOwner":        isOwner,
			"groupConfirmed": groupConfirmed,
			"userConfirmed":  userConfirmed,
		})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	if !result.Next() {
		return nil, exceptions.ErrUserOrGroupNotFound
	}
	return mapMembership(result.Record(), "m")

}

func (g *GroupStore) MarkInvitationAsAccepted(ctx context.Context, membershipKey model.MembershipKey, decisionFrom group2.MembershipParty) error {

	var cyphers = []string{`MATCH (u:User {id:$userId})-[m:IsMemberOf]->(g:Group {id:$groupId})`}

	if decisionFrom == group2.PartyUser {
		cyphers = append(cyphers, "SET m += {userConfirmed: true, isMember: m.groupConfirmed}")
	} else if decisionFrom == group2.PartyGroup {
		cyphers = append(cyphers, "SET m += {groupConfirmed: true, isMember: m.userConfirmed}")
	} else {
		return exceptions.ErrUnknownParty
	}

	cyphers = append(cyphers, "RETURN m")

	session, err := g.graphDriver.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	result, err := session.Run(strings.Join(cyphers, "\n"), map[string]interface{}{
		"userId":  membershipKey.UserKey.String(),
		"groupId": membershipKey.GroupKey.String(),
	})

	if err != nil {
		return err
	}
	if result.Err() != nil {
		return err
	}

	return nil
}

func (g *GroupStore) GetMembershipsForUser(ctx context.Context, userKey model.UserKey, membershipStatus *group2.MembershipStatus) (*group2.Memberships, error) {

	session, err := g.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	cypher := "MATCH (u:User {id:$userId})-[m:IsMemberOf]->(g:Group)"
	var wheres []string
	wheres = getMembershipStatusWhereClauses(membershipStatus, wheres)
	if len(wheres) > 0 {
		cypher = cypher + "\nWHERE " + strings.Join(wheres, " AND ") + "\n"
	}
	cypher = cypher + "RETURN m"

	result, err := session.Run(cypher, map[string]interface{}{
		"userId": userKey.String(),
	})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	var memberships []*group2.Membership
	for result.Next() {
		membership, err := mapMembership(result.Record(), "m")
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, membership)
	}

	return group2.NewMemberships(memberships), nil

}

func getMembershipStatusWhereClauses(membershipStatus *group2.MembershipStatus, wheres []string) []string {
	if membershipStatus != nil {
		if *membershipStatus == group2.ApprovedMembershipStatus {
			wheres = append(wheres, "m.groupConfirmed = true")
			wheres = append(wheres, "m.userConfirmed = true")
		} else if *membershipStatus == group2.PendingGroupMembershipStatus {
			wheres = append(wheres, "m.groupConfirmed = false")
			wheres = append(wheres, "m.userConfirmed = true")
		} else if *membershipStatus == group2.PendingUserMembershipStatus {
			wheres = append(wheres, "m.groupConfirmed = true")
			wheres = append(wheres, "m.userConfirmed = false")
		} else if *membershipStatus == group2.PendingStatus {
			wheres = append(wheres, "(m.groupConfirmed = false OR m.userConfirmed = false)")
		}
	}
	return wheres
}

func (g *GroupStore) filterMembershipStatus(chain *gorm.DB, membershipStatus *group2.MembershipStatus) *gorm.DB {
	if membershipStatus != nil {
		if *membershipStatus == group2.ApprovedMembershipStatus {
			chain = chain.Where("group_confirmed = true AND user_confirmed = true")
		} else if *membershipStatus == group2.PendingGroupMembershipStatus {
			chain = chain.Where("group_confirmed = false AND user_confirmed = true")
		} else if *membershipStatus == group2.PendingUserMembershipStatus {
			chain = chain.Where("group_confirmed = true AND user_confirmed = false")
		} else if *membershipStatus == group2.PendingStatus {
			chain = chain.Where("group_confirmed = false OR user_confirmed = false")
		}
	}
	return chain
}

func (g *GroupStore) GetMembershipsForGroup(ctx context.Context, groupKey model.GroupKey, membershipStatus *group2.MembershipStatus) (*group2.Memberships, error) {

	session, err := g.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	cypher := "MATCH (u:User)-[m:IsMemberOf]->(g:Group {id:$groupId})"
	var wheres []string
	wheres = getMembershipStatusWhereClauses(membershipStatus, wheres)
	if len(wheres) > 0 {
		cypher = cypher + "\nWHERE " + strings.Join(wheres, " AND ") + "\n"
	}
	cypher = cypher + "RETURN m"

	result, err := session.Run(cypher, map[string]interface{}{
		"groupId": groupKey.String(),
	})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	var memberships []*group2.Membership
	for result.Next() {
		membership, err := mapMembership(result.Record(), "m")
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, membership)
	}

	return group2.NewMemberships(memberships), nil

}

func (g *GroupStore) GetMembership(ctx context.Context, membershipKey model.MembershipKey) (*group2.Membership, error) {

	session, err := g.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`	
		MATCH (u:User {id:$userId})-[m:IsMemberOf]-(g:Group {id:$groupId})
		RETURN m`,
		map[string]interface{}{
			"groupId": membershipKey.GroupKey.String(),
			"userId":  membershipKey.UserKey.String(),
		})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	if !result.Next() {
		return nil, exceptions.ErrMembershipNotFound
	}

	return mapMembership(result.Record(), "m")

}

func (g *GroupStore) DeleteMembership(ctx context.Context, membershipKey model.MembershipKey) error {

	session, err := g.graphDriver.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	result, err := session.Run(`	
		MATCH (u:User {id:$userId})-[m:IsMemberOf]-(g:Group {id:$groupId})
		DELETE m`,
		map[string]interface{}{
			"groupId": membershipKey.GroupKey.String(),
			"userId":  membershipKey.UserKey.String(),
		})

	if err != nil {
		return err
	}
	if result.Err() != nil {
		return result.Err()
	}

	summary, err := result.Summary()
	if err != nil {
		return err
	}

	if summary.Counters().RelationshipsDeleted() != 1 {
		return exceptions.ErrMembershipNotFound
	}

	return nil

}
