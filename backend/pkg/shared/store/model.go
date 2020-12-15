package store

import (
	"fmt"
	"github.com/commonpool/backend/model"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	store2 "github.com/commonpool/backend/pkg/group/store"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	"github.com/commonpool/backend/pkg/user/store"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func MapTargets(record neo4j.Record, targetsFieldName string) (*model.Targets, error) {
	field, _ := record.Get(targetsFieldName)

	if field == nil {
		return model.NewEmptyTargets(), nil
	}

	intfs := field.([]interface{})
	var targets []*model.Target
	for _, intf := range intfs {
		node := intf.(neo4j.Node)
		target, err := MapOfferItemTarget(node)
		if err != nil {
			return nil, err
		}
		targets = append(targets, target)
	}
	return model.NewTargets(targets), nil
}

func MapOfferItemTarget(node neo4j.Node) (*model.Target, error) {
	if node == nil {
		return nil, fmt.Errorf("node is nil")
	}
	isGroup := store2.IsGroupNode(node)
	isUser := !isGroup && store.IsUserNode(node)
	if !isGroup && !isUser {
		return nil, fmt.Errorf("target is neither user nor group")
	}

	if isGroup {
		groupKey, err := groupmodel.ParseGroupKey(node.Props()["id"].(string))
		if err != nil {
			return nil, err
		}
		return &model.Target{
			GroupKey: &groupKey,
			Type:     model.GroupTarget,
		}, nil
	}
	userKey := usermodel.NewUserKey(node.Props()["id"].(string))
	return &model.Target{
		UserKey: &userKey,
		Type:    model.UserTarget,
	}, nil
}