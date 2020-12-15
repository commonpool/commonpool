package store

import (
	ctx "context"
	"fmt"
	"github.com/commonpool/backend/logging"
	"github.com/commonpool/backend/pkg/exceptions"
	graph2 "github.com/commonpool/backend/pkg/graph"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	"github.com/commonpool/backend/pkg/resource"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	"github.com/commonpool/backend/pkg/shared/store"
	transaction2 "github.com/commonpool/backend/pkg/transaction"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"strings"
	"time"
)

type ResourceStore struct {
	graphDriver        graph2.Driver
	transactionService transaction2.Service
}

var _ resource.Store = &ResourceStore{}

func NewResourceStore(graphDriver graph2.Driver, transactionService transaction2.Service) *ResourceStore {
	return &ResourceStore{
		graphDriver:        graphDriver,
		transactionService: transactionService,
	}
}

func (rs *ResourceStore) GetByKeys(ctx ctx.Context, resourceKeys *resourcemodel.ResourceKeys) (*resource.GetResourceByKeysResponse, error) {

	graphSession, err := rs.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}

	return rs.getByKeys(ctx, graphSession, resourceKeys)
}

func (rs *ResourceStore) getByKeys(ctx ctx.Context, session neo4j.Session, resourceKeys *resourcemodel.ResourceKeys) (*resource.GetResourceByKeysResponse, error) {

	getResult, err := session.Run(`
		MATCH (resource:Resource) 
		WHERE resource.id IN $ids
		
		WITH resource
		OPTIONAL MATCH (resource)-[:SharedWith]->(group:Group)

		WITH resource, collect(DISTINCT group.id) as sharedWithGroupIds
		OPTIONAL MATCH (resource)-[:SharedWith]->(viewer)

		WITH resource, sharedWithGroupIds, collect(distinct viewer) as viewers
		OPTIONAL MATCH (resource)-[:OwnedBy]->(owner)

		WITH resource, sharedWithGroupIds, viewers, collect(DISTINCT owner) as owners
		OPTIONAL MATCH (resource)-[:ManagedBy]->(manager)

		WITH resource, sharedWithGroupIds, viewers, owners, collect(DISTINCT manager) as managers
		RETURN resource, sharedWithGroupIds, viewers, owners, managers`,
		map[string]interface{}{
			"ids": resourceKeys.Strings(),
		})

	if err != nil {
		return nil, err
	}
	if getResult.Err() != nil {
		return nil, getResult.Err()
	}

	var resources []*resourcemodel.Resource

	sharings := resourcemodel.NewEmptyResourceSharings()
	claims := resourcemodel.NewEmptyClaims()

	for getResult.Next() {
		r, err := rs.mapGraphResourceRecord(getResult.Record(), "resource")
		if err != nil {
			return nil, err
		}

		sharingsForResource, err := rs.mapGraphSharingRecord(getResult.Record(), "resource", "sharedWithGroupIds")
		if err != nil {
			return nil, err
		}
		sharings.AppendAll(sharingsForResource)

		ownerTargets, err := store.MapTargets(getResult.Record(), "owners")
		if err != nil {
			return nil, err
		}
		claims.AppendAll(createClaimsForTargets(r.Key, resourcemodel.OwnershipClaim, ownerTargets))

		managerTargets, err := store.MapTargets(getResult.Record(), "managers")
		if err != nil {
			return nil, err
		}
		claims.AppendAll(createClaimsForTargets(r.Key, resourcemodel.ManagerClaim, managerTargets))

		viewerTargets, err := store.MapTargets(getResult.Record(), "viewers")
		if err != nil {
			return nil, err
		}
		claims.AppendAll(createClaimsForTargets(r.Key, resourcemodel.ViewerClaim, viewerTargets))

		resources = append(resources, r)
	}

	return &resource.GetResourceByKeysResponse{
		Sharings:  sharings,
		Resources: resourcemodel.NewResources(resources),
		Claims:    claims,
	}, nil

}

func createClaimsForTargets(resourceKey resourcemodel.ResourceKey, claimType resourcemodel.ClaimType, targets *resourcemodel.Targets) *resourcemodel.Claims {
	var claims []*resourcemodel.Claim
	for _, target := range targets.Items {
		claims = append(claims, &resourcemodel.Claim{
			ResourceKey: resourceKey,
			ClaimType:   claimType,
			For:         target,
		})
	}
	return resourcemodel.NewClaims(claims)
}

// GetByKey Gets a resource by key
func (rs *ResourceStore) GetByKey(ctx ctx.Context, getResourceByKeyQuery *resource.GetResourceByKeyQuery) (*resource.GetResourceByKeyResponse, error) {

	graphSession, err := rs.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}

	return rs.getByKey(ctx, graphSession, getResourceByKeyQuery)

}

