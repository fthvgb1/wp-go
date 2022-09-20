package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type MapCache[K comparable, V any] struct {
	data         map[K]V
	mutex        *sync.Mutex
	setCacheFunc func(...any) (V, error)
	expireTime   time.Duration
	setTime      time.Time
	incr         int
}

func NewMapCache[K comparable, V any](fun func(...any) (V, error), duration time.Duration) *MapCache[K, V] {
	return &MapCache[K, V]{
		mutex:        &sync.Mutex{},
		setCacheFunc: fun,
		expireTime:   duration,
	}
}

func (c *MapCache[K, V]) FlushCache(k any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	key := k.(K)
	delete(c.data, key)
}

func (c *MapCache[K, V]) GetCache(ctx context.Context, key K, timeout time.Duration, params ...any) (V, error) {
	l := len(c.data)
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
			c.data[key] = r
			c.incr++
		}
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			done := make(chan struct{})
			go func() {
				call()
				done <- struct{}{}
			}()
			select {
			case <-ctx.Done():
				err = errors.New(fmt.Sprintf("get cache %v %s", key, ctx.Err().Error()))
			case <-done:
			}
		} else {
			call()
		}

	}
	return c.data[key], err
}
