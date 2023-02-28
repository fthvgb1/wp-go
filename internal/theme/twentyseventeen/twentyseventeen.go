package twentyseventeen

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/maps"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/gin-gonic/gin"
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
	return p
}()

var pipe = common.HandlePipe(common.Render, ready, dispatch)

func Hook(h *common.Handle) {
	pipe(h)
}

func ready(next common.HandleFn[*common.Handle], h *common.Handle) {
	h.WidgetAreaData()
	h.GetPassword()
	h.PushHandleFn(constraints.AllStats, common.NewHandleFn(calClass, 15))
	h.PushHeadScript(
		common.NewComponents(colorScheme, 10),
		common.NewComponents(customHeader, 10),
	)
	h.SetData("HeaderImage", getHeaderImage(h))
	h.SetData("scene", h.Scene())
	next(h)
}

func dispatch(next common.HandleFn[*common.Handle], h *common.Handle) {
	switch h.Scene() {
	case constraints.Detail:
		detail(next, h.Detail)
	default:
		index(next, h.Index)
	}
}

var listPostsPlugins = func() map[string]common.Plugin[models.Posts, *common.Handle] {
	return maps.Merge(common.ListPostPlugins(), map[string]common.Plugin[models.Posts, *common.Handle]{
		"twentyseventeen_postThumbnail": postThumbnail,
	})
}()

func index(next common.HandleFn[*common.Handle], i *common.IndexHandle) {
	err := i.BuildIndexData(common.NewIndexParams(i.C))
	if err != nil {
		i.SetTempl(str.Join(ThemeName, "/posts/error.gohtml"))
		i.Render()
		return
	}
	i.SetPageEle(paginate)
	i.SetPostsPlugins(listPostsPlugins)
	next(i.Handle)
}

func detail(next common.HandleFn[*common.Handle], d *common.DetailHandle) {
	err := d.BuildDetailData()
	if err != nil {
		d.SetTempl(str.Join(ThemeName, "/posts/error.gohtml"))
		d.Render()
		return
	}
	if d.Post.Thumbnail.Path != "" {
		img := wpconfig.Thumbnail(d.Post.Thumbnail.OriginAttachmentData, "full", "", "thumbnail", "post-thumbnail")
		img.Sizes = "100vw"
		img.Srcset = fmt.Sprintf("%s %dw, %s", img.Path, img.Width, img.Srcset)
		d.Post.Thumbnail = img
	}

	d.CommentRender = commentFormat

	next(d.Handle)
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

func postThumbnail(next common.Fn[models.Posts], h *common.Handle, t models.Posts) models.Posts {
	if t.Thumbnail.Path != "" {
		t.Thumbnail.Sizes = "(max-width: 767px) 89vw, (max-width: 1000px) 54vw, (max-width: 1071px) 543px, 580px"
		if h.Scene() == constraints.Detail {
			t.Thumbnail.Sizes = "100vw"
		}
	}
	return next(t)
}

var header = reload.Vars(models.PostThumbnail{})

func getHeaderImage(h *common.Handle) (r models.PostThumbnail) {
	img := header.Load()
	if img.Path != "" {
		return r
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
	r.Sizes = "100vw"
	header.Store(r)
	return
}

func calClass(h *common.Handle) {
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
