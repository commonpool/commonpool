package queries

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource/domain"
	"github.com/commonpool/backend/pkg/resource/readmodel"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"strings"
)

type SearchResources struct {
	db *gorm.DB
}

func NewSearchResources(db *gorm.DB) *SearchResources {
	return &SearchResources{db: db}
}

type SearchResourcesQuery struct {
	ResourceType *domain.ResourceType `query:"resource_type"`
	CallType     *domain.CallType     `query:"call_type"`
	Query        *string              `query:"query"`
	Skip         int                  `query:"skip" validate:"gte=0"`
	Take         int                  `query:"take" validate:"gte=0"`
	CreatedBy    *string              `query:"created_by"`
	GroupKey     *keys.GroupKey       `query:"group_id"`
}

func (k SearchResourcesQuery) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	if k.ResourceType != nil {
		encoder.AddString("resource_type", string(*k.ResourceType))
	}
	if k.CallType != nil {
		encoder.AddString("call_type", string(*k.CallType))
	}
	if k.CreatedBy != nil {
		encoder.AddString("created_by", string(*k.CreatedBy))
	}
	if k.GroupKey != nil {
		encoder.AddString("shared_with_group", k.GroupKey.String())
	}
	encoder.AddInt("skip", k.Skip)
	encoder.AddInt("take", k.Take)
	return nil
}

func (q *SearchResources) Get(ctx context.Context, query *SearchResourcesQuery) ([]*readmodel.ResourceReadModel, error) {

	var sb strings.Builder
	var params []interface{}
	var clauses []string

	if query.CreatedBy != nil {
		clauses = append(clauses, "created_by = ?")
		params = append(params, *query.CreatedBy)
	}

	if query.ResourceType != nil {
		clauses = append(clauses, "resource_type = ?")
		params = append(params, *query.ResourceType)
	}

	if query.CallType != nil {
		clauses = append(clauses, "call_type = ?")
		params = append(params, *query.CallType)
	}

	if query.GroupKey != nil && query.GroupKey.String() != "" {
		clauses = append(clauses, "exists (select null from resource_sharing_read_models s where s.resource_key = resource_read_models.resource_key and s.group_key = ?)")
		params = append(params, *query.GroupKey)
	}

	if query.Query != nil && *query.Query != "" {
		clauses = append(clauses, "name like ?")
		params = append(params, *query.Query+"%")
	}

	sb.WriteString("select * from resource_read_models ")
	if len(clauses) > 0 {
		sb.WriteString("where ")
		sb.WriteString(strings.Join(clauses, " and "))
	}

	sb.WriteString(" limit ?")
	params = append(params, query.Take)
	sb.WriteString(" offset ? ")
	params = append(params, query.Skip)

	var result []*readmodel.DbResourceReadModel
	sql := sb.String()
	err := q.db.Debug().Model(&readmodel.DbResourceReadModel{}).Raw(sql, params...).Scan(&result).Error
	return mapResourceReadModels(result), err

}
