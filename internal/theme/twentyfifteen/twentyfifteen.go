package twentyfifteen

import (
	"embed"
	"encoding/json"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/plugins/wphandle"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
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

func Hook(h *wp.Handle) {
	wp.Run(h, configs)
}

func configs(h *wp.Handle) {
	wphandle.UsePlugins(h)
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
	h.PushRender(constraints.AllStats, wp.NewHandleFn(wp.IndexRender, 50, "wp.IndexRender"))
	h.PushRender(constraints.Detail, wp.NewHandleFn(wp.DetailRender, 50, "wp.DetailRender"))
	h.PushRender(constraints.Detail, wp.NewHandleFn(postThumb, 60, "postThumb"))
	h.PushDataHandler(constraints.Detail, wp.NewHandleFn(wp.Details, 100, "wp.Details"))
	h.PushDataHandler(constraints.AllScene, wp.NewHandleFn(wp.Indexs, 100, "wp.Indexs"))
	h.PushDataHandler(constraints.AllScene, wp.NewHandleFn(wp.PreCodeAndStats, 80, "wp.PreCodeAndStats"))
	h.PushRender(constraints.AllScene, wp.NewHandleFn(wp.PreTemplate, 70, "wp.PreTemplate"))
}

func postThumb(h *wp.Handle) {
	if h.Detail.Post.Thumbnail.Path != "" {
		h.Detail.Post.Thumbnail = wpconfig.Thumbnail(h.Detail.Post.Thumbnail.OriginAttachmentData, "post-thumbnail", "")
	}
}
