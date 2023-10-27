package cachemanager

import (
	"context"
	"errors"
	"github.com/fthvgb1/wp-go/cache"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"time"
)

var ctx = context.Background()

var mapFlush = safety.NewMap[string, func(any)]()
var getSingleFn = safety.NewMap[string, func(context.Context, any, time.Duration, ...any) (any, error)]()
var getBatchFn = safety.NewMap[string, func(context.Context, any, time.Duration, ...any) (any, error)]()
var anyFlush = safety.NewMap[string, func()]()

type flush interface {
	Flush(ctx context.Context)
}

type clearExpired interface {
	ClearExpired(ctx context.Context)
}

var clears []clearExpired

var flushes []flush

func Flush() {
	for _, f := range flushes {
		f.Flush(ctx)
	}
}

func FlushMapVal[T any](name string, keys ...T) {
	v, ok := mapFlush.Load(name)
	if !ok || len(keys) < 1 {
		return
	}
	v(keys)
}

func FlushAnyVal(name ...string) {
	for _, s := range name {
		v, ok := anyFlush.Load(s)
		if ok {
			v()
		}
	}
}

func pushFlushMap[K comparable, V any](m *cache.MapCache[K, V], args ...any) {
	name := parseArgs(args...)
	if name != "" {
		anyFlush.Store(name, func() {
			m.Flush(ctx)
		})
		mapFlush.Store(name, func(a any) {
			k, ok := a.([]K)
			if ok && len(k) > 0 {
				m.Del(ctx, k...)
			}
		})
		getSingleFn.Store(name, func(ct context.Context, k any, t time.Duration, a ...any) (any, error) {
			kk, ok := k.(K)
			if !ok {
				return nil, errors.New(str.Join("cache ", name, " key type err"))
			}
			return m.GetCache(ct, kk, t, a...)
		})
		getBatchFn.Store(name, func(ct context.Context, k any, t time.Duration, a ...any) (any, error) {
			kk, ok := k.([]K)
			if !ok {
				return nil, errors.New(str.Join("cache ", name, " key type err"))
			}
			return m.GetCacheBatch(ct, kk, t, a...)
		})
		FlushPush()
	}
}

func Get[T, K any](name string, ct context.Context, key K, timeout time.Duration, params ...any) (r T, err error) {
	v, ok := getSingleFn.Load(name)
	if !ok {
		err = errors.New(str.Join("cache ", name, " doesn't exist"))
		return
	}
	vv, err := v(ct, key, timeout, params...)
	if err != nil {
		return r, err
	}
	r = vv.(T)
	return
}
func GetMultiple[T, K any](name string, ct context.Context, key []K, timeout time.Duration, params ...any) (r []T, err error) {
	v, ok := getBatchFn.Load(name)
	if !ok {
		err = errors.New(str.Join("cache ", name, " doesn't exist"))
		return
	}
	vv, err := v(ct, key, timeout, params...)
	if err != nil {
		return r, err
	}
	r = vv.([]T)
	return
}

func parseArgs(args ...any) string {
	var name string
	for _, arg := range args {
		v, ok := arg.(string)
		if ok {
			name = v
		}
	}
	return name
}

func NewMapCache[K comparable, V any](data cache.Cache[K, V], batchFn cache.MapBatchFn[K, V], fn cache.MapSingleFn[K, V], args ...any) *cache.MapCache[K, V] {
	m := cache.NewMapCache[K, V](data, fn, batchFn)
	pushFlushMap(m, args...)
	FlushPush(m)
	ClearPush(m)
	return m
}
func NewMemoryMapCache[K comparable, V any](batchFn cache.MapBatchFn[K, V],
	fn cache.MapSingleFn[K, V], expireTime time.Duration, args ...any) *cache.MapCache[K, V] {
	return NewMapCache[K, V](cache.NewMemoryMapCache[K, V](expireTime), batchFn, fn, args...)
}

func FlushPush(f ...flush) {
	flushes = append(flushes, f...)
}
func ClearPush(c ...clearExpired) {
	clears = append(clears, c...)
}

func ClearExpired() {
	for _, c := range clears {
		c.ClearExpired(ctx)
	}
}
