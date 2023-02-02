package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	"sync"
	"time"
)

type MapCache[K comparable, V any] struct {
	data         Cache[K, V]
	mux          sync.Mutex
	cacheFunc    func(...any) (V, error)
	batchCacheFn func(...any) (map[K]V, error)
	expireTime   time.Duration
}

func (m *MapCache[K, V]) SetCacheFunc(fn func(...any) (V, error)) {
	m.cacheFunc = fn
}

func (m *MapCache[K, V]) GetLastSetTime(ctx context.Context, k K) (t time.Time) {
	tt := m.data.Ttl(ctx, k, m.expireTime)
	if tt <= 0 {
		return
	}
	return time.Now().Add(m.data.Ttl(ctx, k, m.expireTime))
}

func (m *MapCache[K, V]) SetCacheBatchFn(fn func(...any) (map[K]V, error)) {
	m.batchCacheFn = fn
	if m.cacheFunc == nil {
		m.setCacheFn(fn)
	}
}

func (m *MapCache[K, V]) setCacheFn(fn func(...any) (map[K]V, error)) {
	m.cacheFunc = func(a ...any) (V, error) {
		var err error
		var r map[K]V
		var id K
		ctx, ok := a[0].(context.Context)
		if ok {
			id = a[1].(K)
			r, err = fn(ctx, []K{id})
		} else {
			id = a[0].(K)
			r, err = fn([]K{id})
		}

		if err != nil {
			var rr V
			return rr, err
		}
		return r[id], err
	}
}

func NewMemoryMapCacheByFn[K comparable, V any](fn func(...any) (V, error), expireTime time.Duration) *MapCache[K, V] {
	return &MapCache[K, V]{
		data:       NewMemoryMapCache[K, V](),
		cacheFunc:  fn,
		expireTime: expireTime,
		mux:        sync.Mutex{},
	}
}
func NewMemoryMapCacheByBatchFn[K comparable, V any](fn func(...any) (map[K]V, error), expireTime time.Duration) *MapCache[K, V] {
	r := &MapCache[K, V]{
		data:         NewMemoryMapCache[K, V](),
		batchCacheFn: fn,
		expireTime:   expireTime,
		mux:          sync.Mutex{},
	}
	r.setCacheFn(fn)
	return r
}

func (m *MapCache[K, V]) Flush(ctx context.Context) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.data.Flush(ctx)
}

func (m *MapCache[K, V]) Get(ctx context.Context, k K) (V, bool) {
	return m.data.Get(ctx, k)
}

func (m *MapCache[K, V]) Set(ctx context.Context, k K, v V) {
	m.data.Set(ctx, k, v, m.expireTime)
}

func (m *MapCache[K, V]) GetCache(c context.Context, key K, timeout time.Duration, params ...any) (V, error) {
	data, ok := m.data.Get(c, key)
	var err error
	if !ok || m.data.Ttl(c, key, m.expireTime) <= 0 {
		ver := m.data.Ver(c, key)
		call := func() {
			m.mux.Lock()
			defer m.mux.Unlock()
			if m.data.Ver(c, key) > ver {
				data, _ = m.data.Get(c, key)
				return
			}
			data, err = m.cacheFunc(params...)
			if err != nil {
				return
			}
			m.Set(c, key, data)
		}
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(c, timeout)
			defer cancel()
			done := make(chan struct{}, 1)
			go func() {
				call()
				done <- struct{}{}
			}()
			select {
			case <-ctx.Done():
				err = errors.New(fmt.Sprintf("get cache %v %s", key, ctx.Err().Error()))
			case <-done:
			}
		} else {
			call()
		}

	}
	return data, err
}

func (m *MapCache[K, V]) GetCacheBatch(c context.Context, key []K, timeout time.Duration, params ...any) ([]V, error) {
	var res []V
	ver := 0
	needFlush := slice.FilterAndMap(key, func(k K) (r K, ok bool) {
		if _, ok := m.data.Get(c, k); !ok {
			return k, true
		}
		ver += m.data.Ver(c, k)
		return
	})

	var err error
	if len(needFlush) > 0 {
		call := func() {
			m.mux.Lock()
			defer m.mux.Unlock()

			vers := slice.Reduce(needFlush, func(t K, r int) int {
				r += m.data.Ver(c, t)
				return r
			}, 0)

			if vers > ver {
				return
			}

			r, er := m.batchCacheFn(params...)
			if err != nil {
				err = er
				return
			}
			for k, v := range r {
				m.Set(c, k, v)
			}
		}
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(c, timeout)
			defer cancel()
			done := make(chan struct{}, 1)
			go func() {
				call()
				done <- struct{}{}
			}()
			select {
			case <-ctx.Done():
				err = errors.New(fmt.Sprintf("get cache %v %s", key, ctx.Err().Error()))
			case <-done:
			}
		} else {
			call()
		}
	}
	res = slice.FilterAndMap(key, func(k K) (V, bool) {
		return m.data.Get(c, k)
	})
	return res, err
}

func (m *MapCache[K, V]) ClearExpired(ctx context.Context) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.data.ClearExpired(ctx, m.expireTime)
}
