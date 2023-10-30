package cachemanager

import (
	"context"
	"errors"
	"github.com/fthvgb1/wp-go/app/cmd/reload"
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

var expiredTime = safety.NewMap[string, expire]()

type expire struct {
	fn          func() time.Duration
	p           *safety.Var[time.Duration]
	isUseManger *safety.Var[bool]
}

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
	name, _ := parseArgs(args...)
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

func parseArgs(args ...any) (string, func() time.Duration) {
	var name string
	var fn func() time.Duration
	for _, arg := range args {
		v, ok := arg.(string)
		if ok {
			name = v
			continue
		}
		vv, ok := arg.(func() time.Duration)
		if ok {
			fn = vv
		}

	}
	return name, fn
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
	name, f := parseArgs(args...)
	var t, tt func() time.Duration
	t = f
	if t == nil {
		t = func() time.Duration {
			return expireTime
		}
	}
	tt = t
	if name != "" {
		expireTime = t()
		p := safety.NewVar(expireTime)
		e := expire{
			fn:          t,
			p:           p,
			isUseManger: safety.NewVar(false),
		}
		expiredTime.Store(name, e)
		reload.Push(func() {
			if !e.isUseManger.Load() {
				e.p.Store(tt())
			}
		}, str.Join("cacheManger-", name, "-expiredTime"))
		t = func() time.Duration {
			return e.p.Load()
		}
	}

	return NewMapCache[K, V](cache.NewMemoryMapCache[K, V](t), batchFn, fn, args...)
}

func SetExpireTime(t time.Duration, name ...string) {
	for _, s := range name {
		v, ok := expiredTime.Load(s)
		if !ok {
			continue
		}
		v.p.Store(t)
		if !v.isUseManger.Load() {
			v.isUseManger.Store(true)
		}
	}
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
