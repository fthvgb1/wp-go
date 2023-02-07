package cache

import (
	"context"
	"time"
)

type Cache[K comparable, V any] interface {
	Get(ctx context.Context, key K) (V, bool)
	Set(ctx context.Context, key K, val V, expire time.Duration)
	Ttl(ctx context.Context, key K, expire time.Duration) time.Duration
	Ver(ctx context.Context, key K) int
	Flush(ctx context.Context)
	Delete(ctx context.Context, key K)
	ClearExpired(ctx context.Context, expire time.Duration)
}
