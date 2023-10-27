package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/safety"
	"sync"
	"time"
)

type MemoryMapCache[K comparable, V any] struct {
	*safety.Map[K, mapVal[V]]
	expireTime time.Duration
}

func NewMemoryMapCacheByFn[K comparable, V any](fn MapSingleFn[K, V], expireTime time.Duration) *MapCache[K, V] {
	return &MapCache[K, V]{
		Cache:     NewMemoryMapCache[K, V](expireTime),
		cacheFunc: fn,
		mux:       sync.Mutex{},
	}
}

func NewMemoryMapCache[K comparable, V any](expireTime time.Duration) *MemoryMapCache[K, V] {
	return &MemoryMapCache[K, V]{
		Map:        safety.NewMap[K, mapVal[V]](),
		expireTime: expireTime,
	}
}

type mapVal[T any] struct {
	setTime time.Time
	ver     int
	data    T
}

func (m *MemoryMapCache[K, V]) GetExpireTime(_ context.Context) time.Duration {
	return m.expireTime
}

func (m *MemoryMapCache[K, V]) Get(_ context.Context, key K) (r V, ok bool) {
	v, ok := m.Load(key)
	if !ok {
		return
	}
	r = v.data
	t := m.expireTime - time.Now().Sub(v.setTime)
	if t <= 0 {
		ok = false
	}
	return
}

func (m *MemoryMapCache[K, V]) Set(_ context.Context, key K, val V) {
	v, ok := m.Load(key)
	t := time.Now()
	if ok {
		v.data = val
		v.setTime = t
		v.ver++
	} else {
		v = mapVal[V]{
			setTime: t,
			ver:     1,
			data:    val,
		}
	}
	m.Store(key, v)
}

func (m *MemoryMapCache[K, V]) Ttl(_ context.Context, key K) time.Duration {
	v, ok := m.Load(key)
	if !ok {
		return time.Duration(-1)
	}
	return m.expireTime - time.Now().Sub(v.setTime)
}

func (m *MemoryMapCache[K, V]) Ver(_ context.Context, key K) int {
	v, ok := m.Load(key)
	if !ok {
		return -1
	}
	return v.ver
}

func (m *MemoryMapCache[K, V]) Flush(context.Context) {
	m.Map.Flush()
}

func (m *MemoryMapCache[K, V]) Del(_ context.Context, keys ...K) {
	for _, key := range keys {
		m.Map.Delete(key)
	}
}

func (m *MemoryMapCache[K, V]) ClearExpired(_ context.Context) {
	now := time.Duration(time.Now().UnixNano())
	m.Range(func(k K, v mapVal[V]) bool {
		if now > time.Duration(v.setTime.UnixNano())+m.expireTime {
			m.Map.Delete(k)
		}
		return true
	})
}
