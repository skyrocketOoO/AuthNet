package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	errors "github.com/rotisserie/eris"
)

func InitDb() (*redis.Client, func(), error) {
	rdsCli := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdsCli.Ping(ctx).Err(); err != nil {
		return nil, nil, errors.Wrap(err, "ping redis failed")
	}

	var Disconnect = func() {
		rdsCli.Close()
	}

	return rdsCli, Disconnect, nil
}
