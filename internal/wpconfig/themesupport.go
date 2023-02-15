package wpconfig

type ThemeSupport struct {
	CoreBlockPatterns                bool                    `json:"core-block-patterns"`
	WidgetsBlockEditor               bool                    `json:"widgets-block-editor"`
	AutomaticFeedLinks               bool                    `json:"automatic-feed-links"`
	TitleTag                         bool                    `json:"title-tag"`
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
