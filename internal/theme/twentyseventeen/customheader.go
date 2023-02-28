package twentyseventeen

import (
	"fmt"
	"github.com/fthvgb1/wp-go/internal/theme/common"
)

func customHeader(h *common.Handle) (r string) {
	headerTextColor := h.ThemeMods.HeaderTextcolor
	if headerTextColor == "" || headerTextColor == h.ThemeMods.ThemeSupport.CustomHeader.DefaultTextColor {
		return
	}
	css := `
		.site-title,
		.site-description {
			position: absolute;
			clip: rect(1px, 1px, 1px, 1px);
		}`
	if headerTextColor != "blank" {
		css = fmt.Sprintf(customHeaderCss, headerTextColor)
	}
	r = fmt.Sprintf(`<style id="twentyseventeen-custom-header-styles" type="text/css">%s</style>`, css)
	return
}

var customHeaderCss = `
		.site-title a,
		.colors-dark .site-title a,
		.colors-custom .site-title a,
		body.has-header-image .site-title a,
		body.has-header-video .site-title a,
		body.has-header-image.colors-dark .site-title a,
		body.has-header-video.colors-dark .site-title a,
		body.has-header-image.colors-custom .site-title a,
		body.has-header-video.colors-custom .site-title a,
		.site-description,
		.colors-dark .site-description,
		.colors-custom .site-description,
		body.has-header-image .site-description,
		body.has-header-video .site-description,
		body.has-header-image.colors-dark .site-description,
		body.has-header-video.colors-dark .site-description,
		body.has-header-image.colors-custom .site-description,
		body.has-header-video.colors-custom .site-description {
			color: #%s;
		}
`
