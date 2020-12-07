package store

import (
	ctx "context"
	errs "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/graph"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"strings"
	"time"
)

type ResourceStore struct {
	graphDriver graph.GraphDriver
}

var _ resource.Store = &ResourceStore{}

func NewResourceStore(graphDriver graph.GraphDriver) *ResourceStore {
	return &ResourceStore{
		graphDriver: graphDriver,
	}
}

func (rs *ResourceStore) GetByKeys(ctx ctx.Context, resourceKeys *model.ResourceKeys) (*resource.Resources, error) {

	graphSession, err := rs.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}

	return rs.getByKeys(ctx, graphSession, resourceKeys)
}

func (rs *ResourceStore) getByKeys(ctx ctx.Context, session neo4j.Session, resourceKeys *model.ResourceKeys) (*resource.Resources, error) {

	getResult, err := session.Run(`
		MATCH (r:Resource) 
		WHERE r.id IN $ids		
		OPTIONAL MATCH (r)-[s:SharedWith]->(g:Group)
		RETURN r, s, g`,
		map[string]interface{}{
			"ids": resourceKeys.Strings(),
		})

	if err != nil {
		return nil, err
	}
	if getResult.Err() != nil {
		return nil, getResult.Err()
	}

	var resources []*resource.Resource
	for getResult.Next() {
		r, err := rs.mapGraphResourceRecord(getResult.Record(), "r")
		if err != nil {
			return nil, err
		}
		resources = append(resources, r)
	}

	return resource.NewResources(resources), nil

}

// GetByKey Gets a resource by key
func (rs *ResourceStore) GetByKey(ctx ctx.Context, getResourceByKeyQuery *resource.GetResourceByKeyQuery) *resource.GetResourceByKeyResponse {

	graphSession, err := rs.graphDriver.GetSession()
	if err != nil {
		return &resource.GetResourceByKeyResponse{
			Error: err,
		}
	}

	return rs.getByKey(ctx, graphSession, getResourceByKeyQuery)

}

func (rs *ResourceStore) getByKey(ctx ctx.Context, session neo4j.Session, getResourceByKeyQuery *resource.GetResourceByKeyQuery) *resource.GetResourceByKeyResponse {

	getResult, err := session.Run(`
		MATCH (r:Resource {id : $id})
		OPTIONAL MATCH (r)-[s:SharedWith]->(g:Group)
		RETURN r, collect(g.id) as groupIds`,
		map[string]interface{}{
			"id": getResourceByKeyQuery.ResourceKey.String(),
		})

	if err != nil {
		return &resource.GetResourceByKeyResponse{
			Error: err,
		}
	}

	if getResult.Err() != nil {
		return &resource.GetResourceByKeyResponse{
			Error: getResult.Err(),
		}
	}

	if !getResult.Next() {
		return &resource.GetResourceByKeyResponse{
			Error: errs.ErrResourceNotFound,
		}
	}

	res, err := rs.mapGraphResourceRecord(getResult.Record(), "r")
	if err != nil {
		return &resource.GetResourceByKeyResponse{
			Error: err,
		}
	}

	sharings, err := rs.mapGraphSharingRecord(getResult.Record(), "r", "groupIds")
	if err != nil {
		return &resource.GetResourceByKeyResponse{
			Error: err,
		}
	}

	return &resource.GetResourceByKeyResponse{
		Resource: res,
		Sharings: sharings,
	}

}

// Delete deletes a resource
func (rs *ResourceStore) Delete(deleteResourceQuery *resource.DeleteResourceQuery) *resource.DeleteResourceResponse {

	graphSession, err := rs.graphDriver.GetSession()
	if err != nil {
		return &resource.DeleteResourceResponse{
			Error: err,
		}
	}

	graphSession.Run(`MATCH (r:Resource {id:$id}) DETACH DELETE`, map[string]interface{}{
		"id": deleteResourceQuery.ResourceKey.String(),
	})

	now := time.Now()
	deleteResult, err := graphSession.Run(`
			MATCH (r:Resource {id:$id})
			SET r += {
				deletedAt:$deletedAt
			}
			RETURN r`,
		map[string]interface{}{
			"id":        deleteResourceQuery.ResourceKey.String(),
			"deletedAt": now,
		})
	if err != nil {
		return &resource.DeleteResourceResponse{
			Error: err,
		}
	}
	if deleteResult.Err() != nil {
		return &resource.DeleteResourceResponse{
			Error: deleteResult.Err(),
		}
	}
	if !deleteResult.Next() {
		return &resource.DeleteResourceResponse{
			Error: errs.ErrResourceNotFound,
		}
	}

	return &resource.DeleteResourceResponse{}

}

