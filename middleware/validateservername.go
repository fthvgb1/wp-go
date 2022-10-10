package middleware

import (
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/vars"
	"net/http"
	"strings"
)

func ValidateServerNames() func(ctx *gin.Context) {
	serverName := helper.SimpleSliceToMap(vars.Conf.TrustServerNames, func(v string) string {
		return v
	})
	return func(c *gin.Context) {
		if len(serverName) > 0 {
			if _, ok := serverName[strings.Split(c.Request.Host, ":")[0]]; !ok {
				c.Status(http.StatusForbidden)
				c.Abort()
			}
		}
	}
}
