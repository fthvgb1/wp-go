package twentyfifteen

import (
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
)

var style = `<style type="text/css" id="twentyfifteen-header-css">`
var defaultTextStyle = `.site-header {
			padding-top: 14px;
			padding-bottom: 14px;
		}

		.site-branding {
			min-height: 42px;
		}

		@media screen and (min-width: 46.25em) {
			.site-header {
				padding-top: 21px;
				padding-bottom: 21px;
			}
			.site-branding {
				min-height: 56px;
			}
		}
		@media screen and (min-width: 55em) {
			.site-header {
				padding-top: 25px;
				padding-bottom: 25px;
			}
			.site-branding {
				min-height: 62px;
			}
		}
		@media screen and (min-width: 59.6875em) {
			.site-header {
				padding-top: 0;
				padding-bottom: 0;
			}
			.site-branding {
				min-height: 0;
			}
		}`
var imgStyle = `.site-header {

			/*
			 * No shorthand so the Customizer can override individual properties.
			 * @see https://core.trac.wordpress.org/ticket/31460
			 */
			background-image: url("%s");
			background-repeat: no-repeat;
			background-position: 50% 50%;
			-webkit-background-size: cover;
			-moz-background-size:    cover;
			-o-background-size:      cover;
			background-size:         cover;
		}

		@media screen and (min-width: 59.6875em) {
			body:before {

				/*
				 * No shorthand so the Customizer can override individual properties.
				 * @see https://core.trac.wordpress.org/ticket/31460
				 */
				background-image: url("%s");
				background-repeat: no-repeat;
				background-position: 100% 50%;
				-webkit-background-size: cover;
				-moz-background-size:    cover;
				-o-background-size:      cover;
				background-size:         cover;
				border-right: 0;
			}

			.site-header {
				background: transparent;
			}
		}`

var header = reload.Vars(constraints.Defaults)

func calCustomHeader(h *wp.Handle) (r string, rand bool) {
	img, rand := h.GetCustomHeader()
	if img.Path == "" && h.DisplayHeaderText() {
		return
	}
	ss := str.NewBuilder()
	ss.WriteString(style)
	if img.Path == "" && !h.DisplayHeaderText() {
		ss.WriteString(defaultTextStyle)
	}
	if img.Path != "" {
		ss.Sprintf(imgStyle, img.Path, img.Path)
	}
	if !h.DisplayHeaderText() {
		ss.WriteString(`.site-title,
		.site-description {
			clip: rect(1px, 1px, 1px, 1px);
			position: absolute;
		}`)
	}
	ss.WriteString("</style>")
	r = ss.String()
	return
}

func customHeader(h *wp.Handle) {
	headers := header.Load()
	if headers == constraints.Defaults {
		headerss, rand := calCustomHeader(h)
		headers = headerss
		if !rand {
			header.Store(headers)
		}
	}
	h.SetData("customHeader", headers)
	return
}
