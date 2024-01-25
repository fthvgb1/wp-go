package twentyseventeen

import (
	"context"
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/plugins"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/theme/wp/components"
	"github.com/fthvgb1/wp-go/app/theme/wp/middleware"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/plugin/pagination"
	"strings"
)

const ThemeName = "twentyseventeen"

var paginate = func() plugins.PageEle {
	p := plugins.TwentyFifteenPagination()
	p.PrevEle = `<a class="prev page-numbers" href="%s"><svg class="icon icon-arrow-left" aria-hidden="true" role="img"> <use href="#icon-arrow-left" xlink:href="#icon-arrow-left"></use> </svg>
<span class="screen-reader-text">上一页</span></a>`
	p.NextEle = strings.Replace(p.NextEle, "下一页", `<span class="screen-reader-text">下一页</span>
<svg class="icon icon-arrow-right" aria-hidden="true" role="img"> <use href="#icon-arrow-right" xlink:href="#icon-arrow-right"></use> 
</svg>`, 1)
	commentPageEle = plugins.PaginationNav{
		Currents: p.Current,
		Prevs:    p.Prev,
		Nexts:    p.Next,
		Dotss:    p.Dots,
		Middles:  p.Middle,
		Steps: func() int {
			return 2
		},
		Urlss: plugins.TwentyFifteenCommentPagination().Urls,
	}

	return p
}()

var commentPageEle pagination.Render

func Hook(h *wp.Handle) {
	wp.Run(h, configs)
}

func configs(h *wp.Handle) {
	wp.InitPipe(h)
	middleware.CommonMiddleware(h)
	h.AddActionFilter("bodyClass", calClass)
	h.PushCacheGroupHeadScript(constraints.AllScene, "colorScheme-customHeader", 10, colorScheme, customHeader)
	components.WidgetArea(h)
	pushScripts(h)
	h.PushRender(constraints.AllStats, wp.NewHandleFn(calCustomHeader, 10.005, "calCustomHeader"))
	wp.SetComponentsArgs(widgets.Widget, map[string]string{
		"{$before_widget}": `<section id="%s" class="%s">`,
		"{$after_widget}":  `</section>`,
	})
	h.PushRender(constraints.AllStats,
		wp.NewHandleFn(wp.PreTemplate, 70.005, "wp.PreTemplate"),
		wp.NewHandleFn(errorsHandle, 80.005, "errorsHandle"),
	)
	videoHeader(h)
	h.SetData("colophon", colophon)
	setPaginationAndRender(h)
	h.CommonComponents()
	h.PushPostPlugin(postThumbnail)
	wp.SetComponentsArgsForMap(widgets.Search, "{$form}", searchForm)
	wp.PushIndexHandler(constraints.PipeRender, h, wp.NewHandleFn(wp.IndexRender, 10.005, "wp.IndexRender"))
	h.PushRender(constraints.Detail, wp.NewHandleFn(wp.DetailRender, 10.005, "wp.DetailRender"))
	h.PushDataHandler(constraints.Detail, wp.NewHandleFn(wp.Detail, 100.005, "wp.Detail"), wp.NewHandleFn(postThumb, 90.005, "{theme}.postThumb"))
	wp.PushIndexHandler(constraints.PipeData, h, wp.NewHandleFn(wp.Index, 100.005, "wp.Index"))
	h.PushDataHandler(constraints.AllScene, wp.NewHandleFn(wp.PreCodeAndStats, 90.005, "wp.PreCodeAndStats"))
}

func setPaginationAndRender(h *wp.Handle) {
	h.PushHandler(constraints.PipeRender, constraints.Detail, wp.NewHandleFn(func(hh *wp.Handle) {
		d := hh.GetDetailHandle()
		d.CommentRender = commentFormat
		d.CommentPageEle = commentPageEle
	}, 150, "setPaginationAndRender"))
	wp.PushIndexHandler(constraints.PipeRender, h, wp.NewHandleFn(func(hh *wp.Handle) {
		i := hh.GetIndexHandle()
		i.SetPageEle(paginate)
	}, 150, "setPaginationAndRender"))
}

var searchForm = `<form role="search" method="get" class="search-form" action="/">
	<label for="search-form-1">
		<span class="screen-reader-text">{$label}：</span>
	</label>
	<input type="search" id="search-form-1" class="search-field" placeholder="{$placeholder}…" value="{$value}" name="s">
	<button type="submit" class="search-submit">
<svg class="icon icon-search" aria-hidden="true" role="img"> <use href="#icon-search" xlink:href="#icon-search"></use> </svg>
<span class="screen-reader-text">{$button}</span>
</button>
</form>`

func errorsHandle(h *wp.Handle) {
	switch h.Stats {
	case constraints.Error404, constraints.InternalErr, constraints.ParamError:
		logs.IfError(h.Err(), "报错：")
		h.SetTempl("twentyseventeen/posts/error.gohtml")
	}
}

