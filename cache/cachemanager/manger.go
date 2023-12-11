package cachemanager

import (
	"context"
	"errors"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"time"
)

var mapFlush = safety.NewMap[string, func(any)]()
var anyFlush = safety.NewMap[string, func()]()

var varCache = safety.NewMap[string, any]()
var mapCache = safety.NewMap[string, any]()

func SetVarCache[T any](name string, v *cache.VarCache[T]) error {
	vv, ok := varCache.Load(name)
	if !ok {
		varCache.Store(name, v)
		return nil
	}
	_, ok = vv.(*cache.VarCache[T])
	if ok {
		varCache.Store(name, v)
		return nil
	}
	return errors.New(str.Join("cache ", name, " type err"))
}

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

type flush interface {
	Flush(ctx context.Context)
}

type clearExpired interface {
	ClearExpired(ctx context.Context)
}

var clears []clearExpired

var flushes []flush

func Flush() {
	ctx := context.WithValue(context.Background(), "execFlushBy", "mangerFlushFn")
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

func PushMangerMap[K comparable, V any](name string, m *cache.MapCache[K, V]) {
	if name == "" {
		return
	}
	mapCache.Store(name, m)
	anyFlush.Store(name, func() {
		ctx := context.WithValue(context.Background(), "ctx", "registerFlush")
		m.Flush(ctx)
	})
	mapFlush.Store(name, func(a any) {
		k, ok := a.([]K)
		if ok && len(k) > 0 {
			ctx := context.WithValue(context.Background(), "ctx", "registerFlush")
			m.Del(ctx, k...)
		}
	})
}

func Get[T any, K comparable](name string, ct context.Context, key K, timeout time.Duration, params ...any) (r T, err error) {
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
func GetMultiple[T any, K comparable](name string, ct context.Context, key []K, timeout time.Duration, params ...any) (r []T, err error) {
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
func GetMultipleToMap[T any, K comparable](name string, ct context.Context, key []K, timeout time.Duration, params ...any) (r map[K]T, err error) {
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

func NewPaginationCache[K comparable, V any](m *cache.MapCache[string, helper.PaginationData[V]], maxNum int,
	dbFn cache.DbFn[K, V], localFn cache.LocalFn[K, V], dbKeyFn, localKeyFn func(K, ...any) string, fetchNum int, name string, a ...any) *cache.Pagination[K, V] {
	fn := helper.ParseArgs([]func() int(nil), a...)
	var ma, fet func() int
	if len(fn) > 0 {
		ma = fn[0]
		if len(fn) > 1 {
			fet = fn[1]
		}
	}
	if ma == nil {
		ma = reload.FnVal(str.Join("paginationCache-", name, "-maxNum"), maxNum, nil)
	}
	if fet == nil {
		fet = reload.FnVal(str.Join("paginationCache-", name, "-fetchNum"), fetchNum, nil)
	}
	return cache.NewPagination(m, ma, dbFn, localFn, dbKeyFn, localKeyFn, fet, name)
}

func NewMapCache[K comparable, V any](data cache.Cache[K, V], batchFn cache.MapBatchFn[K, V], fn cache.MapSingleFn[K, V], args ...any) *cache.MapCache[K, V] {
	inc := helper.ParseArgs((*cache.IncreaseUpdate[K, V])(nil), args...)
	m := cache.NewMapCache[K, V](data, fn, batchFn, inc, args...)
	FlushPush(m)
	ClearPush(m)
	name, f := parseArgs(args...)
	if name != "" {
		PushMangerMap(name, m)
	}
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
	fn := reload.FnVal(str.Join("cacheManger-", name, "-expiredTime"), expireTime, expireTimeFn)
	c.SetExpiredTime(fn)
}

func ChangeExpireTime(t time.Duration, name ...string) {
	for _, s := range name {
		reload.ChangeFnVal(s, t)
	}
}

func FlushPush(f ...flush) {
	flushes = append(flushes, f...)
}
func ClearPush(c ...clearExpired) {
	clears = append(clears, c...)
}

func ClearExpired() {
	ctx := context.WithValue(context.Background(), "execClearExpired", "mangerClearExpiredFn")
	for _, c := range clears {
		c.ClearExpired(ctx)
	}
}

func NewVarCache[T any](c cache.AnyCache[T], fn func(context.Context, ...any) (T, error), a ...any) *cache.VarCache[T] {
	inc := helper.ParseArgs((*cache.IncreaseUpdateVar[T])(nil), a...)
	ref := helper.ParseArgs(cache.RefreshVar[T](nil), a...)
	v := cache.NewVarCache(c, fn, inc, ref, a...)
	FlushPush(v)
	name, _ := parseArgs(a...)
	if name != "" {
		varCache.Store(name, v)
	}
	cc, ok := any(c).(clearExpired)
	if ok {
		ClearPush(cc)
	}
	return v
}

func GetVarVal[T any](name string, ctx context.Context, duration time.Duration, a ...any) (r T, err error) {
	ctx = context.WithValue(ctx, "getCache", name)
	ca, ok := GetVarCache[T](name)
	if !ok {
		err = errors.New(str.Join("cache ", name, " is not exist"))
		return
	}
	v, err := ca.GetCache(ctx, duration, a...)
	if err != nil {
		return
	}
	r = v
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
	vv, err := getMap[K, V](name)
	return vv, err == nil
}
