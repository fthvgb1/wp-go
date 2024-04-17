package cachemanager

import (
	"context"
	"errors"
	"github.com/fthvgb1/wp-go/cache"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/slice/mockmap"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/safety"
	"time"
)

var varCache = safety.NewMap[string, any]()

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

func NewVarCache[T any](c cache.AnyCache[T], fn func(context.Context, ...any) (T, error), a ...any) *cache.VarCache[T] {
	inc := helper.ParseArgs((*cache.IncreaseUpdateVar[T])(nil), a...)
	ref := helper.ParseArgs(cache.RefreshVar[T](nil), a...)
	v := cache.NewVarCache(c, fn, inc, ref, a...)

	name, _ := parseArgs(a...)
	if name != "" {
		varCache.Store(name, v)
	}
	PushOrSetFlush(mockmap.Item[string, Fn]{
		Name:  name,
		Value: v.Flush,
	})
	cc, ok := any(c).(clearExpired)
	if ok {
		PushOrSetClearExpired(mockmap.Item[string, Fn]{
			Name:  name,
			Value: cc.ClearExpired,
		})
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
