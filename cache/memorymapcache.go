package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/helper/number"
	"github.com/fthvgb1/wp-go/safety"
	"time"
)

type MemoryMapCache[K comparable, V any] struct {
	safety.Map[K, mapVal[V]]
}

func NewMemoryMapCache[K comparable, V any]() *MemoryMapCache[K, V] {
	return &MemoryMapCache[K, V]{Map: safety.NewMap[K, mapVal[V]]()}
}

type mapVal[T any] struct {
	setTime time.Time
	ver     int
	data    T
}

func (m *MemoryMapCache[K, V]) Get(_ context.Context, key K) (r V, ok bool) {
	v, ok := m.Load(key)
	if ok {
		return v.data, true
	}
	return
}

func (m *MemoryMapCache[K, V]) Set(_ context.Context, key K, val V, _ time.Duration) {
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

func (m *MemoryMapCache[K, V]) Ttl(_ context.Context, key K, expire time.Duration) time.Duration {
	v, ok := m.Load(key)
	if !ok {
		return -1
	}
	return number.Max(time.Duration(0), expire-time.Duration(time.Now().UnixNano()-v.setTime.UnixNano()))
}

func (m *MemoryMapCache[K, V]) Ver(_ context.Context, key K) int {
	v, ok := m.Load(key)
	if !ok {
		return -1
	}
	return v.ver
}

func (m *MemoryMapCache[K, V]) Flush(context.Context) {
	m.Map = safety.NewMap[K, mapVal[V]]()
}

func (m *MemoryMapCache[K, V]) Delete(_ context.Context, key K) {
	m.Map.Delete(key)
}

func (m *MemoryMapCache[K, V]) ClearExpired(_ context.Context, expire time.Duration) {
	now := time.Duration(time.Now().UnixNano())

	m.Range(func(k K, v mapVal[V]) bool {
		if now > time.Duration(v.setTime.UnixNano())+expire {
			m.Map.Delete(k)
		}
		return true
	})
}
