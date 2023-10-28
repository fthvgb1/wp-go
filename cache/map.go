package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"sync"
	"time"
)

type MapCache[K comparable, V any] struct {
	Cache[K, V]
	mux          sync.Mutex
	cacheFunc    MapSingleFn[K, V]
	batchCacheFn MapBatchFn[K, V]
}

type MapSingleFn[K, V any] func(context.Context, K, ...any) (V, error)
type MapBatchFn[K comparable, V any] func(context.Context, []K, ...any) (map[K]V, error)

func NewMapCache[K comparable, V any](data Cache[K, V], cacheFunc MapSingleFn[K, V], batchCacheFn MapBatchFn[K, V]) *MapCache[K, V] {
	r := &MapCache[K, V]{
		Cache:        data,
		mux:          sync.Mutex{},
		cacheFunc:    cacheFunc,
		batchCacheFn: batchCacheFn,
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
			v, er := fn(ctx, id, a...)
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

func (m *MapCache[K, V]) GetLastSetTime(ctx context.Context, k K) (t time.Time) {
	tt := m.Ttl(ctx, k)
	if tt <= 0 {
		return
	}
	return time.Now().Add(m.Ttl(ctx, k)).Add(-m.GetExpireTime(ctx))
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

func (m *MapCache[K, V]) Flush(ctx context.Context) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.Cache.Flush(ctx)
}

func (m *MapCache[K, V]) GetCache(c context.Context, key K, timeout time.Duration, params ...any) (V, error) {
	data, ok := m.Get(c, key)
	var err error
	if !ok {
		ver := m.Ver(c, key)
		call := func() {
			m.mux.Lock()
			defer m.mux.Unlock()
			if m.Ver(c, key) > ver {
				data, _ = m.Get(c, key)
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
	var res = make([]V, 0, len(key))
	var needIndex = make(map[K]int)
	var ver = make(map[K]int)
	for i, k := range key {
		v, ok := m.Get(c, k)
		if !ok {
			ver[k] = m.Ver(c, k)
			needIndex[k] = i
		}
		res = append(res, v)
	}
	if len(needIndex) < 1 {
		return res, nil
	}

	var err error
	call := func() {
		m.mux.Lock()
		defer m.mux.Unlock()
		needFlushs := maps.FilterToSlice(needIndex, func(k K, v int) (K, bool) {
			return k, ver[k] >= m.Ver(c, k)
		})

		if len(needFlushs) < 1 {
			for k, i := range needIndex {
				res[i], _ = m.Get(c, k)
			}
			return
		}

		r, er := m.batchCacheFn(c, needFlushs, params...)
		if err != nil {
			err = er
			return
		}
		for k, i := range needIndex {
			v, ok := r[k]
			if ok {
				res[i] = v
				m.Set(c, k, v)
			} else {
				v, ok = m.Get(c, k)
				if ok {
					res[i] = v
				}
			}
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

	return res, err
}
