package twentyfifteen

import (
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/app/plugins"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/theme/wp/components"
	"github.com/fthvgb1/wp-go/app/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"strings"
)

const ThemeName = "twentyfifteen"

func Hook(h *wp.Handle) {
	wp.Run(h, configs)
}

func configs(h *wp.Handle) {
	h.AddActionFilter(widgets.Search, func(h *wp.Handle, s string, args ...any) string {
		return strings.ReplaceAll(s, `class="search-submit"`, `class="search-submit screen-reader-text"`)
	})
	wp.InitPipe(h)
	h.PushHandler(constraints.PipeMiddleware, constraints.Home,
		wp.NewHandleFn(widget.CheckCategory, 100, "widget.CheckCategory"))

	h.Index.SetPageEle(plugins.TwentyFifteenPagination())
	h.PushCacheGroupHeadScript(constraints.AllScene, "CalCustomBackGround", 10.005, CalCustomBackGround)
	h.PushCacheGroupHeadScript(constraints.AllScene, "colorSchemeCss", 10.0056, colorSchemeCss)
	h.CommonComponents()
	components.WidgetArea(h)
	wp.ReplyCommentJs(h)
	h.SetData("customHeader", customHeader(h))
	wp.PushIndexHandler(constraints.PipeRender, h, wp.NewHandleFn(wp.IndexRender, 50.005, "wp.IndexRender"))
	h.PushRender(constraints.Detail, wp.NewHandleFn(wp.DetailRender, 50.005, "wp.DetailRender"))
	h.PushRender(constraints.Detail, wp.NewHandleFn(postThumb, 60.005, "postThumb"))
	h.PushDataHandler(constraints.Detail, wp.NewHandleFn(wp.Detail, 100.005, "wp.Detail"))
	wp.PushIndexHandler(constraints.PipeData, h, wp.NewHandleFn(wp.Index, 100.005, "wp.Index"))
	h.PushDataHandler(constraints.AllScene, wp.NewHandleFn(wp.PreCodeAndStats, 80.005, "wp.PreCodeAndStats"))
	h.PushRender(constraints.AllScene, wp.NewHandleFn(wp.PreTemplate, 70.005, "wp.PreTemplate"))
}

func postThumb(h *wp.Handle) {
	if h.Detail.Post.Thumbnail.Path != "" {
		h.Detail.Post.Thumbnail = wpconfig.Thumbnail(h.Detail.Post.Thumbnail.OriginAttachmentData, "post-thumbnail", "")
	}
}
