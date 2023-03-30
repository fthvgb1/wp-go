package widget

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/html"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var searchTemplate = `{$before_widget}
{$title}
{$form}
{$after_widget}`

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

func searchArgs() map[string]string {
	return map[string]string{
		"{$aria_label}":  "",
		"{$title}":       "",
		"{$form}":        "",
		"{$button}":      "搜索",
		"{$placeholder}": "搜索&hellip;",
		"{$label}":       "搜索：",
	}
}

func Search(h *wp.Handle, id string) string {
	form := html5SearchForm
	args := reload.GetAnyValBys("widget-search-args", h, func(h *wp.Handle) map[string]string {
		search := searchArgs()
		commonArgs := wp.GetComponentsArgs(h, widgets.Widget, map[string]string{})
		args := wp.GetComponentsArgs(h, widgets.Search, search)
		args = maps.FilterZeroMerge(search, CommonArgs(), commonArgs, args)
		args["{$before_widget}"] = fmt.Sprintf(args["{$before_widget}"], str.Join("search-", id), str.Join("widget widget_", "search"))
		if args["{$title}"] == "" {
			args["{$title}"] = wpconfig.GetPHPArrayVal("widget_search", "", int64(2), "title")
		}

		if args["{$title}"] != "" {
			args["{$title}"] = str.Join(args["{$before_title}"], args["{$title}"], args["{$after_title}"])
		}
		if args["{$form}"] != "" {
			form = args["{$form}"]
			delete(args, "{$form}")
		}
		if !slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "navigation-widgets") {
			form = xmlSearchForm
		}

		return args
	})
	args = maps.Copy(args)
	s := strings.ReplaceAll(searchTemplate, "{$form}", form)
	val := ""
	if h.Scene() == constraints.Search {
		val = html.SpecialChars(h.Index.Param.Search)
	}
	s = strings.ReplaceAll(s, "{$value}", val)
	return h.ComponentFilterFnHook(widgets.Search, str.Replace(s, args))
}
