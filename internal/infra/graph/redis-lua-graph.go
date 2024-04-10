package graph

import (
	"context"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/redis/go-redis/v9"
	"github.com/skyrocketOoO/AuthNet/domain"
	"github.com/skyrocketOoO/AuthNet/utils"
	"github.com/skyrocketOoO/go-utility/queue"
	"github.com/skyrocketOoO/go-utility/set"
)

type RedisLuaGraphInfra struct {
	dbRepo domain.DbRepository
	rdsCli *redis.Client
}

func NewRedisLuaGraphInfra(dbRepo domain.DbRepository,
	rdsCli *redis.Client) *RedisLuaGraphInfra {
	return &RedisLuaGraphInfra{
		dbRepo: dbRepo,
		rdsCli: rdsCli,
	}
}

func (g *RedisLuaGraphInfra) Check(c context.Context, sbj domain.Vertex,
	obj domain.Vertex, searchCond domain.SearchCond) (bool, error) {
	result := g.rdsCli.Eval(c, `
		local function bfs(start, target)
			local queue = {start}
			local visited = {}
			while #queue > 0 do
				local current = table.remove(queue, 1)
				if not visited[current] then
					visited[current] = true
					local members = redis.call('SMEMBERS', current)
					for _, member in ipairs(members) do
						if member == target then
							return 1
						else
							table.insert(queue, member)
						end
					end
				end
			end
			return 0
		end
		return bfs(KEYS[1], KEYS[2])
	`, []string{vertexTostring(sbj), vertexTostring(obj)})
	if err := result.Err(); err != nil {
		return false, err
	}
	if result.Val().(int64) == 1 {
		return true, nil
	}
	return false, nil
}

func (g *RedisLuaGraphInfra) GetPassedVertices(c context.Context, start domain.Vertex,
	isSbj bool, searchCond domain.SearchCond, collectCond domain.CollectCond,
	maxDepth int) ([]domain.Vertex, error) {
	if isSbj {
		if err := utils.ValidateVertex(start, true); err != nil {
			return nil, err
		}
		depth := 0
		verticesSet := set.NewSet[domain.Vertex]()
		visited := set.NewSet[domain.Vertex]()
		q := queue.NewQueue[domain.Vertex]()
		visited.Add(start)
		q.Push(start)
		for !q.IsEmpty() {
			qLen := q.Len()
			for i := 0; i < qLen; i++ {
				vertex, _ := q.Pop()
				query := domain.Edge{
					SbjNs:   vertex.Ns,
					SbjName: vertex.Name,
					SbjRel:  vertex.Rel,
				}
				qEdges, err := g.dbRepo.Get(c, query, true)
				if err != nil {
					return nil, err
				}

				for _, edge := range qEdges {
					child := domain.Vertex{
						Ns:   edge.ObjNs,
						Name: edge.ObjName,
						Rel:  edge.ObjRel,
					}
					if collectCond.ShouldCollect(child) {
						verticesSet.Add(child)
					}
					if !searchCond.ShouldStop(child) &&
						!visited.Exist(child) {
						visited.Add(child)
						q.Push(child)
					}
				}
			}
			depth++
			if depth >= maxDepth {
				break
			}
		}

		return verticesSet.ToSlice(), nil
	} else {
		return nil, domain.ErrNotImplemented{}
	}
}

