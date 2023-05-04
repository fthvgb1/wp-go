package twentyfifteen

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/helper/slice"
	"strconv"
	"strings"
)

func colorSchemeCss(h *wp.Handle) string {
	s := slice.Filter([]string{calColorSchemeCss(h), calSidebarTextColorCss(h), calHeaderBackgroundColorCss(h)}, func(s string, i int) bool {
		return s != ""
	})
	if len(s) < 1 {
		return ""
	}
	return fmt.Sprintf(`<style id='%s-inline-css'%s>\n%s\n</style>`, "twentyfifteen-style", "", strings.Join(s, "\n"))
}
func calColorSchemeCss(h *wp.Handle) (r string) {
	color := getColorScheme(h)
	if "default" == h.CommonThemeMods().ColorScheme || len(color) < 1 {
		return
	}
	textColorRgb := slice.ToAnySlice(Hex2RgbUint8(color[3]))
	sidebarTextColorRgb := Hex2RgbUint8(color[4])
	sidebarTextColorRgbs := slice.ToAnySlice(sidebarTextColorRgb)
	colors := map[string]string{
		"background_color":            color[0],
		"header_background_color":     color[1],
		"box_background_color":        color[2],
		"textcolor":                   color[3],
		"secondary_textcolor":         fmt.Sprintf("rgba(%d, %d, %d, 0.7)", textColorRgb...),
		"border_color":                fmt.Sprintf("rgba(%d, %d, %d, 0.7)", textColorRgb...),
		"border_focus_color":          fmt.Sprintf("rgba(%d, %d, %d, 0.7)", textColorRgb...),
		"sidebar_textcolor":           color[4],
		"sidebar_border_color":        fmt.Sprintf("rgba(%d, %d, %d, 0.7)", sidebarTextColorRgbs...),
		"sidebar_border_focus_color":  fmt.Sprintf("rgba(%d, %d, %d, 0.7)", sidebarTextColorRgbs...),
		"secondary_sidebar_textcolor": fmt.Sprintf("rgba(%d, %d, %d, 0.7)", sidebarTextColorRgbs...),
		"meta_box_background_color":   color[5],
	}
	r = cssTemplate
	for k, v := range colors {
		r = strings.ReplaceAll(r, fmt.Sprintf(`{$colors['%s']}`, k), v)
	}
	return
}

func calSidebarTextColorCss(h *wp.Handle) (r string) {
	colors := getColorScheme(h)
	themeMods := h.CommonThemeMods()
	if themeMods.SidebarTextcolor == "" || themeMods.SidebarTextcolor == colors[4] {
		return
	}
	linkColorRgb := Hex2RgbUint8(themeMods.SidebarTextcolor)
	color := slice.ToAnySlice(linkColorRgb)
	textColor := fmt.Sprintf(`rgba( %[1]v, %[2]v, %[3]v, 0.7)`, color...)
	borderColor := fmt.Sprintf(`rgba( %[1]v, %[2]v, %[3]v, 0.1)`, color...)
	borderFocusColor := fmt.Sprintf(`rgba( %[1]v, %[2]v, %[3]v, 0.3)`, color...)
	r = fmt.Sprintf(sidebarTextColorTemplate, themeMods.SidebarTextcolor, textColor, borderColor, borderFocusColor)
	return
}

func calHeaderBackgroundColorCss(h *wp.Handle) (r string) {
	colors := getColorScheme(h)
	themeMods := h.CommonThemeMods()
	if themeMods.HeaderBackgroundColor == "" || themeMods.HeaderBackgroundColor == colors[1] {
		return
	}
	r = fmt.Sprintf(headerBackgroundColorCssTemplate, themeMods.HeaderBackgroundColor, themeMods.HeaderBackgroundColor)
	return
}

func getColorScheme(h *wp.Handle) (r []string) {
	x, ok := colorscheme[h.CommonThemeMods().ColorScheme]
	if ok {
		r = x.Colors
	}
	return
}

type ColorScheme struct {
	Label  string   `json:"label,omitempty"`
	Colors []string `json:"colors,omitempty"`
}

func Hex2RgbUint8(color string) []uint8 {
	var r []uint8
	color = strings.TrimLeft(color, "#")
	fn := func(s string) uint8 {
		n, _ := strconv.ParseInt(s, 16, 0)
		return uint8(n)
	}
	switch len(color) {
	case 3:
		r = []uint8{color[0], color[1], color[2]}
	case 6:
		r = []uint8{fn(color[:2]), fn(color[2:4]), fn(color[4:])}
	}
	return r
}

