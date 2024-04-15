package cachemanager

import (
	"context"
	"errors"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/helper"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"time"
)

var mapDelFuncs = safety.NewMap[string, func(any)]()

var mapCache = safety.NewMap[string, any]()

func SetMapCache[K comparable, V any](name string, ca *cache.MapCache[K, V]) error {
	v, ok := mapCache.Load(name)
	if !ok {
		mapCache.Store(name, ca)
		return nil
	}
	_, ok = v.(*cache.MapCache[K, V])
	if !ok {
		return errors.New(str.Join("cache ", name, " type err"))
	}
	mapCache.Store(name, ca)
	return nil
}

// PushMangerMap will del mapCache val with name When call DelMapCacheVal
func PushMangerMap[K comparable, V any](name string, m *cache.MapCache[K, V]) {
	if name == "" {
		return
	}
	mapCache.Store(name, m)
	mapDelFuncs.Store(name, func(a any) {
		k, ok := a.([]K)
		if ok && len(k) > 0 {
			mm, ok := mapCache.Load(name)
			if !ok {
				return
			}
			c, ok := mm.(*cache.MapCache[K, V])
			if !ok {
				return
			}
			ctx := context.WithValue(context.Background(), "ctx", "registerFlush")
			c.Del(ctx, k...)
		}
	})
}

func GetBy[T any, K comparable](name string, ct context.Context, key K, timeout time.Duration, params ...any) (r T, err error) {
	ct = context.WithValue(ct, "getCache", name)
	ca, err := getMap[K, T](name)
	if err != nil {
		return r, err
	}
	vv, err := ca.GetCache(ct, key, timeout, params...)
	if err != nil {
		return r, err
	}
	r = vv
	return
}

func getMap[K comparable, T any](name string) (*cache.MapCache[K, T], error) {
	m, ok := mapCache.Load(name)
	if !ok {
		return nil, errors.New(str.Join("cache ", name, " doesn't exist"))
	}
	vk, ok := m.(*cache.MapCache[K, T])
	if !ok {
		return nil, errors.New(str.Join("cache ", name, " type error"))
	}
	return vk, nil
}
func GetBatchBy[T any, K comparable](name string, ct context.Context, key []K, timeout time.Duration, params ...any) (r []T, err error) {
	ct = context.WithValue(ct, "getCache", name)
	ca, err := getMap[K, T](name)
	if err != nil {
		return r, err
	}
	vv, err := ca.GetCacheBatch(ct, key, timeout, params...)
	if err != nil {
		return r, err
	}
	r = vv
	return
}
func GetBatchByToMap[T any, K comparable](name string, ct context.Context, key []K, timeout time.Duration, params ...any) (r map[K]T, err error) {
	ct = context.WithValue(ct, "getCache", name)
	ca, err := getMap[K, T](name)
	if err != nil {
		return r, err
	}
	vv, err := ca.GetBatchToMap(ct, key, timeout, params...)
	if err != nil {
		return r, err
	}
	r = vv
	return
}

func NewMapCache[K comparable, V any](data cache.Cache[K, V], batchFn cache.MapBatchFn[K, V], fn cache.MapSingleFn[K, V], args ...any) *cache.MapCache[K, V] {
	inc := helper.ParseArgs((*cache.IncreaseUpdate[K, V])(nil), args...)
	m := cache.NewMapCache[K, V](data, fn, batchFn, inc, buildLockFn[K](args...), args...)
	name, f := parseArgs(args...)
	if name != "" {
		PushMangerMap(name, m)
	}
	PushOrSetFlush(Queue{
		Name: name,
		Fn:   m.Flush,
	})
	PushOrSetClearExpired(Queue{
		Name: name,
		Fn:   m.ClearExpired,
	})
	if f != nil && name != "" {
		SetExpireTime(any(data).(cache.SetTime), name, 0, f)
	}
	return m
}

func NewMemoryMapCache[K comparable, V any](batchFn cache.MapBatchFn[K, V],
	fn cache.MapSingleFn[K, V], expireTime time.Duration, args ...any) *cache.MapCache[K, V] {

	c := NewMapCache[K, V](cache.NewMemoryMapCache[K, V](func() time.Duration {
		return expireTime
	}), batchFn, fn, args...)
	return c
}

func GetMapCache[K comparable, V any](name string) (*cache.MapCache[K, V], bool) {
	vv, err := getMap[K, V](name)
	return vv, err == nil
}

func DelMapCacheVal[T any](name string, keys ...T) {
	v, ok := mapDelFuncs.Load(name)
	if !ok || len(keys) < 1 {
		return
	}
	v(keys)
}
