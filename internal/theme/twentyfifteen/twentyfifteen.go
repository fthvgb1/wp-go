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

var pipe = wp.HandlePipe(wp.NothingToDo, widget.MiddleWare(ready,
	wp.PipeHandle(constraints.PipeData, wp.PipeKey, wp.PipeDataHandle),
	wp.PipeHandle(constraints.PipeRender, wp.PipeKey, wp.PipeRender),
)...)

func Hook(h *wp.Handle) {
	pipe(h)
}

func ready(next wp.HandleFn[*wp.Handle], h *wp.Handle) {
	wp.InitThemeArgAndConfig(configs, h)
	h.GetPassword()
	next(h)
}

func configs(h *wp.Handle) {
	conf := config.GetConfig()
	h.PushComponentFilterFn(widgets.Search, func(h *wp.Handle, s string, args ...any) string {
		return strings.ReplaceAll(s, `class="search-submit"`, `class="search-submit screen-reader-text"`)
	})
	h.Index.SetPageEle(plugins.TwentyFifteenPagination())
	wphandle.RegisterPlugins(h, conf.Plugins...)
	h.PushCacheGroupHeadScript("CalCustomBackGround", 10, CalCustomBackGround, colorSchemeCss)
	h.CommonComponents()
	h.Index.SetListPlugin(wp.PostsPlugins(wp.PostPlugin(), wp.GetListPostPlugins(conf.ListPagePlugins, wp.ListPostPlugins())...))
	components.WidgetArea(h)
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
