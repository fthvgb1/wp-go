package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/maps"
	"sync"
	"time"
)

type MapCache[K comparable, V any] struct {
	Cache[K, V]
	mux                sync.Mutex
	cacheFunc          MapSingleFn[K, V]
	batchCacheFn       MapBatchFn[K, V]
	getCacheBatch      func(c context.Context, key []K, timeout time.Duration, params ...any) ([]V, error)
	getCacheBatchToMap func(c context.Context, key []K, timeout time.Duration, params ...any) (map[K]V, error)
	increaseUpdate     IncreaseUpdate[K, V]
	refresh            Refresh[K, V]
}
type IncreaseUpdate[K comparable, V any] struct {
	CycleTime func() time.Duration
	Fn        IncreaseFn[K, V]
}

func NewIncreaseUpdate[K comparable, V any](name string, fn IncreaseFn[K, V], cycleTime time.Duration, tFn func() time.Duration) IncreaseUpdate[K, V] {
	tFn = reload.FnVal(name, cycleTime, tFn)
	return IncreaseUpdate[K, V]{CycleTime: tFn, Fn: fn}
}

type MapSingleFn[K, V any] func(context.Context, K, ...any) (V, error)
type MapBatchFn[K comparable, V any] func(context.Context, []K, ...any) (map[K]V, error)
type IncreaseFn[K comparable, V any] func(c context.Context, currentData V, k K, t time.Time, a ...any) (data V, save bool, refresh bool, err error)

func NewMapCache[K comparable, V any](ca Cache[K, V], cacheFunc MapSingleFn[K, V], batchCacheFn MapBatchFn[K, V], inc IncreaseUpdate[K, V]) *MapCache[K, V] {
	r := &MapCache[K, V]{
		Cache:          ca,
		mux:            sync.Mutex{},
		cacheFunc:      cacheFunc,
		batchCacheFn:   batchCacheFn,
		increaseUpdate: inc,
	}
	if cacheFunc == nil && batchCacheFn != nil {
		r.setDefaultCacheFn(batchCacheFn)
	} else if batchCacheFn == nil && cacheFunc != nil {
		r.SetDefaultBatchFunc(cacheFunc)
	}
	ex, ok := any(ca).(Expend[K, V])
	if !ok {
		r.getCacheBatch = r.getCacheBatchs
		r.getCacheBatchToMap = r.getBatchToMapes
	} else {
		r.getCacheBatch = r.getBatches(ex)
		r.getCacheBatchToMap = r.getBatchToMap(ex)
	}
	re, ok := any(ca).(Refresh[K, V])
	if ok {
		r.refresh = re
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
	if ok {
		if m.increaseUpdate.Fn == nil || m.refresh == nil {
			return data, err
		}
		nowTime := time.Now()
		if nowTime.Sub(m.GetLastSetTime(c, key)) < m.increaseUpdate.CycleTime() {
			return data, err
		}
		fn := func() {
			m.mux.Lock()
			defer m.mux.Unlock()
			if nowTime.Sub(m.GetLastSetTime(c, key)) < m.increaseUpdate.CycleTime() {
				return
			}
			dat, save, refresh, er := m.increaseUpdate.Fn(c, data, key, m.GetLastSetTime(c, key), params...)
			if er != nil {
				err = er
				return
			}
			if refresh {
				m.refresh.Refresh(c, key, params...)
			}
			if save {
				m.Set(c, key, dat)
				data = dat
			}
		}
		if timeout > 0 {
			er := helper.RunFnWithTimeout(c, timeout, fn, fmt.Sprintf("increateUpdate cache %v err", key))
			if err == nil && er != nil {
				return data, er
			}
		} else {
			fn()
		}
		return data, err
	}
	call := func() {
		m.mux.Lock()
		defer m.mux.Unlock()
		if data, ok = m.Get(c, key); ok {
			return
		}
		data, err = m.cacheFunc(c, key, params...)
		if err != nil {
			return
		}
		m.Set(c, key, data)
	}
	if timeout > 0 {
		er := helper.RunFnWithTimeout(c, timeout, call, fmt.Sprintf("get cache %v ", key))
		if err == nil && er != nil {
			err = er
		}
	} else {
		call()
	}
	return data, err
}

func (m *MapCache[K, V]) GetCacheBatch(c context.Context, key []K, timeout time.Duration, params ...any) ([]V, error) {
	return m.getCacheBatch(c, key, timeout, params...)
}

func (m *MapCache[K, V]) GetBatchToMap(c context.Context, key []K, timeout time.Duration, params ...any) (map[K]V, error) {
	return m.getCacheBatchToMap(c, key, timeout, params...)
}
func (m *MapCache[K, V]) getBatchToMap(e Expend[K, V]) func(c context.Context, key []K, timeout time.Duration, params ...any) (map[K]V, error) {
	return func(ctx context.Context, key []K, timeout time.Duration, params ...any) (map[K]V, error) {
		var res map[K]V
		var err error
		mm, err := e.Gets(ctx, key)
		if err != nil || len(key) == len(mm) {
			return mm, err
		}
		var needIndex = make(map[K]int)
		res = mm
		var flushKeys []K
		for i, k := range key {
			_, ok := mm[k]
			if !ok {
				flushKeys = append(flushKeys, k)
				needIndex[k] = i
			}
		}

		call := func() {
			m.mux.Lock()
			defer m.mux.Unlock()
			mmm, er := e.Gets(ctx, maps.FilterToSlice(needIndex, func(k K, v int) (K, bool) {
				return k, true
			}))
			if er != nil {
				err = er
				return
			}
			for k, v := range mmm {
				res[k] = v
				delete(needIndex, k)
			}

			if len(needIndex) < 1 {
				return
			}

			r, er := m.batchCacheFn(ctx, maps.FilterToSlice(needIndex, func(k K, v int) (K, bool) {
				return k, true
			}), params...)
			if er != nil {
				err = er
				return
			}
			e.Sets(ctx, r)

			for k := range needIndex {
				v, ok := r[k]
				if ok {
					res[k] = v
				}
			}
		}
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			done := make(chan struct{}, 1)
			go func() {
				call()
				done <- struct{}{}
			}()
			select {
			case <-ctx.Done():
				err = errors.New(fmt.Sprintf("get cache %v %s", key, ctx.Err().Error()))
				return nil, err
			case <-done:
			}
		} else {
			call()
		}

		return res, err
	}
}
func (m *MapCache[K, V]) getBatchToMapes(c context.Context, key []K, timeout time.Duration, params ...any) (r map[K]V, err error) {
	r = make(map[K]V)
	var needIndex = make(map[K]int)
	for i, k := range key {
		v, ok := m.Get(c, k)
		if !ok {
			needIndex[k] = i
		} else {
			r[k] = v
		}
	}
	if len(needIndex) < 1 {
		return
	}

	call := func() {
		m.mux.Lock()
		defer m.mux.Unlock()
		needFlushs := maps.FilterToSlice(needIndex, func(k K, v int) (K, bool) {
			vv, ok := m.Get(c, k)
			if ok {
				r[k] = vv
				delete(needIndex, k)
				return k, false
			}
			return k, true
		})

		if len(needFlushs) < 1 {
			return
		}

		rr, er := m.batchCacheFn(c, needFlushs, params...)
		if er != nil {
			err = er
			return
		}
		for k := range needIndex {
			v, ok := rr[k]
			if ok {
				r[k] = v
			}
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
			return nil, err
		case <-done:
		}
	} else {
		call()
	}
	return
}

