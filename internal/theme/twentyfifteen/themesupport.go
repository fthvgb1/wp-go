package twentyfifteen

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

var themesupport themeSupport
var colorscheme map[string]ColorScheme
