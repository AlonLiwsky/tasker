package executionDB

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string) error
}

type DB interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type repository struct {
	db DB
}

func NewRepository(db DB) Repository {
	return &repository{db: db}
}
