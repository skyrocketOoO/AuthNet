package mongo2

import (
	"context"
	"time"

	errors "github.com/rotisserie/eris"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitDb() (*mongo.Client, func(), error) {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, nil, err
	}

	collection := client.Database("zanzibar-dag").Collection("edges")
	_, err = collection.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "obj_ns", Value: 1}, {Key: "obj_name", Value: 1}, {Key: "obj_rel", Value: 1}},
			Options: options.Index().SetName("object_index"),
		},
		{
			Keys:    bson.D{{Key: "sbj_ns", Value: 1}, {Key: "sbj_name", Value: 1}, {Key: "sbj_rel", Value: 1}},
			Options: options.Index().SetName("subject_index"),
		},
	})
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, nil, errors.Wrap(err, "Unable to connect to MongoDB")
	}

	var Disconnect = func() {
		client.Disconnect(ctx)
	}
	return client, Disconnect, nil
}
