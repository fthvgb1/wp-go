package cache

import (
	"context"
	"log"
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

func (c *SliceCache[T]) GetCache(ctx context.Context, timeout time.Duration) []T {
	l := len(c.data)
	expired := time.Duration(c.setTime.Unix())+c.expireTime/time.Second < time.Duration(time.Now().Unix())
	if l < 1 || (l > 0 && c.expireTime >= 0 && expired) {
		t := c.incr
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		done := make(chan struct{})
		go func() {
			c.mutex.Lock()
			defer c.mutex.Unlock()
			if c.incr > t {
				return
			}
			r, err := c.setCacheFunc()
			if err != nil {
				log.Printf("set cache err[%s]", err)
				return
			}
			c.setTime = time.Now()
			c.data = r
			c.incr++
			done <- struct{}{}
		}()
		select {
		case <-ctx.Done():
			log.Printf("get cache timeout")
		case <-done:

		}
	}
	return c.data
}