// Create creates a resource
func (rs *ResourceStore) Create(createResourceQuery *resource.CreateResourceQuery) *resource.CreateResourceResponse {

	graphSession, err := rs.graphDriver.GetSession()
	if err != nil {
		return &resource.CreateResourceResponse{
			Error: err,
		}
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
				valueInHoursFrom:$valueInHoursFrom,
				valueInHoursTo:$valueInHoursTo
			})-[c:CreatedBy]->(u)
			`

	if createResourceQuery.SharedWith != nil && len(createResourceQuery.SharedWith.Items) > 0 {
		params["groupIds"] = createResourceQuery.SharedWith.Strings()
		cypher = cypher + `
			WITH u, r
			OPTIONAL MATCH (g:Group)
			WHERE g.id IN $groupIds
			CALL apoc.do.when(
				g IS NOT NULL,
				'MERGE (r)-[s:SharedWith]->(g) return s, g',
				'',
				{r: r, g: g})
			YIELD value
			RETURN r, g`
	} else {
		cypher = cypher + `
			RETURN r, null as g`
	}

	createResult, err := graphSession.Run(cypher, params)

	if err != nil {
		return &resource.CreateResourceResponse{
			Error: err,
		}
	}
	if createResult.Err() != nil {
		return &resource.CreateResourceResponse{
			Error: createResult.Err(),
		}
	}

	createResult.Next()

	record := createResult.Record()
	_, err = rs.mapGraphResourceRecord(record, "r")
	if err != nil {
		return &resource.CreateResourceResponse{
			Error: err,
		}
	}

	return &resource.CreateResourceResponse{}

}

func (rs *ResourceStore) mapGraphResourceRecord(record neo4j.Record, key string) (*resource.Resource, error) {
	resourceRecord, _ := record.Get(key)
	node := resourceRecord.(neo4j.Node)
	return MapResourceNode(node)
}

func IsResourceNode(node neo4j.Node) bool {
	return NodeHasLabel(node, "Resource")
}

func MapResourceNode(node neo4j.Node) (*resource.Resource, error) {
	var graphResource = GraphResource{}
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

func (rs *ResourceStore) mapGraphSharingRecord(record neo4j.Record, resourceFieldKey string, groupIdsFieldKey string) (*resource.Sharings, error) {
	resourceField, _ := record.Get(resourceFieldKey)
	resourceNode := resourceField.(neo4j.Node)
	resourceId := resourceNode.Props()["id"].(string)
	resourceKey, err := model.ParseResourceKey(resourceId)
	if err != nil {
		return nil, err
	}

	groupIdsField, _ := record.Get(groupIdsFieldKey)
	if groupIdsField == nil {
		return resource.NewEmptyResourceSharings(), nil
	}
	groupIds := groupIdsField.([]interface{})
	var sharings []*resource.Sharing
	for _, groupId := range groupIds {
		groupKey, err := group.ParseGroupKey(groupId.(string))
		if err != nil {
			return nil, err
		}
		sharing := &resource.Sharing{
			ResourceKey: resourceKey,
			GroupKey:    groupKey,
		}
		sharings = append(sharings, sharing)
	}

	return resource.NewResourceSharings(sharings), nil
}

// Update updates a resource
func (rs *ResourceStore) Update(request *resource.UpdateResourceQuery) *resource.UpdateResourceResponse {

	session, err := rs.graphDriver.GetSession()
	if err != nil {
		return &resource.UpdateResourceResponse{
			Error: err,
		}
	}
	defer session.Close()

	now := time.Now().UTC().UnixNano() / 1e6
	resourceKey := request.Resource.GetKey()
	updateResult, err := session.Run(`
			MATCH (r:Resource {id:$id})
			OPTIONAL MATCH (r)-[notSharedWith:SharedWith]-(g1:Group)
			WHERE NOT (g1.id IN $groupIds)			
			OPTIONAL MATCH (g:Group)
			WHERE g.id in $groupIds
			SET r += {
				updatedAt:datetime({epochMillis:$updatedAt}),
				summary:$summary,
				description:$description,
				valueInHoursFrom:$valueInHoursFrom,
				valueInHoursTo:$valueInHoursTo
			}
			DELETE notSharedWith
			WITH r, g
			CALL apoc.do.when(
				g IS NOT NULL,
				'MERGE (r)-[s:SharedWith]->(g) return s, g',
				'',
				{r: r, g: g})
			YIELD value
			RETURN r`,
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
		return &resource.UpdateResourceResponse{
			Error: err,
		}
	}

	if updateResult.Err() != nil {
		return &resource.UpdateResourceResponse{
			Error: updateResult.Err(),
		}
	}

	if !updateResult.Next() {
		return &resource.UpdateResourceResponse{
			Error: errs.ErrResourceNotFound,
		}
	}

	return resource.NewUpdateResourceResponse(nil)

}

// Search search for resources
func (rs *ResourceStore) Search(request *resource.SearchResourcesQuery) *resource.SearchResourcesResponse {

	session, err := rs.graphDriver.GetSession()
	if err != nil {
		return &resource.SearchResourcesResponse{
			Error: err,
		}
	}
	defer session.Close()

	propertyValues := map[string]interface{}{}
	var matchClauses = []string{
		"(r:Resource)",
	}
	var whereClauses []string

	if request.CreatedBy != "" {
		matchClauses = append([]string{"(createdBy:User {id:$createdById})"}, matchClauses...)
		whereClauses = append(whereClauses, "(r)<-[:CreatedBy]-(createdBy)")
		propertyValues["createdById"] = request.CreatedBy
	}

	if request.Type != nil {
		whereClauses = append(whereClauses, "r.type = $type")
		propertyValues["type"] = *request.Type
	}

	if request.Query != nil && *request.Query != "" {
		whereClauses = append(whereClauses, "r.summary =~ $query")
		propertyValues["query"] = ".*" + *request.Query + ".*"
	}

	var cyper = "MATCH "
	cyper = cyper + strings.Join(matchClauses, ",")

	if len(whereClauses) > 0 {
		cyper = cyper + "\n WHERE "
		cyper = cyper + strings.Join(whereClauses, " AND ")
	}

	if request.SharedWithGroup != nil {
		propertyValues["groupId"] = request.SharedWithGroup.String()
		cyper = cyper + "\nMATCH (r)-[s:SharedWith]->(g:Group {id:$groupId})"
	} else {
		cyper = cyper + "\nOPTIONAL MATCH (r)-[s:SharedWith]->(g:Group)"
	}

	countCypher := cyper + `
