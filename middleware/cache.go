package middleware

import (
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/helper"
	"strings"
)

func SetStaticFileCache(c *gin.Context) {
	f := strings.Split(strings.TrimLeft(c.FullPath(), "/"), "/")
	if len(f) > 1 && helper.IsContainInArr(f[0], []string{"wp-includes", "wp-content"}) {
		c.Header("Cache-Control", "private, max-age=86400")
	}
}
