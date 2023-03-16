package widget

import (
	"github.com/fthvgb1/wp-go/helper/html"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/safety"
	"strings"
)

var searchTemplate = `{$before_widget}
{$title}
{$form}
{$after_widget}`

var searchArgs = func() safety.Var[map[string]string] {
	v := safety.Var[map[string]string]{}
	v.Store(map[string]string{
		"{$id}":            "2",
		"{$before_widget}": `<aside id="search-2" class="widget widget_search">`,
		"{$after_widget}":  `</aside>`,
		"{$aria_label}":    "",
		"{$title}":         "",
		"{$before_title}":  `<h2 class="widget-title">`,
		"{$after_title}":   `</h2>`,
		"{$form}":          "",
		"{$button}":        "搜索",
		"{$placeholder}":   "搜索&hellip;",
		"{$label}":         "搜索：",
	})
	return v
}()

var html5SearchForm = `<form role="search" {$aria_label} method="get" class="search-form" action="/">
				<label>
					<span class="screen-reader-text">{$label}</span>
					<input type="search" class="search-field" placeholder="{$placeholder}" value="{$value}" name="s" />
				</label>
				<input type="submit" class="search-submit" value="{$button}" />
			</form>`
var xmlSearchForm = `<form role="search" {$aria_label} method="get" id="searchform" class="searchform" action="/">
				<div>
					<label class="screen-reader-text" for="s">{$label}</label>
					<input type="text" value="{$value}" name="s" id="s" />
					<input type="submit" id="searchsubmit" value="{$button}" />
				</div>
			</form>`

func SearchForm(h *wp.Handle) string {
	args := wp.GetComponentsArgs(h, widgets.Search, searchArgs.Load())
	args = maps.FilterZeroMerge(searchArgs.Load(), args)
	if args["{$title}"] == "" {
		args["{$title}"] = wpconfig.GetPHPArrayVal("widget_search", "", int64(2), "title")
	}
	if id, ok := args["{$id}"]; ok && id != "" {
		args["{$before_widget}"] = strings.ReplaceAll(args["{$before_widget}"], "2", args["{$id}"])
	}
	if args["{$title}"] != "" {
		args["{$title}"] = str.Join(args["{$before_title}"], args["{$title}"], args["{$after_title}"])
	}
	args["{$value}"] = ""
	if h.Scene() == constraints.Search {
		args["{$value}"] = html.SpecialChars(h.Index.Param.Search)
	}
	form := html5SearchForm
	if args["{$form}"] != "" {
		form = args["{$form}"]
		delete(args, "{$form}")
	}

	if !slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "navigation-widgets") {
		form = xmlSearchForm
	}
	s := strings.ReplaceAll(searchTemplate, "{$form}", form)
	return h.ComponentFilterFnHook(widgets.Search, str.Replace(s, args))
}
