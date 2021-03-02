package graph

import "github.com/neo4j/neo4j-go-driver/v4/neo4j"

func NodeHasLabel(node neo4j.Node, nodeLabel string) bool {
	for _, label := range node.Labels {
		if nodeLabel == label {
			return true
		}
	}
	return false
}
