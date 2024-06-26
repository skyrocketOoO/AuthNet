package mongo

import (
	"context"

	"github.com/skyrocketOoO/AuthNet/domain"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoRepository struct {
	client    *mongo.Client
	colClient *mongo.Collection
}

func NewMongoRepository(client *mongo.Client) (*MongoRepository, error) {
	return &MongoRepository{
		client:    client,
		colClient: client.Database(viper.GetString("db")).Collection("edges"),
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
			cursor, err := r.colClient.Find(c, bson.D{})
			if err != nil {
				return nil, err
			}
			defer cursor.Close(c)
			if err := cursor.All(c, &edges); err != nil {
				return nil, err
			}
		} else {
			cursor, err := r.colClient.Find(c, edgeToBSONDWithoutZeroVal(edge))
			if err != nil {
				return nil, err
			}
			defer cursor.Close(c)
			if err := cursor.All(c, &edges); err != nil {
				return nil, err
			}
		}
	} else {
		cursor, err := r.colClient.Find(c, edge)
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
	_, err := r.colClient.InsertOne(
		c,
		edge,
	)
	return err
}

func (r *MongoRepository) Delete(c context.Context, edge domain.Edge,
	queryMode bool) error {
	if queryMode {
		_, err := r.colClient.DeleteMany(c, edgeToBSONDWithoutZeroVal(edge))
		return err
	} else {
		if _, err := r.Get(c, edge, false); err != nil {
			return err
		}
		_, err := r.colClient.DeleteOne(c, edge)
		return err
	}
}

func (r *MongoRepository) ClearAll(c context.Context) error {
	_, err := r.colClient.DeleteMany(c, bson.M{})
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
