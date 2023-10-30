package cache

import (
	"context"
	"time"
)

type Cache[K comparable, V any] interface {
	Get(ctx context.Context, key K) (V, bool)
	Set(ctx context.Context, key K, val V)
	GetExpireTime(ctx context.Context) time.Duration
	Ttl(ctx context.Context, key K) time.Duration
	Flush(ctx context.Context)
	Del(ctx context.Context, key ...K)
	ClearExpired(ctx context.Context)
}

type Expend[K comparable, V any] interface {
	Gets(ctx context.Context, k []K) (map[K]V, error)
	Sets(ctx context.Context, m map[K]V)
}
