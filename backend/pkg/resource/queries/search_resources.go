package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"gorm.io/gorm"
)

type SearchResources struct {
	db *gorm.DB
}

func NewSearchResources(db *gorm.DB) *SearchResources {
	return &SearchResources{db: db}
}

func (q *SearchResources) Get(ctx context.Context, query *SearchResourcesQuery) ([]*readmodel.ResourceReadModel, error) {

	qry := q.db.Model(&readmodel.ResourceReadModel{}).Offset(query.Skip).Limit(query.Take)

	if query.CreatedBy != nil {
		qry = qry.Where("created_by_id = ?", *query.CreatedBy)
	}

	if query.Type != nil {
		qry = qry.Where("resource_type = ?", *query.Type)
	}

	if query.CallType != nil {
		qry = qry.Where("call_type = ?", *query.CallType)
	}

	if query.SharedWithGroup != nil {
		sql := "exists (select null from resource_sharing_read_models s where s.resource_key = resource_read_models.resource_key)"
		qry = qry.Where(q.db.Raw(sql))
	}

	if query.Query != nil {
		qry = qry.Where("resource_name like ?", *query.Query+"%")
	}

	var result []*readmodel.ResourceReadModel
	err := qry.Find(&result).Error
	return result, err

}

type SearchResourcesQuery struct {
	Type            *domain.ResourceType
	CallType        *domain.CallType
	Query           *string
	Skip            int
	Take            int
	CreatedBy       *string
	SharedWithGroup *keys.GroupKey
}

func NewSearchResourcesQuery(
	query *string,
	resourceType *domain.ResourceType,
	resourceSubType *domain.CallType,
	skip int,
	take int,
	createdBy *string,
	sharedWithGroup *keys.GroupKey) *SearchResourcesQuery {
	return &SearchResourcesQuery{
		Type:            resourceType,
		CallType:        resourceSubType,
		Query:           query,
		Skip:            skip,
		Take:            take,
		CreatedBy:       createdBy,
		SharedWithGroup: sharedWithGroup,
	}
}
