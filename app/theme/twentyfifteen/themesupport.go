package twentyfifteen

import "github.com/fthvgb1/wp-go/app/wpconfig"

type themeSupport struct {
	CustomBackground      customBackground        `json:"custom-background"`
	EditorColorPalette    []EditorColorPalette    `json:"editor-color-palette"`
	EditorGradientPresets []EditorGradientPresets `json:"editor-gradient-presets"`
}
type customBackground struct {
	DefaultImage         string `json:"default-image"`
	DefaultPreset        string `json:"default-preset"`
	DefaultPositionX     string `json:"default-position-x"`
	DefaultPositionY     string `json:"default-position-y"`
	DefaultSize          string `json:"default-size"`
	DefaultRepeat        string `json:"default-repeat"`
	DefaultAttachment    string `json:"default-attachment"`
	DefaultColor         string `json:"default-color"`
	WpHeadCallback       string `json:"wp-head-callback"`
	AdminHeadCallback    string `json:"admin-head-callback"`
	AdminPreviewCallback string `json:"admin-preview-callback"`
}

type EditorColorPalette struct {
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Color string `json:"color"`
}
type EditorGradientPresets struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	Gradient string `json:"gradient"`
}

var themesupport = themeSupport{
	CustomBackground: customBackground{
		DefaultImage:         "",
		DefaultPreset:        "default",
		DefaultPositionX:     "left",
		DefaultPositionY:     "top",
		DefaultSize:          "auto",
		DefaultRepeat:        "repeat",
		DefaultAttachment:    "fixed",
		DefaultColor:         "f1f1f1",
		WpHeadCallback:       "_custom_background_cb",
		AdminHeadCallback:    "",
		AdminPreviewCallback: "",
	},
	EditorColorPalette: []EditorColorPalette{
		{
			Name:  "暗灰色",
			Slug:  "dark-gray",
			Color: "#111",
		},
		{
			Name:  "亮灰色",
			Slug:  "light-gray",
			Color: "#f1f1f1",
		},
		{
			Name:  "白色",
			Slug:  "white",
			Color: "#fff",
		},
		{
			Name:  "黄色",
			Slug:  "yellow",
			Color: "#f4ca16",
		},
		{
			Name:  "暗棕色",
			Slug:  "dark-brown",
			Color: "#352712",
		},
		{
			Name:  "粉色",
			Slug:  "medium-pink",
			Color: "#e53b51",
		},
		{
			Name:  "浅粉色",
			Slug:  "light-pink",
			Color: "#ffe5d1",
		},
		{
			Name:  "暗紫色",
			Slug:  "dark-purple",
			Color: "#2e2256",
		},
		{
			Name:  "紫色",
			Slug:  "purple",
			Color: "#674970",
		},
		{
			Name:  "蓝灰色",
			Slug:  "blue-gray",
			Color: "#22313f",
		},
		{
			Name:  "亮蓝色",
			Slug:  "bright-blue",
			Color: "#55c3dc",
		},
		{
			Name:  "浅蓝色",
			Slug:  "light-blue",
			Color: "#e9f2f9",
		},
	},
	EditorGradientPresets: []EditorGradientPresets{
		{
			Name:     "Dark Gray Gradient",
			Slug:     "dark-gray-gradient-gradient",
			Gradient: "linear-gradient(90deg, rgba(17,17,17,1) 0%, rgba(42,42,42,1) 100%)",
		},
		{
			Name:     "Light Gray Gradient",
			Slug:     "light-gray-gradient",
			Gradient: "linear-gradient(90deg, rgba(241,241,241,1) 0%, rgba(215,215,215,1) 100%)",
		},
		{
			Name:     "White Gradient",
			Slug:     "white-gradient",
			Gradient: "linear-gradient(90deg, rgba(255,255,255,1) 0%, rgba(230,230,230,1) 100%)",
		},
		{
			Name:     "Yellow Gradient",
			Slug:     "yellow-gradient",
			Gradient: "linear-gradient(90deg, rgba(244,202,22,1) 0%, rgba(205,168,10,1) 100%)",
		},
		{
			Name:     "Dark Brown Gradient",
			Slug:     "dark-brown-gradient",
			Gradient: "linear-gradient(90deg, rgba(53,39,18,1) 0%, rgba(91,67,31,1) 100%)",
		},
		{
			Name:     "Medium Pink Gradient",
			Slug:     "medium-pink-gradient",
			Gradient: "linear-gradient(90deg, rgba(229,59,81,1) 0%, rgba(209,28,51,1) 100%)",
		},
		{
			Name:     "Light Pink Gradient",
			Slug:     "light-pink-gradient",
			Gradient: "linear-gradient(90deg, rgba(255,229,209,1) 0%, rgba(255,200,158,1) 100%)",
		},
		{
			Name:     "Dark Purple Gradient",
			Slug:     "dark-purple-gradient",
			Gradient: "linear-gradient(90deg, rgba(46,34,86,1) 0%, rgba(66,48,123,1) 100%)",
		},
		{
			Name:     "Purple Gradient",
			Slug:     "purple-gradient",
			Gradient: "linear-gradient(90deg, rgba(103,73,112,1) 0%, rgba(131,93,143,1) 100%)",
		},
		{
			Name:     "Blue Gray Gradient",
			Slug:     "blue-gray-gradient",
			Gradient: "linear-gradient(90deg, rgba(34,49,63,1) 0%, rgba(52,75,96,1) 100%)",
		},
		{
			Name:     "Bright Blue Gradient",
			Slug:     "bright-blue-gradient",
			Gradient: "linear-gradient(90deg, rgba(85,195,220,1) 0%, rgba(43,180,211,1) 100%)",
		},
		{
			Name:     "Light Blue Gradient",
			Slug:     "light-blue-gradient",
			Gradient: "linear-gradient(90deg, rgba(233,242,249,1) 0%, rgba(193,218,238,1) 100%)",
		},
	},
}
var colorscheme = map[string]ColorScheme{
	"default": {
		Label: "Default",
		Colors: []string{
			"#f1f1f1",
			"#ffffff",
			"#ffffff",
			"#333333",
			"#333333",
			"#f7f7f7",
		},
	},
	"dark": {
		Label: "Dark",
		Colors: []string{
			"#111111",
			"#202020",
			"#202020",
			"#bebebe",
			"#bebebe",
			"#1b1b1b",
		},
	},

	"pink": {
		Label: "Pink",
		Colors: []string{
			"#ffe5d1",
			"#e53b51",
			"#ffffff",
			"#352712",
			"#ffffff",
			"#f1f1f1",
		},
	},
	"purple": {
		Label: "Purple",
		Colors: []string{
			"#674970",
			"#2e2256",
			"#ffffff",
			"#2e2256",
			"#ffffff",
			"#f1f1f1",
		},
	},
	"blue": {
		Label: "Blue",
		Colors: []string{
			"#e9f2f9",
			"#55c3dc",
			"#ffffff",
			"#22313f",
			"#ffffff",
			"#f1f1f1",
		},
	},
}