func (rs *ResourceStore) getByKey(ctx ctx.Context, session neo4j.Session, getResourceByKeyQuery *resource.GetResourceByKeyQuery) (*resource.GetResourceByKeyResponse, error) {

	key := getResourceByKeyQuery.ResourceKey
	response, err := rs.GetByKeys(ctx, resourcemodel.NewResourceKeys([]resourcemodel.ResourceKey{key}))
	if err != nil {
		return nil, err
	}

	r, err := response.Resources.GetResource(key)
	if err != nil {
		return nil, err
	}

	return &resource.GetResourceByKeyResponse{
		Resource: r,
		Sharings: response.Sharings,
		Claims:   response.Claims,
	}, nil

}

// Delete deletes a resource
func (rs *ResourceStore) Delete(ctx context.Context, resourceKey resourcemodel.ResourceKey) error {

	graphSession, err := rs.graphDriver.GetSession()
	if err != nil {
		return err
	}

	graphSession.Run(`MATCH (r:Resource {id:$id}) DETACH DELETE`, map[string]interface{}{
		"id": resourceKey.String(),
	})

	now := time.Now()
	deleteResult, err := graphSession.Run(`
			MATCH (r:Resource {id:$id})
			SET r += {
				deletedAt:$deletedAt
			}
			RETURN r`,
		map[string]interface{}{
			"id":        resourceKey.String(),
			"deletedAt": now,
		})
	if err != nil {
		return err
	}
	if deleteResult.Err() != nil {
		return err
	}
	if !deleteResult.Next() {
		return err
	}

	return nil

}

// Create creates a resource
func (rs *ResourceStore) Create(ctx context.Context, createResourceQuery *resource.CreateResourceQuery) error {

	graphSession, err := rs.graphDriver.GetSession()
	if err != nil {
		return err
	}

	resourceKey := createResourceQuery.Resource.GetKey()
	now := time.Now().UTC().UnixNano() / 1e6

	params := map[string]interface{}{
		"userId":           createResourceQuery.Resource.CreatedBy,
		"id":               resourceKey.String(),
		"createdAt":        now,
		"updatedAt":        now,
		"deletedAt":        nil,
		"summary":          createResourceQuery.Resource.Summary,
		"description":      createResourceQuery.Resource.Description,
		"createdBy":        createResourceQuery.Resource.CreatedBy,
		"type":             createResourceQuery.Resource.Type,
		"subType":          createResourceQuery.Resource.SubType,
		"valueInHoursFrom": createResourceQuery.Resource.ValueInHoursFrom,
		"valueInHoursTo":   createResourceQuery.Resource.ValueInHoursTo,
	}

	cypher := `
			MATCH(u:User {id:$userId})
			CREATE (r:Resource {
				id:$id,
				createdAt:datetime({epochMillis:$createdAt}),
				updatedAt:datetime({epochMillis:$updatedAt}),
				deletedAt:null,
				summary:$summary,
				description:$description,
				createdBy:$createdBy,
				type:$type,
				subType:$subType,
				valueInHoursFrom:$valueInHoursFrom,
				valueInHoursTo:$valueInHoursTo
			})-[c:CreatedBy]->(u),
			(r)-[:OwnedBy]->(u)
			`

	if createResourceQuery.SharedWith != nil && len(createResourceQuery.SharedWith.Items) > 0 {
		params["groupIds"] = createResourceQuery.SharedWith.Strings()
		cypher = cypher + `

			WITH u, r

			CALL {
				
				WITH u, r
				MATCH (g:Group)
				WHERE g.id IN $groupIds
				CREATE (r)-[s:SharedWith {createdAt:datetime({epochMillis:$createdAt})}]->(g) 
				RETURN collect(distinct g.id) as groupIds
			
			}

			RETURN r, groupIds`
	} else {
		cypher = cypher + `
			RETURN r`
	}

	createResult, err := graphSession.Run(cypher, params)

	if err != nil {
		return err
	}
	if createResult.Err() != nil {
		return err
	}

	if !createResult.Next() {
		return fmt.Errorf("unexpected result count")

	}

	record := createResult.Record()

	groupIdsField, ok := record.Get("groupIds")
	if ok {
		groupIdsIntfs := groupIdsField.([]interface{})
		for _, groupIdIntf := range groupIdsIntfs {
			groupId := groupIdIntf.(string)
			groupKey, err := groupmodel.ParseGroupKey(groupId)
			if err != nil {
				return err
			}
			_, err = rs.transactionService.UserSharedResourceWithGroup(groupKey, createResourceQuery.Resource.Key)
			return err
		}
	}

	return nil

}

