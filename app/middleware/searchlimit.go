package middleware

import (
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/gin-gonic/gin"
)

func SearchLimit(num int64, clearNum ...int64) func(ctx *gin.Context) {
	fn, reFn := IpLimit(num, clearNum...)
	reload.Append(func() {
		reFn(config.GetConfig().SingleIpSearchNum)
	}, "search-ip-limit-number")
	return func(c *gin.Context) {
		if c.Query("s") != "" {
			fn(c)
		} else {
			c.Next()
		}

	}
}
