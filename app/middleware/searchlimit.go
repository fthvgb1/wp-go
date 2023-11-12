package middleware

import (
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/gin-gonic/gin"
)

func SearchLimit(num int64) func(ctx *gin.Context) {
	fn, reFn := IpLimit(num)
	reload.Push(func() {
		reFn(config.GetConfig().SingleIpSearchNum)
	})
	return func(c *gin.Context) {
		if c.Query("s") != "" {
			fn(c)
		} else {
			c.Next()
		}

	}
}
