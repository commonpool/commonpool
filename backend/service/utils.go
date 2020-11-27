package service

import (
	"context"
	"github.com/commonpool/backend/logging"
	"go.uber.org/zap"
)

func GetCtx(ctx context.Context, serviceName string, actionName string) (context.Context, *zap.Logger) {
	l := logging.WithContext(ctx).
		Named("service."+serviceName).
		With(
			zap.String("service", serviceName),
			zap.String("action", actionName))
	return ctx, l
}
