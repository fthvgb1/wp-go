package twentyfifteen

import (
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/plugins"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/theme/wp/components"
	"github.com/fthvgb1/wp-go/app/theme/wp/middleware"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache/reload"
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
	middleware.CommonMiddleware(h)
	setPaginationAndRender(h)
	h.PushCacheGroupHeadScript(constraints.AllScene, "CalCustomBackGround", 10.005, CalCustomBackGround)
	h.PushCacheGroupHeadScript(constraints.AllScene, "colorSchemeCss", 10.0056, colorSchemeCss)
	h.CommonComponents()
	components.WidgetArea(h)
	h.PushRender(constraints.AllScene, wp.NewHandleFn(renderCustomHeader, 20.5, "renderCustomHeader"))
	wp.PushIndexHandler(constraints.PipeRender, h, wp.NewHandleFn(wp.IndexRender, 50.005, "wp.IndexRender"))
	h.PushRender(constraints.Detail, wp.NewHandleFn(wp.DetailRender, 50.005, "wp.DetailRender"))
	h.PushRender(constraints.Detail, wp.NewHandleFn(postThumb, 60.005, "postThumb"))
	h.PushDataHandler(constraints.Detail, wp.NewHandleFn(wp.Detail, 100.005, "wp.Detail"))
	wp.PushIndexHandler(constraints.PipeData, h, wp.NewHandleFn(wp.Index, 100.005, "wp.Index"))
	h.PushDataHandler(constraints.AllScene, wp.NewHandleFn(wp.PreCodeAndStats, 80.005, "wp.PreCodeAndStats"))
	h.PushRender(constraints.AllScene, wp.NewHandleFn(wp.PreTemplate, 70.005, "wp.PreTemplate"))
}

func setPaginationAndRender(h *wp.Handle) {
	h.PushHandler(constraints.PipeRender, constraints.Detail, wp.NewHandleFn(func(hh *wp.Handle) {
		d := hh.GetDetailHandle()
		d.CommentRender = plugins.CommentRender()
		d.CommentPageEle = plugins.TwentyFifteenCommentPagination()
	}, 150, "setPaginationAndRender"))

	wp.PushIndexHandler(constraints.PipeRender, h, wp.NewHandleFn(func(hh *wp.Handle) {
		i := hh.GetIndexHandle()
		i.SetPageEle(plugins.TwentyFifteenPagination())
	}, 150, "setPaginationAndRender"))
}

func postThumb(h *wp.Handle) {
	d := h.GetDetailHandle()
	if d.Post.Thumbnail.Path != "" {
		d.Post.Thumbnail = getPostThumbs(d.Post.Id, d.Post)
	}
}

var getPostThumbs = reload.BuildMapFn[uint64]("twentyFifteen-post-thumbnail", postThumbs)

func postThumbs(post models.Posts) models.PostThumbnail {
	return wpconfig.Thumbnail(post.Thumbnail.OriginAttachmentData, "post-thumbnail", "")
}

func renderCustomHeader(h *wp.Handle) {
	h.SetData("customHeader", customHeader(h))
}
