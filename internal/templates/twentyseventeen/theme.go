package twentyseventeen

import (
	"github.com/elliotchance/phpserialize"
	"github.com/gin-gonic/gin"
	"github/fthvgb1/wp-go/helper"
	"github/fthvgb1/wp-go/internal/pkg/cache"
	"github/fthvgb1/wp-go/internal/pkg/logs"
	"github/fthvgb1/wp-go/internal/pkg/models"
	"github/fthvgb1/wp-go/internal/plugins"
	"github/fthvgb1/wp-go/internal/wpconfig"
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

func Hook(c *gin.Context, h gin.H, scene int) (r string) {
	if _, ok := plugins.IndexSceneMap[scene]; ok {
		r = "twentyseventeen/posts/index.gohtml"
		h["HeaderImage"] = getHeaderImage(c)
	} else if _, ok := plugins.DetailSceneMap[scene]; ok {
		r = "twentyseventeen/posts/detail.gohtml"
	}
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
