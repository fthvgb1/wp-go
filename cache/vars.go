package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/safety"
	"sync"
	"time"
)

type VarCache[T any] struct {
	AnyCache[T]
	setCacheFunc func(context.Context, ...any) (T, error)
	mutex        sync.Mutex
}

func (t *VarCache[T]) GetCache(ctx context.Context, timeout time.Duration, params ...any) (T, error) {
	data, ok := t.Get(ctx)
	if ok {
		return data, nil
	}
	var err error
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
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		done := make(chan struct{}, 1)
		go func() {
			call()
			done <- struct{}{}
			close(done)
		}()
		select {
		case <-ctx.Done():
			err = errors.New(fmt.Sprintf("get cache %s", ctx.Err().Error()))
		case <-done:

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

func NewVarCache[T any](cache AnyCache[T], fn func(context.Context, ...any) (T, error)) *VarCache[T] {
	return &VarCache[T]{
		AnyCache: cache, setCacheFunc: fn, mutex: sync.Mutex{},
	}
}

func (c *VarMemoryCache[T]) Flush(_ context.Context) {
	c.v.Flush()
}
