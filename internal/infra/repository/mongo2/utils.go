package mongo2

import (
	"context"

	"github.com/skyrocketOoO/AuthNet/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Vertex2 struct {
	Id   primitive.ObjectID `bson:"_id"`
	Ns   string             `bson:"ns"`
	Name string             `bson:"name"`
	Rel  string             `bson:"rel"`
}

func insertIfNotExist(c context.Context, cli *mongo.Collection,
	vertex domain.Vertex) (primitive.ObjectID, error) {
	res := cli.FindOne(c, bson.M{
		"ns":   vertex.Ns,
		"name": vertex.Name,
		"rel":  vertex.Rel,
	})
	if err := res.Err(); err != nil {
		if err != mongo.ErrNoDocuments {
			return [12]byte{}, res.Err()
		}
	}
	var v struct {
		Id primitive.ObjectID `bson:"_id"`
	}
	res.Decode(&v)
	if v.Id.IsZero() {
		res, err := cli.InsertOne(c, bson.M{
			"ns":   vertex.Ns,
			"name": vertex.Name,
			"rel":  vertex.Rel,
		})
		if err != nil {
			return [12]byte{}, err
		}
		return res.InsertedID.(primitive.ObjectID), nil
	}
	return v.Id, nil
}

func queryVertex(c context.Context, cli *mongo.Collection,
	vertex domain.Vertex) ([]Vertex2, error) {
	m := bson.M{}
	if vertex.Ns != "" {
		m["ns"] = vertex.Ns
	}
	if vertex.Name != "" {
		m["name"] = vertex.Name
	}
	if vertex.Rel != "" {
		m["rel"] = vertex.Rel
	}

	res, err := cli.Find(c, m)
	if err != nil {
		return nil, err
	}

	vs := []Vertex2{}
	res.All(c, &vs)

	return vs, nil
}
