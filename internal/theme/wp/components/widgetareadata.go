package components

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var widgetFn = map[string]widgetComponent{
	"search":          {fn: widget.Search},
	"recent-posts":    {fn: widget.RecentPosts},
	"recent-comments": {fn: widget.RecentComments},
	"archives":        {fn: widget.Archive},
	"categories":      {fn: widget.Category},
	"meta":            {fn: widget.Meta, cacheKey: "widgetMeta"},
}

type widgetComponent struct {
	fn       func(h *wp.Handle, id string) string
	cacheKey string
}

func WidgetArea(h *wp.Handle) {
	sidebar := reload.GetAnyValBys("sidebarWidgets", h, sidebars)
	h.PushComponents(constraints.SidebarsWidgets, sidebar...)
}

func sidebars(*wp.Handle) []wp.Components[string] {
	v := wpconfig.GetPHPArrayVal("sidebars_widgets", []any{}, "sidebar-1")
	return slice.FilterAndMap(v, func(t any) (wp.Components[string], bool) {
		vv := t.(string)
		ss := strings.Split(vv, "-")
		id := ss[len(ss)-1]
		name := strings.Join(ss[0:len(ss)-1], "-")
		widgetComponents, ok := widgetFn[name]
		if name != "block" && !ok {
			return wp.Components[string]{}, false
		}
		var component wp.Components[string]
		if name == "block" {
			fn := Block(id)
			if fn == nil {
				return component, false
			}
			component.Fn = fn
		} else {
			component.Fn = widget.Fn(id, widgetComponents.fn)
			component.CacheKey = widgetComponents.cacheKey
		}
		component.Order = 10
		return component, true
	})
}
