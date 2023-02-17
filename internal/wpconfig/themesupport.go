package wpconfig

type ThemeSupport struct {
	CoreBlockPatterns                bool         `json:"core-block-patterns"`
	WidgetsBlockEditor               bool         `json:"widgets-block-editor"`
	AutomaticFeedLinks               bool         `json:"automatic-feed-links"`
	TitleTag                         bool         `json:"title-tag"`
	PostThumbnails                   bool         `json:"post-thumbnails"`
	Menus                            bool         `json:"menus"`
	HTML5                            []string     `json:"html5"`
	PostFormats                      []string     `json:"post-formats"`
	CustomLogo                       CustomLogo   `json:"custom-logo"`
	CustomizeSelectiveRefreshWidgets bool         `json:"customize-selective-refresh-widgets"`
	EditorStyle                      bool         `json:"editor-style"`
	EditorStyles                     bool         `json:"editor-styles"`
	WpBlockStyles                    bool         `json:"wp-block-styles"`
	ResponsiveEmbeds                 bool         `json:"responsive-embeds"`
	CustomHeader                     CustomHeader `json:"custom-header"`
	Widgets                          bool         `json:"widgets"`
}
type CustomLogo struct {
	Width              int    `json:"width"`
	Height             int    `json:"height"`
	FlexWidth          bool   `json:"flex-width"`
	FlexHeight         bool   `json:"flex-height"`
	HeaderText         string `json:"header-text"`
	UnlinkHomepageLogo bool   `json:"unlink-homepage-logo"`
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
