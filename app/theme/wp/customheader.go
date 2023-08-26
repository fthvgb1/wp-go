package wp

import (
	"encoding/json"
	"fmt"
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/model"
	"regexp"
)

func (h *Handle) DisplayHeaderText() bool {
	return h.themeMods.ThemeSupport.CustomHeader.HeaderText && "blank" != h.themeMods.HeaderTextcolor
}

func (h *Handle) GetCustomHeaderImg() (r models.PostThumbnail, isRand bool) {
	var err error
	img := reload.GetAnyValBys("headerImages", h.theme, func(theme string) []models.PostThumbnail {
		hs, er := h.GetHeaderImages(h.theme)
		if er != nil {
			err = er
			return nil
		}
		return hs
	})
	if err != nil {
		logs.Error(err, "获取页眉背景图失败")
		return
	}
	hs := slice.Copy(img)

	if len(hs) < 1 {
		return
	}
	if len(hs) > 1 {
		isRand = true
	}
	r, _ = slice.RandPop(&hs)
	return
}

type VideoPlay struct {
	Pause      string `json:"pause,omitempty"`
	Play       string `json:"play,omitempty"`
	PauseSpeak string `json:"pauseSpeak,omitempty"`
	PlaySpeak  string `json:"playSpeak,omitempty"`
}

type VideoSetting struct {
	MimeType  string    `json:"mimeType,omitempty"`
	PosterUrl string    `json:"posterUrl,omitempty"`
	VideoUrl  string    `json:"videoUrl,omitempty"`
	Width     int       `json:"width,omitempty"`
	Height    int       `json:"height,omitempty"`
	MinWidth  int       `json:"minWidth,omitempty"`
	MinHeight int       `json:"minHeight,omitempty"`
	L10n      VideoPlay `json:"l10n"`
}

var videoReg = regexp.MustCompile(`^https?://(?:www\.)?(?:youtube\.com/watch|youtu\.be/)`)

func GetVideoSetting(h *Handle, u string) (string, error) {

	img, _ := h.GetCustomHeaderImg()
	v := VideoSetting{
		MimeType:  GetMimeType(u),
		PosterUrl: img.Path,
		VideoUrl:  u,
		Width:     img.Width,
		Height:    img.Height,
		MinWidth:  900,
		MinHeight: 500,
		L10n: VideoPlay{
			Pause:      "暂停",
			Play:       "播放",
			PauseSpeak: "视频已暂停",
			PlaySpeak:  "视频正在播放",
		},
	}
	if is := videoReg.FindString(u); is != "" {
		v.MimeType = "video/x-youtube"
	}
	_ = h.DoActionFilter("videoSetting", "", &v)
	s, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	setting := fmt.Sprintf(`var %s = %s`, "_wpCustomHeaderSettings", string(s))
	script := str.Join(`<script id="wp-custom-header-js-extra">`, setting, "</script>\n")
	return script, nil
}

