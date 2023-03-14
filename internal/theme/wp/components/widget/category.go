package widget

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/safety"
	"strings"
)

var categoryArgs safety.Var[map[string]string]
var categoryConfig = func() safety.Var[map[any]any] {
	v := safety.Var[map[any]any]{}
	v.Store(map[any]any{
		"count":        int64(0),
		"dropdown":     int64(0),
		"hierarchical": int64(0),
		"title":        "分类",
	})
	categoryArgs.Store(map[string]string{
		"{$before_widget}":  `<aside id="categories-2" class="widget widget_categories">`,
		"{$after_widget}":   "</aside>",
		"{$before_title}":   `<h2 class="widget-title">`,
		"{$after_title}":    "</h2>",
		"{$before_sidebar}": "",
		"{$after_sidebar}":  "",
		"{$nav}":            "",
		"{$navCloser}":      "",
		"{$title}":          "",
		"{$dropdown_id}":    "archives-dropdown-2",
		"{$dropdown_type}":  "monthly",
		"{$dropdown_label}": "选择月份",
	})
	return v
}()

var categoryTemplate = `{$before_widget}
{$title}
{$nav}
{$html}
{$navCloser}
{$after_widget}
`

func Category(h *wp.Handle) string {
	args := wp.GetComponentsArgs(h, widgets.ArchiveArgs, categoryArgs.Load())
	args = maps.FilterZeroMerge(categoryArgs.Load(), args)
	conf := wpconfig.GetPHPArrayVal("widget_categories", categoryConfig.Load(), int64(2))
	conf = maps.FilterZeroMerge(categoryConfig.Load(), conf)
	args["{$title}"] = str.Join(args["{$before_title}"], conf["title"].(string), args["{$after_title}"])
	t := categoryTemplate
	dropdown := conf["dropdown"].(int64)
	categories := cache.CategoriesTags(h.C, constraints.Category)
	if dropdown == 1 {

	} else {
		t = strings.ReplaceAll(t, "{$html}", categoryUL(h, args, conf, categories))
	}
	return str.Replace(t, args)
}

func categoryUL(h *wp.Handle, args map[string]string, conf map[any]any, categories []models.TermsMy) string {
	if slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "navigation-widgets") {
		args["{$nav}"] = fmt.Sprintf(`<nav aria-label="%s">`, args["{title}"])
		args["{$navCloser}"] = "</nav>"
	}
	s := str.NewBuilder()
	s.WriteString("<ul>\n")
	isCount := conf["count"].(int64)
	for _, category := range categories {
		count := ""
		if isCount != 0 {
			count = fmt.Sprintf("(%d)", category.Count)
		}
		s.Sprintf(`	<li class="cat-item cat-item-%d">
		<a href="/p/category/%s">%s %s</a>
	</li>
`, category.Terms.TermId, category.Name, category.Name, count)
	}
	s.WriteString("</ul>")
	return s.String()
}
