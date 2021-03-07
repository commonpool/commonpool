package store

import (
	graph2 "github.com/commonpool/backend/pkg/graph"
	"github.com/commonpool/backend/pkg/resource"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func MapResourceNode(node neo4j.Node) (*resource.Resource, error) {
	var graphResource = Resource{}
	err := mapstructure.Decode(node.Props, &graphResource)
	if err != nil {
		return nil, err
	}
	mappedResource, err := mapGraphResourceToResource(&graphResource)
	if err != nil {
		return nil, err
	}
	return mappedResource, nil
}

func IsResourceNode(node neo4j.Node) bool {
	return graph2.NodeHasLabel(node, "Resource")
}
