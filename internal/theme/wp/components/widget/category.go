package widget

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/helper/tree"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/pkg/models"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"github.com/fthvgb1/wp-go/safety"
	"net/http"
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
		"{$before_widget}":    `<aside id="categories-2" class="widget widget_categories">`,
		"{$after_widget}":     "</aside>",
		"{$before_title}":     `<h2 class="widget-title">`,
		"{$after_title}":      "</h2>",
		"{$before_sidebar}":   "",
		"{$after_sidebar}":    "",
		"{$class}":            "postform",
		"{$show_option_none}": "选择分类",
		"{$name}":             "cat",
		"{$selectId}":         "cat",
		"{$required}":         "",
		"{$nav}":              "",
		"{$navCloser}":        "",
		"{$title}":            "",
		"{$dropdown_id}":      "archives-dropdown-2",
		"{$dropdown_type}":    "monthly",
		"{$dropdown_label}":   "选择月份",
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
	args := wp.GetComponentsArgs(h, widgets.Categories, categoryArgs.Load())
	args = maps.FilterZeroMerge(categoryArgs.Load(), args)
	conf := wpconfig.GetPHPArrayVal("widget_categories", categoryConfig.Load(), int64(2))
	conf = maps.FilterZeroMerge(categoryConfig.Load(), conf)
	args["{$title}"] = str.Join(args["{$before_title}"], conf["title"].(string), args["{$after_title}"])
	t := categoryTemplate
	dropdown := conf["dropdown"].(int64)
	if id, ok := args["{$id}"]; ok && id != "" {
		args["{$before_widget}"] = strings.ReplaceAll(args["{$before_widget}"], "2", args["{$id}"])
	}
	categories := cache.CategoriesTags(h.C, constraints.Category)
	if dropdown == 1 {
		t = strings.ReplaceAll(t, "{$html}", categoryDropdown(h, args, conf, categories))

	} else {
		t = strings.ReplaceAll(t, "{$html}", categoryUL(h, args, conf, categories))
	}
	return h.ComponentFilterFnHook(widgets.Categories, str.Replace(t, args))
}

func categoryUL(h *wp.Handle, args map[string]string, conf map[any]any, categories []models.TermsMy) string {
	if slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "navigation-widgets") {
		args["{$nav}"] = fmt.Sprintf(`<nav aria-label="%s">`, args["{title}"])
		args["{$navCloser}"] = "</nav>"
	}
	s := str.NewBuilder()
	s.WriteString("<ul>\n")
	isCount := conf["count"].(int64)
	if conf["hierarchical"].(int64) == 0 {
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
	} else {

		m := tree.Roots(categories, 0, func(cate models.TermsMy) (child, parent uint64) {
			return cate.TermTaxonomyId, cate.Parent
		})
		cate := &tree.Node[models.TermsMy, uint64]{Data: models.TermsMy{}}
		if h.Scene() == constraints.Category {
			cat := h.C.Param("category")
			i, ca := slice.SearchFirst(categories, func(my models.TermsMy) bool {
				return cat == my.Name
			})
			if i > 0 {
				cate = m[ca.TermTaxonomyId]
			}
		}

		r := m[0]
		categoryLi(r, cate, tree.Ancestor(m, 0, cate), isCount, s)
	}

	s.WriteString("</ul>")
	return s.String()
}

func categoryLi(root *tree.Node[models.TermsMy, uint64], cate, roots *tree.Node[models.TermsMy, uint64], isCount int64, s *str.Builder) {
	for _, child := range *root.Children {
		category := child.Data
		count := ""
		if isCount != 0 {
			count = fmt.Sprintf("(%d)", category.Count)
		}
		var class []string

		if len(*child.Children) > 0 && cate.Data.TermTaxonomyId > 0 {
			if category.TermTaxonomyId == cate.Parent {
				class = append(class, "current-cat-parent")
			}

			if cate.Parent > 0 && category.TermTaxonomyId == roots.Data.TermTaxonomyId {
				class = append(class, "current-cat-ancestor")
			}
		}
		aria := ""
		if category.TermTaxonomyId == cate.Data.TermTaxonomyId {
			class = append(class, "current-cat")
			aria = `aria-current="page"`
		}
		s.Sprintf(`	<li class="cat-item cat-item-%d %s">
		<a %s href="/p/category/%s">%s %s</a>
	</li>
`, category.Terms.TermId, strings.Join(class, " "), aria, category.Name, category.Name, count)

		if len(*child.Children) > 0 {
			s.WriteString(`<ul class="children">
`)
			categoryLi(&child, cate, roots, isCount, s)
			s.WriteString(`</ul>`)
		}
	}

}

