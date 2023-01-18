package theme

import (
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/plugin/pagination"
	"github.com/gin-gonic/gin"
)

var themeMap = map[string]func(int, *gin.Context, gin.H, int){}

func AddThemeHookFunc(name string, fn func(int, *gin.Context, gin.H, int)) {
	if _, ok := themeMap[name]; ok {
		panic("exists same name theme")
	}
	themeMap[name] = fn
}

func Hook(themeName string, status int, c *gin.Context, h gin.H, scene int) {
	fn, ok := themeMap[themeName]
	if ok && fn != nil {
		fn(status, c, h, scene)
		return
	}
	if _, ok := plugins.IndexSceneMap[scene]; ok {
		p, ok := h["pagination"]
		if ok {
			pp, ok := p.(pagination.ParsePagination)
			if ok {
				h["pagination"] = pagination.Paginate(plugins.TwentyFifteenPagination(), pp)
			}
		}
		c.HTML(status, "twentyfifteen/posts/index.gohtml", h)
		return
	} else if _, ok := plugins.DetailSceneMap[scene]; ok {
		c.HTML(status, "twentyfifteen/posts/detail.gohtml", h)
		return
	}
	c.HTML(status, "twentyfifteen/posts/index.gohtml", h)
}
