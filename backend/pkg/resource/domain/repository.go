package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

type ResourceRepository interface {
	Load(ctx context.Context, resourceKey keys.ResourceKey) (*Resource, error)
	Save(ctx context.Context, resource *Resource) error
}
