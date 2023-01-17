package theme

import (
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/internal/plugins"
	"github/fthvgb1/wp-go/internal/templates/twentyseventeen"
)

var themeMap = map[string]func(*gin.Context, gin.H, int) string{}

func InitTheme() {
	HookFunc(twentyseventeen.ThemeName, twentyseventeen.Hook)
}

func HookFunc(themeName string, fn func(*gin.Context, gin.H, int) string) {
	themeMap[themeName] = fn
}

func Hook(themeName string, c *gin.Context, h gin.H, scene int) string {
	fn, ok := themeMap[themeName]
	if ok && fn != nil {
		return fn(c, h, scene)
	}
	if _, ok := plugins.IndexSceneMap[scene]; ok {
		return "twentyfifteen/posts/index.gohtml"
	} else if _, ok := plugins.DetailSceneMap[scene]; ok {
		return "twentyfifteen/posts/detail.gohtml"
	}
	return "twentyfifteen/posts/detail.gohtml"
}
