package mongo

import (
	"context"
	"time"

	errors "github.com/rotisserie/eris"
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
