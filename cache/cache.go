package cache

import (
	"context"
	"database/sql/driver"
	"time"
)

// Item represents a single item in cache and will contain the results of a
// single SQL query.
type Item struct {
	Cols              []string
	DatabaseTypeNames []string
	Rows              [][]driver.Value
	Query             string
	QueryAt           int64
}

// Cacher represents a backend cache that can be used by sqlcache package.
type Cacher interface {
	// Get must return a pointer to the item, a boolean representing whether
	// item is present or not, and an error (must be nil when key is not
	// present).
	Get(ctx context.Context, key string) (*Item, bool, error)
	// Set sets the item into cache with the given TTL.
	Set(ctx context.Context, key string, item *Item, ttl time.Duration) error
	// Del delete item from cache
	Del(ctx context.Context, keys ...string) error
}
