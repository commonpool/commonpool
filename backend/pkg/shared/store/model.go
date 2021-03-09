package store

import (
	"fmt"
	store3 "github.com/commonpool/backend/pkg/auth/store"
	store2 "github.com/commonpool/backend/pkg/group/store"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func MapTargets(record *neo4j.Record, targetsFieldName string) (*keys.Targets, error) {
	field, _ := record.Get(targetsFieldName)

	if field == nil {
		return keys.NewEmptyTargets(), nil
	}

	intfs := field.([]interface{})
	var targets []*keys.Target
	for _, intf := range intfs {
		node := intf.(neo4j.Node)
		target, err := MapOfferItemTarget(node)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}
	return keys.NewTargets(targets), nil
}

func MapOfferItemTarget(node neo4j.Node) (*keys.Target, error) {
	isGroup := store2.IsGroupNode(node)
	isUser := !isGroup && store3.IsUserNode(node)
	if !isGroup && !isUser {
		return nil, fmt.Errorf("target is neither user nor group")
	}

	if isGroup {
		groupKey, err := keys.ParseGroupKey(node.Props["id"].(string))
		if err != nil {
			return nil, err
		}
		return &keys.Target{
			GroupKey: &groupKey,
			Type:     keys.GroupTarget,
		}, nil
	}
	userKey := keys.NewUserKey(node.Props["id"].(string))
	return &keys.Target{
		UserKey: &userKey,
		Type:    keys.UserTarget,
	}, nil
}
