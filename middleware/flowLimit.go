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
	m   map[string]int
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
		m:   make(map[string]int),
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
		if m.searchLimit(true, c, f) {
			c.Abort()
			return
		}
		atomic.AddInt64(&flow, 1)
		if flow >= vars.Conf.MaxRequestSleepNum && flow <= vars.Conf.MaxRequestNum {
			t := randFn(vars.Conf.SleepTime[0], vars.Conf.SleepTime[1])
			time.Sleep(t)
		} else if flow > vars.Conf.MaxRequestNum {
			c.String(http.StatusForbidden, "请求太多了，服务器君压力山大中==!, 请稍后访问")
			c.Abort()
			atomic.AddInt64(&flow, -1)
			m.searchLimit(false, c, f)
			return
		}

		c.Next()
		m.searchLimit(false, c, f)
		atomic.AddInt64(&flow, -1)
	}
}

func (m *IpLimitMap) set(k string, n int) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.m[k] = n
}

func (m *IpLimitMap) searchLimit(a bool, c *gin.Context, f []string) (isForbid bool) {
	ip := c.ClientIP()
	if f[0] == "" && c.Query("s") != "" {
		if a {
			i, ok := m.m[ip]
			if ok {
				num := vars.Conf.SingleIpSearchNum
				if num < 1 {
					num = 10
				}
				if i > num {
					return true
				}
			} else {
				i = 0
			}
			i++
			m.set(ip, i)
		} else {
			m.set(ip, m.m[ip]-1)
			if m.m[ip] == 0 {
				m.mux.Lock()
				delete(m.m, ip)
				m.mux.Unlock()
			}
		}
	}
	return
}
