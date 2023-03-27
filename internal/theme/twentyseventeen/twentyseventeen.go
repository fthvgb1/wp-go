package twentyseventeen

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/maps"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/plugins/wphandle"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/gin-gonic/gin"
	"strings"
)

const ThemeName = "twentyseventeen"

func Init(fs embed.FS) {
	b, err := fs.ReadFile(str.Join(ThemeName, "/themesupport.json"))
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &themesupport)
}

var paginate = func() plugins.PageEle {
	p := plugins.TwentyFifteenPagination()
	p.PrevEle = `<a class="prev page-numbers" href="%s"><svg class="icon icon-arrow-left" aria-hidden="true" role="img"> <use href="#icon-arrow-left" xlink:href="#icon-arrow-left"></use> </svg>
<span class="screen-reader-text">上一页</span></a>`
	p.NextEle = strings.Replace(p.NextEle, "下一页", `<span class="screen-reader-text">下一页</span>
<svg class="icon icon-arrow-right" aria-hidden="true" role="img"> <use href="#icon-arrow-right" xlink:href="#icon-arrow-right"></use> 
</svg>`, 1)
	return p
}()

var pipe = wp.HandlePipe(wp.Render, widget.MiddleWare(ready)...)

func Hook(h *wp.Handle) {
	pipe(h)
}

func ready(next wp.HandleFn[*wp.Handle], h *wp.Handle) {
	h.GetPassword()
	wphandle.RegisterPlugins(h, config.GetConfig().Plugins...)
	h.PushHandleFn(constraints.AllStats, wp.NewHandleFn(calClass, 20))
	h.PushHandleFn(constraints.AllStats, wp.NewHandleFn(index, 100))
	h.PushHandleFn(constraints.Detail, wp.NewHandleFn(detail, 100))
	h.PushGroupHandleFn(constraints.AllStats, 90, wp.PreCodeAndStats, wp.PreTemplate, errorsHandle)
	h.PushCacheGroupHeadScript("colorScheme-customHeader", 10, colorScheme, customHeader)
	h.PushHandleFn(constraints.Ok, wp.NewHandleFn(components.WidgetArea, 20))
	pushScripts(h)
	h.SetData("HeaderImage", getHeaderImage(h))
	h.SetComponentsArgs(widgets.Widget, map[string]string{
		"{$before_widget}": `<section id="%s" class="%s">`,
		"{$after_widget}":  `</section>`,
	})
	wp.SetComponentsArgsForMap(h, widgets.Search, "{$form}", searchForm)

	next(h)
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

var listPostsPlugins = func() map[string]wp.Plugin[models.Posts, *wp.Handle] {
	return maps.Merge(wp.ListPostPlugins(), map[string]wp.Plugin[models.Posts, *wp.Handle]{
		"twentyseventeen_postThumbnail": postThumbnail,
	})
}()

func errorsHandle(h *wp.Handle) {
	switch h.Stats {
	case constraints.Error404, constraints.InternalErr, constraints.ParamError:
		logs.ErrPrintln(h.Err(), "报错：")
		h.SetTempl("twentyseventeen/posts/error.gohtml")
	}
}

func index(h *wp.Handle) {
	if h.Scene() == constraints.Detail {
		return
	}
	i := h.Index
	err := i.BuildIndexData(wp.NewIndexParams(i.C))
	if err != nil {
		i.SetErr(err)
	}
	h.SetData("scene", h.Scene())
	i.SetPageEle(paginate)
	i.SetPostsPlugins(listPostsPlugins)
}

func detail(h *wp.Handle) {
	d := h.Detail
	err := d.BuildDetailData()
	if err != nil {
		d.SetErr(err)
	}
	if d.Post.Thumbnail.Path != "" {
		img := wpconfig.Thumbnail(d.Post.Thumbnail.OriginAttachmentData, "full", "", "thumbnail", "post-thumbnail")
		img.Sizes = "100vw"
		img.Srcset = fmt.Sprintf("%s %dw, %s", img.Path, img.Width, img.Srcset)
		d.Post.Thumbnail = img
	}
	d.CommentRender = commentFormat
}

var commentFormat = comment{}

type comment struct {
	plugins.CommonCommentFormat
}

func (c comment) FormatLi(ctx *gin.Context, m models.Comments, depth int, isTls bool, eo, parent string) string {
	templ := plugins.CommonLi()
	templ = strings.ReplaceAll(templ, `<a rel="nofollow" class="comment-reply-link"
               href="/p/{{PostId}}?replytocom={{CommentId}}#respond" data-commentid="{{CommentId}}" data-postid="{{PostId}}"
               data-belowelement="div-comment-{{CommentId}}" data-respondelement="respond"
               data-replyto="回复给{{CommentAuthor}}"
               aria-label="回复给{{CommentAuthor}}">回复</a>`, `<a rel="nofollow" class="comment-reply-link"
               href="/p/{{PostId}}?replytocom={{CommentId}}#respond" data-commentid="{{CommentId}}" data-postid="{{PostId}}"
               data-belowelement="div-comment-{{CommentId}}" data-respondelement="respond"
               data-replyto="回复给{{CommentAuthor}}"
               aria-label="回复给{{CommentAuthor}}"><svg class="icon icon-mail-reply" aria-hidden="true" role="img"> <use href="#icon-mail-reply" xlink:href="#icon-mail-reply"></use> </svg>回复</a>`)
	return plugins.FormatLi(templ, ctx, m, depth, isTls, eo, parent)
}

func postThumbnail(next wp.Fn[models.Posts], h *wp.Handle, t models.Posts) models.Posts {
	if t.Thumbnail.Path != "" {
		t.Thumbnail.Sizes = "(max-width: 767px) 89vw, (max-width: 1000px) 54vw, (max-width: 1071px) 543px, 580px"
		if h.Scene() == constraints.Detail {
			t.Thumbnail.Sizes = "100vw"
		}
	}
	return next(t)
}

var header = reload.Vars(models.PostThumbnail{})

func getHeaderImage(h *wp.Handle) (r models.PostThumbnail) {
	img := header.Load()
	if img.Path != "" {
		return img
	}
	image, rand := h.GetCustomHeader()
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

func calClass(h *wp.Handle) {
	themeMods := h.CommonThemeMods()
	u := wpconfig.GetThemeModsVal(ThemeName, "header_image", themeMods.ThemeSupport.CustomHeader.DefaultImage)
	var class []string
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
	h.PushClass(class...)
}
