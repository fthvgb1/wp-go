package route

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/actions"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/middleware"
	"github/fthvgb1/wp-go/static"
	"github/fthvgb1/wp-go/templates"
	"github/fthvgb1/wp-go/vars"
	"html/template"
	"net/http"
	"time"
)

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.New()
	if len(vars.Conf.TrustIps) > 0 {
		err := r.SetTrustedProxies(vars.Conf.TrustIps)
		if err != nil {
			panic(err)
		}
	}

	r.HTMLRender = templates.NewFsTemplate(template.FuncMap{
		"unescaped": func(s string) any {
			return template.HTML(s)
		},
		"dateCh": func(t time.Time) any {
			return t.Format("2006年 01月 02日")
		},
	}).SetTemplate()
	r.Use(
		middleware.ValidateServerNames(),
		gin.Logger(),
		gin.Recovery(),
		middleware.FlowLimit(vars.Conf.MaxRequestSleepNum, vars.Conf.MaxRequestNum, vars.Conf.SleepTime),
		middleware.SetStaticFileCache,
	)
	//gzip 因为一般会用nginx做反代时自动使用gzip,所以go这边本身可以不用
	if vars.Conf.Gzip {
		r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{
			"/wp-includes/", "/wp-content/",
		})))
	}
	f := static.Fs{FS: static.FsEx, Path: "wp-includes"}
	r.StaticFileFS("/favicon.ico", "favicon.ico", http.FS(static.FsEx))
	r.StaticFS("/wp-includes", http.FS(f))
	r.StaticFS("/wp-content", http.FS(static.Fs{
		FS:   static.FsEx,
		Path: "wp-content",
	}))
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("go-wp", store))
	r.GET("/", middleware.SearchLimit(vars.Conf.SingleIpSearchNum), actions.Index)
	r.GET("/page/:page", actions.Index)
	r.GET("/p/category/:category", actions.Index)
	r.GET("/p/category/:category/page/:page", actions.Index)
	r.GET("/p/tag/:tag", actions.Index)
	r.GET("/p/tag/:tag/page/:page", actions.Index)
	r.GET("/p/date/:year/:month", actions.Index)
	r.GET("/p/date/:year/:month/page/:page", actions.Index)
	r.POST("/login", actions.Login)
	r.GET("/p/:id", actions.Detail)
	r.GET("/p/:id/feed", actions.PostFeed)
	r.GET("/feed", actions.Feed)
	r.GET("/comments/feed", actions.CommentsFeed)
	r.POST("/comment", actions.PostComment)
	if helper.IsContainInArr(gin.Mode(), []string{gin.DebugMode, gin.TestMode}) {
		pprof.Register(r, "dev/pprof")
	}
	return r
}
