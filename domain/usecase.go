package domain

import (
	"context"

	"github.com/go-echarts/go-echarts/v2/charts"
)

type Usecase interface {
	Healthy(ctx context.Context) error

	Get(c context.Context, edge Edge, queryMode bool) (edges []Edge, err error)
	Create(c context.Context, edge Edge) error
	Delete(c context.Context, edge Edge, queryMode bool) error
	ClearAll(c context.Context) error
	// BatchOperation(operations []Operation) error
	// GetNamespaces(c context.Context) (namespaces []string, err error)
	CheckAuth(c context.Context, sbj Vertex, obj Vertex,
		searchCond SearchCond) (fond bool, err error)

	GetObjAuths(c context.Context, sbj Vertex, searchCond SearchCond,
		collectCond CollectCond, maxDepth int) (
		vertices []Vertex, err error)
	GetSbjsWhoHasAuth(c context.Context, obj Vertex, searchCond SearchCond,
		collectCond CollectCond, maxDepth int) (
		vertices []Vertex, err error)
	GetTree(c context.Context, subject Vertex, maxDepth int) (*TreeNode, error)
	SeeTree(c context.Context, sbj Vertex, maxDepth int) (*charts.Tree,
		error)

	Batch(c context.Context, operations []Operation) error
}
