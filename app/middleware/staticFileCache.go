package middleware

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

var path = map[string]struct{}{
	"wp-includes": {},
	"wp-content":  {},
	"favicon.ico": {},
}

func SetStaticFileCache(c *gin.Context) {
	f := strings.Split(strings.TrimLeft(c.FullPath(), "/"), "/")
	if _, ok := path[f[0]]; ok {
		if ".php" == filepath.Ext(c.Request.URL.Path) {
			c.Abort()
			c.Status(http.StatusForbidden)
			return
		}
		t := config.GetConfig().CacheTime.CacheControl
		if t > 0 {
			c.Header("Cache-Control", fmt.Sprintf("private, max-age=%d", int(t.Seconds())))
		}
	}

	c.Next()
}
