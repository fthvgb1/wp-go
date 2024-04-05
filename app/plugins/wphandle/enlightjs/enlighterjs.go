package enlightjs

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/phphelper"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/goccy/go-json"
)

type Config struct {
	Selectors Selectors `json:"selectors"`
	Options   Options   `json:"options"`
}

type Options struct {
	Indent           int64  `json:"indent,omitempty"`
	AmpersandCleanup bool   `json:"ampersandCleanup,omitempty"`
	Linehover        bool   `json:"linehover,omitempty"`
	RawcodeDbclick   bool   `json:"rawcodeDbclick,omitempty"`
	TextOverflow     string `json:"textOverflow,omitempty"`
	Linenumbers      bool   `json:"linenumbers,omitempty"`
	Theme            string `json:"theme,omitempty"`
	Language         string `json:"language,omitempty"`
	RetainCssClasses bool   `json:"retainCssClasses,omitempty"`
	Collapse         bool   `json:"collapse,omitempty"`
	ToolbarOuter     string `json:"toolbarOuter,omitempty"`
	ToolbarTop       string `json:"toolbarTop,omitempty"`
	ToolbarBottom    string `json:"toolbarBottom,omitempty"`
}

type Selectors struct {
	Block  string `json:"block,omitempty"`
	Inline string `json:"inline,omitempty"`
}

func EnlighterJS(h *wp.Handle) {
	h.PushGroupHeadScript(constraints.AllScene, "enlighterjs-css", 20, `<link rel='stylesheet' id='enlighterjs-css'  href='/wp-content/plugins/enlighter/cache/enlighterjs.min.css' media='all' />`)

	h.PushCacheGroupFooterScript(constraints.AllScene, "enlighterJs", 10, func(h *wp.Handle) string {
		op := wpconfig.GetOption("enlighter-options")
		opp, err := phphelper.UnPHPSerializeToStrAnyMap(op)
		if err != nil {
			logs.Error(err, "获取enlighter-option失败", op)
			return ""
		}
		v := Config{
			Selectors: Selectors{
				Block:  maps.GetStrAnyValWithDefaults(opp, "enlighterjs-selector-block", "pre.EnlighterJSRAW"),
				Inline: maps.GetStrAnyValWithDefaults(opp, "enlighterjs-selector-inline", "code.EnlighterJSRAW"),
			},
			Options: Options{
				Indent:           maps.GetStrAnyValWithDefaults[int64](opp, "enlighterjs-indent", 4),
				AmpersandCleanup: maps.GetStrAnyValWithDefaults(opp, "enlighterjs-ampersandcleanup", true),
				Linehover:        maps.GetStrAnyValWithDefaults(opp, "enlighterjs-linehover", true),
				RawcodeDbclick:   maps.GetStrAnyValWithDefaults(opp, "enlighterjs-rawcodedbclick", true),
				TextOverflow:     maps.GetStrAnyValWithDefaults(opp, "enlighterjs-textoverflow", "break"),
				Linenumbers:      maps.GetStrAnyValWithDefaults[bool](opp, "enlighterjs-linenumbers", true),
				Theme:            maps.GetStrAnyValWithDefaults(opp, "enlighterjs-theme", "enlighter"),
				Language:         maps.GetStrAnyValWithDefaults(opp, "enlighterjs-language", "generic"),
				RetainCssClasses: maps.GetStrAnyValWithDefaults(opp, "enlighterjs-retaincss", false),
				Collapse:         false,
				ToolbarOuter:     "",
				ToolbarTop:       "{BTN_RAW}{BTN_COPY}{BTN_WINDOW}{BTN_WEBSITE}",
				ToolbarBottom:    "",
			},
		}
		conf, err := json.Marshal(v)
		if err != nil {
			logs.Error(err, "json化enlighterjs配置失败")
			return ""
		}
		return fmt.Sprintf(enlighterjs, conf)
	})
}

var enlighterjs = `<script src='/wp-content/plugins/enlighter/cache/enlighterjs.min.js?ver=0A0B0C' id='enlighterjs-js'></script>
<script id='enlighterjs-js-after'>
        !function(e,n){if("undefined"!=typeof EnlighterJS){var o=%s;(e.EnlighterJSINIT=function(){EnlighterJS.init(o.selectors.block,o.selectors.inline,o.options)})()}else{(n&&(n.error||n.log)||function(){})("Error: EnlighterJS resources not loaded yet!")}}(window,console);
    </script>`
