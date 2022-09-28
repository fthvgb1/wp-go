package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type MapCache[K comparable, V any] struct {
	data            map[K]mapCacheStruct[V]
	mutex           *sync.Mutex
	setCacheFunc    func(...any) (V, error)
	setBatchCacheFn func(...any) (map[K]V, error)
	expireTime      time.Duration
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
	m.setCacheFunc = fn
}

func (m *MapCache[K, V]) SetCacheBatchFunc(fn func(...any) (map[K]V, error)) {
	m.setBatchCacheFn = fn
}

func NewMapCacheByFn[K comparable, V any](fun func(...any) (V, error), expireTime time.Duration) *MapCache[K, V] {
	return &MapCache[K, V]{
		mutex:        &sync.Mutex{},
		setCacheFunc: fun,
		expireTime:   expireTime,
		data:         make(map[K]mapCacheStruct[V]),
	}
}
func NewMapCacheByBatchFn[K comparable, V any](fn func(...any) (map[K]V, error), expireTime time.Duration) *MapCache[K, V] {
	return &MapCache[K, V]{
		mutex:           &sync.Mutex{},
		setBatchCacheFn: fn,
		expireTime:      expireTime,
		data:            make(map[K]mapCacheStruct[V]),
	}
}

func (m *MapCache[K, V]) Flush() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data = make(map[K]mapCacheStruct[V])
}

func (m *MapCache[K, V]) Get(k K) V {
	return m.data[k].data
}

func (m *MapCache[K, V]) Set(k K, v V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.set(k, v)
}

func (m *MapCache[K, V]) SetByBatchFn(params ...any) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	r, err := m.setBatchCacheFn(params...)
	if err != nil {
		return err
	}
	for k, v := range r {
		m.set(k, v)
	}
	return nil
}

func (m *MapCache[K, V]) set(k K, v V) {
	data, ok := m.data[k]
	t := time.Now()
	if !ok {
		data.data = v
		data.setTime = t
		data.incr++
		m.data[k] = data
	} else {
		m.data[k] = mapCacheStruct[V]{
			data:    v,
			setTime: t,
		}
	}
}

func (m *MapCache[K, V]) GetCache(c context.Context, key K, timeout time.Duration, params ...any) (V, error) {
	data, ok := m.data[key]
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
			r, er := m.setCacheFunc(params...)
			if err != nil {
				err = er
				return
			}
			data.setTime = time.Now()
			data.data = r
			m.data[key] = data
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
		d, ok := m.data[k]
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
				if ddd, ok := m.data[dd]; ok {
					tt = tt + ddd.incr
				}
			}
			if tt > t {
				return
			}
			r, er := m.setBatchCacheFn(params...)
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
		d := m.data[k]
		res = append(res, d.data)
	}
	return res, err
}

func (m *MapCache[K, V]) ClearExpired() {
	now := time.Duration(time.Now().UnixNano())
	m.mutex.Lock()
	defer m.mutex.Unlock()
	for k, v := range m.data {
		if now > time.Duration(v.setTime.UnixNano())+m.expireTime {
			delete(m.data, k)
		}
	}
}