func (g *RedisLuaGraphInfra) GetTree(c context.Context, sbj domain.Vertex, maxDepth int) (
	*domain.TreeNode, error) {
	if res, err := g.dbRepo.Get(
		c,
		domain.Edge{
			SbjNs:   sbj.Ns,
			SbjName: sbj.Name,
			SbjRel:  sbj.Rel,
		}, true); err != nil {
		return nil, err
	} else if len(res) == 0 {
		return nil, domain.ErrRecordNotFound{}
	}

	root := &domain.TreeNode{
		Ns:       sbj.Ns,
		Name:     sbj.Name,
		Rel:      sbj.Rel,
		Children: []*domain.TreeNode{},
	}
	visited := map[domain.Vertex]*domain.TreeNode{}
	visited[sbj] = root
	q := queue.NewQueue[*domain.TreeNode]()
	q.Push(root)
	for depth := 1; depth <= maxDepth && !q.IsEmpty(); depth++ {
		for i := 0; i < q.Len(); i++ {
			v, err := q.Pop()
			if err != nil {
				return nil, err
			}
			edges, err := g.dbRepo.Get(c,
				domain.Edge{
					SbjNs:   v.Ns,
					SbjName: v.Name,
					SbjRel:  v.Rel,
				},
				true,
			)
			if err != nil {
				return nil, err
			}
			for _, edge := range edges {
				obj := domain.Vertex{
					Ns:   edge.ObjNs,
					Name: edge.ObjName,
					Rel:  edge.ObjRel,
				}
				if node, ok := visited[obj]; !ok {
					newNode := &domain.TreeNode{
						Ns:       obj.Ns,
						Name:     obj.Name,
						Rel:      obj.Rel,
						Children: []*domain.TreeNode{},
					}
					q.Push(newNode)
					visited[obj] = newNode
					v.Children = append(v.Children, newNode)
				} else {
					v.Children = append(v.Children, node)
				}
			}
		}
	}
	return root, nil
}

func (g *RedisLuaGraphInfra) SeeTree(c context.Context, sbj domain.Vertex, maxDepth int) (
	*charts.Tree, error) {
	if res, err := g.dbRepo.Get(
		c,
		domain.Edge{
			SbjNs:   sbj.Ns,
			SbjName: sbj.Name,
			SbjRel:  sbj.Rel,
		}, true); err != nil {
		return nil, err
	} else if len(res) == 0 {
		return nil, domain.ErrRecordNotFound{}
	}

	root := domain.Vertex{
		Ns:   sbj.Ns,
		Name: sbj.Name,
		Rel:  sbj.Rel,
	}
	visited := map[string]*opts.TreeData{}
	rootTreeData := &opts.TreeData{
		Name:     vertexTostring(root),
		Children: []*opts.TreeData{},
	}
	visited[vertexTostring(root)] = rootTreeData
	q := queue.NewQueue[domain.Vertex]()
	q.Push(root)
	for depth := 1; depth <= maxDepth && !q.IsEmpty(); depth++ {
		for i := 0; i < q.Len(); i++ {
			v, err := q.Pop()
			if err != nil {
				return nil, err
			}
			edges, err := g.dbRepo.Get(c,
				domain.Edge{
					SbjNs:   v.Ns,
					SbjName: v.Name,
					SbjRel:  v.Rel,
				},
				true,
			)
			if err != nil {
				return nil, err
			}
			sbj := visited[vertexTostring(v)]
			for _, edge := range edges {
				obj := domain.Vertex{
					Ns:   edge.ObjNs,
					Name: edge.ObjName,
					Rel:  edge.ObjRel,
				}
				if node, ok := visited[vertexTostring(obj)]; !ok {
					newNode := &opts.TreeData{
						Name:     vertexTostring(obj),
						Children: []*opts.TreeData{},
					}
					q.Push(obj)
					visited[vertexTostring(obj)] = newNode
					sbj.Children = append(sbj.Children, newNode)
				} else {
					sbj.Children = append(sbj.Children, node)
				}
			}
		}
	}

	graph := charts.NewTree()
	graph.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "95vh"}),
		charts.WithTitleOpts(opts.Title{Title: "basic tree example"}),
	)
	graph.AddSeries("tree", []opts.TreeData{*rootTreeData}).
		SetSeriesOptions(
			charts.WithTreeOpts(
				opts.TreeChart{
					Layout:           "orthogonal",
					Orient:           "LR",
					InitialTreeDepth: -1,
					Leaves: &opts.TreeLeaves{
						Label: &opts.Label{Show: true, Position: "right", Color: "Black"},
					},
				},
			),
			charts.WithLabelOpts(opts.Label{Show: true, Position: "top", Color: "Black"}),
		)

	return graph, nil

}
