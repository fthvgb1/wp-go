package twentyseventeen

import (
	"fmt"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

func pushScripts(h *wp.Handle) {
	h.PushCacheGroupHeadScript(constraints.AllScene, "{theme}.head", 30, func(h *wp.Handle) string {
		head := headScript
		if "dark" == wpconfig.GetThemeModsVal(ThemeName, "colorscheme", "light") {
			head = fmt.Sprintf("%s\n%s", headScript, ` <link rel="stylesheet" id="twentyseventeen-colors-dark-css" href="/wp-content/themes/twentyseventeen/assets/css/colors-dark.css?ver=20191025" media="all">`)
		}
		return head
	})
	h.PushGroupFooterScript(constraints.AllScene, "{theme}.footer", 20, footerScript)

}

var headScript = `<link rel='stylesheet' id='twentyseventeen-style-css' href='/wp-content/themes/twentyseventeen/style.css?ver=20221101' media='all' />
    <link rel='stylesheet' id='twentyseventeen-block-style-css' href='/wp-content/themes/twentyseventeen/assets/css/blocks.css?ver=20220912' media='all' />
    <!--[if lt IE 9]>
    <link rel='stylesheet' id='twentyseventeen-ie8-css' href='/wp-content/themes/twentyseventeen/assets/css/ie8.css?ver=20161202' media='all' />
    <![endif]-->
    <!--[if lt IE 9]>
    <script src='/wp-content/themes/twentyseventeen/assets/js/html5.js?ver=20161020' id='html5-js'></script>
    <![endif]-->
    <script src='/wp-includes/js/jquery/jquery.min.js?ver=3.6.0' id='jquery-core-js'></script>
    <script src='/wp-includes/js/jquery/jquery-migrate.min.js?ver=3.3.2' id='jquery-migrate-js'></script>
    <link rel="https://api.w.org/" href="/wp-json/" /><link rel="EditURI" type="application/rsd+xml" title="RSD" href="/xmlrpc.php?rsd" />
    <link rel="wlwmanifest" type="application/wlwmanifest+xml" href="/wp-includes/wlwmanifest.xml" />
    <meta name="generator" content="WordPress 6.1.1" />
    <style>.recentcomments a{display:inline !important;padding:0 !important;margin:0 !important;}</style>`

var footerScript = `<script id="twentyseventeen-skip-link-focus-fix-js-extra">
        var twentyseventeenScreenReaderText = {"quote":"<svg class=\"icon icon-quote-right\" aria-hidden=\"true\" role=\"img\"> <use href=\"#icon-quote-right\" xlink:href=\"#icon-quote-right\"><\/use> <\/svg>"};
    </script>

    <script src="/wp-content/themes/twentyseventeen/assets/js/skip-link-focus-fix.js?ver=20161114" id="twentyseventeen-skip-link-focus-fix-js"></script>
    <script src="/wp-content/themes/twentyseventeen/assets/js/global.js?ver=20211130" id="twentyseventeen-global-js"></script>
    <script src="/wp-content/themes/twentyseventeen/assets/js/jquery.scrollTo.js?ver=2.1.3" id="jquery-scrollto-js"></script>`