func postThumb(h *wp.Handle) {
	d := h.GetDetailHandle()
	if d.Post.Thumbnail.Path != "" {
		img := wpconfig.Thumbnail(d.Post.Thumbnail.OriginAttachmentData, "full", "", "thumbnail", "post-thumbnail")
		img.Sizes = "100vw"
		img.Srcset = fmt.Sprintf("%s %dw, %s", img.Path, img.Width, img.Srcset)
		d.Post.Thumbnail = img
	}
}

var commentFormat = comment{}

type comment struct {
	plugins.CommonCommentFormat
}

var commentLi = plugins.CommonLi()

var respondFn = plugins.Responds(respondStr)

func (c comment) FormatLi(_ context.Context, m models.Comments, depth, maxDepth, page int, isTls, isThreadComments bool, eo, parent string) string {
	return plugins.FormatLi(commentLi, m, respondFn, depth, maxDepth, page, isTls, isThreadComments, eo, parent)
}

var colophon = `<footer id="colophon" class="site-footer">
            <div class="wrap">
                <div class="site-info">
                    <a href="https://github.com/fthvgb1/wp-go" class="imprint">自豪地采用 wp-go</a>
                </div>
            </div>
        </footer>`

var respondStr = `<a rel="nofollow" class="comment-reply-link"
               href="/p/{{PostId}}?replytocom={{CommentId}}#respond" data-commentid="{{CommentId}}" data-postid="{{PostId}}"
               data-belowelement="div-comment-{{CommentId}}" data-respondelement="respond"
               data-replyto="回复给{{CommentAuthor}}"
               aria-label="回复给{{CommentAuthor}}"><svg class="icon icon-mail-reply" aria-hidden="true" role="img"> <use href="#icon-mail-reply" xlink:href="#icon-mail-reply"></use> </svg>回复</a>`

func postThumbnail(h *wp.Handle, posts *models.Posts) {
	if posts.Thumbnail.Path != "" {
		posts.Thumbnail.Sizes = "(max-width: 767px) 89vw, (max-width: 1000px) 54vw, (max-width: 1071px) 543px, 580px"
		if h.Scene() == constraints.Detail {
			posts.Thumbnail.Sizes = "100vw"
		}
	}
}

var header = reload.Vars(models.PostThumbnail{}, "twentyseventeen-headerImage")

func calCustomHeader(h *wp.Handle) {
	h.SetData("HeaderImage", getHeaderImage(h))
}

func getHeaderImage(h *wp.Handle) (r models.PostThumbnail) {
	img := header.Load()
	if img.Path != "" {
		return img
	}
	image, rand := h.GetCustomHeaderImg()
	if image.Path != "" {
		r = image
		r.Sizes = "100vw"
		if !rand {
			header.Store(r)
		}
		return
	}
	r.Path = helper.CutUrlHost(h.CommonThemeMods().ThemeSupport.CustomHeader.DefaultImage)
	r.Width = 2000
	r.Height = 1200
	header.Store(r)
	return
}

func calClass(h *wp.Handle, s string, _ ...any) string {
	class := strings.Split(s, " ")
	themeMods := h.CommonThemeMods()
	u := wpconfig.GetThemeModsVal(ThemeName, "header_image", themeMods.ThemeSupport.CustomHeader.DefaultImage)
	if u != "" && u != "remove-header" {
		class = append(class, "has-header-image")
	}
	if len(themeMods.SidebarsWidgets.Data.Sidebar1) > 0 {
		class = append(class, "has-sidebar")
	}
	if themeMods.HeaderTextcolor == "blank" {
		class = append(class, "title-tagline-hidden")
	}
	class = append(class, "hfeed")
	class = append(class, str.Join("colors-", wpconfig.GetThemeModsVal(ThemeName, "colorscheme", "light")))
	if h.Scene() == constraints.Archive {
		if "one-column" == wpconfig.GetThemeModsVal(ThemeName, "page_layout", "") {
			class = append(class, "page-one-column")
		} else {
			class = append(class, "page-two-column")
		}
	}
	return strings.Join(class, " ")
}

func videoHeader(h *wp.Handle) {
	h.AddActionFilter("videoSetting", videoPlay)
	wp.CustomVideo(h, constraints.Home)
}

func videoPlay(h *wp.Handle, _ string, a ...any) string {
	if len(a) < 1 {
		return ""
	}
	v, ok := a[0].(*wp.VideoSetting)
	if !ok {
		return ""
	}
	img := getHeaderImage(h)
	v.Width = img.Width
	v.Height = img.Height
	v.PosterUrl = img.Path
	v.L10n.Play = `<span class="screen-reader-text">播放背景视频</span><svg class="icon icon-play" aria-hidden="true" role="img"> <use href="#icon-play" xlink:href="#icon-play"></use> </svg>`
	v.L10n.Pause = `<span class="screen-reader-text">暂停背景视频</span><svg class="icon icon-pause" aria-hidden="true" role="img"> <use href="#icon-pause" xlink:href="#icon-pause"></use> </svg>`
	return ""
}
