package twentyseventeen

import (
	"github.com/elliotchance/phpserialize"
	"github.com/fthvgb1/wp-go/helper"
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

func Hook(status int, c *gin.Context, h gin.H, scene int) {
	templ := "twentyseventeen/posts/index.gohtml"
	if _, ok := plugins.IndexSceneMap[scene]; ok {
		h["HeaderImage"] = getHeaderImage(c)
		p, ok := h["pagination"]
		if ok {
			pp, ok := p.(pagination.ParsePagination)
			if ok {
				h["pagination"] = pagination.Paginate(paginate, pp)
			}
		}
	} else if _, ok := plugins.DetailSceneMap[scene]; ok {
		templ = "twentyseventeen/posts/detail.gohtml"
	}
	c.HTML(status, templ, h)
	return
}

func getHeaderImage(c *gin.Context) (r models.PostThumbnail) {
	r.Path = "/wp-content/themes/twentyseventeen/assets/images/header.jpg"
	r.Width = 2000
	r.Height = 1200
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
		if m.Thumbnail.Path != "" {
			r = m.Thumbnail
		}
	}
	return
}

func getHeaderMarkup() (r HeaderImageMeta, err error) {
	mods, ok := wpconfig.Options.Load("theme_mods_twentyseventeen")
	var rr map[any]any
	if ok {
		err = phpserialize.Unmarshal([]byte(mods), &rr)
		if err == nil {
			rx := helper.MapAnyAnyToStrAny(rr)
			r, err = helper.StrAnyMapToStruct[HeaderImageMeta](rx)
		}
	}
	return
}
