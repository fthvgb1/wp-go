package twentyfifteen

import (
	"github.com/fthvgb1/wp-go/app/pkg/config"
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
	conf := config.GetConfig()
	h.PushComponentFilterFn(widgets.Search, func(h *wp.Handle, s string, args ...any) string {
		return strings.ReplaceAll(s, `class="search-submit"`, `class="search-submit screen-reader-text"`)
	})
	wp.InitPipe(h)
	h.PushHandler(constraints.PipeMiddleware, constraints.Home,
		wp.NewHandleFn(widget.IsCategory, 100, "widget.IsCategory"))

	h.Index.SetPageEle(plugins.TwentyFifteenPagination())
	h.PushCacheGroupHeadScript(constraints.AllScene, "CalCustomBackGround", 10, CalCustomBackGround)
	h.PushCacheGroupHeadScript(constraints.AllScene, "colorSchemeCss", 10, colorSchemeCss)
	h.CommonComponents()
	h.Index.SetListPlugin(wp.PostsPlugins(wp.PostPlugin(), wp.GetListPostPlugins(conf.ListPagePlugins, wp.ListPostPlugins())...))
	components.WidgetArea(h)
	wp.ReplyCommentJs(h)
	h.SetData("customHeader", customHeader(h))
	wp.PushIndexHandler(constraints.PipeRender, h, wp.NewHandleFn(wp.IndexRender, 50, "wp.IndexRender"))
	h.PushRender(constraints.Detail, wp.NewHandleFn(wp.DetailRender, 50, "wp.DetailRender"))
	h.PushRender(constraints.Detail, wp.NewHandleFn(postThumb, 60, "postThumb"))
	h.PushDataHandler(constraints.Detail, wp.NewHandleFn(wp.Detail, 100, "wp.Detail"))
	wp.PushIndexHandler(constraints.PipeData, h, wp.NewHandleFn(wp.Index, 100, "wp.Index"))
	h.PushDataHandler(constraints.AllScene, wp.NewHandleFn(wp.PreCodeAndStats, 80, "wp.PreCodeAndStats"))
	h.PushRender(constraints.AllScene, wp.NewHandleFn(wp.PreTemplate, 70, "wp.PreTemplate"))
}

func postThumb(h *wp.Handle) {
	if h.Detail.Post.Thumbnail.Path != "" {
		h.Detail.Post.Thumbnail = wpconfig.Thumbnail(h.Detail.Post.Thumbnail.OriginAttachmentData, "post-thumbnail", "")
	}
}
