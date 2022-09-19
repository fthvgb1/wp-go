package cache

import (
	"log"
	"sync"
	"time"
)

type SliceCache[T any] struct {
	data         []T
	mutex        *sync.Mutex
	setCacheFunc func() ([]T, error)
	expireTime   time.Duration
	setTime      time.Time
}

func NewSliceCache[T any](fun func() ([]T, error), duration time.Duration) *SliceCache[T] {
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

func (c *SliceCache[T]) GetCache() []T {
	l := len(c.data)
	expired := time.Duration(c.setTime.Unix())+c.expireTime/time.Second < time.Duration(time.Now().Unix())
	if l < 1 || (l > 0 && c.expireTime > 0 && expired) {
		r, err := c.setCacheFunc()
		if err != nil {
			log.Printf("set cache err[%s]", err)
			return nil
		}
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.setTime = time.Now()
		c.data = r
	}
	return c.data
}
