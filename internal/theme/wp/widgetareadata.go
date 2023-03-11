package wp

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
)

var widgets = map[string]func(*Handle) string{
	"search-2":          SearchForm,
	"recent-posts-2":    RecentPosts,
	"recent-comments-2": RecentComments,
	"archives-2":        Archive,
}

func (h *Handle) WidgetArea() {
	v := wpconfig.GetPHPArrayVal("sidebars_widgets", []any{}, "sidebar-1")
	sidebar := slice.FilterAndMap(v, func(t any) (func(*Handle) string, bool) {
		widget := t.(string)
		fn, ok := widgets[widget]
		if ok {
			return fn, true
		}
		return nil, false
	})
	h.PushHandleFn(constraints.Ok, NewHandleFn(func(h *Handle) {
		h.PushGroupCacheComponentFn(constraints.SidebarsWidgets, constraints.SidebarsWidgets, 10, sidebar...)
	}, 30))
	h.ginH["categories"] = cache.CategoriesTags(h.C, constraints.Category)
}
