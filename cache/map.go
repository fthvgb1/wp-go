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

type mapCacheStruct[T any] struct {
	setTime time.Time
	incr    int
	data    T
}

func NewMapCache[K comparable, V any](fun func(...any) (V, error), expireTime time.Duration) *MapCache[K, V] {
	return &MapCache[K, V]{
		mutex:        &sync.Mutex{},
		setCacheFunc: fun,
		expireTime:   expireTime,
		data:         make(map[K]mapCacheStruct[V]),
	}
}
func NewMapBatchCache[K comparable, V any](fn func(...any) (map[K]V, error), expireTime time.Duration) *MapCache[K, V] {
	return &MapCache[K, V]{
		mutex:           &sync.Mutex{},
		setBatchCacheFn: fn,
		expireTime:      expireTime,
		data:            make(map[K]mapCacheStruct[V]),
	}
}

func (m *MapCache[K, V]) FlushCache(k any) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	key := k.(K)
	delete(m.data, key)
}

func (m *MapCache[K, V]) Get(k K) V {
	return m.data[k].data
}

func (m *MapCache[K, V]) Set(k K, v V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.set(k, v)
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
	var err error
	expired := time.Duration(data.setTime.Unix())+m.expireTime/time.Second < time.Duration(time.Now().Unix())
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

func (m *MapCache[K, V]) GetCacheBatch(c context.Context, key K, timeout time.Duration, params ...any) (V, error) {
	data, ok := m.data[key]
	if !ok {
		data = mapCacheStruct[V]{}
	}
	var err error
	expired := time.Duration(data.setTime.Unix())+m.expireTime/time.Second < time.Duration(time.Now().Unix())
	//todo 这里应该判断下取出的值是否为零值，不过怎么操作呢？
	if !ok || (ok && m.expireTime >= 0 && expired) {
		t := data.incr
		call := func() {
			m.mutex.Lock()
			defer m.mutex.Unlock()
			if data.incr > t {
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
			data.data = m.data[key].data
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
