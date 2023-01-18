package middleware

import (
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/safety"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func ValidateServerNames() (func(ctx *gin.Context), func()) {
	var serverName safety.Map[string, struct{}]
	fn := func() {
		r := config.Conf.Load().TrustServerNames
		if len(r) > 0 {
			for _, name := range r {
				serverName.Store(name, struct{}{})
			}
		} else {
			serverName.Flush()
		}

	}
	fn()
	return func(c *gin.Context) {
		if serverName.Len() > 0 {
			if _, ok := serverName.Load(strings.Split(c.Request.Host, ":")[0]); !ok {
				c.Status(http.StatusForbidden)
				c.Abort()
				return
			}
		}
		c.Next()
	}, fn
}
