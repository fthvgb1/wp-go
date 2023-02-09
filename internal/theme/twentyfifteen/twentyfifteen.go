package twentyfifteen

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/fthvgb1/wp-go/plugin/pagination"
)

const ThemeName = "twentyfifteen"

type handle struct {
	common.Handle
	templ string
}

func Hook(h2 common.Handle) {
	h := handle{
		Handle: h2,
		templ:  "twentyfifteen/posts/index.gohtml",
	}
	//h.GinH["HeaderImage"] = h.getHeaderImage(h.C)
	if h.Stats == plugins.Empty404 {
		h.C.HTML(h.Code, h.templ, h.GinH)
		return
	}
	if h.Scene == plugins.Detail {
		h.Detail()
		return
	}
	h.Index()
}

var plugin = common.Plugins()

func (h handle) Index() {
	if h.Stats != plugins.Empty404 {

		h.GinH["posts"] = slice.Map(
			h.GinH["posts"].([]models.Posts),
			common.PluginFn[models.Posts](plugin, h.Handle, common.DigestsAndOthers(h.C)))

		p, ok := h.GinH["pagination"]
		if ok {
			pp, ok := p.(pagination.ParsePagination)
			if ok {
				h.GinH["pagination"] = pagination.Paginate(plugins.TwentyFifteenPagination(), pp)
			}
		}
	}

	//h.GinH["bodyClass"] = h.bodyClass()
	h.C.HTML(h.Code, h.templ, h.GinH)
}

func (h handle) Detail() {
	//h.GinH["bodyClass"] = h.bodyClass()
	//host, _ := wpconfig.Options.Load("siteurl")
	if h.GinH["comments"] != nil {
		comments := h.GinH["comments"].([]models.Comments)
		dep := h.GinH["maxDep"].(int)
		h.GinH["comments"] = plugins.FormatComments(h.C, plugins.CommentRender(), comments, dep)
	}

	h.templ = "twentyfifteen/posts/detail.gohtml"
	h.C.HTML(h.Code, h.templ, h.GinH)
}