func (rs *ResourceStore) mapGraphResourceRecord(record neo4j.Record, key string) (*resourcemodel.Resource, error) {
	resourceRecord, _ := record.Get(key)
	node := resourceRecord.(neo4j.Node)
	return MapResourceNode(node)
}

func IsResourceNode(node neo4j.Node) bool {
	return graph2.NodeHasLabel(node, "Resource")
}

func MapResourceNode(node neo4j.Node) (*resourcemodel.Resource, error) {
	var graphResource = Resource{}
	err := mapstructure.Decode(node.Props(), &graphResource)
	if err != nil {
		return nil, err
	}
	mappedResource, err := mapGraphResourceToResource(&graphResource)
	if err != nil {
		return nil, err
	}
	return mappedResource, nil
}

func (rs *ResourceStore) mapGraphSharingRecord(record neo4j.Record, resourceFieldKey string, groupIdsFieldKey string) (*resourcemodel.Sharings, error) {
	resourceField, _ := record.Get(resourceFieldKey)
	resourceNode := resourceField.(neo4j.Node)
	resourceId := resourceNode.Props()["id"].(string)
	resourceKey, err := resourcemodel.ParseResourceKey(resourceId)
	if err != nil {
		return nil, err
	}

	groupIdsField, _ := record.Get(groupIdsFieldKey)
	if groupIdsField == nil {
		return resourcemodel.NewEmptyResourceSharings(), nil
	}
	groupIds := groupIdsField.([]interface{})
	var sharings []*resourcemodel.Sharing
	for _, groupId := range groupIds {
		groupKey, err := groupmodel.ParseGroupKey(groupId.(string))
		if err != nil {
			return nil, err
		}
		sharing := &resourcemodel.Sharing{
			ResourceKey: resourceKey,
			GroupKey:    groupKey,
		}
		sharings = append(sharings, sharing)
	}

	return resourcemodel.NewResourceSharings(sharings), nil
}

// Update updates a resource
func (rs *ResourceStore) Update(ctx context.Context, request *resource.UpdateResourceQuery) error {

	session, err := rs.graphDriver.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	now := time.Now().UTC().UnixNano() / 1e6
	resourceKey := request.Resource.GetKey()
	updateResult, err := session.Run(`
			MATCH (resource:Resource {id:$id})

			WITH resource

			OPTIONAL MATCH (resource)-[notSharedWithRel:SharedWith]->(notSharedWith:Group)
			WHERE NOT (notSharedWith.id IN $groupIds)		
			
			WITH resource, notSharedWithRel, notSharedWith

			OPTIONAL MATCH (sharedWithGroup:Group)
			WHERE sharedWithGroup.id IN $groupIds and NOT (resource)-[:SharedWith]->(sharedWithGroup)

			WITH resource, notSharedWithRel, notSharedWith, sharedWithGroup
	
			WITH 
				resource,
				notSharedWithRel,
				notSharedWith,
				sharedWithGroup,
				collect(distinct notSharedWith.id) as deletedSharingGroupIds,
				collect(distinct sharedWithGroup.id) as createdSharingGroupIds

			SET resource += {
				updatedAt:datetime({epochMillis:$updatedAt}),
				summary:$summary,
				description:$description,
				valueInHoursFrom:$valueInHoursFrom,
				valueInHoursTo:$valueInHoursTo
			}

			WITH resource, sharedWithGroup, notSharedWithRel, createdSharingGroupIds, deletedSharingGroupIds

			call apoc.do.when(
				notSharedWithRel IS NOT NULL,
				'DELETE notSharedWithRel RETURN "bla" as a',	
				'',
				{notSharedWithRel: notSharedWithRel}
			)
			YIELD value as a

			WITH resource, sharedWithGroup, createdSharingGroupIds, deletedSharingGroupIds

			call apoc.do.when(
				sharedWithGroup IS NOT NULL, 
				'CREATE (resource)-[sharedWith:SharedWith]->(sharedWithGroup) RETURN "bla" as b',
				'',
				{resource: resource, sharedWithGroup: sharedWithGroup})
			YIELD value as b

			RETURN resource, createdSharingGroupIds, deletedSharingGroupIds`,
		map[string]interface{}{
			"id":               resourceKey.String(),
			"updatedAt":        now,
			"summary":          request.Resource.Summary,
			"description":      request.Resource.Description,
			"valueInHoursFrom": request.Resource.ValueInHoursFrom,
			"valueInHoursTo":   request.Resource.ValueInHoursTo,
			"groupIds":         request.SharedWith.Strings(),
		})

	if err != nil {
		return err
	}

	if updateResult.Err() != nil {
		return updateResult.Err()
	}

	if !updateResult.Next() {
		return exceptions.ErrResourceNotFound
	}

	deletedSharingField, _ := updateResult.Record().Get("deletedSharingGroupIds")
	deletedSharingIntfs := deletedSharingField.([]interface{})
	for _, deletedSharingIntf := range deletedSharingIntfs {
		groupId := deletedSharingIntf.(string)
		groupKey, err := groupmodel.ParseGroupKey(groupId)
		if err != nil {
			return err
		}
		_, err = rs.transactionService.UserRemovedResourceFromGroup(groupKey, request.Resource.Key)
		if err != nil {
			return err
		}
	}

	createdSharingsField, _ := updateResult.Record().Get("createdSharingGroupIds")
	createdSharingIntfs := createdSharingsField.([]interface{})
	for _, createdSharingIntf := range createdSharingIntfs {
		groupId := createdSharingIntf.(string)
		groupKey, err := groupmodel.ParseGroupKey(groupId)
		if err != nil {
			return err
		}
		_, err = rs.transactionService.UserSharedResourceWithGroup(groupKey, request.Resource.Key)
		if err != nil {
			return err
		}
	}

	return nil

}

