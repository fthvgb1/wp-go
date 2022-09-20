package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type SliceCache[T any] struct {
	data         []T
	mutex        *sync.Mutex
	setCacheFunc func(...any) ([]T, error)
	expireTime   time.Duration
	setTime      time.Time
	incr         int
}

func NewSliceCache[T any](fun func(...any) ([]T, error), duration time.Duration) *SliceCache[T] {
	return &SliceCache[T]{
		mutex:        &sync.Mutex{},
		setCacheFunc: fun,
		expireTime:   duration,
	}
}

func (c *SliceCache[T]) FlushCache() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data = nil
}

func (c *SliceCache[T]) GetCache(ctx context.Context, timeout time.Duration, params ...any) ([]T, error) {
	l := len(c.data)
	data := c.data
	var err error
	expired := time.Duration(c.setTime.Unix())+c.expireTime/time.Second < time.Duration(time.Now().Unix())
	if l < 1 || (l > 0 && c.expireTime >= 0 && expired) {
		t := c.incr
		call := func() {
			c.mutex.Lock()
			defer c.mutex.Unlock()
			if c.incr > t {
				return
			}
			r, er := c.setCacheFunc(params...)
			if err != nil {
				err = er
				return
			}
			c.setTime = time.Now()
			c.data = r
			data = r
			c.incr++
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
