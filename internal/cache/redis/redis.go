package redis

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	cacheMetrics "homework-1/internal/metrics"
	"time"
)

type Cache struct {
	client  redis.Client
	metrics *cacheMetrics.Metrics
}

func New(opts *redis.Options, metrics *cacheMetrics.Metrics) *Cache {
	client := redis.NewClient(opts)
	return &Cache{
		client:  *client,
		metrics: metrics,
	}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	log.Debugf("Get from redis: %s", key)
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			c.metrics.CacheMissCounter.Inc()
			return "", redis.Nil
		}
		return "", errors.Wrap(err, "failed to get value from redis")
	}
	c.metrics.CacheHitCounter.Inc()
	return val, nil
}

func (c *Cache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	log.Debugf("Set to redis. Key: %s val: %s ex: %d", key, value, expiration)
	err := c.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return errors.Wrap(err, "failed to set value to redis")
	}
	return nil
}

func (c *Cache) Del(ctx context.Context, key string) error {
	log.Debugf("Del from redis: %s", key)
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return errors.Wrap(err, "failed to delete value from redis")
	}
	return nil
}
