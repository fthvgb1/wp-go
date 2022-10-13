package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type MapCache[K comparable, V any] struct {
	data         atomic.Value
	mutex        *sync.Mutex
	cacheFunc    func(...any) (V, error)
	batchCacheFn func(...any) (map[K]V, error)
	expireTime   time.Duration
}

func NewMapCache[K comparable, V any](expireTime time.Duration) *MapCache[K, V] {
	var v atomic.Value
	v.Store(make(map[K]mapCacheStruct[V]))
	return &MapCache[K, V]{expireTime: expireTime, data: v}
}

type mapCacheStruct[T any] struct {
	setTime time.Time
	incr    int
	data    T
}

func (m *MapCache[K, V]) SetCacheFunc(fn func(...any) (V, error)) {
	m.cacheFunc = fn
}

func (m *MapCache[K, V]) GetSetTime(k K) (t time.Time) {
	r, ok := m.data.Load().(map[K]mapCacheStruct[V])[k]
	if ok {
		t = r.setTime
	}
	return
}

func (m *MapCache[K, V]) SetCacheBatchFunc(fn func(...any) (map[K]V, error)) {
	m.batchCacheFn = fn
	if m.cacheFunc == nil {
		m.setCacheFn(fn)
	}
}

func (m *MapCache[K, V]) setCacheFn(fn func(...any) (map[K]V, error)) {
	m.cacheFunc = func(a ...any) (V, error) {
		id := a[0].(K)
		r, err := fn([]K{id})
		if err != nil {
			var rr V
			return rr, err
		}
		return r[id], err
	}
}

func NewMapCacheByFn[K comparable, V any](fn func(...any) (V, error), expireTime time.Duration) *MapCache[K, V] {
	var d atomic.Value
	d.Store(make(map[K]mapCacheStruct[V]))
	return &MapCache[K, V]{
		mutex:      &sync.Mutex{},
		cacheFunc:  fn,
		expireTime: expireTime,
		data:       d,
	}
}
func NewMapCacheByBatchFn[K comparable, V any](fn func(...any) (map[K]V, error), expireTime time.Duration) *MapCache[K, V] {
	var d atomic.Value
	d.Store(make(map[K]mapCacheStruct[V]))
	r := &MapCache[K, V]{
		mutex:        &sync.Mutex{},
		batchCacheFn: fn,
		expireTime:   expireTime,
		data:         d,
	}
	r.setCacheFn(fn)
	return r
}

func (m *MapCache[K, V]) Flush() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	var d atomic.Value
	d.Store(make(map[K]mapCacheStruct[V]))
	m.data = d
}

func (m *MapCache[K, V]) Get(k K) V {
	return m.data.Load().(map[K]mapCacheStruct[V])[k].data
}

func (m *MapCache[K, V]) Set(k K, v V) {
	m.set(k, v)
}

func (m *MapCache[K, V]) SetByBatchFn(params ...any) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	r, err := m.batchCacheFn(params...)
	if err != nil {
		return err
	}
	for k, v := range r {
		m.set(k, v)
	}
	return nil
}

func (m *MapCache[K, V]) set(k K, v V) {
	d, ok := m.data.Load().(map[K]mapCacheStruct[V])
	t := time.Now()
	data := d[k]
	if !ok {
		data.data = v
		data.setTime = t
		data.incr++
	} else {
		data = mapCacheStruct[V]{
			data:    v,
			setTime: t,
		}
	}
	d[k] = data
	m.data.Store(d)
}

func (m *MapCache[K, V]) GetCache(c context.Context, key K, timeout time.Duration, params ...any) (V, error) {
	d := m.data.Load().(map[K]mapCacheStruct[V])
	data, ok := d[key]
	if !ok {
		data = mapCacheStruct[V]{}
	}
	now := time.Duration(time.Now().UnixNano())
	var err error
	expired := time.Duration(data.setTime.UnixNano())+m.expireTime < now
	//todo 这里应该判断下取出的值是否为零值，不过怎么操作呢？
	if !ok || (ok && m.expireTime >= 0 && expired) {
		t := data.incr
		call := func() {
			tmp, o := m.data.Load().(map[K]mapCacheStruct[V])[key]
			if o && tmp.incr > t {
				return
			}
			m.mutex.Lock()
			defer m.mutex.Unlock()
			r, er := m.cacheFunc(params...)
			if err != nil {
				err = er
				return
			}
			data.setTime = time.Now()
			data.data = r
			data.incr++
			d[key] = data
			m.data.Store(d)
		}
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(c, timeout)
			defer cancel()
			done := make(chan struct{})
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
	return data.data, err
}

func (m *MapCache[K, V]) GetCacheBatch(c context.Context, key []K, timeout time.Duration, params ...any) ([]V, error) {
	var needFlush []K
	var res []V
	t := 0
	now := time.Duration(time.Now().UnixNano())
	data := m.data.Load().(map[K]mapCacheStruct[V])
	for _, k := range key {
		d, ok := data[k]
		if !ok {
			needFlush = append(needFlush, k)
			continue
		}
		expired := time.Duration(d.setTime.UnixNano())+m.expireTime < now
		if expired {
			needFlush = append(needFlush, k)
		}
		t = t + d.incr
	}
	var err error
	//todo 这里应该判断下取出的值是否为零值，不过怎么操作呢？
	if len(needFlush) > 0 {
		call := func() {
			tt := 0
			for _, dd := range needFlush {
				if ddd, ok := data[dd]; ok {
					tt = tt + ddd.incr
				}
			}
			if tt > t {
				return
			}
			m.mutex.Lock()
			defer m.mutex.Unlock()
			r, er := m.batchCacheFn(params...)
			if err != nil {
				err = er
				return
			}
			for k, v := range r {
				m.set(k, v)
			}
		}
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(c, timeout)
			defer cancel()
			done := make(chan struct{})
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
	for _, k := range key {
		d := data[k]
		res = append(res, d.data)
	}
	return res, err
}

func (m *MapCache[K, V]) ClearExpired() {
	now := time.Duration(time.Now().UnixNano())
	m.mutex.Lock()
	defer m.mutex.Unlock()
	data := m.data.Load().(map[K]mapCacheStruct[V])
	for k, v := range data {
		if now > time.Duration(v.setTime.UnixNano())+m.expireTime {
			delete(data, k)
		}
	}
	m.data.Store(data)
}
