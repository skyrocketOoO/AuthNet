package usecase

import (
	"context"
	"math"
	"sync"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/skyrocketOoO/AuthNet/domain"
	"github.com/skyrocketOoO/AuthNet/utils"
)

type Usecase struct {
	sqlRepo        domain.DbRepository
	graphInfra     domain.GraphInfra
	PageStatesLock sync.RWMutex
}

func NewUsecase(sqlRepo domain.DbRepository,
	graphInfra domain.GraphInfra) *Usecase {
	usecase := &Usecase{
		sqlRepo:        sqlRepo,
		graphInfra:     graphInfra,
		PageStatesLock: sync.RWMutex{},
	}

	return usecase
}

func (u *Usecase) Healthy(c context.Context) error {
	// do something check like db connection is established
	if err := u.sqlRepo.Ping(c); err != nil {
		return err
	}

	return nil
}

func (u *Usecase) Get(c context.Context, edge domain.Edge, queryMode bool) (
	[]domain.Edge, error) {
	return u.sqlRepo.Get(c, edge, queryMode)
}

func (u *Usecase) Create(c context.Context, edge domain.Edge) error {
	if err := utils.ValidateRel(edge); err != nil {
		return err
	}
	ok, err := u.graphInfra.Check(
		c,
		domain.Vertex{
			Ns:   edge.SbjNs,
			Name: edge.SbjName,
			Rel:  edge.SbjRel,
		},
		domain.Vertex{
			Ns:   edge.ObjNs,
			Name: edge.ObjName,
			Rel:  edge.ObjRel,
		},
		domain.SearchCond{},
	)
	if err != nil {
		return err
	}
	if ok {
		return domain.ErrGraphCycle{}
	}

	return u.sqlRepo.Create(c, edge)
}

func (u *Usecase) Delete(c context.Context, edge domain.Edge,
	queryMode bool) error {
	if !queryMode {
		if err := utils.ValidateRel(edge); err != nil {
			return err
		}
	}
	return u.sqlRepo.Delete(c, edge, queryMode)
}

func (u *Usecase) ClearAll(c context.Context) error {
	return u.sqlRepo.ClearAll(c)
}

func (u *Usecase) CheckAuth(c context.Context, sbj domain.Vertex,
	obj domain.Vertex, searchCond domain.SearchCond) (fond bool, err error) {
	if err := utils.ValidateVertex(obj, false); err != nil {
		return false, err
	}
	if err := utils.ValidateVertex(sbj, true); err != nil {
		return false, err
	}
	return u.graphInfra.Check(c, sbj, obj, searchCond)
}

func (u *Usecase) GetObjAuths(c context.Context, sbj domain.Vertex,
	searchCond domain.SearchCond, collectCond domain.CollectCond,
	maxDepth int) (vertices []domain.Vertex, err error) {
	return u.graphInfra.GetPassedVertices(c, sbj, true, searchCond,
		collectCond, maxDepth)
}

func (u *Usecase) GetSbjsWhoHasAuth(c context.Context, obj domain.Vertex,
	searchCond domain.SearchCond, collectCond domain.CollectCond,
	maxDepth int) (vertices []domain.Vertex, err error) {
	return u.graphInfra.GetPassedVertices(c, obj, false, searchCond,
		collectCond, maxDepth)
}

func (u *Usecase) GetTree(c context.Context, sbj domain.Vertex, maxDepth int) (
	*domain.TreeNode, error) {
	if maxDepth == 0 {
		maxDepth = math.MaxInt
	}
	return u.graphInfra.GetTree(c, sbj, maxDepth)
}

func (u *Usecase) SeeTree(c context.Context, sbj domain.Vertex, maxDepth int) (
	*charts.Tree, error) {
	if maxDepth == 0 {
		maxDepth = math.MaxInt
	}
	return u.graphInfra.SeeTree(c, sbj, maxDepth)
}
