package twentyseventeen

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/gin-gonic/gin"
	"net/http"
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

type handle struct {
	*common.Handle
}

func newHandle(h *common.Handle) *handle {
	return &handle{Handle: h}
}

type indexHandle struct {
	*common.IndexHandle
	h *handle
}

func newIndexHandle(iHandle *common.IndexHandle) *indexHandle {
	return &indexHandle{IndexHandle: iHandle, h: newHandle(iHandle.Handle)}
}

type detailHandle struct {
	*common.DetailHandle
	h *handle
}

func newDetailHandle(dHandle *common.DetailHandle) *detailHandle {
	return &detailHandle{DetailHandle: dHandle, h: newHandle(dHandle.Handle)}
}

func Hook(h *common.Handle) {
	h.WidgetAreaData()
	h.GetPassword()
	h.GinH["HeaderImage"] = getHeaderImage(h.C)
	if h.Scene == constraints.Detail {
		newDetailHandle(common.NewDetailHandle(h)).Detail()
		return
	}
	newIndexHandle(common.NewIndexHandle(h)).Index()
}

var pluginFns = func() map[string]common.Plugin[models.Posts, *common.Handle] {
	return maps.Merge(common.ListPostPlugins(), map[string]common.Plugin[models.Posts, *common.Handle]{
		"twentyseventeen_postThumbnail": postThumbnail,
	})
}()

func (i *indexHandle) Index() {
	i.Templ = "twentyseventeen/posts/index.gohtml"
	err := i.BuildIndexData(common.NewIndexParams(i.C))
	if err != nil {
		i.Stats = constraints.Error404
		i.Code = http.StatusNotFound
		i.GinH["bodyClass"] = i.h.bodyClass()
		i.C.HTML(i.Code, i.Templ, i.GinH)
		return
	}
	i.PostsPlugins = pluginFns
	i.PageEle = paginate
	i.Render()
}

func (d *detailHandle) Detail() {
	err := d.BuildDetailData()
	if err != nil {
		d.Code = http.StatusNotFound
		d.Stats = constraints.Error404
		d.GinH["bodyClass"] = d.h.bodyClass()
		d.C.HTML(d.Code, d.Templ, d.GinH)
		return
	}
	d.GinH["bodyClass"] = d.h.bodyClass()
	img := wpconfig.Thumbnail(d.Post.Thumbnail.OriginAttachmentData, "thumbnail", "", "thumbnail", "post-thumbnail")
	img.Width = img.OriginAttachmentData.Width
	img.Height = img.OriginAttachmentData.Height
	img.Sizes = "100vw"
	img.Srcset = fmt.Sprintf("%s %dw, %s", img.Path, img.Width, img.Srcset)
	d.Post.Thumbnail = img
	d.CommentRender = commentFormat
	d.GinH["post"] = d.Post
	d.Render()

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
		if h.Scene == constraints.Detail {
			t.Thumbnail.Sizes = "100vw"
		}
	}
	return next(t)
}

func getHeaderImage(c *gin.Context) (r models.PostThumbnail) {
	r.Path = "/wp-content/themes/twentyseventeen/assets/images/header.jpg"
	r.Width = 2000
	r.Height = 1200
	hs, err := cache.GetHeaderImages(c, ThemeName)
	if err != nil {
		logs.ErrPrintln(err, "获取页眉背景图失败")
	} else if len(hs) > 0 && err == nil {
		_, r = slice.Rand(hs)

	}
	r.Sizes = "100vw"
	return
}

func (i *handle) bodyClass() string {
	s := ""
	if constraints.Ok != i.Stats {
		return "error404"
	}
	switch i.Scene {
	case constraints.Search:
		s = "search-no-results"
		if len(i.GinH["posts"].([]models.Posts)) > 0 {
			s = "search-results"
		}
	case constraints.Category, constraints.Tag:
		cat := i.C.Param("category")
		if cat == "" {
			cat = i.C.Param("tag")
		}
		_, cate := slice.SearchFirst(cache.CategoriesTags(i.C, i.Scene), func(my models.TermsMy) bool {
			return my.Name == cat
		})
		if cate.Slug[0] != '%' {
			s = cate.Slug
		}
		s = fmt.Sprintf("category-%d %v", cate.Terms.TermId, s)
	case constraints.Detail:
		s = fmt.Sprintf("postid-%d", i.GinH["post"].(models.Posts).Id)
	}
	return str.Join(class[i.Scene], s)
}

var class = map[int]string{
	constraints.Home:     "home blog ",
	constraints.Archive:  "archive date page-two-column",
	constraints.Category: "archive category page-two-column",
	constraints.Tag:      "archive category page-two-column ",
	constraints.Search:   "search ",
	constraints.Detail:   "post-template-default single single-post single-format-standard ",
}
