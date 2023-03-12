package wp

import (
	"github.com/fthvgb1/wp-go/helper/html"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/components"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var searchTemplate = `{$before_widget}
{$title}
{$form}
{$after_widget}`

var searchArgs = map[string]string{
	"{$before_widget}": `<aside id="search-2" class="widget widget_search">`,
	"{$after_widget}":  `</aside>`,
	"{$aria_label}":    "",
	"{$title}":         "",
	"{$before_title}":  `<h2 class="widget-title">`,
	"{$after_title}":   `</h2>`,
	"{$button}":        "搜索",
	"{$placeholder}":   "搜索&hellip;",
	"{$label}":         "搜索：",
}

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

func SearchForm(h *Handle) string {
	args := GetComponentsArgs(h, components.SearchFormArgs, searchArgs)
	args = maps.Merge(searchArgs, args)
	if args["{$title}"] == "" {
		args["{$title}"] = wpconfig.GetPHPArrayVal("widget_search", "", int64(2), "title")
	}
	if args["{$title}"] != "" {
		args["{$title}"] = str.Join(args["{$before_title}"], args["{$title}"], args["{$after_title}"])
	}
	args["{$value}"] = ""
	if h.Scene() == constraints.Search {
		args["{$value}"] = html.SpecialChars(h.Index.Param.Search)
	}
	form := html5SearchForm
	if !slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "search-form") {
		form = xmlSearchForm
	}
	s := strings.ReplaceAll(searchTemplate, "{$form}", form)
	return h.ComponentFilterFnHook(components.SearchFormArgs, str.Replace(s, args))
}
