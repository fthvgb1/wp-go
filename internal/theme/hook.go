package theme

import (
	"github.com/fthvgb1/wp-go/internal/theme/common"
)

var themeMap = map[string]func(handle common.Handle){}

func addThemeHookFunc(name string, fn func(handle common.Handle)) {
	if _, ok := themeMap[name]; ok {
		panic("exists same name theme")
	}
	themeMap[name] = fn
}

func Hook(themeName string, handle common.Handle) {
	fn, ok := themeMap[themeName]
	if ok && fn != nil {
		fn(handle)
		return
	}
	/*if _, ok := plugins.IndexSceneMap[scene]; ok {
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
	logs.ErrPrintln(errors.New("what happening"), " how reached here", themeName, code, h, scene, status)*/
}