var categoryDropdownJs = `/* <![CDATA[ */
(function() {
	var dropdown = document.getElementById( "%s" );
	function onCatChange() {
		if ( dropdown.options[ dropdown.selectedIndex ].value > 0 ) {
			dropdown.parentNode.submit();
		}
	}
	dropdown.onchange = onCatChange;
})();
/* ]]> */
`

func categoryDropdown(h *wp.Handle, args map[string]string, conf map[any]any, categories []models.TermsMy) string {
	s := str.NewBuilder()
	s.WriteString(`<form action="/" method="get">
`)
	s.Sprintf(`	<label class="screen-reader-text" for="%s">%s</label>
`, args["{$selectId}"], args["{$title}"])
	if len(categories) > 0 {
		s.Sprintf(`	<select %s name="%s" id="%s" class="%s">
`, args["{$required}"], args["{$name}"], args["{$selectId}"], args["{$class}"])
		s.Sprintf(`		<option value="%[1]s">%[1]s</option>
`, args["{$show_option_none}"])
		currentCategory := ""
		if h.Scene() == constraints.Category {
			currentCategory = h.Index.Param.Category
		}
		showCount := conf["count"].(int64)
		fn := func(category models.TermsMy, deep int) {
			lv := fmt.Sprintf("level-%d", deep+1)
			sep := strings.Repeat("&nbsp;", deep*2)
			selected := ""
			if category.Name == currentCategory {
				selected = "selected"
			}
			count := ""
			if showCount != 0 {
				count = fmt.Sprintf("(%d)", category.Count)
			}
			s.Sprintf(`		<option class="%s" %s value="%d">%s%s %s</option>
`, lv, selected, category.Terms.TermId, sep, category.Name, count)
		}
		if conf["hierarchical"].(int64) == 0 {
			for _, category := range categories {
				fn(category, 0)
			}
		} else {
			tree.Root(categories, 0, func(t models.TermsMy) (child, parent uint64) {
				return t.TermTaxonomyId, t.Parent
			}).Loop(func(category models.TermsMy, deep int) {
				fn(category, deep)
			})
		}
		s.WriteString("	</select>\n")
	}
	s.WriteString("</form>\n")
	attr := ""
	if !slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "script") {
		attr = ` type="text/javascript"`
	}
	s.Sprintf(`<script%s>
`, attr)
	s.Sprintf(categoryDropdownJs, args["{$selectId}"])
	s.WriteString("</script>\n")
	return s.String()
}

func IsCategory(next wp.HandleFn[*wp.Handle], h *wp.Handle) {
	if h.Scene() != constraints.Home {
		next(h)
		return
	}
	name, ok := parseDropdownCate(h)
	if !ok {
		next(h)
		return
	}
	h.C.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/p/category/%s", name))
}

func parseDropdownCate(h *wp.Handle) (cateName string, r bool) {
	cate := wp.GetComponentsArgs[map[string]string](h, widgets.Categories, categoryArgs.Load())
	name, ok := cate["{$name}"]
	if !ok || name == "" {
		return
	}
	cat := h.C.Query(name)
	if cat == "" {
		return
	}
	id := str.ToInteger[uint64](cat, 0)
	if id < 1 {
		return
	}
	i, cc := slice.SearchFirst(cache.CategoriesTags(h.C, constraints.Category), func(my models.TermsMy) bool {
		return id == my.Terms.TermId
	})
	if i < 0 {
		return
	}
	r = true
	cateName = cc.Name
	return
}
