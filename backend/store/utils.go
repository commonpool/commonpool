package store

import (
	"context"
	"github.com/commonpool/backend/logging"
	"go.uber.org/zap"
)

func GetCtx(ctx context.Context, storeName string, actionName string) (context.Context, *zap.Logger) {
	l := logging.WithContext(ctx).
		Named("store."+storeName).
		With(zap.String("service", storeName), zap.String("action", actionName))
	return ctx, l
}
