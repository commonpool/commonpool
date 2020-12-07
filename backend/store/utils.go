package store

import (
	"context"
	"github.com/commonpool/backend/logging"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.uber.org/zap"
)

func GetCtx(ctx context.Context, storeName string, actionName string) (context.Context, *zap.Logger) {
	l := logging.WithContext(ctx).
		Named("store."+storeName).
		With(zap.String("service", storeName), zap.String("action", actionName))
	return ctx, l
}

func NodeHasLabel(node neo4j.Node, nodeLabel string) bool {
	for _, label := range node.Labels() {
		if nodeLabel == label {
			return true
		}
	}
	return false
}
