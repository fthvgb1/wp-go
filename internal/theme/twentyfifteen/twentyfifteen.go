package twentyfifteen

import (
	"embed"
	"encoding/json"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/plugins/wphandle"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components/widget"
	"strings"
)

const ThemeName = "twentyfifteen"

func Init(fs embed.FS) {
	b, err := fs.ReadFile("twentyfifteen/themesupport.json")
	if err != nil {
		logs.Error(err, "读取themesupport.json失败")
		return
	}
	err = json.Unmarshal(b, &themesupport)
	if err != nil {
		logs.Error(err, "解析themesupport失败")
		return
	}
	bytes, err := fs.ReadFile("twentyfifteen/colorscheme.json")
	if err != nil {
		logs.Error(err, "读取colorscheme.json失败")
		return
	}
	err = json.Unmarshal(bytes, &colorscheme)
	if err != nil {
		logs.Error(err, "解析colorscheme失败")
		return
	}
}

var pipe = wp.HandlePipe(wp.ExecuteHandleFn, widget.MiddleWare(ready, data)...)

func Hook(h *wp.Handle) {
	pipe(h)
}

func ready(next wp.HandleFn[*wp.Handle], h *wp.Handle) {
	wp.InitThemeArgAndConfig(configs, h)
	h.GetPassword()
	next(h)
}

func data(next wp.HandleFn[*wp.Handle], h *wp.Handle) {
	if h.Scene() == constraints.Detail {
		wp.Details(h)
	} else {
		wp.Indexs(h)
	}
	h.DetermineHandleFns()
	next(h)
}

func configs(h *wp.Handle) {
	h.PushComponentFilterFn(widgets.Search, func(h *wp.Handle, s string, args ...any) string {
		return strings.ReplaceAll(s, `class="search-submit"`, `class="search-submit screen-reader-text"`)
	})
	wphandle.RegisterPlugins(h, config.GetConfig().Plugins...)
	h.PushCacheGroupHeadScript("CalCustomBackGround", 10, CalCustomBackGround, colorSchemeCss)
	h.CommonComponents()
	h.PushHandleFn(constraints.Ok, wp.NewHandleFn(components.WidgetArea, 20))
	h.PushHandleFn(constraints.AllStats, wp.NewHandleFn(customHeader, 10))
	h.PushHandleFn(constraints.AllStats, wp.NewHandleFn(wp.IndexRender, 50))
	h.PushHandleFn(constraints.Detail, wp.NewHandleFn(wp.DetailRender, 50))
}