var _ = func() struct{} {
	v := wpconfig.ThemeSupport{
		CoreBlockPatterns:  true,
		WidgetsBlockEditor: true,
		AutomaticFeedLinks: true,
		TitleTag:           true,
		PostThumbnails:     true,
		Menus:              true,
		HTML5: []string{
			"search-form",
			"comment-form",
			"comment-list",
			"gallery",
			"caption",
			"script",
			"style",
			"navigation-widgets",
		},
		PostFormats: []string{
			"aside",
			"image",
			"video",
			"quote",
			"link",
			"gallery",
			"status",
			"audio",
			"chat",
		},
		CustomLogo: wpconfig.CustomLogo{
			Width:              248,
			Height:             248,
			FlexWidth:          false,
			FlexHeight:         true,
			HeaderText:         "",
			UnlinkHomepageLogo: false,
		},
		CustomizeSelectiveRefreshWidgets: true,
		EditorStyle:                      true,
		EditorStyles:                     true,
		WpBlockStyles:                    true,
		ResponsiveEmbeds:                 true,
		CustomHeader: wpconfig.CustomHeader{
			DefaultImage:         "",
			RandomDefault:        false,
			Width:                954,
			Height:               1300,
			FlexHeight:           false,
			FlexWidth:            false,
			DefaultTextColor:     "333333",
			HeaderText:           true,
			Uploads:              true,
			WpHeadCallback:       "twentyfifteen_header_style",
			AdminHeadCallback:    "",
			AdminPreviewCallback: "",
			Video:                false,
			VideoActiveCallback:  "is_front_page",
		},
		Widgets: true,
	}
	wpconfig.SetThemeSupport(ThemeName, v)
	return struct{}{}
}()
