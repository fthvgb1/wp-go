package route

import (
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/static"
	"html/template"
	"net/http"
	"strings"
)

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(setStaticFileCache)
	r.SetFuncMap(template.FuncMap{"unescaped": func(s string) interface{} {
		return template.HTML(s)
	}})
	f := static.Fs{FS: static.FsEx, Path: "wp-includes"}
	r.StaticFS("/wp-includes", http.FS(f))
	r.StaticFS("/wp-content", http.FS(static.Fs{
		FS:   static.FsEx,
		Path: "wp-content",
	}))
	r.LoadHTMLGlob("templates/*")

	// Ping test
	r.GET("/", index)

	return r
}

func setStaticFileCache(c *gin.Context) {
	f := strings.Split(strings.TrimLeft(c.FullPath(), "/"), "/")
	if len(f) > 1 && helper.IsContainInArr(f[0], []string{"wp-includes", "wp-content"}) {
		c.Header("Cache-Control", "private, max-age=86400")
	}
}
