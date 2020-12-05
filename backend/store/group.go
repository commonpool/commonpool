package store

import (
	"context"
	errs "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/graph"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
	"time"
)

type GroupStore struct {
	graphDriver graph.GraphDriver
}

var _ group.Store = &GroupStore{}

func NewGroupStore(graphDriver graph.GraphDriver) *GroupStore {
	return &GroupStore{
		graphDriver: graphDriver,
	}
}

func (g *GroupStore) GetGroups(take int, skip int) (*group.Groups, int64, error) {

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

	var groups []*group.Group

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

	return group.NewGroups(groups), *count, nil

}

func mapRecordToGroup(record neo4j.Record, key string) (*group.Group, error) {
	field, _ := record.Get(key)
	node := field.(neo4j.Node)
	graphGroup := GraphGroup{}
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

func mapGraphGroupToGroup(graphGroup *GraphGroup) (*group.Group, error) {

	groupKey, err := group.ParseGroupKey(graphGroup.ID)
	if err != nil {
		return nil, err
	}
	userKey := model.NewUserKey(graphGroup.CreatedBy)

	return &group.Group{
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
	description string) (*group.Group, *group.Membership, error) {

	ctx, _ = GetCtx(ctx, "GroupStore", "CreateGroupAndMembership")

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

	membership, err := MapGraphMembership(record, "m")

	return grp, membership, nil
}

func (g *GroupStore) GetGroup(ctx context.Context, groupKey model.GroupKey) (*group.Group, error) {

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
		return nil, errs.ErrGroupNotFound
	}

	grp, err := mapRecordToGroup(result.Record(), "g")
	if err != nil {
		return nil, err
	}

	return grp, nil

}

func (g *GroupStore) GetGroupsByKeys(ctx context.Context, groupKeys *model.GroupKeys) (*group.Groups, error) {

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
	var groups []*group.Group
	for result.Next() {
		grp, err := mapRecordToGroup(result.Record(), "g")
		if err != nil {
			return nil, err
		}
		groups = append(groups, grp)
	}

	return group.NewGroups(groups), nil
}

func (g *GroupStore) CreateMembership(ctx context.Context, membershipKey model.MembershipKey, isMember bool, isAdmin bool, isOwner bool, isDeactivated bool, groupConfirmed bool, userConfirmed bool) (*group.Membership, error) {

	ctx, _ = GetCtx(ctx, "GroupStore", "CreateMembership")

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
		return nil, errs.ErrUserOrGroupNotFound
	}
	return MapGraphMembership(result.Record(), "m")

}

func (g *GroupStore) MarkInvitationAsAccepted(ctx context.Context, membershipKey model.MembershipKey, decisionFrom group.MembershipParty) error {

	ctx, l := GetCtx(ctx, "GroupStore", "MarkInvitationAsAccepted")

	l = l.With(zap.Object("membership", membershipKey))
	l.Debug("marking invitation as accepted")

	var cyphers = []string{`MATCH (u:User {id:$userId})-[m:IsMemberOf]->(g:Group {id:$groupId})`}

	if decisionFrom == group.PartyUser {
		cyphers = append(cyphers, "SET m += {userConfirmed: true, isMember: m.groupConfirmed}")
	} else if decisionFrom == group.PartyGroup {
		cyphers = append(cyphers, "SET m += {groupConfirmed: true, isMember: m.userConfirmed}")
	} else {
		return errs.ErrUnknownParty
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

func (g *GroupStore) GetMembershipsForUser(ctx context.Context, userKey model.UserKey, membershipStatus *group.MembershipStatus) (*group.Memberships, error) {

	ctx, _ = GetCtx(ctx, "GroupStore", "GetMembershipsForUser")

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

	var memberships []*group.Membership
	for result.Next() {
		membership, err := MapGraphMembership(result.Record(), "m")
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, membership)
	}

	return group.NewMemberships(memberships), nil

}

func getMembershipStatusWhereClauses(membershipStatus *group.MembershipStatus, wheres []string) []string {
	if membershipStatus != nil {
		if *membershipStatus == group.ApprovedMembershipStatus {
			wheres = append(wheres, "m.groupConfirmed = true")
			wheres = append(wheres, "m.userConfirmed = true")
		} else if *membershipStatus == group.PendingGroupMembershipStatus {
			wheres = append(wheres, "m.groupConfirmed = false")
			wheres = append(wheres, "m.userConfirmed = true")
		} else if *membershipStatus == group.PendingUserMembershipStatus {
			wheres = append(wheres, "m.groupConfirmed = true")
			wheres = append(wheres, "m.userConfirmed = false")
		} else if *membershipStatus == group.PendingStatus {
			wheres = append(wheres, "(m.groupConfirmed = false OR m.userConfirmed = false)")
		}
	}
	return wheres
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

	ctx, _ = GetCtx(ctx, "GroupStore", "GetMembershipsForGroup")

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

	var memberships []*group.Membership
	for result.Next() {
		membership, err := MapGraphMembership(result.Record(), "m")
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, membership)
	}

	return group.NewMemberships(memberships), nil

}

func (g *GroupStore) GetMembership(ctx context.Context, membershipKey model.MembershipKey) (*group.Membership, error) {

	ctx, _ = GetCtx(ctx, "GroupStore", "GetMembership")

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
		return nil, errs.ErrMembershipNotFound
	}

	return MapGraphMembership(result.Record(), "m")

}

func (g *GroupStore) DeleteMembership(ctx context.Context, membershipKey model.MembershipKey) error {

	ctx, _ = GetCtx(ctx, "GroupStore", "GetMembership")

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
		return errs.ErrMembershipNotFound
	}

	return nil

}