func CustomVideo(h *Handle, scene ...string) (ok bool) {
	mod, err := wpconfig.GetThemeMods(h.theme)
	if err != nil {
		logs.Error(err, "getThemeMods fail", h.theme)
		return
	}
	if !mod.ThemeSupport.CustomHeader.Video || (mod.HeaderVideo < 1 && mod.ExternalHeaderVideo == "") {
		return
	}
	u := ""
	if mod.HeaderVideo > 0 {
		post, err := cache.GetPostById(h.C, uint64(mod.HeaderVideo))
		if err != nil {
			logs.Error(err, "get headerVideo fail", mod.HeaderVideo)
			return
		}
		u = post.Metas["_wp_attached_file"].(string)
		u = str.Join("/wp-content/uploads/", u)
	} else {
		u = mod.ExternalHeaderVideo
	}

	hs, err := GetVideoSetting(h, u)
	if err != nil {
		logs.Error(err, "get headerVideo fail", mod.HeaderVideo)
		return
	}
	scriptss := []string{
		"/wp-includes/js/dist/vendor/wp-polyfill-inert.min.js",
		"/wp-includes/js/dist/vendor/regenerator-runtime.min.js",
		"/wp-includes/js/dist/vendor/wp-polyfill.min.js",
		"/wp-includes/js/dist/dom-ready.min.js",
		"/wp-includes/js/dist/hooks.min.js",
		"/wp-includes/js/dist/i18n.min.js",
		"/wp-includes/js/dist/a11y.min.js",
		"/wp-includes/js/wp-custom-header.min.js",
	}
	scriptss = slice.Map(scriptss, func(t string) string {
		return fmt.Sprintf(`<script src="%s" id="wp-%s-js"></script>
`, t, str.Replaces(t, []string{
			"/wp-includes/js/dist/vendor/", "/wp-includes/js/dist/", "/wp-includes/js/", ".min.js", ".js", "wp-", "",
		}))
	})

	var tr = `<script id="wp-i18n-js-after">
wp.i18n.setLocaleData( { 'text direction\u0004ltr': [ 'ltr' ] } );
</script>
<script id='wp-a11y-js-translations'>
( function( domain, translations ) {
	var localeData = translations.locale_data[ domain ] || translations.locale_data.messages;
	localeData[""].domain = domain;
	wp.i18n.setLocaleData( localeData, domain );
} )( "default", {"translation-revision-date":"2023-04-23 22:48:55+0000","generator":"GlotPress/4.0.0-alpha.4","domain":"messages","locale_data":{"messages":{"":{"domain":"messages","plural-forms":"nplurals=1; plural=0;","lang":"zh_CN"},"Notifications":["u901au77e5"]}},"comment":{"reference":"wp-includes/js/dist/a11y.js"}} );
</script>
<script src='/wp-includes/js/dist/a11y.min.js?ver=ecce20f002eda4c19664' id='wp-a11y-js'></script>
`
	c := []Components[string]{
		NewComponent("wp-a11y-js-translations", tr, true, 10.0065, nil),
		NewComponent("VideoSetting", hs, true, 10.0064, nil),
		NewComponent("header-script", scriptss[len(scriptss)-1], true, 10.0063, nil),
	}
	for _, s := range scene {
		h.PushGroupFooterScript(s, "wp-custom-header", 10.0066, scriptss[0:len(scriptss)-2]...)
		h.PushFooterScript(s, c...)
	}
	ok = true
	return
}

func (h *Handle) GetHeaderImages(theme string) (r []models.PostThumbnail, err error) {
	meta, err := wpconfig.GetThemeMods(theme)
	if err != nil || meta.HeaderImage == "" {
		return
	}
	if "random-uploaded-image" != meta.HeaderImage {
		m, er := cache.GetPostById(h.C, uint64(meta.HeaderImagData.AttachmentId))
		if er != nil {
			err = er
			return
		}
		m.Thumbnail = thumb(m, theme)
		r = []models.PostThumbnail{m.Thumbnail}
		return
	}

	headers, er := model.Finds[models.Posts](h.C, model.Conditions(
		model.Where(model.SqlBuilder{
			{"post_type", "attachment"},
			{"post_status", "inherit"},
			{"meta_value", theme},
			{"meta_key", "_wp_attachment_is_custom_header"},
		}),
		model.Fields("a.ID"),
		model.Group("a.ID"),
		model.Join(model.SqlBuilder{
			{" a", "left join", "wp_postmeta b", "a.ID=b.post_id"},
		}),
	))

	if er != nil {
		err = er
		return
	}
	if len(headers) > 0 {
		posts, er := cache.GetPostsByIds(h.C, slice.Map(headers, func(t models.Posts) uint64 {
			return t.Id
		}))
		if er != nil {
			err = er
			return
		}
		r = slice.Map(posts, func(m models.Posts) models.PostThumbnail {
			return thumb(m, theme)
		})
	}
	return

}

func thumb(m models.Posts, theme string) models.PostThumbnail {
	m.Thumbnail = wpconfig.Thumbnail(m.AttachmentMetadata, "full", "", "thumbnail", "post-thumbnail", fmt.Sprintf("%s-thumbnail-avatar", theme))
	m.Thumbnail.Width = m.AttachmentMetadata.Width
	m.Thumbnail.Height = m.AttachmentMetadata.Height
	if m.Thumbnail.Path != "" {
		if len(m.AttachmentMetadata.Sizes) > 0 {
			m.Thumbnail.Srcset = str.Join(m.Thumbnail.Path, " 2000w, ", m.Thumbnail.Srcset)
		}
	}
	return m.Thumbnail
}
