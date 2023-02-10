package twentyfifteen

import (
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/fthvgb1/wp-go/plugin/pagination"
	"github.com/gin-contrib/sessions"
)

const ThemeName = "twentyfifteen"

type handle struct {
	common.Handle
	templ string
}

func Hook(cHandle common.Handle) {
	h := handle{
		Handle: cHandle,
		templ:  "twentyfifteen/posts/index.gohtml",
	}
	h.WidgetAreaData()
	h.Session = sessions.Default(h.C)

	if h.Stats == constraints.Error404 {
		h.C.HTML(h.Code, h.templ, h.GinH)
		return
	}
	if h.Scene == constraints.Detail {
		h.Detail()
		return
	}
	h.Index()
}

var plugin = common.ListPostPlugins()

func (h handle) Index() {
	err := h.Indexs()
	if err != nil {
		h.C.HTML(h.Code, h.templ, h.GinH)
		return
	}

	h.ExecListPagePlugin(plugin)
	page, ok := maps.GetStrMapAnyVal[pagination.ParsePagination](h.GinH, "pagination")
	if ok {
		h.GinH["pagination"] = pagination.Paginate(plugins.TwentyFifteenPagination(), page)
	}
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
