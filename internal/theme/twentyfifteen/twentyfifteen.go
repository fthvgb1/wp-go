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
		return
	}
	err = json.Unmarshal(b, &themesupport)
	logs.ErrPrintln(err, "解析themesupport失败")
	bytes, err := fs.ReadFile("twentyfifteen/colorscheme.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &colorscheme)
	if err != nil {
		return
	}
	logs.ErrPrintln(err, "解析colorscheme失败")
}

var pipe = wp.HandlePipe(wp.Render, widget.MiddleWare(ready)...)

func Hook(h *wp.Handle) {
	pipe(h)
}

func ready(next wp.HandleFn[*wp.Handle], h *wp.Handle) {
	h.GetPassword()
	h.PushComponentFilterFn(widgets.Search, func(h *wp.Handle, s string) string {
		return strings.ReplaceAll(s, `class="search-submit"`, `class="search-submit screen-reader-text"`)
	})
	wphandle.RegisterPlugins(h, config.GetConfig().Plugins...)

	h.PushCacheGroupHeadScript("CalCustomBackGround", 10, CalCustomBackGround, colorSchemeCss)
	h.PushHandleFn(constraints.Ok, wp.NewHandleFn(components.WidgetArea, 20))
	h.PushHandleFn(constraints.AllStats, wp.NewHandleFn(customHeader, 10))
	h.PushHandleFn(constraints.AllStats, wp.NewHandleFn(wp.Indexs, 100))
	h.PushHandleFn(constraints.Detail, wp.NewHandleFn(wp.Details, 100))
	next(h)
}
