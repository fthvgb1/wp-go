package cache

import (
	"context"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/safety"
	"sync"
	"time"
)

type VarCache[T any] struct {
	AnyCache[T]
	setCacheFunc   func(context.Context, ...any) (T, error)
	mutex          sync.Mutex
	increaseUpdate IncreaseUpdateVar[T]
	refresh        RefreshVar[T]
}

type IncreaseUpdateVar[T any] struct {
	CycleTime func() time.Duration
	Fn        IncreaseVarFn[T]
}

type IncreaseVarFn[T any] func(c context.Context, currentData T, t time.Time, a ...any) (data T, save bool, refresh bool, err error)

func (t *VarCache[T]) GetCache(ctx context.Context, timeout time.Duration, params ...any) (T, error) {
	data, ok := t.Get(ctx)
	var err error
	if ok {
		if t.increaseUpdate.Fn != nil && t.refresh != nil {
			nowTime := time.Now()
			if t.increaseUpdate.CycleTime() > nowTime.Sub(t.GetLastSetTime(ctx)) {
				return data, nil
			}
			fn := func() {
				t.mutex.Lock()
				defer t.mutex.Unlock()
				da, save, refresh, er := t.increaseUpdate.Fn(ctx, data, t.GetLastSetTime(ctx), params...)
				if er != nil {
					err = er
					return
				}
				if save {
					t.Set(ctx, da)
				}
				if refresh {
					t.refresh.Refresh(ctx, params...)
				}
			}
			if timeout > 0 {
				er := helper.RunFnWithTimeout(ctx, timeout, fn, "increaseUpdate cache fail")
				if err == nil && er != nil {
					err = er
				}
			} else {
				fn()
			}
		}
		return data, nil
	}
	call := func() {
		t.mutex.Lock()
		defer t.mutex.Unlock()
		dat, ok := t.Get(ctx)
		if ok {
			data = dat
			return
		}
		r, er := t.setCacheFunc(ctx, params...)
		if er != nil {
			err = er
			return
		}
		t.Set(ctx, r)
		data = r
	}
	if timeout > 0 {
		er := helper.RunFnWithTimeout(ctx, timeout, call, "get cache fail")
		if err == nil && er != nil {
			err = er
		}
	} else {
		call()
	}
	return data, err
}

type VarMemoryCache[T any] struct {
	v          *safety.Var[vars[T]]
	expireTime func() time.Duration
}

func (c *VarMemoryCache[T]) ClearExpired(ctx context.Context) {
	c.Flush(ctx)
}

func NewVarMemoryCache[T any](expireTime func() time.Duration) *VarMemoryCache[T] {
	return &VarMemoryCache[T]{v: safety.NewVar(vars[T]{}), expireTime: expireTime}
}

func (c *VarMemoryCache[T]) Get(_ context.Context) (T, bool) {
	v := c.v.Load()
	return v.data, c.expireTime() >= time.Now().Sub(v.setTime)
}

func (c *VarMemoryCache[T]) Set(_ context.Context, v T) {
	vv := c.v.Load()
	vv.data = v
	vv.setTime = time.Now()
	vv.incr++
	c.v.Store(vv)
}

func (c *VarMemoryCache[T]) SetExpiredTime(f func() time.Duration) {
	c.expireTime = f
}

type vars[T any] struct {
	data    T
	setTime time.Time
	incr    int
}

func (c *VarMemoryCache[T]) GetLastSetTime(_ context.Context) time.Time {
	return c.v.Load().setTime
}

func NewVarCache[T any](cache AnyCache[T], fn func(context.Context, ...any) (T, error), inc IncreaseUpdateVar[T], ref RefreshVar[T]) *VarCache[T] {
	return &VarCache[T]{
		AnyCache: cache, setCacheFunc: fn, mutex: sync.Mutex{},
		increaseUpdate: inc,
		refresh:        ref,
	}
}

func (c *VarMemoryCache[T]) Flush(_ context.Context) {
	c.v.Flush()
}
