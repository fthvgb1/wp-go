package twentyseventeen

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/fthvgb1/wp-go/plugin/pagination"
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

type handle struct {
	common.Handle
	templ string
}

func Hook(cHandle common.Handle) {
	h := handle{
		Handle: cHandle,
		templ:  "twentyseventeen/posts/index.gohtml",
	}
	h.GinH["HeaderImage"] = h.getHeaderImage(h.C)
	if h.Scene == plugins.Detail {
		h.Detail()
		return
	}
	h.Index()
}

var plugin = func() []common.Plugin[models.Posts] {
	return append(common.Plugins(), postThumbnail)
}()

func (h handle) Index() {
	if h.Stats != plugins.Empty404 {

		h.GinH["posts"] = slice.Map(
			h.GinH["posts"].([]models.Posts),
			common.PluginFn[models.Posts](plugin, h.Handle, common.DigestsAndOthers(h.C)))

		p, ok := h.GinH["pagination"]
		if ok {
			pp, ok := p.(pagination.ParsePagination)
			if ok {
				h.GinH["pagination"] = pagination.Paginate(paginate, pp)
			}
		}
	}

	h.GinH["bodyClass"] = h.bodyClass()
	h.C.HTML(h.Code, h.templ, h.GinH)
}

func (h handle) Detail() {
	post := h.GinH["post"].(models.Posts)
	h.GinH["bodyClass"] = h.bodyClass()
	if h.Stats == plugins.Empty404 {
		h.C.HTML(h.Code, h.templ, h.GinH)
		return
	}
	//host, _ := wpconfig.Options.Load("siteurl")
	host := ""
	img := plugins.Thumbnail(post.Thumbnail.OriginAttachmentData, "thumbnail", host, "thumbnail", "post-thumbnail")
	img.Width = img.OriginAttachmentData.Width
	img.Height = img.OriginAttachmentData.Height
	img.Sizes = "100vw"
	img.Srcset = fmt.Sprintf("%s %dw, %s", img.Path, img.Width, img.Srcset)
	post.Thumbnail = img
	h.GinH["post"] = post
	if h.GinH["comments"] != nil {
		comments := h.GinH["comments"].([]models.Comments)
		dep := h.GinH["maxDep"].(int)
		h.GinH["comments"] = plugins.FormatComments(h.C, commentFormat, comments, dep)
	}
	h.templ = "twentyseventeen/posts/detail.gohtml"
	h.C.HTML(h.Code, h.templ, h.GinH)
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

func postThumbnail(next common.Fn[models.Posts], h common.Handle, t models.Posts) models.Posts {
	if t.Thumbnail.Path != "" {
		t.Thumbnail.Sizes = "(max-width: 767px) 89vw, (max-width: 1000px) 54vw, (max-width: 1071px) 543px, 580px"
		if h.Scene == plugins.Detail {
			t.Thumbnail.Sizes = "100vw"
		}
	}
	return next(t)
}

func (h handle) getHeaderImage(c *gin.Context) (r models.PostThumbnail) {
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

func (h handle) bodyClass() string {
	s := ""
	if h.Stats == plugins.Empty404 {
		return "error404"
	}
	switch h.Scene {
	case plugins.Search:
		s = "search-no-results"
		if len(h.GinH["posts"].([]models.Posts)) > 0 {
			s = "search-results"
		}
	case plugins.Category, plugins.Tag:
		cat := h.C.Param("category")
		if cat == "" {
			cat = h.C.Param("tag")
		}
		_, cate := slice.SearchFirst(cache.CategoriesTags(h.C, h.Scene), func(my models.TermsMy) bool {
			return my.Name == cat
		})
		if cate.Slug[0] != '%' {
			s = cate.Slug
		}
		s = fmt.Sprintf("category-%d %v", cate.Terms.TermId, s)
	case plugins.Detail:
		s = fmt.Sprintf("postid-%d", h.GinH["post"].(models.Posts).Id)
	}
	return str.Join(class[h.Scene], s)
}

var class = map[int]string{
	plugins.Home:     "home blog ",
	plugins.Archive:  "archive date page-two-column",
	plugins.Category: "archive category page-two-column",
	plugins.Tag:      "archive category page-two-column ",
	plugins.Search:   "search ",
	plugins.Detail:   "post-template-default single single-post single-format-standard ",
}
