package twentyseventeen

import (
	"fmt"
	"github.com/elliotchance/phpserialize"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/logs"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/plugins"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/plugin/pagination"
	"github.com/gin-gonic/gin"
	"strings"
)

const ThemeName = "twentyseventeen"

type HeaderImageMeta struct {
	CustomCssPostId  int       `json:"custom_css_post_id,omitempty"`
	NavMenuLocations []string  `json:"nav_menu_locations,omitempty"`
	HeaderImage      string    `json:"header_image,omitempty"`
	HeaderImagData   ImageData `json:"header_image_data,omitempty"`
}

type ImageData struct {
	AttachmentId int64  `json:"attachment_id,omitempty"`
	Url          string `json:"url,omitempty"`
	ThumbnailUrl string `json:"thumbnail_url,omitempty"`
	Height       int64  `json:"height,omitempty"`
	Width        int64  `json:"width,omitempty"`
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

func Hook(status int, c *gin.Context, h gin.H, scene, stats int) {
	templ := "twentyseventeen/posts/index.gohtml"
	if _, ok := plugins.IndexSceneMap[scene]; ok {
		h["HeaderImage"] = getHeaderImage(c)
		posts := h["posts"].([]models.Posts)
		p, ok := h["pagination"]
		if ok {
			pp, ok := p.(pagination.ParsePagination)
			if ok {
				h["pagination"] = pagination.Paginate(paginate, pp)
			}
		}
		d := 0
		s := ""
		if scene == plugins.Search {
			if len(posts) > 0 {
				d = 1
			} else {
				d = 0
			}
		} else if scene == plugins.Category {
			cate := slice.Filter(cache.Categories(c), func(my models.TermsMy) bool {
				return my.Name == c.Param("category")
			})[0]
			d = int(cate.Terms.TermId)
			if cate.Slug[0] != '%' {
				s = cate.Slug
			}
		}
		h["bodyClass"] = bodyClass(scene, d, s)
		h["posts"] = postThumbnail(posts, scene)
	} else if scene == plugins.Detail {
		h["HeaderImage"] = getHeaderImage(c)
		post := h["post"].(models.Posts)
		h["bodyClass"] = bodyClass(scene, int(post.Id))
		host, _ := wpconfig.Options.Load("siteurl")
		img := plugins.Thumbnail(post.Thumbnail.OriginAttachmentData, "thumbnail", host, "thumbnail", "post-thumbnail")
		img.Width = img.OriginAttachmentData.Width
		img.Height = img.OriginAttachmentData.Height
		img.Sizes = "100vw"
		img.Srcset = fmt.Sprintf("%s %dw, %s", img.Path, img.Width, img.Srcset)
		post.Thumbnail = img
		h["post"] = post
		comments := h["comments"].([]models.Comments)
		dep := h["maxDep"].(int)
		h["comments"] = plugins.FormatComments(c, comment{}, comments, dep)
		templ = "twentyseventeen/posts/detail.gohtml"
	}
	c.HTML(status, templ, h)
	return
}

type comment struct {
	plugins.CommonCommentFormat
}

func (c comment) FormatLi(ctx *gin.Context, m models.Comments, depth int, eo, parent string) string {
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
	return plugins.FormatLi(templ, ctx, m, depth, eo, parent)
}

func postThumbnail(posts []models.Posts, scene int) []models.Posts {
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

func getHeaderImage(c *gin.Context) (r models.PostThumbnail) {
	r.Path = "/wp-content/themes/twentyseventeen/assets/images/header.jpg"
	meta, err := getHeaderMarkup()
	if err != nil {
		logs.ErrPrintln(err, "解析主题背景图设置错误")
		return
	}
	if meta.HeaderImagData.AttachmentId > 0 {
		m, err := cache.GetPostById(c, uint64(meta.HeaderImagData.AttachmentId))
		if err != nil {
			logs.ErrPrintln(err, "获取主题背景图信息错误")
			return
		}
		host, _ := wpconfig.Options.Load("siteurl")
		m.Thumbnail = plugins.Thumbnail(m.AttachmentMetadata, "thumbnail", host, "thumbnail", "post-thumbnail", "twentyseventeen-thumbnail-avatar")
		if m.Thumbnail.Path != "" {
			r = m.Thumbnail
			if len(m.AttachmentMetadata.Sizes) > 0 {
				r.Srcset = str.Join(r.Path, " 2000vw, ", r.Srcset)
			}
		}
	}
	r.Width = 2000
	r.Height = 1200
	r.Sizes = "100vw"
	return
}

func getHeaderMarkup() (r HeaderImageMeta, err error) {
	mods, ok := wpconfig.Options.Load("theme_mods_twentyseventeen")
	var rr map[any]any
	if ok {
		err = phpserialize.Unmarshal([]byte(mods), &rr)
		if err == nil {
			rx := maps.AnyAnyToStrAny(rr)
			r, err = maps.StrAnyMapToStruct[HeaderImageMeta](rx)
		}
	}
	return
}

func bodyClass(scene, d int, a ...any) string {
	s := ""
	if scene == plugins.Search {
		if d > 0 {
			s = "search-results"
		} else {
			s = "search-no-results"
		}
	} else if scene == plugins.Category {
		s = fmt.Sprintf("category-%d %v", d, a[0])
	} else if scene == plugins.Detail {
		s = fmt.Sprintf("postid-%d", d)
	}
	return map[int]string{
		plugins.Home:     "home blog ",
		plugins.Archive:  "archive date page-two-column",
		plugins.Category: str.Join("archive category page-two-column ", s),
		plugins.Search:   str.Join("search ", s),
		plugins.Detail:   str.Join("post-template-default single single-post single-format-standard ", s),
	}[scene]
}
