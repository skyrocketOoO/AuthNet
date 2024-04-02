package redis

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/skyrocketOoO/AuthNet/domain"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) (*RedisRepository, error) {
	return &RedisRepository{
		client: client,
	}, nil
}

func (r *RedisRepository) Ping(c context.Context) error {
	return r.client.Ping(c).Err()
}

func (r *RedisRepository) Get(c context.Context, edge domain.Edge,
	queryMode bool) ([]domain.Edge, error) {
	if queryMode {
		if edge == (domain.Edge{}) {
			keys, err := r.getKeysFromPattern(c, "[^$]*")
			if err != nil {
				return nil, err
			}
			edges := []domain.Edge{}
			for _, key := range keys {
				values, err := r.getValues(c, key)
				if err != nil {
					return nil, err
				}
				fromSplit := strings.Split(key, "%")
				for _, to := range values {
					toSplit := strings.Split(to, "%")
					edges = append(edges, domain.Edge{
						ObjNs:   toSplit[0],
						ObjName: toSplit[1],
						ObjRel:  toSplit[2],
						SbjNs:   fromSplit[0],
						SbjName: fromSplit[1],
						SbjRel:  fromSplit[2],
					})
				}
			}
			return edges, nil
		} else {
			fromStr, toStr := edgeToKeyValue(edge)
			if toStr != "%%" {
				toPattern := vertexToPattern(domain.Vertex{
					Ns:   edge.ObjNs,
					Name: edge.ObjName,
					Rel:  edge.ObjRel,
				})
				keys, err := r.getKeysFromPattern(c, "$"+toPattern)
				if err != nil {
					return nil, err
				}
				edges := []domain.Edge{}
				for _, key := range keys {
					values, err := r.getValues(c, key)
					if err != nil {
						return nil, err
					}
					fromSplit := strings.Split(key, "%")
					for _, to := range values {
						toSplit := strings.Split(to, "%")
						edges = append(edges, domain.Edge{
							ObjNs:   toSplit[0],
							ObjName: toSplit[1],
							ObjRel:  toSplit[2],
							SbjNs:   fromSplit[0],
							SbjName: fromSplit[1],
							SbjRel:  fromSplit[2],
						})
					}
				}
				return edges, nil
			} else {
				values, err := r.getValues(c, fromStr)
				if err != nil {
					return nil, err
				}
				edges := []domain.Edge{}
				for _, val := range values {
					strSplit := strings.Split(val, "%")
					edges = append(edges, domain.Edge{
						ObjNs:   strSplit[0],
						ObjName: strSplit[1],
						ObjRel:  strSplit[2],
						SbjNs:   edge.SbjNs,
						SbjName: edge.SbjName,
						SbjRel:  edge.SbjRel,
					})
				}
				return edges, nil
			}
		}
	} else {
		from, to := edgeToKeyValue(edge)
		rdsBoolCmd := r.client.SIsMember(c, from, to)
		if rdsBoolCmd.Err() != nil {
			return nil, rdsBoolCmd.Err()
		}
		if rdsBoolCmd.Val() {
			return []domain.Edge{edge}, nil
		} else {
			return nil, domain.ErrRecordNotFound{}
		}
	}
}

func (r *RedisRepository) Create(c context.Context, edge domain.Edge) error {
	from, to := edgeToKeyValue(edge)
	if err := r.client.SAdd(c, from, to).Err(); err != nil {
		return err
	}
	return r.client.SAdd(c, addReverse(to), from).Err()
}

func (r *RedisRepository) Delete(c context.Context, edge domain.Edge,
	queryMode bool) error {
	if queryMode {
		fromP := vertexToPattern(domain.Vertex{
			Ns:   edge.SbjNs,
			Name: edge.SbjName,
			Rel:  edge.SbjRel,
		})
		toP := vertexToPattern(domain.Vertex{
			Ns:   edge.ObjNs,
			Name: edge.ObjName,
			Rel:  edge.ObjRel,
		})
		keys, err := r.getKeysFromPattern(c, fromP)
		if err != nil {
			return err
		}
		for _, key := range keys {
			values, err := r.getValues(c, key)
			if err != nil {
				return err
			}
			for _, val := range values {
				match, err := filepath.Match(toP, val)
				if err != nil {
					return err
				}
				if match {
					r.client.SRem(c, key, val)
					r.client.SRem(c, addReverse(val), key)
				}
			}
		}
	} else {
		if _, err := r.Get(c, edge, false); err != nil {
			return err
		}
		from, to := edgeToKeyValue(edge)
		return r.client.SRem(c, from, to).Err()
	}
	return nil
}

func (r *RedisRepository) ClearAll(c context.Context) error {
	return r.client.FlushDB(c).Err()
}

func vertexToString(v domain.Vertex) string {
	return v.Ns + "%" + v.Name + "%" + v.Rel
}

func edgeToKeyValue(edge domain.Edge) (from string, to string) {
	from = vertexToString(domain.Vertex{
		Ns:   edge.SbjNs,
		Name: edge.SbjName,
		Rel:  edge.SbjRel,
	})
	to = vertexToString(domain.Vertex{
		Ns:   edge.ObjNs,
		Name: edge.ObjName,
		Rel:  edge.ObjRel,
	})
	return
}

func vertexToPattern(vertex domain.Vertex) string {
	if vertex == (domain.Vertex{}) {
		return "*"
	}
	res := ""
	if vertex.Ns == "" {
		res += "*"
	} else {
		res += vertex.Ns
	}
	res += "%"
	if vertex.Name == "" {
		res += "*"
	} else {
		res += vertex.Name
	}
	res += "%"
	if vertex.Rel == "" {
		res += "*"
	} else {
		res += vertex.Rel
	}
	return res
}

func (r *RedisRepository) getKeysFromPattern(c context.Context,
	pattern string) ([]string, error) {

	result := r.client.Keys(c, pattern)
	if err := result.Err(); err != nil {
		return nil, err
	}
	keys := []string{}
	if err := result.ScanSlice(&keys); err != nil {
		return nil, err
	}
	return keys, nil
}

func addReverse(in string) string {
	return "$" + in
}

func (r *RedisRepository) getValues(c context.Context, key string) (
	[]string, error) {
	res := r.client.SMembers(c, key)
	if err := res.Err(); err != nil {
		return nil, err
	}
	values := []string{}
	if err := res.ScanSlice(&values); err != nil {
		return nil, err
	}
	return values, nil
}
