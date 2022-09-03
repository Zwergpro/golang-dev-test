package redis

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	"homework-1/internal/metrics"
	"time"
)

type Cache struct {
	client  redis.Client
	metrics *metrics.Metrics
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			c.metrics.CacheMissCounter.Inc()
			return "", nil
		}
		return "", errors.Wrap(err, "failed to get value from redis")
	}
	c.metrics.CacheHitCounter.Inc()
	return val, nil
}

func (c *Cache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := c.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return errors.Wrap(err, "failed to set value to redis")
	}
	return nil
}
