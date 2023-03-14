package components

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

var widgets = map[string]func(*wp.Handle) string{
	"search-2":          widget.SearchForm,
	"recent-posts-2":    widget.RecentPosts,
	"recent-comments-2": widget.RecentComments,
	"archives-2":        widget.Archive,
	"categories-2":      widget.Category,
}

func WidgetArea(h *wp.Handle) {
	v := wpconfig.GetPHPArrayVal("sidebars_widgets", []any{}, "sidebar-1")
	sidebar := slice.FilterAndMap(v, func(t any) (func(*wp.Handle) string, bool) {
		vv := t.(string)
		fn, ok := widgets[vv]
		if ok {
			return fn, true
		}
		return nil, false
	})
	h.PushHandleFn(constraints.Ok, wp.NewHandleFn(func(h *wp.Handle) {
		h.PushGroupComponentFns(constraints.SidebarsWidgets, 10, sidebar...)
	}, 30))
	h.SetData("categories", cache.CategoriesTags(h.C, constraints.Category))
}
