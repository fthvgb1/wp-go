package route

import (
	"github.com/fthvgb1/wp-go/app/actions"
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/middleware"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/static"
	"github.com/fthvgb1/wp-go/app/theme"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

var hooker []func(r *gin.Engine)

// Hook 方便插件在init时使用
func Hook(fn ...func(r *gin.Engine)) {
	hooker = append(hooker, fn...)
}

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.New()
	c := config.GetConfig()
	if len(c.TrustIps) > 0 {
		err := r.SetTrustedProxies(c.TrustIps)
		if err != nil {
			panic(err)
		}
	}

	r.HTMLRender = theme.Template()
	wpconfig.SetTemplateFs(theme.TemplateFs)
	siteFlowLimitMiddleware, siteFlow := middleware.FlowLimit(c.MaxRequestSleepNum, c.MaxRequestNum, c.CacheTime.SleepTime)
	r.Use(
		gin.Logger(),
		middleware.ValidateServerNames(),
		middleware.RecoverAndSendMail(gin.DefaultErrorWriter),
		siteFlowLimitMiddleware,
		middleware.SetStaticFileCache,
	)
	//gzip 因为一般会用nginx做反代时自动使用gzip,所以go这边本身可以不用
	if c.Gzip {
		r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{
			"/wp-includes/", "/wp-content/",
		})))
	}

	f := static.Fs{FS: static.FsDir, Path: "wp-includes"}
	r.StaticFileFS("/favicon.ico", "favicon.ico", http.FS(static.FsDir))
	r.StaticFS("/wp-includes", http.FS(f))
	r.StaticFS("/wp-content/plugins", http.FS(static.Fs{
		FS:   static.FsDir,
		Path: "wp-content/plugins",
	}))
	r.StaticFS("/wp-content/themes", http.FS(static.Fs{
		FS:   static.FsDir,
		Path: "wp-content/themes",
	}))
	if c.UploadDir != "" {
		r.Static("/wp-content/uploads", c.UploadDir)
	}
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("go-wp", store))
	r.GET("/", middleware.SearchLimit(c.SingleIpSearchNum), actions.ThemeHook(constraints.Home))
	r.GET("/page/:page", actions.ThemeHook(constraints.Home))
	r.GET("/p/category/:category", actions.ThemeHook(constraints.Category))
	r.GET("/p/category/:category/page/:page", actions.ThemeHook(constraints.Category))
	r.GET("/p/tag/:tag", actions.ThemeHook(constraints.Tag))
	r.GET("/p/tag/:tag/page/:page", actions.ThemeHook(constraints.Tag))
	r.GET("/p/date/:year/:month", actions.ThemeHook(constraints.Archive))
	r.GET("/p/date/:year/:month/page/:page", actions.ThemeHook(constraints.Archive))
	r.GET("/p/author/:author", actions.ThemeHook(constraints.Author))
	r.GET("/p/author/:author/page/:page", actions.ThemeHook(constraints.Author))
	r.POST("/login", actions.Login)
	r.GET("/p/:id", actions.ThemeHook(constraints.Detail))
	r.GET("/p/:id/feed", actions.PostFeed)
	r.GET("/feed", actions.Feed)
	r.GET("/comments/feed", actions.CommentsFeed)
	r.NoRoute(actions.ThemeHook(constraints.NoRoute))
	commentMiddleWare, _ := middleware.FlowLimit(c.MaxRequestSleepNum, 5, c.CacheTime.SleepTime)
	r.POST("/comment", commentMiddleWare, actions.PostComment)
	if c.Pprof != "" {
		pprof.Register(r, c.Pprof)
	}
	for _, fn := range hooker {
		fn(r)
	}

	reload.Push(func() {
		c := config.GetConfig()
		siteFlow(c.MaxRequestSleepNum, c.MaxRequestNum, c.CacheTime.SleepTime)
	})
	return r
}
