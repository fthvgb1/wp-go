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
	v *safety.Var[vars[T]]
}

type vars[T any] struct {
	data         T
	mutex        *sync.Mutex
	setCacheFunc func(...any) (T, error)
	expireTime   time.Duration
	setTime      time.Time
	incr         int
}

func (c *VarCache[T]) GetLastSetTime() time.Time {
	return c.v.Load().setTime
}

func NewVarCache[T any](fun func(...any) (T, error), duration time.Duration) *VarCache[T] {
	return &VarCache[T]{
		v: safety.NewVar(vars[T]{
			mutex:        &sync.Mutex{},
			setCacheFunc: fun,
			expireTime:   duration,
		}),
	}
}

func (c *VarCache[T]) IsExpired() bool {
	v := c.v.Load()
	return time.Duration(v.setTime.UnixNano())+v.expireTime < time.Duration(time.Now().UnixNano())
}

func (c *VarCache[T]) Flush() {
	v := c.v.Load()
	mu := v.mutex
	mu.Lock()
	defer mu.Unlock()
	var vv T
	v.data = vv
	c.v.Store(v)
}

func (c *VarCache[T]) GetCache(ctx context.Context, timeout time.Duration, params ...any) (T, error) {
	v := c.v.Load()
	data := v.data
	var err error
	if v.expireTime <= 0 || ((time.Duration(v.setTime.UnixNano()) + v.expireTime) < time.Duration(time.Now().UnixNano())) {
		t := v.incr
		call := func() {
			v.mutex.Lock()
			defer v.mutex.Unlock()
			vv := c.v.Load()
			if vv.incr > t {
				return
			}
			r, er := vv.setCacheFunc(params...)
			if err != nil {
				err = er
				return
			}
			vv.setTime = time.Now()
			vv.data = r
			data = r
			vv.incr++
			c.v.Store(vv)
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

	}
	return data, err
}
