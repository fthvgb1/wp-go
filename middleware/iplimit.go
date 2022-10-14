package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"sync/atomic"
)

type IpLimitMap struct {
	mux      *sync.RWMutex
	m        map[string]*int64
	limitNum int64
}

func IpLimit(num int64) func(ctx *gin.Context) {
	m := IpLimitMap{
		mux:      &sync.RWMutex{},
		m:        make(map[string]*int64),
		limitNum: num,
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		s := false
		m.mux.RLock()
		i, ok := m.m[ip]
		m.mux.RUnlock()
		defer func() {
			ii := atomic.LoadInt64(i)
			if s && ii > 0 {
				atomic.AddInt64(i, -1)
				if atomic.LoadInt64(i) == 0 {
					m.mux.Lock()
					delete(m.m, ip)
					m.mux.Unlock()
				}
			}
		}()

		if !ok {
			m.mux.Lock()
			i = new(int64)
			m.m[ip] = i
			m.mux.Unlock()
		}

		if m.limitNum > 0 && atomic.LoadInt64(i) >= m.limitNum {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}
		atomic.AddInt64(i, 1)
		s = true
		c.Next()
	}
}
