package route

import (
	"github.com/fthvgb1/wp-go/internal/actions"
	"github.com/fthvgb1/wp-go/internal/middleware"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/static"
	"github.com/fthvgb1/wp-go/internal/theme"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
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

	r.HTMLRender = theme.NewFsTemplate(theme.FuncMap()).SetTemplate()
	validServerName, reloadValidServerNameFn := middleware.ValidateServerNames()
	fl, flReload := middleware.FlowLimit(c.MaxRequestSleepNum, c.MaxRequestNum, c.SleepTime)
	r.Use(
		gin.Logger(),
		validServerName,
		middleware.RecoverAndSendMail(gin.DefaultErrorWriter),
		fl,
		middleware.SetStaticFileCache,
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
	sl, slRload := middleware.SearchLimit(c.SingleIpSearchNum)
	r.GET("/", sl, actions.Index)
	r.GET("/page/:page", actions.Index)
	r.GET("/p/category/:category", actions.Index)
	r.GET("/p/category/:category/page/:page", actions.Index)
	r.GET("/p/tag/:tag", actions.Index)
	r.GET("/p/tag/:tag/page/:page", actions.Index)
	r.GET("/p/date/:year/:month", actions.Index)
	r.GET("/p/date/:year/:month/page/:page", actions.Index)
	r.GET("/p/author/:author", actions.Index)
	r.GET("/p/author/:author/page/:page", actions.Index)
	r.POST("/login", actions.Login)
	r.GET("/p/:id", actions.Detail)
	r.GET("/p/:id/feed", actions.PostFeed)
	r.GET("/feed", actions.Feed)
	r.GET("/comments/feed", actions.CommentsFeed)
	cfl, _ := middleware.FlowLimit(c.MaxRequestSleepNum, 5, c.SleepTime)
	r.POST("/comment", cfl, actions.PostComment)
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
