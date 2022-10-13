package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"sync/atomic"
)

type IpLimitMap struct {
	mux      *sync.Mutex
	m        map[string]*int64
	limitNum int64
}

func IpLimit(num int64) func(ctx *gin.Context) {
	m := IpLimitMap{
		mux:      &sync.Mutex{},
		m:        make(map[string]*int64),
		limitNum: num,
	}

	return func(c *gin.Context) {
		ip := c.ClientIP()
		s := false
		defer func() {
			i, ok := m.m[ip]
			if ok && s && *i > 0 {
				//time.Sleep(time.Second * 3)
				atomic.AddInt64(i, -1)
				if *i == 0 {
					m.mux.Lock()
					delete(m.m, ip)
					m.mux.Unlock()
				}
			}
		}()
		i, ok := m.m[ip]
		if !ok {
			m.mux.Lock()
			i = new(int64)
			m.m[ip] = i
			m.mux.Unlock()
		}
		if m.limitNum > 0 && *i >= m.limitNum {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}
		s = true
		atomic.AddInt64(i, 1)
	}
}
