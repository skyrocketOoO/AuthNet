package redis

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/skyrocketOoO/AuthNet/domain"
)

type Redis2Repository struct {
	client *redis.Client
}

func NewRedis2Repository(client *redis.Client) (*Redis2Repository, error) {
	log.Info().Msg("initialize redis2Repository")
	return &Redis2Repository{
		client: client,
	}, nil
}

func (r *Redis2Repository) Ping(c context.Context) error {
	return r.client.Ping(c).Err()
}

func (r *Redis2Repository) Get(c context.Context, edge domain.Edge,
	queryMode bool) ([]domain.Edge, error) {
	if queryMode {
		if edge == (domain.Edge{}) {
			keys, err := r.getKeysFromPattern(c, "*")
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
			_, toStr := edgeToKeyValue(edge)
			if toStr != "%%" {
				return nil, domain.ErrNotImplemented{}
			} else {
				fromPattern := vertexToPattern(
					domain.Vertex{
						Ns:   edge.SbjNs,
						Name: edge.SbjName,
						Rel:  edge.SbjRel,
					},
				)
				keys, err := r.getKeysFromPattern(c, fromPattern)
				if err != nil {
					return nil, err
				}
				edges := []domain.Edge{}
				for _, key := range keys {
					v, err := stringToVertex(key)
					if err != nil {
						return nil, err
					}
					values, err := r.getValues(c, key)
					if err != nil {
						return nil, err
					}
					for _, val := range values {
						strSplit := strings.Split(val, "%")
						edges = append(edges, domain.Edge{
							ObjNs:   strSplit[0],
							ObjName: strSplit[1],
							ObjRel:  strSplit[2],
							SbjNs:   v.Ns,
							SbjName: v.Name,
							SbjRel:  v.Rel,
						})
					}
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

func (r *Redis2Repository) Create(c context.Context, edge domain.Edge) error {
	from, to := edgeToKeyValue(edge)
	return r.client.SAdd(c, from, to).Err()
}

func (r *Redis2Repository) Delete(c context.Context, edge domain.Edge,
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
					if err := r.client.SRem(c, key, val).Err(); err != nil {
						return err
					}
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

func (r *Redis2Repository) ClearAll(c context.Context) error {
	return r.client.FlushDB(c).Err()
}

func (r *Redis2Repository) getKeysFromPattern(c context.Context,
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

func (r *Redis2Repository) getValues(c context.Context, key string) (
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
