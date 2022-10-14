package middleware

import "github.com/gin-gonic/gin"

func SearchLimit(num int64) func(ctx *gin.Context) {
	fn := IpLimit(num)
	return func(c *gin.Context) {
		if "/" == c.FullPath() && c.Query("s") != "" {
			fn(c)
		} else {
			c.Next()
		}

	}
}
