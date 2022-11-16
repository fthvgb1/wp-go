package middleware

import "github.com/gin-gonic/gin"

func SearchLimit(num int64) (func(ctx *gin.Context), func(int64)) {
	fn, reFn := IpLimit(num)
	return func(c *gin.Context) {
		if c.Query("s") != "" {
			fn(c)
		} else {
			c.Next()
		}

	}, reFn
}
