package theme

import (
	"errors"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/plugin/pagination"
	"github.com/gin-gonic/gin"
)

var themeMap = map[string]func(int, *gin.Context, gin.H, int, int){}

func addThemeHookFunc(name string, fn func(int, *gin.Context, gin.H, int, int)) {
	if _, ok := themeMap[name]; ok {
		panic("exists same name theme")
	}
	themeMap[name] = fn
}

func Hook(themeName string, code int, c *gin.Context, h gin.H, scene, status int) {
	fn, ok := themeMap[themeName]
	if ok && fn != nil {
		fn(code, c, h, scene, status)
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
		c.HTML(code, "twentyfifteen/posts/index.gohtml", h)
		return
	} else if scene == plugins.Detail {
		h["comments"] = plugins.FormatComments(c, plugins.CommentRender(), h["comments"].([]models.Comments), h["maxDep"].(int))
		c.HTML(code, "twentyfifteen/posts/detail.gohtml", h)
		return
	}
	logs.ErrPrintln(errors.New("what happening"), " how reached here", themeName, code, h, scene, status)
}
