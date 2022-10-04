package middleware

import (
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/vars"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type IpLimitMap struct {
	mux *sync.Mutex
	m   map[string]*int64
}

func FlowLimit() func(ctx *gin.Context) {
	var flow int64
	rand.Seed(time.Now().UnixNano())
	randFn := func(start, end time.Duration) time.Duration {
		end++
		return time.Duration(rand.Intn(int(end-start)) + int(start))
	}
	m := IpLimitMap{
		mux: &sync.Mutex{},
		m:   make(map[string]*int64),
	}
	statPath := map[string]struct{}{
		"wp-includes": {},
		"wp-content":  {},
		"favicon.ico": {},
	}
	return func(c *gin.Context) {
		f := strings.Split(strings.TrimLeft(c.FullPath(), "/"), "/")
		_, ok := statPath[f[0]]
		if len(f) > 0 && ok {
			c.Next()
			return
		}
		s := false
		ip := c.ClientIP()
		defer m.searchLimit(false, c, ip, f, &s)
		if m.searchLimit(true, c, ip, f, &s) {
			c.Abort()
			return
		}
		atomic.AddInt64(&flow, 1)
		defer func() {
			atomic.AddInt64(&flow, -1)
		}()
		if flow >= vars.Conf.MaxRequestSleepNum && flow <= vars.Conf.MaxRequestNum {
			t := randFn(vars.Conf.SleepTime[0], vars.Conf.SleepTime[1])
			time.Sleep(t)
		} else if flow > vars.Conf.MaxRequestNum {
			c.String(http.StatusForbidden, "请求太多了，服务器君表示压力山大==!, 请稍后访问")
			c.Abort()

			return
		}
		c.Next()

	}
}

func (m *IpLimitMap) searchLimit(start bool, c *gin.Context, ip string, f []string, s *bool) (isForbid bool) {
	if f[0] == "" && c.Query("s") != "" {
		if start {
			i, ok := m.m[ip]
			num := vars.Conf.SingleIpSearchNum
			if !ok {
				m.mux.Lock()
				i = new(int64)
				m.m[ip] = i
				m.mux.Unlock()
			}
			if num > 0 && *i >= num {
				isForbid = true
				return
			}
			*s = true
			atomic.AddInt64(i, 1)
			return
		}
		i, ok := m.m[ip]
		if ok && *s && *i > 0 {
			atomic.AddInt64(i, -1)
			if *i == 0 {
				m.mux.Lock()
				delete(m.m, ip)
				m.mux.Unlock()
			}
		}
	}
	return
}
