package enlightjs

import (
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
)

func EnlighterJS(h *wp.Handle) {
	h.PushGroupHeadScript(20, func(h *wp.Handle) string {
		return `<link rel='stylesheet' id='enlighterjs-css'  href='/wp-content/plugins/enlighter/cache/enlighterjs.min.css' media='all' />`
	})

	h.PushGroupFooterScript(10, func(h *wp.Handle) string {
		return str.Join(`<script src='/wp-content/plugins/enlighter/cache/enlighterjs.min.js?ver=0A0B0C' id='enlighterjs-js'></script>`, "\n", enlighterjs)
	})
}

var enlighterjs = `<script id='enlighterjs-js-after'>
        !function(e,n){if("undefined"!=typeof EnlighterJS){var o={"selectors":{"block":"pre.EnlighterJSRAW","inline":"code.EnlighterJSRAW"},"options":{"indent":4,"ampersandCleanup":true,"linehover":true,"rawcodeDbclick":false,"textOverflow":"break","linenumbers":true,"theme":"enlighter","language":"generic","retainCssClasses":false,"collapse":false,"toolbarOuter":"","toolbarTop":"{BTN_RAW}{BTN_COPY}{BTN_WINDOW}{BTN_WEBSITE}","toolbarBottom":""}};(e.EnlighterJSINIT=function(){EnlighterJS.init(o.selectors.block,o.selectors.inline,o.options)})()}else{(n&&(n.error||n.log)||function(){})("Error: EnlighterJS resources not loaded yet!")}}(window,console);
    </script>`