var cssTemplate = `
/* Color Scheme */

	/* Background Color */
	body {
		background-color: {$colors['background_color']};
	}

	/* Sidebar Background Color */
	body:before,
	.site-header {
		background-color: {$colors['header_background_color']};
	}

	/* Box Background Color */
	.post-navigation,
	.pagination,
	.secondary,
	.site-footer,
	.hentry,
	.page-header,
	.page-content,
	.comments-area,
	.widecolumn {
		background-color: {$colors['box_background_color']};
	}

	/* Box Background Color */
	button,
	input[type="button"],
	input[type="reset"],
	input[type="submit"],
	.pagination .prev,
	.pagination .next,
	.widget_calendar tbody a,
	.widget_calendar tbody a:hover,
	.widget_calendar tbody a:focus,
	.page-links a,
	.page-links a:hover,
	.page-links a:focus,
	.sticky-post {
		color: {$colors['box_background_color']};
	}

	/* Main Text Color */
	button,
	input[type="button"],
	input[type="reset"],
	input[type="submit"],
	.pagination .prev,
	.pagination .next,
	.widget_calendar tbody a,
	.page-links a,
	.sticky-post {
		background-color: {$colors['textcolor']};
	}

	/* Main Text Color */
	body,
	blockquote cite,
	blockquote small,
	a,
	.dropdown-toggle:after,
	.image-navigation a:hover,
	.image-navigation a:focus,
	.comment-navigation a:hover,
	.comment-navigation a:focus,
	.widget-title,
	.entry-footer a:hover,
	.entry-footer a:focus,
	.comment-metadata a:hover,
	.comment-metadata a:focus,
	.pingback .edit-link a:hover,
	.pingback .edit-link a:focus,
	.comment-list .reply a:hover,
	.comment-list .reply a:focus,
	.site-info a:hover,
	.site-info a:focus {
		color: {$colors['textcolor']};
	}

	/* Main Text Color */
	.entry-content a,
	.entry-summary a,
	.page-content a,
	.comment-content a,
	.pingback .comment-body > a,
	.author-description a,
	.taxonomy-description a,
	.textwidget a,
	.entry-footer a:hover,
	.comment-metadata a:hover,
	.pingback .edit-link a:hover,
	.comment-list .reply a:hover,
	.site-info a:hover {
		border-color: {$colors['textcolor']};
	}

	/* Secondary Text Color */
	button:hover,
	button:focus,
	input[type="button"]:hover,
	input[type="button"]:focus,
	input[type="reset"]:hover,
	input[type="reset"]:focus,
	input[type="submit"]:hover,
	input[type="submit"]:focus,
	.pagination .prev:hover,
	.pagination .prev:focus,
	.pagination .next:hover,
	.pagination .next:focus,
	.widget_calendar tbody a:hover,
	.widget_calendar tbody a:focus,
	.page-links a:hover,
	.page-links a:focus {
		background-color: {$colors['textcolor']}; /* Fallback for IE7 and IE8 */
		background-color: {$colors['secondary_textcolor']};
	}

	/* Secondary Text Color */
	blockquote,
	a:hover,
	a:focus,
	.main-navigation .menu-item-description,
	.post-navigation .meta-nav,
	.post-navigation a:hover .post-title,
	.post-navigation a:focus .post-title,
	.image-navigation,
	.image-navigation a,
	.comment-navigation,
	.comment-navigation a,
	.widget,
	.author-heading,
	.entry-footer,
	.entry-footer a,
	.taxonomy-description,
	.page-links > .page-links-title,
	.entry-caption,
	.comment-author,
	.comment-metadata,
	.comment-metadata a,
	.pingback .edit-link,
	.pingback .edit-link a,
	.post-password-form label,
	.comment-form label,
	.comment-notes,
	.comment-awaiting-moderation,
	.logged-in-as,
	.form-allowed-tags,
	.no-comments,
	.site-info,
	.site-info a,
	.wp-caption-text,
	.gallery-caption,
	.comment-list .reply a,
	.widecolumn label,
	.widecolumn .mu_register label {
		color: {$colors['textcolor']}; /* Fallback for IE7 and IE8 */
		color: {$colors['secondary_textcolor']};
	}

	/* Secondary Text Color */
	blockquote,
	.logged-in-as a:hover,
	.comment-author a:hover {
		border-color: {$colors['textcolor']}; /* Fallback for IE7 and IE8 */
		border-color: {$colors['secondary_textcolor']};
	}

	/* Border Color */
	hr,
	.dropdown-toggle:hover,
	.dropdown-toggle:focus {
		background-color: {$colors['textcolor']}; /* Fallback for IE7 and IE8 */
		background-color: {$colors['border_color']};
	}

	/* Border Color */
	pre,
	abbr[title],
	table,
	th,
	td,
	input,
	textarea,
	.main-navigation ul,
	.main-navigation li,
	.post-navigation,
	.post-navigation div + div,
	.pagination,
	.comment-navigation,
	.widget li,
	.widget_categories .children,
	.widget_nav_menu .sub-menu,
	.widget_pages .children,
	.site-header,
	.site-footer,
	.hentry + .hentry,
	.author-info,
	.entry-content .page-links a,
	.page-links > span,
	.page-header,
	.comments-area,
	.comment-list + .comment-respond,
	.comment-list article,
	.comment-list .pingback,
	.comment-list .trackback,
	.comment-list .reply a,
	.no-comments {
		border-color: {$colors['textcolor']}; /* Fallback for IE7 and IE8 */
		border-color: {$colors['border_color']};
	}

	/* Border Focus Color */
	a:focus,
	button:focus,
	input:focus {
		outline-color: {$colors['textcolor']}; /* Fallback for IE7 and IE8 */
		outline-color: {$colors['border_focus_color']};
	}

	input:focus,
	textarea:focus {
		border-color: {$colors['textcolor']}; /* Fallback for IE7 and IE8 */
		border-color: {$colors['border_focus_color']};
	}

	/* Sidebar Link Color */
	.secondary-toggle:before {
		color: {$colors['sidebar_textcolor']};
	}

	.site-title a,
	.site-description {
		color: {$colors['sidebar_textcolor']};
	}

	/* Sidebar Text Color */
	.site-title a:hover,
	.site-title a:focus {
		color: {$colors['secondary_sidebar_textcolor']};
	}

	/* Sidebar Border Color */
	.secondary-toggle {
		border-color: {$colors['sidebar_textcolor']}; /* Fallback for IE7 and IE8 */
		border-color: {$colors['sidebar_border_color']};
	}

	/* Sidebar Border Focus Color */
	.secondary-toggle:hover,
	.secondary-toggle:focus {
		border-color: {$colors['sidebar_textcolor']}; /* Fallback for IE7 and IE8 */
		border-color: {$colors['sidebar_border_focus_color']};
	}

	.site-title a {
		outline-color: {$colors['sidebar_textcolor']}; /* Fallback for IE7 and IE8 */
		outline-color: {$colors['sidebar_border_focus_color']};
	}

	/* Meta Background Color */
	.entry-footer {
		background-color: {$colors['meta_box_background_color']};
	}

	@media screen and (min-width: 38.75em) {
		/* Main Text Color */
		.page-header {
			border-color: {$colors['textcolor']};
		}
	}

	@media screen and (min-width: 59.6875em) {
		/* Make sure its transparent on desktop */
		.site-header,
		.secondary {
			background-color: transparent;
		}

		/* Sidebar Background Color */
		.widget button,
		.widget input[type="button"],
		.widget input[type="reset"],
		.widget input[type="submit"],
		.widget_calendar tbody a,
		.widget_calendar tbody a:hover,
		.widget_calendar tbody a:focus {
			color: {$colors['header_background_color']};
		}

		/* Sidebar Link Color */
		.secondary a,
		.dropdown-toggle:after,
		.widget-title,
		.widget blockquote cite,
		.widget blockquote small {
			color: {$colors['sidebar_textcolor']};
		}

		.widget button,
		.widget input[type="button"],
		.widget input[type="reset"],
		.widget input[type="submit"],
		.widget_calendar tbody a {
			background-color: {$colors['sidebar_textcolor']};
		}

		.textwidget a {
			border-color: {$colors['sidebar_textcolor']};
		}

		/* Sidebar Text Color */
		.secondary a:hover,
		.secondary a:focus,
		.main-navigation .menu-item-description,
		.widget,
		.widget blockquote,
		.widget .wp-caption-text,
		.widget .gallery-caption {
			color: {$colors['secondary_sidebar_textcolor']};
		}

		.widget button:hover,
		.widget button:focus,
		.widget input[type="button"]:hover,
		.widget input[type="button"]:focus,
		.widget input[type="reset"]:hover,
		.widget input[type="reset"]:focus,
		.widget input[type="submit"]:hover,
		.widget input[type="submit"]:focus,
		.widget_calendar tbody a:hover,
		.widget_calendar tbody a:focus {
			background-color: {$colors['secondary_sidebar_textcolor']};
		}

		.widget blockquote {
			border-color: {$colors['secondary_sidebar_textcolor']};
		}

		/* Sidebar Border Color */
		.main-navigation ul,
		.main-navigation li,
		.widget input,
		.widget textarea,
		.widget table,
		.widget th,
		.widget td,
		.widget pre,
		.widget li,
		.widget_categories .children,
		.widget_nav_menu .sub-menu,
		.widget_pages .children,
		.widget abbr[title] {
			border-color: {$colors['sidebar_border_color']};
		}

		.dropdown-toggle:hover,
		.dropdown-toggle:focus,
		.widget hr {
			background-color: {$colors['sidebar_border_color']};
		}

		.widget input:focus,
		.widget textarea:focus {
			border-color: {$colors['sidebar_border_focus_color']};
		}

		.sidebar a:focus,
		.dropdown-toggle:focus {
			outline-color: {$colors['sidebar_border_focus_color']};
		}
	}
`

