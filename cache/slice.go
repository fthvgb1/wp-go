package cache

import (
	"context"
	"errors"
	"fmt"
	"github/fthvgb1/wp-go/safety"
	"sync"
	"time"
)

type SliceCache[T any] struct {
	v safety.Var[slice[T]]
}

type slice[T any] struct {
	data         []T
	mutex        *sync.Mutex
	setCacheFunc func(...any) ([]T, error)
	expireTime   time.Duration
	setTime      time.Time
	incr         int
}

func (c *SliceCache[T]) GetLastSetTime() time.Time {
	return c.v.Load().setTime
}

func NewSliceCache[T any](fun func(...any) ([]T, error), duration time.Duration) *SliceCache[T] {
	return &SliceCache[T]{
		v: safety.NewVar(slice[T]{
			mutex:        &sync.Mutex{},
			setCacheFunc: fun,
			expireTime:   duration,
		}),
	}
}

func (c *SliceCache[T]) FlushCache() {
	mu := c.v.Load().mutex
	mu.Lock()
	defer mu.Unlock()
	c.v.Delete()
}

func (c *SliceCache[T]) GetCache(ctx context.Context, timeout time.Duration, params ...any) ([]T, error) {
	v := c.v.Load()
	l := len(v.data)
	data := v.data
	var err error
	expired := time.Duration(v.setTime.UnixNano())+v.expireTime < time.Duration(time.Now().UnixNano())
	if l < 1 || (l > 0 && v.expireTime >= 0 && expired) {
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
			done := make(chan struct{})
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
