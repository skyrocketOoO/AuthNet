package mongo2

import (
	"context"

	"github.com/skyrocketOoO/AuthNet/domain"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoRepository struct {
	client    *mongo.Client
	edgeCli   *mongo.Collection
	vertexCli *mongo.Collection
}

func NewMongoRepository(client *mongo.Client) (*MongoRepository, error) {
	db := viper.GetString("db")
	return &MongoRepository{
		client:    client,
		edgeCli:   client.Database(db).Collection("edges"),
		vertexCli: client.Database(db).Collection("vertices"),
	}, nil
}

func (r *MongoRepository) Ping(c context.Context) error {
	return r.client.Ping(c, readpref.Primary())
}

func (r *MongoRepository) Get(c context.Context, edge domain.Edge,
	queryMode bool) ([]domain.Edge, error) {
	edges := []domain.Edge{}
	if queryMode {
		if edge == (domain.Edge{}) {
			cursor, err := r.edgeCli.Find(c, bson.D{})
			if err != nil {
				return nil, err
			}
			defer cursor.Close(c)
			var modelEdges []struct {
				U primitive.ObjectID
				V primitive.ObjectID
			}
			if err := cursor.All(c, &modelEdges); err != nil {
				return nil, err
			}
			for _, edge := range modelEdges {
				u := r.vertexCli.FindOne(c, bson.M{"_id": edge.U})
				uv := domain.Vertex{}
				u.Decode(&uv)
				v := r.vertexCli.FindOne(c, bson.M{"_id": edge.V})
				vv := domain.Vertex{}
				v.Decode(&vv)
				edges = append(edges, domain.Edge{
					SbjNs:   uv.Ns,
					SbjName: uv.Name,
					SbjRel:  uv.Rel,
					ObjNs:   vv.Ns,
					ObjName: vv.Name,
					ObjRel:  vv.Rel,
				})
			}

		} else {
			sIds, err := queryVertex(c, r.vertexCli, domain.Vertex{
				Ns: edge.SbjNs, Name: edge.SbjName, Rel: edge.SbjRel,
			})
			if err != nil {
				return nil, err
			}
			oIds, err := queryVertex(c, r.vertexCli, domain.Vertex{
				Ns: edge.ObjNs, Name: edge.ObjName, Rel: edge.ObjRel,
			})
			if err != nil {
				return nil, err
			}
			for _, sId := range sIds {
				for _, oId := range oIds {
					res := r.edgeCli.FindOne(c, bson.M{
						"u": sId.Id, "v": oId.Id,
					})
					if err := res.Err(); err != nil {
						if err != mongo.ErrNoDocuments {
							return nil, res.Err()
						}
						continue
					}
					edges = append(edges, domain.Edge{
						SbjNs:   sId.Ns,
						SbjName: sId.Name,
						SbjRel:  sId.Rel,
						ObjNs:   oId.Ns,
						ObjName: oId.Name,
						ObjRel:  oId.Rel,
					})
				}
			}
		}
	} else {
		cursor, err := r.edgeCli.Find(c, edge)
		if err != nil {
			return nil, err
		}
		defer cursor.Close(c)
		if err := cursor.All(c, &edges); err != nil {
			return nil, err
		}
		if len(edges) == 0 {
			return nil, domain.ErrRecordNotFound{}
		} else if len(edges) > 1 {
			return nil, domain.ErrDuplicateRecord{}
		}
	}
	return edges, nil
}

func (r *MongoRepository) Create(c context.Context, edge domain.Edge) error {
	subId, err := insertIfNotExist(c, r.vertexCli, domain.Vertex{
		Ns:   edge.SbjNs,
		Name: edge.SbjName,
		Rel:  edge.SbjRel,
	})
	if err != nil {
		return err
	}
	objId, err := insertIfNotExist(c, r.vertexCli, domain.Vertex{
		Ns:   edge.ObjNs,
		Name: edge.ObjName,
		Rel:  edge.ObjRel,
	})
	if err != nil {
		return err
	}

	_, err = r.edgeCli.InsertOne(c, bson.M{
		"u": subId,
		"v": objId,
	})
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepository) Delete(c context.Context, edge domain.Edge,
	queryMode bool) error {
	if queryMode {
		_, err := r.edgeCli.DeleteMany(c, edgeToBSONDWithoutZeroVal(edge))
		return err
	} else {
		if _, err := r.Get(c, edge, false); err != nil {
			return err
		}
		_, err := r.edgeCli.DeleteOne(c, edge)
		return err
	}
}

func (r *MongoRepository) ClearAll(c context.Context) error {
	_, err := r.edgeCli.DeleteMany(c, bson.M{})
	return err
}

func edgeToBSONDWithoutZeroVal(e domain.Edge) bson.D {
	doc := bson.D{}
	if e.ObjNs != "" {
		doc = append(doc, bson.E{Key: "obj_ns", Value: e.ObjNs})
	}
	if e.ObjName != "" {
		doc = append(doc, bson.E{Key: "obj_name", Value: e.ObjName})
	}
	if e.ObjRel != "" {
		doc = append(doc, bson.E{Key: "obj_rel", Value: e.ObjRel})
	}
	if e.SbjNs != "" {
		doc = append(doc, bson.E{Key: "sbj_ns", Value: e.SbjNs})
	}
	if e.SbjName != "" {
		doc = append(doc, bson.E{Key: "sbj_name", Value: e.SbjName})
	}
	if e.SbjRel != "" {
		doc = append(doc, bson.E{Key: "sbj_rel", Value: e.SbjRel})
	}
	return doc
}
