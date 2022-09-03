package cache

import "time"

type KVCache interface {
	Get(key string)
	Set(key string, value string, expiration time.Duration)
}
