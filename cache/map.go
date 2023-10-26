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
	handle       Cache[K, V]
	mux          sync.Mutex
	cacheFunc    MapSingleFn[K, V]
	batchCacheFn MapBatchFn[K, V]
	expireTime   time.Duration
}

type MapSingleFn[K, V any] func(context.Context, K, ...any) (V, error)
type MapBatchFn[K comparable, V any] func(context.Context, []K, ...any) (map[K]V, error)

func NewMapCache[K comparable, V any](data Cache[K, V], cacheFunc MapSingleFn[K, V], batchCacheFn MapBatchFn[K, V], expireTime time.Duration) *MapCache[K, V] {
	r := &MapCache[K, V]{
		handle:       data,
		mux:          sync.Mutex{},
		cacheFunc:    cacheFunc,
		batchCacheFn: batchCacheFn,
		expireTime:   expireTime,
	}
	if cacheFunc == nil && batchCacheFn != nil {
		r.setDefaultCacheFn(batchCacheFn)
	} else if batchCacheFn == nil && cacheFunc != nil {
		r.SetDefaultBatchFunc(cacheFunc)
	}
	return r
}

func (m *MapCache[K, V]) SetDefaultBatchFunc(fn MapSingleFn[K, V]) {
	m.batchCacheFn = func(ctx context.Context, ids []K, a ...any) (map[K]V, error) {
		var err error
		rr := make(map[K]V)
		for _, id := range ids {
			v, er := fn(ctx, id)
			if er != nil {
				err = errors.Join(er)
				continue
			}
			rr[id] = v
		}
		return rr, err
	}
}

func (m *MapCache[K, V]) SetCacheFunc(fn MapSingleFn[K, V]) {
	m.cacheFunc = fn
}
func (m *MapCache[K, V]) GetHandle() Cache[K, V] {
	return m.handle
}

func (m *MapCache[K, V]) Ttl(ctx context.Context, k K) time.Duration {
	return m.handle.Ttl(ctx, k, m.expireTime)
}

func (m *MapCache[K, V]) GetLastSetTime(ctx context.Context, k K) (t time.Time) {
	tt := m.handle.Ttl(ctx, k, m.expireTime)
	if tt <= 0 {
		return
	}
	return time.Now().Add(m.handle.Ttl(ctx, k, m.expireTime)).Add(-m.expireTime)
}

func (m *MapCache[K, V]) SetCacheBatchFn(fn MapBatchFn[K, V]) {
	m.batchCacheFn = fn
	if m.cacheFunc == nil {
		m.setDefaultCacheFn(fn)
	}
}

func (m *MapCache[K, V]) setDefaultCacheFn(fn MapBatchFn[K, V]) {
	m.cacheFunc = func(ctx context.Context, k K, a ...any) (V, error) {
		var err error
		var r map[K]V
		r, err = fn(ctx, []K{k}, a...)

		if err != nil {
			var rr V
			return rr, err
		}
		return r[k], err
	}
}

func NewMapCacheByFn[K comparable, V any](cacheType Cache[K, V], fn MapSingleFn[K, V], expireTime time.Duration) *MapCache[K, V] {
	r := &MapCache[K, V]{
		mux:        sync.Mutex{},
		cacheFunc:  fn,
		expireTime: expireTime,
		handle:     cacheType,
	}
	r.SetDefaultBatchFunc(fn)
	return r
}
func NewMapCacheByBatchFn[K comparable, V any](cacheType Cache[K, V], fn MapBatchFn[K, V], expireTime time.Duration) *MapCache[K, V] {
	r := &MapCache[K, V]{
		mux:          sync.Mutex{},
		batchCacheFn: fn,
		expireTime:   expireTime,
		handle:       cacheType,
	}
	r.setDefaultCacheFn(fn)
	return r
}

func (m *MapCache[K, V]) Flush(ctx context.Context) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.handle.Flush(ctx)
}

func (m *MapCache[K, V]) Get(ctx context.Context, k K) (V, bool) {
	return m.handle.Get(ctx, k)
}

func (m *MapCache[K, V]) Set(ctx context.Context, k K, v V) {
	m.handle.Set(ctx, k, v, m.expireTime)
}

func (m *MapCache[K, V]) Delete(ctx context.Context, k K) {
	m.handle.Delete(ctx, k)
}
func (m *MapCache[K, V]) Ver(ctx context.Context, k K) int {
	return m.handle.Ver(ctx, k)
}

func (m *MapCache[K, V]) GetCache(c context.Context, key K, timeout time.Duration, params ...any) (V, error) {
	data, ok := m.handle.Get(c, key)
	var err error
	if !ok || m.handle.Ttl(c, key, m.expireTime) <= 0 {
		ver := m.handle.Ver(c, key)
		call := func() {
			m.mux.Lock()
			defer m.mux.Unlock()
			if m.handle.Ver(c, key) > ver {
				data, _ = m.handle.Get(c, key)
				return
			}
			data, err = m.cacheFunc(c, key, params...)
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
		if _, ok := m.handle.Get(c, k); !ok {
			return k, true
		}
		ver += m.handle.Ver(c, k)
		return
	})

	var err error
	if len(needFlush) > 0 {
		call := func() {
			m.mux.Lock()
			defer m.mux.Unlock()

			vers := slice.Reduce(needFlush, func(t K, r int) int {
				return r + m.handle.Ver(c, t)
			}, 0)

			if vers > ver {
				return
			}

			r, er := m.batchCacheFn(c, key, params...)
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
		return m.handle.Get(c, k)
	})
	return res, err
}

func (m *MapCache[K, V]) ClearExpired(ctx context.Context) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.handle.ClearExpired(ctx, m.expireTime)
}
