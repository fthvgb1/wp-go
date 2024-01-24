package route

import (
	"github.com/fthvgb1/wp-go/app/actions"
	"github.com/fthvgb1/wp-go/app/middleware"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/theme"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache/reload"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
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

	r.HTMLRender = theme.BuildTemplate()
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

	if c.WpDir == "" {
		panic("wordpress path can't be empty")
	}
	r.Static("/wp-content/uploads", str.Join(c.WpDir, "/wp-content/uploads"))
	r.Static("/wp-content/themes", str.Join(c.WpDir, "/wp-content/themes"))
	r.Static("/wp-content/plugins", str.Join(c.WpDir, "/wp-content/plugins"))
	r.Static("/wp-includes/css", str.Join(c.WpDir, "/wp-includes/css"))
	r.Static("/wp-includes/fonts", str.Join(c.WpDir, "/wp-includes/fonts"))
	r.Static("/wp-includes/js", str.Join(c.WpDir, "/wp-includes/js"))

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("go-wp", store))
	r.GET("/", actions.Feed, middleware.SearchLimit(c.SingleIpSearchNum),
		actions.ThemeHook(constraints.Home))
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
	r.GET("/p/:id/comment-page-:page", actions.ThemeHook(constraints.Detail))
	r.GET("/p/:id/feed", actions.PostFeed)
	r.GET("/feed", actions.SiteFeed)
	r.GET("/comments/feed", actions.CommentsFeed)
	//r.NoRoute(actions.ThemeHook(constraints.NoRoute))
	commentMiddleWare, _ := middleware.FlowLimit(c.MaxRequestSleepNum, 5, c.CacheTime.SleepTime)
	r.POST("/comment", commentMiddleWare, actions.PostComment)
	if c.Pprof != "" {
		pprof.Register(r, c.Pprof)
	}
	for _, fn := range hooker {
		fn(r)
	}

	reload.Append(func() {
		c := config.GetConfig()
		siteFlow(c.MaxRequestSleepNum, c.MaxRequestNum, c.CacheTime.SleepTime)
	}, "site-flowLimit-config")
	return r
}
