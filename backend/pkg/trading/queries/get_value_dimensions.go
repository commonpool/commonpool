package queries

import (
	"context"
	domain2 "github.com/commonpool/backend/pkg/resource/domain"
)

type GetValueDimensions struct {
}

func NewGetValueDimensions() *GetValueDimensions {
	return &GetValueDimensions{}
}

func (q *GetValueDimensions) Get(ctx context.Context) domain2.ValueDimensions {
	return domain2.AllDimensions
}
