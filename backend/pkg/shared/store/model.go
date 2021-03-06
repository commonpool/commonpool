package store

import (
	"fmt"
	store2 "github.com/commonpool/backend/pkg/group/store"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/commonpool/backend/pkg/user/store"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func MapTargets(record *neo4j.Record, targetsFieldName string) (*domain.Targets, error) {
	field, _ := record.Get(targetsFieldName)

	if field == nil {
		return domain.NewEmptyTargets(), nil
	}

	intfs := field.([]interface{})
	var targets []*domain.Target
	for _, intf := range intfs {
		node := intf.(neo4j.Node)
		target, err := MapOfferItemTarget(node)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}
	return domain.NewTargets(targets), nil
}

func MapOfferItemTarget(node neo4j.Node) (*domain.Target, error) {
	isGroup := store2.IsGroupNode(node)
	isUser := !isGroup && store.IsUserNode(node)
	if !isGroup && !isUser {
		return nil, fmt.Errorf("target is neither user nor group")
	}

	if isGroup {
		groupKey, err := keys.ParseGroupKey(node.Props["id"].(string))
		if err != nil {
			return nil, err
		}
		return &domain.Target{
			GroupKey: &groupKey,
			Type:     domain.GroupTarget,
		}, nil
	}
	userKey := keys.NewUserKey(node.Props["id"].(string))
	return &domain.Target{
		UserKey: &userKey,
		Type:    domain.UserTarget,
	}, nil
}
