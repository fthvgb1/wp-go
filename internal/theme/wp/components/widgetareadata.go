package components

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var widgetFn = map[string]func(*wp.Handle) string{
	"search":          widget.SearchForm,
	"recent-posts":    widget.RecentPosts,
	"recent-comments": widget.RecentComments,
	"archives":        widget.Archive,
	"categories":      widget.Category,
	"meta":            widget.Meta,
}

func WidgetArea(h *wp.Handle) {
	args := wp.GetComponentsArgs(h, widgets.Widget, map[string]string{})
	beforeWidget, ok := args["{$before_widget}"]
	if !ok {
		beforeWidget = ""
	} else {
		delete(args, "{$before_widget}")
	}
	v := wpconfig.GetPHPArrayVal("sidebars_widgets", []any{}, "sidebar-1")
	sidebar := slice.FilterAndMap(v, func(t any) (func(*wp.Handle) string, bool) {
		vv := t.(string)
		ss := strings.Split(vv, "-")
		id := ss[len(ss)-1]
		name := strings.Join(ss[0:len(ss)-1], "-")
		fn, ok := widgetFn[name]
		if ok {
			if id != "2" {
				wp.SetComponentsArgsForMap(h, name, "{$id}", id)
			}
			if beforeWidget != "" {
				n := strings.ReplaceAll(name, "-", "_")
				if name == "recent-posts" {
					n = "recent_entries"
				}
				wp.SetComponentsArgsForMap(h, name, "{$before_widget}", fmt.Sprintf(beforeWidget, vv, n))
			}
			if len(args) > 0 {
				for k, val := range args {
					wp.SetComponentsArgsForMap(h, name, k, val)
				}
			}
			return fn, true
		}
		return nil, false
	})
	h.PushGroupComponentFns(constraints.SidebarsWidgets, 10, sidebar...)
	h.SetData("categories", cache.CategoriesTags(h.C, constraints.Category))
}
