package store

import (
	"fmt"
	store2 "github.com/commonpool/backend/pkg/group/store"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/resource"
	"github.com/commonpool/backend/pkg/user/store"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func MapTargets(record neo4j.Record, targetsFieldName string) (*resource.Targets, error) {
	field, _ := record.Get(targetsFieldName)

	if field == nil {
		return resource.NewEmptyTargets(), nil
	}

	intfs := field.([]interface{})
	var targets []*resource.Target
	for _, intf := range intfs {
		node := intf.(neo4j.Node)
		target, err := MapOfferItemTarget(node)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}
	return resource.NewTargets(targets), nil
}

func MapOfferItemTarget(node neo4j.Node) (*resource.Target, error) {
	if node == nil {
		return nil, fmt.Errorf("node is nil")
	}
	isGroup := store2.IsGroupNode(node)
	isUser := !isGroup && store.IsUserNode(node)
	if !isGroup && !isUser {
		return nil, fmt.Errorf("target is neither user nor group")
	}

	if isGroup {
		groupKey, err := keys.ParseGroupKey(node.Props()["id"].(string))
		if err != nil {
			return nil, err
		}
		return &resource.Target{
			GroupKey: &groupKey,
			Type:     resource.GroupTarget,
		}, nil
	}
	userKey := keys.NewUserKey(node.Props()["id"].(string))
	return &resource.Target{
		UserKey: &userKey,
		Type:    resource.UserTarget,
	}, nil
}
