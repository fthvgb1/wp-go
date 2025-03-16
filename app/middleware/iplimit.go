package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
	"sync/atomic"
)

type ipLimitMap struct {
	mux      *sync.RWMutex
	m        map[string]*int64
	limitNum *int64
	clearNum *int64
}

func IpLimit(num int64, clearNum ...int64) (func(ctx *gin.Context), func(int64, ...int64)) {
	m := ipLimitMap{
		mux:      &sync.RWMutex{},
		m:        make(map[string]*int64),
		limitNum: new(int64),
		clearNum: new(int64),
	}
	fn := func(num int64, clearNum ...int64) {
		atomic.StoreInt64(m.limitNum, num)
		if len(clearNum) > 0 {
			atomic.StoreInt64(m.clearNum, clearNum[0])
		}
	}
	fn(num, clearNum...)

	return func(c *gin.Context) {
		if atomic.LoadInt64(m.limitNum) <= 0 {
			c.Next()
			return
		}
		ip := c.ClientIP()
		m.mux.RLock()
		i, ok := m.m[ip]
		m.mux.RUnlock()

		if !ok {
			m.mux.Lock()
			i = new(int64)
			m.m[ip] = i
			m.mux.Unlock()
		}

		defer func() {
			atomic.AddInt64(i, -1)
			if atomic.LoadInt64(i) <= 0 {
				cNum := int(atomic.LoadInt64(m.clearNum))
				if cNum <= 0 {
					m.mux.Lock()
					delete(m.m, ip)
					m.mux.Unlock()
					return
				}

				m.mux.RLock()
				l := len(m.m)
				m.mux.RUnlock()
				if l < cNum {
					m.mux.Lock()
					delete(m.m, ip)
					m.mux.Unlock()
				}
			}
		}()

		if atomic.LoadInt64(i) >= atomic.LoadInt64(m.limitNum) {
			c.String(http.StatusForbidden, "请求太多了，服务器君表示压力山大==!, 请稍后访问")
			c.Abort()
			return
		}
		atomic.AddInt64(i, 1)
		c.Next()
	}, fn
}