// Search search for resources
func (rs *ResourceStore) Search(context ctx.Context, request *resource.SearchResourcesQuery) (*resource.SearchResourcesResponse, error) {

	l := logging.WithContext(context)

	session, err := rs.graphDriver.GetSession()
	if err != nil {
		l.Error("could not get graph session", zap.Error(err))
		return nil, err
	}
	defer session.Close()

	propertyValues := map[string]interface{}{}
	var matchClauses = []string{
		"(r:Resource)",
	}
	var whereClauses []string
	var optionalMatchClauses []string

	if request.CreatedBy != "" {
		matchClauses = append(matchClauses, "(r)-[:CreatedBy]->(createdBy:User {id:$createdById})")
		propertyValues["createdById"] = request.CreatedBy
	}

	if request.Type != nil {
		whereClauses = append(whereClauses, "r.type = $type")
		propertyValues["type"] = *request.Type
	}

	if request.SubType != nil {
		whereClauses = append(whereClauses, "r.subType = $subType")
		propertyValues["subType"] = *request.SubType
	}

	if request.Query != nil && *request.Query != "" {
		whereClauses = append(whereClauses, "r.summary =~ $query")
		propertyValues["query"] = ".*" + *request.Query + ".*"
	}

	if request.SharedWithGroup != nil {
		matchClauses = append(matchClauses, "(r)-[:SharedWith]->(g:Group {id:$groupId})")
		propertyValues["groupId"] = request.SharedWithGroup.String()
	} else {
		optionalMatchClauses = append(optionalMatchClauses, "(r)-[:SharedWith]->(g:Group)")
	}

	var cyper = "MATCH "

	cyper = cyper + strings.Join(matchClauses, ",")

	if len(whereClauses) > 0 {
		cyper = cyper + "\nWHERE "
		cyper = cyper + strings.Join(whereClauses, " AND ")
	}

	if len(optionalMatchClauses) > 0 {
		cyper = cyper + "\nOPTIONAL MATCH"
		cyper = cyper + strings.Join(optionalMatchClauses, "\nOPTIONAL MATCH")
	}

	countCypher := cyper + `
RETURN count(r) as totalCount
`

	countResult, err := session.Run(countCypher, propertyValues)
	if err != nil {
		l.Error("could not execute count query", zap.Error(err))
		return nil, err
	}

	if countResult.Err() != nil {
		l.Error("could not execute count query", zap.Error(countResult.Err()))
		return nil, err
	}

	countResult.Next()

	countField, _ := countResult.Record().Get("totalCount")
	totalCount := countField.(int64)

	propertyValues["take"] = request.Take
	propertyValues["skip"] = request.Skip
	cyper = cyper + `
RETURN r, g
ORDER BY r.summary
SKIP $skip 
LIMIT $take
`

	searchResult, err := session.Run(cyper, propertyValues)

	if err != nil {
		l.Error("could not execute search query", zap.Error(err))
		return nil, err
	}

	if searchResult.Err() != nil {
		l.Error("could not execute search query", zap.Error(searchResult.Err()))
		return nil, err
	}

	var resources []*resourcemodel.Resource

	for searchResult.Next() {

		res, err := rs.mapGraphResourceRecord(searchResult.Record(), "r")
		if err != nil {
			l.Error("could not map resource record", zap.Error(err))
			return nil, err
		}

		resources = append(resources, res)

	}

	return &resource.SearchResourcesResponse{
		Resources:  resourcemodel.NewResources(resources),
		Sharings:   resourcemodel.NewEmptyResourceSharings(),
		Skip:       request.Skip,
		Take:       request.Take,
		TotalCount: int(totalCount),
	}, nil

}
