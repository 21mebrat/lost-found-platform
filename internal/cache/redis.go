package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func New(ctx context.Context, addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		// connection size
		PoolSize: 10,
		// connection timeout
		DialTimeout: 5 * time.Second,

		// operation timeout
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	redisContext, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(redisContext).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return client, nil
}
