package executionDB

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	NotFound error = errors.New("key not found on execution storage")
)

func (r repository) Get(ctx context.Context, key string) (string, error) {
	result, err := r.db.Get(ctx, key).Result()
	if errors.As(err, redis.Nil) {
		return "", NotFound
	}
	return result, err
}

func (r repository) Set(ctx context.Context, key, value string) error {
	return r.db.Set(ctx, key, value, 0).Err()
}
