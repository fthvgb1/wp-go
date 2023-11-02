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

var getVar = safety.NewMap[string, func(context.Context, time.Duration, ...any) (any, error)]()

var expiredTime = safety.NewMap[string, expire]()

var varCache = safety.NewMap[string, any]()
var mapCache = safety.NewMap[string, any]()

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
	if name == "" {
		return
	}
	mapCache.Store(name, m)
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

func Get[T, K any](name string, ct context.Context, key K, timeout time.Duration, params ...any) (r T, err error) {
	ct = context.WithValue(ct, "getCache", name)
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
	ct = context.WithValue(ct, "getCache", name)
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
	name, f := parseArgs(args...)
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

func SetExpireTime(c cache.SetTime, name string, expireTime time.Duration, expireTimeFn func() time.Duration) {
	if name == "" {
		return
	}
	var t, tt func() time.Duration
	t = expireTimeFn
	if t == nil {
		t = func() time.Duration {
			return expireTime
		}
	}
	tt = t
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
	c.SetExpiredTime(t)
}

func ChangeExpireTime(t time.Duration, name ...string) {
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

func NewVarCache[T any](c cache.AnyCache[T], fn func(context.Context, ...any) (T, error), a ...any) *cache.VarCache[T] {
	v := cache.NewVarCache(c, fn)
	FlushPush(v)
	name, _ := parseArgs(a...)
	if name != "" {
		varCache.Store(name, v)
		getVar.Store(name, func(c context.Context, duration time.Duration, a ...any) (any, error) {
			return v.GetCache(c, duration, a...)
		})
	}
	cc, ok := any(c).(clearExpired)
	if ok {
		ClearPush(cc)
	}
	return v
}

func GetVarVal[T any](name string, ctx context.Context, duration time.Duration, a ...any) (r T, err error) {
	ctx = context.WithValue(ctx, "getCache", name)
	fn, ok := getVar.Load(name)
	if !ok {
		err = errors.New(str.Join("cache ", name, " is not exist"))
		return
	}
	v, err := fn(ctx, duration, a...)
	if err != nil {
		return
	}
	vv, ok := v.(T)
	if !ok {
		err = errors.New(str.Join("cache ", name, " value wanted can't match got"))
		return
	}
	r = vv
	return
}

func NewVarMemoryCache[T any](fn func(context.Context, ...any) (T, error), expired time.Duration, a ...any) *cache.VarCache[T] {
	c := cache.NewVarMemoryCache[T](nil)
	name, e := parseArgs(a...)
	SetExpireTime(c, name, expired, e)
	v := NewVarCache[T](c, fn, a...)
	return v
}

func GetVarCache[T any](name string) (*cache.VarCache[T], bool) {
	v, ok := varCache.Load(name)
	if !ok {
		return nil, false
	}
	vv, ok := v.(*cache.VarCache[T])
	return vv, ok
}

func GetMapCache[K comparable, V any](name string) (*cache.MapCache[K, V], bool) {
	v, ok := mapCache.Load(name)
	if !ok {
		return nil, false
	}
	vv, ok := v.(*cache.MapCache[K, V])
	return vv, ok
}
