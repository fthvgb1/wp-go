package wpconfig

import (
	"embed"
	"fmt"
	"github.com/fthvgb1/wp-go/app/phphelper"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/safety"
	"path/filepath"
	"strings"
)

var templateFs embed.FS

func SetTemplateFs(fs embed.FS) {
	templateFs = fs
}

// ThemeMods 只有部分公共的参数，其它的参数调用 GetThemeModsVal 函数获取
type ThemeMods struct {
	CustomCssPostId       int            `json:"custom_css_post_id,omitempty"`
	NavMenuLocations      map[string]int `json:"nav_menu_locations,omitempty"`
	CustomLogo            int            `json:"custom_logo,omitempty"`
	HeaderImage           string         `json:"header_image,omitempty"`
	BackgroundImage       string         `json:"background_image,omitempty"`
	BackgroundSize        string         `json:"background_size,omitempty"`
	BackgroundRepeat      string         `json:"background_repeat,omitempty"`
	BackgroundColor       string         `json:"background_color,omitempty"`
	BackgroundPreset      string         `json:"background_preset"`
	BackgroundPositionX   string         `json:"background_position_x,omitempty"`
	BackgroundPositionY   string         `json:"background_position_y"`
	BackgroundAttachment  string         `json:"background_attachment"`
	ColorScheme           string         `json:"color_scheme"`
	SidebarTextcolor      string         `json:"sidebar_textcolor,omitempty"`
	HeaderBackgroundColor string         `json:"header_background_color,omitempty"`
	HeaderTextcolor       string         `json:"header_textcolor,omitempty"`
	HeaderVideo           int            `json:"header_video,omitempty"`
	ExternalHeaderVideo   string         `json:"external_header_video,omitempty"`
	HeaderImagData        ImageData      `json:"header_image_data,omitempty"`
	SidebarsWidgets       Sidebars       `json:"sidebars_widgets,omitempty"`
	ThemeSupport          ThemeSupport
}

type Sidebars struct {
	Time int          `json:"time,omitempty"`
	Data SidebarsData `json:"data,omitempty"`
}

type ColorScheme struct {
	Label  string   `json:"label,omitempty"`
	Colors []string `json:"colors,omitempty"`
}

type SidebarsData struct {
	WpInactiveWidgets []string `json:"wp_inactive_widgets,omitempty"`
	Sidebar1          []string `json:"sidebar-1,omitempty"`
	Sidebar2          []string `json:"sidebar-2,omitempty"`
	Sidebar3          []string `json:"sidebar-3,omitempty"`
}

type ImageData struct {
	AttachmentId int64  `json:"attachment_id,omitempty"`
	Url          string `json:"url,omitempty"`
	ThumbnailUrl string `json:"thumbnail_url,omitempty"`
	Height       int64  `json:"height,omitempty"`
	Width        int64  `json:"width,omitempty"`
}

func Thumbnail(metadata models.WpAttachmentMetadata, Type, host string, except ...string) (r models.PostThumbnail) {
	up := strings.Split(metadata.File, "/")
	if metadata.File != "" && Type == "full" {
		mimeType := metadata.Sizes["thumbnail"].MimeType
		metadata.Sizes["full"] = models.MetaDataFileSize{
			File:     filepath.Base(metadata.File),
			Width:    metadata.Width,
			Height:   metadata.Height,
			MimeType: mimeType,
			FileSize: metadata.FileSize,
		}
	}
	if siz, ok := metadata.Sizes[Type]; ok {
		r.Path = fmt.Sprintf("%s/wp-content/uploads/%s", host, strings.ReplaceAll(metadata.File, filepath.Base(metadata.File), siz.File))
		r.Width = metadata.Sizes[Type].Width
		r.Height = metadata.Sizes[Type].Height

		r.Srcset = strings.Join(maps.FilterToSlice[string](metadata.Sizes, func(s string, size models.MetaDataFileSize) (r string, ok bool) {
			up[len(up)-1] = size.File
			for _, s2 := range except {
				if s == s2 {
					return
				}
			}
			r = fmt.Sprintf("%s/wp-content/uploads/%s %dw", host, strings.Join(up, "/"), size.Width)
			ok = true
			return
		}), ", ")
		r.Sizes = fmt.Sprintf("(max-width: %dpx) 100vw, %dpx", r.Width, r.Width)
		if r.Width >= 740 && r.Width < 767 {
			r.Sizes = "(max-width: 706px) 89vw, (max-width: 767px) 82vw, 740px"
		} else if r.Width >= 767 {
			r.Sizes = "(max-width: 767px) 89vw, (max-width: 1000px) 54vw, (max-width: 1071px) 543px, 580px"
		}
		r.OriginAttachmentData = metadata
	}
	return
}

var themeModes = func() *safety.Map[string, ThemeMods] {
	m := safety.NewMap[string, ThemeMods]()
	themeModsRaw = safety.NewMap[string, map[string]any]()
	reload.Push(func() {
		m.Flush()
		themeModsRaw.Flush()
	})

	return m
}()

var themeModsRaw *safety.Map[string, map[string]any]

var themeSupport = map[string]ThemeSupport{}

func SetThemeSupport(theme string, support ThemeSupport) {
	themeSupport[theme] = support
}

func GetThemeModsVal[T any](theme, k string, defaults T) (r T) {
	m, ok := themeModsRaw.Load(theme)
	if !ok {
		r = defaults
		return
	}
	r = maps.GetStrAnyValWithDefaults(m, k, defaults)
	return
}

func GetThemeMods(theme string) (r ThemeMods, err error) {
	r, ok := themeModes.Load(theme)
	if ok {
		return
	}
	mods := GetOption(fmt.Sprintf("theme_mods_%s", theme))
	if mods == "" {
		return
	}
	m, err := phphelper.UnPHPSerializeToStrAnyMap(mods)
	if err != nil {
		return
	}
	themeModsRaw.Store(theme, m)
	//这里在的err可以不用处理，因为php的默认值和有设置过的类型可能不一样，直接按有设置的类型处理就行
	r, err = maps.StrAnyMapToStruct[ThemeMods](m)
	if err != nil {
		logs.Error(err, "解析thememods错误(可忽略)")
		err = nil
	}
	r.setThemeSupport(theme)
	themeModes.Store(theme, r)
	return
}

func IsCustomBackground(theme string) bool {
	mods, err := GetThemeMods(theme)
	if err != nil {
		return false
	}
	if mods.BackgroundColor != "" && mods.BackgroundColor != "default-color" || mods.BackgroundImage != "" && mods.BackgroundImage != "default-image" {
		return true
	}

	return false
}

func (m *ThemeMods) setThemeSupport(themeName string) {
	var v ThemeSupport
	vv, ok := themeSupport[themeName]
	if ok {
		m.ThemeSupport = vv
	} else {
		m.ThemeSupport = v
	}
}
