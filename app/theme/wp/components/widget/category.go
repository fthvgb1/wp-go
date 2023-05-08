package widget

import (
	"fmt"
	"github.com/fthvgb1/wp-go/app/cmd/reload"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/helper/tree"
	"net/http"
	"strings"
)

var categoryTemplate = `{$before_widget}
{$title}
{$nav}
{$html}
{$navCloser}
{$after_widget}
`
var categoryConfig = map[any]any{
	"count":        int64(0),
	"dropdown":     int64(0),
	"hierarchical": int64(0),
	"title":        "分类",
}

func categoryArgs() map[string]string {
	return map[string]string{
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
	}
}

func Category(h *wp.Handle, id string) string {
	conf := configs(categoryConfig, "widget_categories", int64(2))

	args := reload.GetAnyValBys("widget-category-args", h, func(h *wp.Handle) map[string]string {
		commonArgs := wp.GetComponentsArgs(h, widgets.Widget, map[string]string{})
		args := wp.GetComponentsArgs(h, widgets.Categories, categoryArgs())
		args = maps.FilterZeroMerge(categoryArgs(), CommonArgs(), commonArgs, args)
		args["{$before_widget}"] = fmt.Sprintf(args["{$before_widget}"], str.Join("categories-", id), str.Join("widget widget_", "categories"))
		args["{$title}"] = str.Join(args["{$before_title}"], conf["title"].(string), args["{$after_title}"])
		if conf["dropdown"].(int64) == 0 && slice.IsContained(h.CommonThemeMods().ThemeSupport.HTML5, "navigation-widgets") {
			args["{$nav}"] = fmt.Sprintf(`<nav aria-label="%s">`, args["{title}"])
			args["{$navCloser}"] = "</nav>"
		}
		return args
	})

	t := categoryTemplate
	dropdown := conf["dropdown"].(int64)
	categories := cache.CategoriesTags(h.C, constraints.Category)
	if dropdown == 1 {
		t = strings.ReplaceAll(t, "{$html}", CategoryDropdown(h, args, conf, categories))
	} else {
		t = strings.ReplaceAll(t, "{$html}", categoryUL(h, args, conf, categories))
	}
	return h.ComponentFilterFnHook(widgets.Categories, str.Replace(t, args))
}

func CategoryLi(h *wp.Handle, conf map[any]any, categories []models.TermsMy) string {
	s := str.NewBuilder()
	isCount := conf["count"].(int64)
	currentCate := models.TermsMy{}
	if h.Scene() == constraints.Category {
		cat := h.C.Param("category")
		_, currentCate = slice.SearchFirst(categories, func(my models.TermsMy) bool {
			return cat == my.Name
		})
	}
	if conf["hierarchical"].(int64) == 0 {
		for _, category := range categories {
			count := ""
			if isCount != 0 {
				count = fmt.Sprintf("(%d)", category.Count)
			}
			current := ""
			if category.TermTaxonomyId == currentCate.TermTaxonomyId {
				current = "current-cat"
			}
			s.Sprintf(`	<li class="cat-item cat-item-%d %s">
		<a href="/p/category/%s">%s %s</a>
	</li>
`, category.Terms.TermId, current, category.Name, category.Name, count)
		}
	} else {

		m := tree.Roots(categories, 0, func(cate models.TermsMy) (child, parent uint64) {
			return cate.TermTaxonomyId, cate.Parent
		})
		cate := &tree.Node[models.TermsMy, uint64]{Data: models.TermsMy{}}
		if currentCate.TermTaxonomyId > 0 {
			cate = m[currentCate.TermTaxonomyId]
		}
		r := m[0]
		categoryLi(r, cate, tree.Ancestor(m, 0, cate), isCount, s)
	}
	return s.String()
}

func categoryUL(h *wp.Handle, args map[string]string, conf map[any]any, categories []models.TermsMy) string {
	s := str.NewBuilder()
	s.WriteString("<ul>\n")
	s.WriteString(CategoryLi(h, conf, categories))
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
	
`, category.Terms.TermId, strings.Join(class, " "), aria, category.Name, category.Name, count)

		if len(*child.Children) > 0 {
			s.WriteString(`	<ul class="children">
`)
			categoryLi(&child, cate, roots, isCount, s)
			s.WriteString(`</ul>
`)
		}
		s.Sprintf(`</li>`)
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

func CategoryDropdown(h *wp.Handle, args map[string]string, conf map[any]any, categories []models.TermsMy) string {
	s := str.NewBuilder()
	s.WriteString(`<form action="/" method="get">
`)
	s.Sprintf(`	<label class="screen-reader-text" for="%s">%s</label>
`, args["{$selectId}"], args["{$title}"])
	s.WriteString(DropdownCategories(h, args, conf, categories))
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

func DropdownCategories(h *wp.Handle, args map[string]string, conf map[any]any, categories []models.TermsMy) string {
	if len(categories) < 1 {
		return ""
	}
	s := str.NewBuilder()
	s.Sprintf(`	<select %s name="%s" id="%s" class="%s">
`, args["{$required}"], args["{$name}"], args["{$selectId}"], args["{$class}"])
	s.Sprintf(`		<option value="-1">%s</option>
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
	return h.ComponentFilterFnHook("wp_dropdown_cats", s.String())
}

func IsCategory(h *wp.Handle) {
	name, ok := parseDropdownCate(h)
	if ok {
		h.C.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/p/category/%s", name))
		h.Abort()
	}
}

func parseDropdownCate(h *wp.Handle) (cateName string, r bool) {
	cate := wp.GetComponentsArgs[map[string]string](h, widgets.Categories, categoryArgs())
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
