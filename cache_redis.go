package sqlcache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v4"

	"github.com/prashanthpai/sqlcache/cache"
)

// Redis implements cache.Cacher interface to use redis as backend with
// go-redis as the redis client library.
type Redis struct {
	c         redis.UniversalClient
	keyPrefix string
}

// Get gets a cache item from redis. Returns pointer to the item, a boolean
// which represents whether key exists or not and an error.
func (r *Redis) Get(ctx context.Context, key string) (*cache.Item, bool, error) {
	b, err := r.c.Get(ctx, r.keyPrefix+key).Bytes()
	switch err {
	case nil:
		var item cache.Item
		if err := msgpack.Unmarshal(b, &item); err != nil {
			return nil, true, err
		}
		return &item, true, nil
	case redis.Nil:
		return nil, false, nil
	default:
		return nil, false, err
	}
}

// Set sets the given item into redis with provided TTL duration.
func (r *Redis) Set(ctx context.Context, key string, item *cache.Item, ttl time.Duration) error {
	b, err := msgpack.Marshal(item)
	if err != nil {
		return err
	}

	_, err = r.c.Set(ctx, r.keyPrefix+key, b, ttl).Result()
	return err
}

// Del delete item from cache
func (r *Redis) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	newKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		newKeys = append(newKeys, r.keyPrefix+key)
	}
	return r.c.Del(ctx, newKeys...).Err()
}

// NewRedis creates a new instance of redis backend using go-redis client.
// All keys created in redis by sqlcache will have start with prefix.
func NewRedis(c redis.UniversalClient, keyPrefix string) *Redis {
	return &Redis{
		c:         c,
		keyPrefix: keyPrefix,
	}
}
