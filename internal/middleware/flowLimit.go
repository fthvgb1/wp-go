package middleware

import (
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/safety"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

func FlowLimit(maxRequestSleepNum, maxRequestNum int64, sleepTime []time.Duration) (func(ctx *gin.Context), func(int64, int64, []time.Duration)) {
	var flow int64
	statPath := map[string]struct{}{
		"wp-includes": {},
		"wp-content":  {},
		"favicon.ico": {},
	}
	s := safety.Var[[]time.Duration]{}
	s.Store(sleepTime)
	fn := func(msn, mn int64, st []time.Duration) {
		atomic.StoreInt64(&maxRequestSleepNum, msn)
		atomic.StoreInt64(&maxRequestNum, mn)
		s.Store(st)
	}
	return func(c *gin.Context) {
		f := strings.Split(strings.TrimLeft(c.FullPath(), "/"), "/")
		_, ok := statPath[f[0]]
		if len(f) > 0 && ok {
			c.Next()
			return
		}

		n := atomic.LoadInt64(&flow)
		if n >= atomic.LoadInt64(&maxRequestSleepNum) && n <= atomic.LoadInt64(&maxRequestNum) {
			ss := s.Load()
			t := helper.RandNum(ss[0], ss[1])
			time.Sleep(t)
		} else if n > atomic.LoadInt64(&maxRequestNum) {
			c.String(http.StatusForbidden, "请求太多了，服务器君表示压力山大==!, 请稍后访问")
			c.Abort()
			return
		}
		atomic.AddInt64(&flow, 1)
		defer func() {
			atomic.AddInt64(&flow, -1)
		}()

		c.Next()
	}, fn
}
