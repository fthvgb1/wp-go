package wpconfig

type ThemeSupport struct {
	CoreBlockPatterns                bool `json:"core-block-patterns"`
	WidgetsBlockEditor               bool `json:"widgets-block-editor"`
	AutomaticFeedLinks               bool `json:"automatic-feed-links"`
	TitleTag                         bool
	CustomLineHeight                 bool                    `json:"title-tag"`
	PostThumbnails                   bool                    `json:"post-thumbnails"`
	Menus                            bool                    `json:"menus"`
	HTML5                            []string                `json:"html5"`
	PostFormats                      []string                `json:"post-formats"`
	CustomLogo                       CustomLogo              `json:"custom-logo"`
	CustomBackground                 CustomBackground        `json:"custom-background"`
	EditorStyle                      bool                    `json:"editor-style"`
	EditorStyles                     bool                    `json:"editor-styles"`
	WpBlockStyles                    bool                    `json:"wp-block-styles"`
	ResponsiveEmbeds                 bool                    `json:"responsive-embeds"`
	EditorColorPalette               []EditorColorPalette    `json:"editor-color-palette"`
	EditorGradientPresets            []EditorGradientPresets `json:"editor-gradient-presets"`
	CustomizeSelectiveRefreshWidgets bool                    `json:"customize-selective-refresh-widgets"`
	StarterContent                   StarterContent          `json:"starter-content"`
	CustomHeader                     CustomHeader            `json:"custom-header"`
	Widgets                          bool                    `json:"widgets"`
}
type CustomLogo struct {
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	FlexWidth          bool   `json:"flex-width"`
	FlexHeight         bool   `json:"flex-height"`
	HeaderText         string `json:"header-text"`
	UnlinkHomepageLogo bool   `json:"unlink-homepage-logo"`
}
type Widgets struct {
	Sidebar1 []string `json:"sidebar-1"`
	Sidebar2 []string `json:"sidebar-2"`
	Sidebar3 []string `json:"sidebar-3"`
}
type About struct {
	Thumbnail string `json:"thumbnail"`
}
type Contact struct {
	Thumbnail string `json:"thumbnail"`
}
type Blog struct {
	Thumbnail string `json:"thumbnail"`
}
type HomepageSection struct {
	Thumbnail string `json:"thumbnail"`
}
type CustomBackground struct {
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

type Posts struct {
	Num0            string          `json:"0"`
	About           About           `json:"about"`
	Contact         Contact         `json:"contact"`
	Blog            Blog            `json:"blog"`
	HomepageSection HomepageSection `json:"homepage-section"`
}
type ImageEspresso struct {
	PostTitle string `json:"post_title"`
	File      string `json:"file"`
}
type ImageSandwich struct {
	PostTitle string `json:"post_title"`
	File      string `json:"file"`
}
type ImageCoffee struct {
	PostTitle string `json:"post_title"`
	File      string `json:"file"`
}
type Attachments struct {
	ImageEspresso ImageEspresso `json:"image-espresso"`
	ImageSandwich ImageSandwich `json:"image-sandwich"`
	ImageCoffee   ImageCoffee   `json:"image-coffee"`
}
type Option struct {
	ShowOnFront  string `json:"show_on_front"`
	PageOnFront  string `json:"page_on_front"`
	PageForPosts string `json:"page_for_posts"`
}
type ThemeMods struct {
	Panel1 string `json:"panel_1"`
	Panel2 string `json:"panel_2"`
	Panel3 string `json:"panel_3"`
	Panel4 string `json:"panel_4"`
}
type Top struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}
type Social struct {
	Name  string   `json:"name"`
	Items []string `json:"items"`
}
type NavMenus struct {
	Top    Top    `json:"top"`
	Social Social `json:"social"`
}
type StarterContent struct {
	Widgets     Widgets     `json:"widgets"`
	Posts       Posts       `json:"posts"`
	Attachments Attachments `json:"attachments"`
	Options     Option      `json:"options"`
	ThemeMods   ThemeMods   `json:"theme_mods"`
	NavMenus    NavMenus    `json:"nav_menus"`
}
type CustomHeader struct {
	DefaultImage         string `json:"default-image"`
	RandomDefault        bool   `json:"random-default"`
	Width                int    `json:"width"`
	Height               int    `json:"height"`
	FlexHeight           bool   `json:"flex-height"`
	FlexWidth            bool   `json:"flex-width"`
	DefaultTextColor     string `json:"default-text-color"`
	HeaderText           bool   `json:"header-text"`
	Uploads              bool   `json:"uploads"`
	WpHeadCallback       string `json:"wp-head-callback"`
	AdminHeadCallback    string `json:"admin-head-callback"`
	AdminPreviewCallback string `json:"admin-preview-callback"`
	Video                bool   `json:"video"`
	VideoActiveCallback  string `json:"video-active-callback"`
}
