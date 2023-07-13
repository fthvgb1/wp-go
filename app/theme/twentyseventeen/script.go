package twentyseventeen

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/theme/wp/scriptloader"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/number"
	str "github.com/fthvgb1/wp-go/helper/strings"
)

func pushScripts(h *wp.Handle) {
	scriptloader.EnqueueStyle("twentyseventeen-style", scriptloader.GetStylesheetUri(), nil, "20230328", "")
	scriptloader.EnqueueStyles("twentyseventeen-block-style", "/assets/css/blocks.css", []string{"twentyseventeen-style"}, "20220912", "")

	if "dark" == wpconfig.GetThemeModsVal(ThemeName, "colorscheme", "light") {
		scriptloader.EnqueueStyles("twentyseventeen-colors-dark", "/assets/css/colors-dark.css",
			[]string{"twentyseventeen-style"}, "20191025", "")
	}

	scriptloader.AddData("twentyseventeen-ie8", "conditional", "lt IE 9")
	scriptloader.EnqueueScripts("html5", "/assets/js/html5.js", nil, "20161020", false)
	scriptloader.AddData("html5", "conditional", "lt IE 9")

	scriptloader.EnqueueScripts("twentyseventeen-skip-link-focus-fix", "/assets/js/skip-link-focus-fix.js",
		nil, "20161114", true)

	l10n := map[string]any{
		"quote": svg(h, map[string]string{"icon": "quote-right"}),
	}

	scriptloader.EnqueueScripts("twentyseventeen-global", "/assets/js/global.js",
		[]string{"jquery"}, "20211130", true)

	scriptloader.EnqueueScripts("jquery-scrollto", "/assets/js/jquery.scrollTo.js",
		[]string{"jquery"}, "2.1.3", true)
	scriptloader.EnqueueScripts("comment-reply", "", nil, "", false)

	//todo  menu top

	scriptloader.AddStaticLocalize("twentyseventeen-skip-link-focus-fix", "twentyseventeenScreenReaderText", l10n)
	scriptloader.AddStaticLocalize("wp-custom-header", "_wpCustomHeaderSettings", map[string]any{
		"mimeType":  `video/mp4`,
		"posterUrl": `/wp-content/uploads/2023/01/cropped-wallhaven-9dm7dd-1.png`,
		"videoUrl":  `/wp-content/uploads/2023/06/BloodMoon_GettyRM_495644264_1080_HD_ZH-CN.mp4`,
		"width":     `2000`,
		"height":    `1199`,
		"minWidth":  `900`,
		"minHeight": `500`,
		"l10n": map[string]any{
			"pause": `<span class="screen-reader-text">暂停背景视频</span><svg class="icon icon-pause" aria-hidden="true" role="img"> <use href="#icon-pause" xlink:href="#icon-pause"></use> </svg>`, "play": `<span class="screen-reader-text">播放背景视频</span><svg class="icon icon-play" aria-hidden="true" role="img"> <use href="#icon-play" xlink:href="#icon-play"></use> </svg>`, "pauseSpeak": `视频已暂停。`, "playSpeak": `视频正在播放。`,
		},
	})
	h.PushCacheGroupHeadScript(constraints.AllScene, "{theme}.head", 30, func(h *wp.Handle) string {
		head := headScript
		if "dark" == wpconfig.GetThemeModsVal(ThemeName, "colorscheme", "light") {
			head = fmt.Sprintf("%s\n%s", headScript, ` <link rel="stylesheet" id="twentyseventeen-colors-dark-css" href="/wp-content/themes/twentyseventeen/assets/css/colors-dark.css?ver=20191025" media="all">`)
		}
		return head
	})
	h.PushGroupFooterScript(constraints.AllScene, "{theme}.footer", 20.005, footerScript)

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

func svg(h *wp.Handle, m map[string]string) string {
	if !maps.IsExists(m, "icon") {
		return ""
	}
	ariaHidden := ` aria-hidden="true"`
	ariaLabelledby := ""
	uniqueId := ""
	if m["title"] != "" {
		ariaHidden = ""
		id := helper.GetContextVal(h.C, "svg", 0)
		uniqueId = number.IntToString(id)
		id++
		h.C.Set("svg", id)
		ariaLabelledby = str.Join(" aria-labelledby=\"title-", uniqueId, "\"")
		if m["desc"] != "" {
			ariaLabelledby = str.Join(" aria-labelledby=\"title-", uniqueId, " desc-", uniqueId, "\"")
		}
	}
	s := str.NewBuilder()
	s.WriteString("<svg class=\"icon icon-", m["icon"], "\"", ariaHidden, ariaLabelledby, " role=\"img\">")
	if m["title"] != "" {
		s.WriteString(`<title id="title-`, uniqueId, `">`, m["title"], "</title>")
		if m["desc"] != "" {
			s.WriteString(`<desc id="desc-`, uniqueId, `">`, m["desc"], `</desc>`)
		}
	}
	s.WriteString(` <use href="#icon-`, m["icon"], `" xlink:href="#icon-`, m["icon"], `"></use> `)
	if m["fallback"] != "" {
		s.WriteString(`<span class="svg-fallback icon-' . esc_attr( $args['icon'] ) . '"></span>`)
	}
	s.WriteString(`<span class="svg-fallback icon-`, m["icon"], `"></span></svg>`)
	return s.String()
}
