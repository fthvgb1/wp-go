package route

import (
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/middleware"
	"github/fthvgb1/wp-go/static"
	"html/template"
	"net/http"
	"time"
)

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(middleware.SetStaticFileCache)
	r.SetFuncMap(template.FuncMap{
		"unescaped": func(s string) interface{} {
			return template.HTML(s)
		},
		"dateCh": func(t time.Time) interface{} {
			return t.Format("2006年01月02日")
		},
	})
	f := static.Fs{FS: static.FsEx, Path: "wp-includes"}
	r.StaticFS("/wp-includes", http.FS(f))
	r.StaticFS("/wp-content", http.FS(static.Fs{
		FS:   static.FsEx,
		Path: "wp-content",
	}))
	r.LoadHTMLGlob("templates/**/*")
	// Ping test
	r.GET("/", index)

	return r
}
