package components

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var widgets = map[string]func(*wp.Handle) string{
	"search":          widget.SearchForm,
	"recent-posts":    widget.RecentPosts,
	"recent-comments": widget.RecentComments,
	"archives":        widget.Archive,
	"categories":      widget.Category,
	"meta":            widget.Meta,
}

func WidgetArea(h *wp.Handle) {
	v := wpconfig.GetPHPArrayVal("sidebars_widgets", []any{}, "sidebar-1")
	sidebar := slice.FilterAndMap(v, func(t any) (func(*wp.Handle) string, bool) {
		vv := t.(string)
		ss := strings.Split(vv, "-")
		id := ss[len(ss)-1]
		name := strings.Join(ss[0:len(ss)-1], "-")
		fn, ok := widgets[name]
		if ok {
			if id != "2" {
				wp.SetComponentsArgsForMap(h, name, "{$id}", id)
			}
			return fn, true
		}
		return nil, false
	})
	h.PushHandleFn(constraints.Ok, wp.NewHandleFn(func(h *wp.Handle) {
		h.PushGroupComponentFns(constraints.SidebarsWidgets, 10, sidebar...)
	}, 30))
	h.SetData("categories", cache.CategoriesTags(h.C, constraints.Category))
}
