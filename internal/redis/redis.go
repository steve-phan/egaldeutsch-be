package redis

import (
	"context"
	"egaldeutsch-be/internal/config"
	"fmt"
	"log/slog"
	"time"

	rawRedis "github.com/redis/go-redis/v9"
)

type RedisClient struct {
	*rawRedis.Client
}

func NewRedisClient(redisConfig config.RedisConfig) (*RedisClient, error) {
	client := rawRedis.NewClient(
		&rawRedis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
			Password: redisConfig.Password,
			DB:       redisConfig.DB,
		},
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("connect to Redis fail: %w", err)
	}
	slog.Info("Connected to Redis successfully")
	return &RedisClient{Client: client}, nil
}

func (r *RedisClient) Close() error {
	return r.Client.Close()
}
