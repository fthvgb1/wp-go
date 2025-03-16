package block

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fthvgb1/wp-go/app/pkg/cache"
	constraints2 "github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/app/pkg/logs"
	"github.com/fthvgb1/wp-go/app/pkg/models"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/cache/reload"
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/helper/number"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"strings"
)

func categoryConf() map[any]any {
	return map[any]any{
		"count":        int64(0),
		"dropdown":     int64(0),
		"hierarchical": int64(0),
		"title":        "分类",
	}
}

func categoryDefaultArgs() map[string]string {
	return map[string]string{
		"{$before_widget}":    `<aside id="%s" class="%s">`,
		"{$after_widget}":     `</aside>`,
		"{$name}":             "cat",
		"{$class}":            "postform",
		"{$selectId}":         "cat",
		"{$required}":         "",
		"{$show_option_none}": "选择分类",
		"{$title}":            "",
	}
}

func parseAttr(attr map[any]any) string {
	var attrs []string
	class := maps.GetAnyAnyValWithDefaults(attr, "", "className")
	classes := strings.Split(class, " ")
	fontsize := maps.GetAnyAnyValWithDefaults(attr, "", "fontSize")
	if fontsize != "" {
		classes = append(classes, fmt.Sprintf("has-%s-font-size", fontsize))
	}
	style := maps.GetAnyAnyValWithDefaults[map[any]any](attr, nil, "style", "typography")
	if len(style) > 0 {
		styless := maps.AnyAnyMapTo(style, func(k, v any) (string, string, bool) {
			kk, ok := k.(string)
			if !ok {
				return "", "", false
			}
			vv, ok := v.(string)
			if !ok {
				return "", "", false
			}
			return kk, vv, true
		})
		styles := maps.FilterToSlice(styless, func(k string, v string) (string, bool) {
			k = str.CamelCaseTo(k, '-')
			return str.Join(k, ":", v), true
		})
		attrs = append(attrs, fmt.Sprintf(`style="%s;"`, strings.Join(styles, ";")))
	}
	attrs = append(attrs, fmt.Sprintf(`class="%s"`, strings.Join(classes, " ")))
	return strings.Join(attrs, " ")
}

var GetCategoryAttr = reload.BuildValFn("block-category-attr", parseAttr)

var GetCategoryConf = reload.BuildValFnWithConfirm("block-category-conf", categoryConfFn, 5)

func categoryConfFn(blockParser ParserBlock) (map[any]any, bool) {
	var con any
	err := json.Unmarshal([]byte(blockParser.Attrs), &con)
	if err != nil {
		logs.Error(err, "解析category attr错误", blockParser.Attrs)
		return nil, false
	}
	var conf map[any]any
	switch con.(type) {
	case map[any]any:
		conf = con.(map[any]any)
	case map[string]any:
		conf = maps.StrAnyToAnyAny(con.(map[string]any))
	}
	conf = maps.FilterZeroMerge(categoryConf(), conf)

	if maps.GetAnyAnyValWithDefaults(conf, false, "showPostCounts") {
		conf["count"] = int64(1)
	}

	if maps.GetAnyAnyValWithDefaults(conf, false, "displayAsDropdown") {
		conf["dropdown"] = int64(1)
	}
	if maps.GetAnyAnyValWithDefaults(conf, false, "showHierarchy") {
		conf["hierarchical"] = int64(1)
	}

	class := maps.GetAnyAnyValWithDefaults(conf, "", "className")
	classes := strings.Split(class, " ")
	classes = append(classes, "wp-block-categories")
	if conf["dropdown"].(int64) == 1 {
		classes = append(classes, "wp-block-categories-dropdown")
		conf["className"] = strings.Join(classes, " ")
	} else {
		classes = append(classes, "wp-block-categories-list")
		conf["className"] = strings.Join(classes, " ")
	}
	return conf, true
}

var GetCategoryArgs = reload.BuildValFnWithAnyParams("block-category-args", categoryArgs)

func categoryArgs(_ ...any) map[string]string {
	args := wp.GetComponentsArgs(widgets.Widget, map[string]string{})
	return maps.FilterZeroMerge(categoryDefaultArgs(), args)
}

func Category(h *wp.Handle, id string, blockParser ParserBlock) (func() string, error) {
	counter := number.Counters[int]()
	conf := GetCategoryConf(blockParser)

	if conf == nil {
		return nil, errors.New("解析block-category配置错误")
	}

	if maps.GetAnyAnyValWithDefaults(conf, false, "showEmpty") {
		h.C.Set("showEmpty", true)
	}
	if maps.GetAnyAnyValWithDefaults(conf, false, "showOnlyTopLevel") {
		h.C.Set("showOnlyTopLevel", true)
	}
	args := GetCategoryArgs()

	return func() string {
		return category(h, id, counter, args, conf)
	}, nil
}

func category(h *wp.Handle, id string, counter number.Counter[int], args map[string]string, conf map[any]any) string {
	var out = ""
	categories := cache.CategoriesTags(h.C, constraints2.Category)
	class := []string{"widget", "widget_block", "widget_categories"}
	if conf["dropdown"].(int64) == 1 {
		out = dropdown(h, categories, counter(), args, conf)
	} else {
		out = categoryUl(h, categories, conf)
	}
	before := fmt.Sprintf(args["{$before_widget}"], str.Join("block-", id), strings.Join(class, " "))
	return str.Join(before, out, args["{$after_widget}"])
}

func categoryUl(h *wp.Handle, categories []models.TermsMy, confAttr map[any]any) string {
	s := str.NewBuilder()
	li := widget.CategoryLi(h, confAttr, categories)
	attrs := GetCategoryAttr(confAttr)
	s.Sprintf(`<ul %s>%s</ul>`, attrs, li)
	return s.String()
}

func dropdown(h *wp.Handle, categories []models.TermsMy, id int, args map[string]string, confAttr map[any]any) string {
	s := str.NewBuilder()
	ids := fmt.Sprintf(`wp-block-categories-%v`, id)
	args = maps.Copy(args)
	args["{$selectId}"] = ids
	attrs := GetCategoryAttr(confAttr)
	selects := widget.DropdownCategories(h, args, confAttr, categories)
	s.Sprintf(`<div %s><label class="screen-reader-text" for="%s">%s</label>%s%s</div>`, attrs, ids, args["{$title}"], selects, strings.ReplaceAll(categoryDropdownScript, "{$id}", ids))
	return s.String()
}

var categoryDropdownScript = `
<script type='text/javascript'>
	/* <![CDATA[ */
	( function() {
		const dropdown = document.getElementById( '{$id}' );
		function onCatChange() {
			if ( dropdown.options[ dropdown.selectedIndex ].value > 0 ) {
				location.href = "/?cat=" + dropdown.options[ dropdown.selectedIndex ].value;
			}
		}
		dropdown.onchange = onCatChange;
	})();
	/* ]]> */
	</script>
`
