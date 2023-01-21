package middleware

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/gin-gonic/gin"
	"strings"
)

func SetStaticFileCache(c *gin.Context) {
	f := strings.Split(strings.TrimLeft(c.FullPath(), "/"), "/")
	if len(f) > 0 && slice.IsContained(f[0], []string{"wp-includes", "wp-content", "favicon.ico"}) {
		c.Header("Cache-Control", "private, max-age=86400")
	}
	c.Next()
}
