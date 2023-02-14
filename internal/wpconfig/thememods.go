package wpconfig

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/internal/phphelper"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/safety"
	"strings"
)

type ThemeMods struct {
	CustomCssPostId       int       `json:"custom_css_post_id,omitempty"`
	NavMenuLocations      []string  `json:"nav_menu_locations,omitempty"`
	CustomLogo            int       `json:"custom_logo"`
	HeaderImage           string    `json:"header_image,omitempty"`
	BackgroundImage       string    `json:"background_image"`
	BackgroundSize        string    `json:"background_size"`
	BackgroundRepeat      string    `json:"background_repeat"`
	BackgroundColor       string    `json:"background_color"`
	ColorScheme           string    `json:"color_scheme"`
	SidebarTextcolor      string    `json:"sidebar_textcolor"`
	HeaderBackgroundColor string    `json:"header_background_color"`
	HeaderTextcolor       string    `json:"header_textcolor"`
	HeaderImagData        ImageData `json:"header_image_data,omitempty"`
	SidebarsWidgets       Sidebars  `json:"sidebars_widgets"`
}

type Sidebars struct {
	Time int          `json:"time,omitempty"`
	Data SidebarsData `json:"data"`
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
	if _, ok := metadata.Sizes[Type]; ok {
		r.Path = fmt.Sprintf("%s/wp-content/uploads/%s", host, metadata.File)
		r.Width = metadata.Sizes[Type].Width
		r.Height = metadata.Sizes[Type].Height
		up := strings.Split(metadata.File, "/")
		r.Srcset = strings.Join(maps.FilterToSlice[string](metadata.Sizes, func(s string, size models.MetaDataFileSize) (r string, ok bool) {
			up[2] = size.File
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

var themeModes = safety.Map[string, ThemeMods]{}

func FlushModes() {
	themeModes.Flush()
}

func GetThemeMods(theme string) (r ThemeMods, err error) {
	r, ok := themeModes.Load(theme)
	if ok {
		return
	}

	mods, ok := Options.Load(fmt.Sprintf("theme_mods_%s", theme))
	if !ok || mods == "" {
		return
	}
	r, err = phphelper.UnPHPSerialize[ThemeMods](mods)
	if err != nil {
		return
	}
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
func IsCustomLogo(theme string) bool {
	mods, err := GetThemeMods(theme)
	if err != nil {
		return false
	}
	if mods.CustomLogo > 0 {
		return true
	}

	return false
}
