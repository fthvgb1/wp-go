package cache

import (
	"context"
	"errors"
	"fmt"
	"github/fthvgb1/wp-go/safeMap"
	"sync"
	"time"
)

type MapCache[K comparable, V any] struct {
	data         safeMap.Map[K, mapCacheStruct[V]]
	mutex        *sync.Mutex
	cacheFunc    func(...any) (V, error)
	batchCacheFn func(...any) (map[K]V, error)
	expireTime   time.Duration
}

func NewMapCache[K comparable, V any](expireTime time.Duration) *MapCache[K, V] {
	return &MapCache[K, V]{expireTime: expireTime}
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
	r, ok := m.data.Load(k)
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
	return &MapCache[K, V]{
		mutex:      &sync.Mutex{},
		cacheFunc:  fn,
		expireTime: expireTime,
		data:       safeMap.NewMap[K, mapCacheStruct[V]](),
	}
}
func NewMapCacheByBatchFn[K comparable, V any](fn func(...any) (map[K]V, error), expireTime time.Duration) *MapCache[K, V] {
	r := &MapCache[K, V]{
		mutex:        &sync.Mutex{},
		batchCacheFn: fn,
		expireTime:   expireTime,
		data:         safeMap.NewMap[K, mapCacheStruct[V]](),
	}
	r.setCacheFn(fn)
	return r
}

func (m *MapCache[K, V]) Flush() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = safeMap.NewMap[K, mapCacheStruct[V]]()
}

func (m *MapCache[K, V]) Get(k K) V {
	r, ok := m.data.Load(k)
	if ok {
		return r.data
	}
	var rr V
	return rr
}

func (m *MapCache[K, V]) Set(k K, v V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
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
	data, ok := m.data.Load(k)
	t := time.Now()
	if !ok {
		data.data = v
		data.setTime = t
		data.incr++
		m.data.Store(k, data)
	} else {
		m.data.Store(k, mapCacheStruct[V]{
			data:    v,
			setTime: t,
		})

	}
}

func (m *MapCache[K, V]) GetCache(c context.Context, key K, timeout time.Duration, params ...any) (V, error) {
	data, ok := m.data.Load(key)
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
			m.mutex.Lock()
			defer m.mutex.Unlock()
			if data.incr > t {
				return
			}
			r, er := m.cacheFunc(params...)
			if err != nil {
				err = er
				return
			}
			data.setTime = time.Now()
			data.data = r
			m.data.Store(key, data)
			data.incr++
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
	for _, k := range key {
		d, ok := m.data.Load(k)
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
			m.mutex.Lock()
			defer m.mutex.Unlock()
			tt := 0
			for _, dd := range needFlush {
				if ddd, ok := m.data.Load(dd); ok {
					tt = tt + ddd.incr
				}
			}
			if tt > t {
				return
			}
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
		d, ok := m.data.Load(k)
		if ok {
			res = append(res, d.data)
		}
	}
	return res, err
}

func (m *MapCache[K, V]) ClearExpired() {
	now := time.Duration(time.Now().UnixNano())
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data.Range(func(k K, v mapCacheStruct[V]) bool {
		if now > time.Duration(v.setTime.UnixNano())+m.expireTime {
			m.data.Delete(k)
		}
		return true
	})
}