RETURN count(r) as totalCount
`

	countResult, err := session.Run(countCypher, propertyValues)
	if err != nil {
		return &resource.SearchResourcesResponse{
			Error: err,
		}
	}

	if countResult.Err() != nil {
		return &resource.SearchResourcesResponse{
			Error: countResult.Err(),
		}
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
		return &resource.SearchResourcesResponse{
			Error: err,
		}
	}

	if searchResult.Err() != nil {
		return &resource.SearchResourcesResponse{
			Error: searchResult.Err(),
		}
	}

	var resources []*resource.Resource

	for searchResult.Next() {

		res, err := rs.mapGraphResourceRecord(searchResult.Record(), "r")
		if err != nil {
			return &resource.SearchResourcesResponse{
				Error: err,
			}
		}

		resources = append(resources, res)

	}

	return &resource.SearchResourcesResponse{
		Resources:  resource.NewResources(resources),
		Sharings:   resource.NewEmptyResourceSharings(),
		Skip:       request.Skip,
		Take:       request.Take,
		TotalCount: int(totalCount),
	}

}

func mapGraphResourceToResource(dbResultItem *GraphResource) (*resource.Resource, error) {

	key, err := model.ParseResourceKey(dbResultItem.ID)
	if err != nil {
		return nil, err
	}
	return &resource.Resource{
		Key:              key,
		CreatedAt:        dbResultItem.CreatedAt,
		UpdatedAt:        dbResultItem.UpdatedAt,
		DeletedAt:        dbResultItem.DeletedAt,
		Summary:          dbResultItem.Summary,
		Description:      dbResultItem.Description,
		CreatedBy:        dbResultItem.CreatedBy,
		Type:             dbResultItem.Type,
		ValueInHoursFrom: dbResultItem.ValueInHoursFrom,
		ValueInHoursTo:   dbResultItem.ValueInHoursTo,
	}, nil
}
