package middleware

import (
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/helper"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

func FlowLimit(maxRequestSleepNum, maxRequestNum int64, sleepTime []time.Duration) func(ctx *gin.Context) {
	var flow int64
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

		n := atomic.LoadInt64(&flow)
		if n >= maxRequestSleepNum && n <= maxRequestNum {
			t := helper.RandNum(sleepTime[0], sleepTime[1])
			time.Sleep(t)
		} else if n > maxRequestNum {
			c.String(http.StatusForbidden, "请求太多了，服务器君表示压力山大==!, 请稍后访问")
			c.Abort()
			return
		}
		atomic.AddInt64(&flow, 1)
		defer func() {
			atomic.AddInt64(&flow, -1)
		}()

		c.Next()
	}
}