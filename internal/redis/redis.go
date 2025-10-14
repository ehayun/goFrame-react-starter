package redis

import (
	"fmt"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() error {
	ctx := gctx.New()
	cfg := g.Cfg()

	// Read Redis configuration (sensitive from env, static from config.yaml)
	host := cfg.MustGet(ctx, "redis.host").String()
	port := cfg.MustGet(ctx, "redis.port").Int()
	password := os.Getenv("REDIS_PASSWORD")
	database := cfg.MustGet(ctx, "redis.database").Int()
	poolSize := cfg.MustGet(ctx, "redis.poolSize").Int()

	// Create Redis client
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       database,
		PoolSize: poolSize,
	})

	// Test connection
	if err := Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	g.Log().Info(ctx, "Redis connection established")

	return nil
}

func Close() error {
	if Client != nil {
		return Client.Close()
	}
	return nil
}
