package route

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/actions"
	"github/fthvgb1/wp-go/middleware"
	"github/fthvgb1/wp-go/static"
	"github/fthvgb1/wp-go/templates"
	"html/template"
	"net/http"
	"time"
)

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.HTMLRender = templates.NewFsTemplate(template.FuncMap{
		"unescaped": func(s string) any {
			return template.HTML(s)
		},
		"dateCh": func(t time.Time) any {
			return t.Format("2006年 01月 02日")
		},
	}).SetTemplate()
	r.Use(middleware.SetStaticFileCache)
	//gzip 因为一般会用nginx做反代时自动使用gzip,所以go这边本身可以不用
	/*r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{
		"/wp-includes/", "/wp-content/",
	})))*/

	f := static.Fs{FS: static.FsEx, Path: "wp-includes"}
	r.StaticFileFS("/favicon.ico", "favicon.ico", http.FS(static.FsEx))
	r.StaticFS("/wp-includes", http.FS(f))
	r.StaticFS("/wp-content", http.FS(static.Fs{
		FS:   static.FsEx,
		Path: "wp-content",
	}))
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("go-wp", store))
	r.GET("/", actions.Index)
	r.GET("/page/:page", actions.Index)
	r.GET("/p/category/:category", actions.Index)
	r.GET("/p/category/:category/page/:page", actions.Index)
	r.GET("/p/tag/:tag", actions.Index)
	r.GET("/p/tag/:tag/page/:page", actions.Index)
	r.GET("/p/date/:year/:month", actions.Index)
	r.GET("/p/date/:year/:month/page/:page", actions.Index)
	r.POST("/login", actions.Login)
	r.GET("/p/:id", actions.Detail)

	return r
}
