package cache

import (
	"context"
	"sync"
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

type SetTime interface {
	SetExpiredTime(func() time.Duration)
}

type AnyCache[T any] interface {
	Get(ctx context.Context) (T, bool)
	Set(ctx context.Context, v T)
	Flush(ctx context.Context)
	GetLastSetTime(ctx context.Context) time.Time
}

type Refresh[K comparable, V any] interface {
	Refresh(ctx context.Context, k K, a ...any)
}
type RefreshVar[T any] interface {
	Refresh(ctx context.Context, a ...any)
}

type Lockss[K comparable] interface {
	GetLock(ctx context.Context, gMut *sync.Mutex, k ...K) *sync.Mutex
}

type LockFn[K comparable] func(ctx context.Context, gMut *sync.Mutex, k ...K) *sync.Mutex

type LocksNum interface {
	SetLockNum(num int)
}
