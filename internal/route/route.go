package route

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	actions2 "github/fthvgb1/wp-go/internal/actions"
	"github/fthvgb1/wp-go/internal/config"
	middleware2 "github/fthvgb1/wp-go/internal/middleware"
	"github/fthvgb1/wp-go/internal/static"
	"github/fthvgb1/wp-go/internal/templates"
	"github/fthvgb1/wp-go/internal/wpconfig"
	"html/template"
	"net/http"
	"time"
)

func SetupRouter() (*gin.Engine, func()) {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.New()
	c := config.Conf.Load()
	if len(c.TrustIps) > 0 {
		err := r.SetTrustedProxies(c.TrustIps)
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
		"getOption": func(k string) string {
			return wpconfig.Options.Value(k)
		},
	}).SetTemplate()
	validServerName, reloadValidServerNameFn := middleware2.ValidateServerNames()
	fl, flReload := middleware2.FlowLimit(c.MaxRequestSleepNum, c.MaxRequestNum, c.SleepTime)
	r.Use(
		gin.Logger(),
		validServerName,
		middleware2.RecoverAndSendMail(gin.DefaultErrorWriter),
		fl,
		middleware2.SetStaticFileCache,
	)
	//gzip 因为一般会用nginx做反代时自动使用gzip,所以go这边本身可以不用
	if c.Gzip {
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
	sl, slRload := middleware2.SearchLimit(c.SingleIpSearchNum)
	r.GET("/", sl, actions2.Index)
	r.GET("/page/:page", actions2.Index)
	r.GET("/p/category/:category", actions2.Index)
	r.GET("/p/category/:category/page/:page", actions2.Index)
	r.GET("/p/tag/:tag", actions2.Index)
	r.GET("/p/tag/:tag/page/:page", actions2.Index)
	r.GET("/p/date/:year/:month", actions2.Index)
	r.GET("/p/date/:year/:month/page/:page", actions2.Index)
	r.GET("/p/author/:author", actions2.Index)
	r.GET("/p/author/:author/page/:page", actions2.Index)
	r.POST("/login", actions2.Login)
	r.GET("/p/:id", actions2.Detail)
	r.GET("/p/:id/feed", actions2.PostFeed)
	r.GET("/feed", actions2.Feed)
	r.GET("/comments/feed", actions2.CommentsFeed)
	cfl, _ := middleware2.FlowLimit(c.MaxRequestSleepNum, 5, c.SleepTime)
	r.POST("/comment", cfl, actions2.PostComment)
	if gin.Mode() != gin.ReleaseMode {
		pprof.Register(r, "dev/pprof")
	}
	fn := func() {
		reloadValidServerNameFn()
		c := config.Conf.Load()
		flReload(c.MaxRequestSleepNum, c.MaxRequestNum, c.SleepTime)
		slRload(c.SingleIpSearchNum)
	}
	return r, fn
}
