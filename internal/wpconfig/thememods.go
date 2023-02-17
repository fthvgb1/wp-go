package wpconfig

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/phphelper"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/safety"
	"path/filepath"
	"strings"
)

var templateFs embed.FS

func SetTemplateFs(fs embed.FS) {
	templateFs = fs
}

type ThemeMod struct {
	CustomCssPostId       int      `json:"custom_css_post_id,omitempty"`
	NavMenuLocations      []string `json:"nav_menu_locations,omitempty"`
	CustomLogo            int      `json:"custom_logo,omitempty"`
	HeaderImage           string   `json:"header_image,omitempty"`
	BackgroundImage       string   `json:"background_image,omitempty"`
	BackgroundSize        string   `json:"background_size,omitempty"`
	BackgroundRepeat      string   `json:"background_repeat,omitempty"`
	BackgroundColor       string   `json:"background_color,omitempty"`
	BackgroundPreset      string   `json:"background_preset"`
	BackgroundPositionX   string   `json:"background_position_x,omitempty"`
	BackgroundPositionY   string   `json:"background_position_y"`
	BackgroundAttachment  string   `json:"background_attachment"`
	ColorScheme           map[string]ColorScheme
	SidebarTextcolor      string    `json:"sidebar_textcolor,omitempty"`
	HeaderBackgroundColor string    `json:"header_background_color,omitempty"`
	HeaderTextcolor       string    `json:"header_textcolor,omitempty"`
	HeaderImagData        ImageData `json:"header_image_data,omitempty"`
	SidebarsWidgets       Sidebars  `json:"sidebars_widgets,omitempty"`
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
	if Type == "full" {
		metadata.Sizes["full"] = models.MetaDataFileSize{
			File:     filepath.Base(metadata.File),
			Width:    metadata.Width,
			Height:   metadata.Height,
			MimeType: metadata.Sizes["thumbnail"].MimeType,
			FileSize: metadata.FileSize,
		}
	}
	if _, ok := metadata.Sizes[Type]; ok {
		r.Path = fmt.Sprintf("%s/wp-content/uploads/%s", host, metadata.File)
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

var themeModes = func() *safety.Map[string, ThemeMod] {
	m := safety.NewMap[string, ThemeMod]()
	reload.Push(func() {
		m.Flush()
	})
	return m
}()

func GetThemeMods(theme string) (r ThemeMod, err error) {
	r, ok := themeModes.Load(theme)
	if ok {
		return
	}

	mods, ok := Options.Load(fmt.Sprintf("theme_mods_%s", theme))
	if !ok || mods == "" {
		return
	}
	r, err = phphelper.UnPHPSerialize[ThemeMod](mods)
	if err != nil {
		return
	}
	r.setThemeColorScheme(theme)
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

func (m *ThemeMod) setThemeColorScheme(themeName string) {
	bytes, err := templateFs.ReadFile(filepath.Join(themeName, "colorscheme.json"))
	if err != nil {
		return
	}
	var scheme map[string]ColorScheme
	err = json.Unmarshal(bytes, &scheme)
	if err != nil {
		return
	}
	m.ColorScheme = scheme
}
func (m *ThemeMod) setThemeSupport(themeName string) {
	bytes, err := templateFs.ReadFile(filepath.Join(themeName, "themesupport.json"))
	if err != nil {
		return
	}
	var themeSupport ThemeSupport
	err = json.Unmarshal(bytes, &themeSupport)
	if err != nil {
		return
	}
	m.ThemeSupport = themeSupport
}
