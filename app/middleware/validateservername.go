package middleware

import (
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ValidateServerNames() func(ctx *gin.Context) {
	sites := reload.VarsBy(func() map[string]struct{} {
		r := config.GetConfig().TrustServerNames
		m := map[string]struct{}{}
		if len(r) > 0 {
			for _, name := range r {
				m[name] = struct{}{}
			}
		}
		return m
	})

	return func(c *gin.Context) {
		m := sites.Load()
		if len(m) > 0 && !maps.IsExists(m, strings.Split(c.Request.Host, ":")[0]) {
			c.Status(http.StatusForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
