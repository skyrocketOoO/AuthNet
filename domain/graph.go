package domain

import "context"

type GraphInfra interface {
	Check(c context.Context, sbj Vertex, obj Vertex, searchCond SearchCond) (
		found bool, err error)
	GetPassedVertices(c context.Context, start Vertex, isSbj bool,
		searchCond SearchCond, collectCond CollectCond, maxDepth int) (
		vertices []Vertex, err error)
	GetTree(subject Vertex, maxDepth int) (*TreeNode, error)
	// GetShortestPath(subject Vertex, object Vertex, searchCond SearchCond) ([]Edge, error)
	// GetAllPaths(subject Vertex, object Vertex, searchCond SearchCond) ([][]Edge, error)
}