func (m *MapCache[K, V]) getCacheBatchs(c context.Context, key []K, timeout time.Duration, params ...any) ([]V, error) {
	var res = make([]V, 0, len(key))
	var needIndex = make(map[K]int)
	for i, k := range key {
		v, ok := m.Get(c, k)
		if !ok {
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
			vv, ok := m.Get(c, k)
			if ok {
				res[needIndex[k]] = vv
				delete(needIndex, k)
				return k, false
			}
			return k, true
		})

		if len(needFlushs) < 1 {
			return
		}

		r, er := m.batchCacheFn(c, needFlushs, params...)
		if er != nil {
			err = er
			return
		}
		for k, i := range needIndex {
			v, ok := r[k]
			if ok {
				res[i] = v
				m.Set(c, k, v)
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
			return nil, err
		case <-done:
		}
	} else {
		call()
	}

	return res, err
}

func (m *MapCache[K, V]) getBatches(e Expend[K, V]) func(ctx context.Context, key []K, timeout time.Duration, params ...any) ([]V, error) {
	cc := e
	return func(ctx context.Context, key []K, timeout time.Duration, params ...any) ([]V, error) {
		var res = make([]V, 0, len(key))
		var needIndex = make(map[K]int)
		var err error
		mm, err := cc.Gets(ctx, key)
		if err != nil {
			return nil, err
		}
		var flushKeys []K
		for i, k := range key {
			v, ok := mm[k]
			if !ok {
				flushKeys = append(flushKeys, k)
				needIndex[k] = i
				var vv V
				v = vv
			}
			res = append(res, v)
		}
		if len(needIndex) < 1 {
			return res, nil
		}

		call := func() {
			m.mux.Lock()
			defer m.mux.Unlock()
			mmm, er := cc.Gets(ctx, maps.FilterToSlice(needIndex, func(k K, v int) (K, bool) {
				return k, true
			}))
			if er != nil {
				err = er
				return
			}
			for k, v := range mmm {
				res[needIndex[k]] = v
				delete(needIndex, k)
			}

			if len(needIndex) < 1 {
				return
			}

			r, er := m.batchCacheFn(ctx, maps.FilterToSlice(needIndex, func(k K, v int) (K, bool) {
				return k, true
			}), params...)
			if er != nil {
				err = er
				return
			}
			cc.Sets(ctx, r)

			for k, i := range needIndex {
				v, ok := r[k]
				if ok {
					res[i] = v
				}
			}
		}
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			done := make(chan struct{}, 1)
			go func() {
				call()
				done <- struct{}{}
			}()
			select {
			case <-ctx.Done():
				err = errors.New(fmt.Sprintf("get cache %v %s", key, ctx.Err().Error()))
				return nil, err
			case <-done:
			}
		} else {
			call()
		}

		return res, err
	}
}
