package route

import (
	"github.com/fthvgb1/wp-go/app/actions"
	"github.com/fthvgb1/wp-go/app/middleware"
	"github.com/fthvgb1/wp-go/app/pkg/config"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/theme"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/helper/slice/mockmap"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type GinSetter func(*gin.Engine)

var setters mockmap.Map[string, GinSetter]

var setterHooks []func(item mockmap.Item[string, GinSetter]) (mockmap.Item[string, GinSetter], bool)

// SetGinAction 方便插件在init时使用
func SetGinAction(name string, hook GinSetter, orders ...float64) {
	setters.Set(name, hook, orders...)
}

func HookGinSetter(fn func(item mockmap.Item[string, GinSetter]) (mockmap.Item[string, GinSetter], bool)) {
	setterHooks = append(setterHooks, fn)
}

// DelGinSetter 方便插件在init时使用
func DelGinSetter(name string) {
	setterHooks = append(setterHooks, func(item mockmap.Item[string, GinSetter]) (mockmap.Item[string, GinSetter], bool) {
		return item, item.Name != name
	})
}

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.New()
	c := config.GetConfig()
	SetGinAction("initTrustIp", func(r *gin.Engine) {
		if len(c.TrustIps) > 0 {
			err := r.SetTrustedProxies(c.TrustIps)
			if err != nil {
				panic(err)
			}
		}
	}, 99.5)

	SetGinAction("setTemplate", func(r *gin.Engine) {
		r.HTMLRender = theme.BuildTemplate()
		wpconfig.SetTemplateFs(theme.TemplateFs)
	}, 90.5)

	siteFlowLimitMiddleware, siteFlow := middleware.WPFlowLimit(c.MaxRequestSleepNum, c.MaxRequestNum, c.CacheTime.SleepTime)
	reload.Append(func() {
		c = config.GetConfig()
		siteFlow(c.MaxRequestSleepNum, c.MaxRequestNum, c.CacheTime.SleepTime)
	}, "site-flowLimit-config")

	SetGinAction("setGlobalMiddleware", func(r *gin.Engine) {
		r.Use(
			gin.Logger(),
			middleware.ValidateServerNames(),
			middleware.RecoverAndSendMail(gin.DefaultErrorWriter),
			siteFlowLimitMiddleware,
			middleware.SetStaticFileCache,
		)
	}, 88.5)

	SetGinAction("setGzip", func(r *gin.Engine) {
		//gzip 因为一般会用nginx做反代时自动使用gzip,所以go这边本身可以不用
		if c.Gzip {
			r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedPaths([]string{
				"/wp-includes/", "/wp-content/",
			})))
		}
	}, 87.6)

	SetGinAction("setWpDir", func(r *gin.Engine) {
		if c.WpDir == "" {
			panic("wordpress path can't be empty")
		}
		r.Static("/wp-content/uploads", str.Join(c.WpDir, "/wp-content/uploads"))
		r.Static("/wp-content/themes", str.Join(c.WpDir, "/wp-content/themes"))
		r.Static("/wp-content/plugins", str.Join(c.WpDir, "/wp-content/plugins"))
		r.Static("/wp-includes/css", str.Join(c.WpDir, "/wp-includes/css"))
		r.Static("/wp-includes/fonts", str.Join(c.WpDir, "/wp-includes/fonts"))
		r.Static("/wp-includes/js", str.Join(c.WpDir, "/wp-includes/js"))
	}, 86.1)

	SetGinAction("setSession", func(r *gin.Engine) {
		store := cookie.NewStore([]byte("secret"))
		r.Use(sessions.Sessions("go-wp", store))
	}, 85.1)

	SetGinAction("setRoute", func(r *gin.Engine) {
		r.GET("/", actions.Feed, middleware.SearchLimit(c.SingleIpSearchNum, 5),
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
		commentMiddleWare, _ := middleware.FlowLimit(5, c.SingleIpSearchNum, c.CacheTime.SleepTime)
		commentIpMiddleware, _ := middleware.IpLimit(5, 2)
		r.POST("/comment", commentMiddleWare, commentIpMiddleware, actions.PostComment)

		r.NoRoute(actions.ThemeHook(constraints.NoRoute))
	}, 84.6)

	SetGinAction("setpprof", func(r *gin.Engine) {
		if c.Pprof != "" {
			pprof.Register(r, c.Pprof)
		}
	}, 80.8)

	for _, hook := range setterHooks {
		setters = slice.FilterAndMap(setters, hook)
	}

	slice.SimpleSort(setters, slice.DESC, func(t mockmap.Item[string, GinSetter]) float64 {
		return t.Order
	})

	for _, fn := range setters {
		fn.Value(r)
	}

	return r
}
