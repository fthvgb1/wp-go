package twentyseventeen

type themeSupport struct {
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

var themesupport = themeSupport{}
