package redis

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

func NewRedisClient() (*redis.Client, error) {
	r := redis.NewClient(&redis.Options{})
	ctx := context.Background()
	err := r.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return r, nil
}
