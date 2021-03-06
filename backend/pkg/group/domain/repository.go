package domain

import (
	"context"
	"github.com/commonpool/backend/pkg/keys"
)

type GroupRepository interface {
	Load(ctx context.Context, groupKey keys.GroupKey) (*Group, error)
	Save(ctx context.Context, group *Group) error
}