var headerBackgroundColorCssTemplate = `
/* Custom Header Background Color */
		body:before,
		.site-header {
			background-color: %s;
		}

		@media screen and (min-width: 59.6875em) {
			.site-header,
			.secondary {
				background-color: transparent;
			}

			.widget button,
			.widget input[type="button"],
			.widget input[type="reset"],
			.widget input[type="submit"],
			.widget_calendar tbody a,
			.widget_calendar tbody a:hover,
			.widget_calendar tbody a:focus {
				color: %s;
			}
		}
`

var sidebarTextColorTemplate = `
/* Custom Sidebar Text Color */
		.site-title a,
		.site-description,
		.secondary-toggle:before {
			color: %[1]v;
		}

		.site-title a:hover,
		.site-title a:focus {
			color: %[1]v; /* Fallback for IE7 and IE8 */
			color: %[2]v;
		}

		.secondary-toggle {
			border-color: %[1]v; /* Fallback for IE7 and IE8 */
			border-color: %[3]v;
		}

		.secondary-toggle:hover,
		.secondary-toggle:focus {
			border-color: %[1]v; /* Fallback for IE7 and IE8 */
			border-color: %[4]v;
		}

		.site-title a {
			outline-color: %[1]v; /* Fallback for IE7 and IE8 */
			outline-color: %[4]v;
		}

		@media screen and (min-width: 59.6875em) {
			.secondary a,
			.dropdown-toggle:after,
			.widget-title,
			.widget blockquote cite,
			.widget blockquote small {
				color: %[1]v;
			}

			.widget button,
			.widget input[type="button"],
			.widget input[type="reset"],
			.widget input[type="submit"],
			.widget_calendar tbody a {
				background-color: %[1]v;
			}

			.textwidget a {
				border-color: %[1]v;
			}

			.secondary a:hover,
			.secondary a:focus,
			.main-navigation .menu-item-description,
			.widget,
			.widget blockquote,
			.widget .wp-caption-text,
			.widget .gallery-caption {
				color: %[2]v;
			}

			.widget button:hover,
			.widget button:focus,
			.widget input[type="button"]:hover,
			.widget input[type="button"]:focus,
			.widget input[type="reset"]:hover,
			.widget input[type="reset"]:focus,
			.widget input[type="submit"]:hover,
			.widget input[type="submit"]:focus,
			.widget_calendar tbody a:hover,
			.widget_calendar tbody a:focus {
				background-color: %[2]v;
			}

			.widget blockquote {
				border-color: %[2]v;
			}

			.main-navigation ul,
			.main-navigation li,
			.secondary-toggle,
			.widget input,
			.widget textarea,
			.widget table,
			.widget th,
			.widget td,
			.widget pre,
			.widget li,
			.widget_categories .children,
			.widget_nav_menu .sub-menu,
			.widget_pages .children,
			.widget abbr[title] {
				border-color: %[3]v;
			}

			.dropdown-toggle:hover,
			.dropdown-toggle:focus,
			.widget hr {
				background-color: %[3]v;
			}

			.widget input:focus,
			.widget textarea:focus {
				border-color: %[4]v;
			}

			.sidebar a:focus,
			.dropdown-toggle:focus {
				outline-color: %[4]v;
			}
		}
`
