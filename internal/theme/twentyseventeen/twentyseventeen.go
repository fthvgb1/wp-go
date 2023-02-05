package twentyseventeen

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
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
	c      *gin.Context
	ginH   gin.H
	scene  int
	status int
	stats  int
	templ  string
}

func Hook(status int, c *gin.Context, ginH gin.H, scene, stats int) {
	h := handle{
		c:      c,
		ginH:   ginH,
		scene:  scene,
		status: status,
		stats:  stats,
		templ:  "twentyseventeen/posts/index.gohtml",
	}
	ginH["HeaderImage"] = h.getHeaderImage(c)
	if stats == plugins.Empty404 {
		c.HTML(status, h.templ, ginH)
		return
	}
	if scene == plugins.Detail {
		h.detail()
		return
	}
	h.index()

}

func (h handle) index() {
	posts := h.ginH["posts"].([]models.Posts)
	p, ok := h.ginH["pagination"]
	if ok {
		pp, ok := p.(pagination.ParsePagination)
		if ok {
			h.ginH["pagination"] = pagination.Paginate(paginate, pp)
		}
	}
	h.ginH["bodyClass"] = h.bodyClass()
	h.ginH["posts"] = h.postThumbnail(posts, h.scene)
	h.c.HTML(h.status, h.templ, h.ginH)
}

func (h handle) detail() {
	post := h.ginH["post"].(models.Posts)
	h.ginH["bodyClass"] = h.bodyClass()
	//host, _ := wpconfig.Options.Load("siteurl")
	host := ""
	img := plugins.Thumbnail(post.Thumbnail.OriginAttachmentData, "thumbnail", host, "thumbnail", "post-thumbnail")
	img.Width = img.OriginAttachmentData.Width
	img.Height = img.OriginAttachmentData.Height
	img.Sizes = "100vw"
	img.Srcset = fmt.Sprintf("%s %dw, %s", img.Path, img.Width, img.Srcset)
	post.Thumbnail = img
	h.ginH["post"] = post
	comments := h.ginH["comments"].([]models.Comments)
	dep := h.ginH["maxDep"].(int)
	h.ginH["comments"] = plugins.FormatComments(h.c, comment{}, comments, dep)
	h.templ = "twentyseventeen/posts/detail.gohtml"
	h.c.HTML(h.status, h.templ, h.ginH)
}

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

func (h handle) postThumbnail(posts []models.Posts, scene int) []models.Posts {
	return slice.Map(posts, func(t models.Posts) models.Posts {
		if t.Thumbnail.Path != "" {
			if slice.IsContained(scene, []int{plugins.Home, plugins.Archive, plugins.Search}) {
				t.Thumbnail.Sizes = "(max-width: 767px) 89vw, (max-width: 1000px) 54vw, (max-width: 1071px) 543px, 580px"
			} else {
				t.Thumbnail.Sizes = "100vw"
			}
		}
		return t
	})
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
	switch h.scene {
	case plugins.Search:
		if len(h.ginH["posts"].([]models.Posts)) > 0 {
			s = "search-results"
		} else {
			s = "search-no-results"
		}
	case plugins.Category:
		cat := h.c.Param("category")
		_, cate := slice.SearchFirst(cache.Categories(h.c), func(my models.TermsMy) bool {
			return my.Name == cat
		})
		if cate.Slug[0] != '%' {
			s = cate.Slug
		}
		s = fmt.Sprintf("category-%d %v", cate.Terms.TermId, s)
	case plugins.Tag:
		cat := h.c.Param("tag")
		_, cate := slice.SearchFirst(cache.Tags(h.c), func(my models.TermsMy) bool {
			return my.Name == cat
		})
		if cate.Slug[0] != '%' {
			s = cate.Slug
		}
		s = fmt.Sprintf("category-%d %v", cate.Terms.TermId, s)
	case plugins.Detail:
		s = fmt.Sprintf("postid-%d", h.ginH["post"].(models.Posts).Id)
	}
	return str.Join(class[h.scene], s)
}

var class = map[int]string{
	plugins.Home:     "home blog ",
	plugins.Archive:  "archive date page-two-column",
	plugins.Category: "archive category page-two-column",
	plugins.Tag:      "archive category page-two-column ",
	plugins.Search:   "search ",
	plugins.Detail:   "post-template-default single single-post single-format-standard ",
}
