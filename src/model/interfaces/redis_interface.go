package interfaces

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisInterface define os métodos necessários do Redis para o service
type RedisInterface interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

// RedisWrapper implementa RedisInterface usando redis.Client
type RedisWrapper struct {
	client *redis.Client
}

func NewRedisWrapper(client *redis.Client) RedisInterface {
	return &RedisWrapper{client: client}
}

func (r *RedisWrapper) Ping(ctx context.Context) *redis.StatusCmd {
	return r.client.Ping(ctx)
}

func (r *RedisWrapper) Get(ctx context.Context, key string) *redis.StringCmd {
	return r.client.Get(ctx, key)
}

func (r *RedisWrapper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return r.client.Set(ctx, key, value, expiration)
}
