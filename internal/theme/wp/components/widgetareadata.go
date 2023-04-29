package components

import (
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var widgetFn = map[string]widgetComponent{
	"search":          {fn: widget.Search, name: "search"},
	"recent-posts":    {fn: widget.RecentPosts, name: "recent-posts"},
	"recent-comments": {fn: widget.RecentComments, name: "recent-comments"},
	"archives":        {fn: widget.Archive, name: "archives"},
	"categories":      {fn: widget.Category, name: "categories"},
	"meta":            {fn: widget.Meta, name: "meta", cached: true},
}

type widgetComponent struct {
	fn     func(h *wp.Handle, id string) string
	cached bool
	name   string
}

func WidgetArea(h *wp.Handle) {
	h.PushComponents(constraints.SidebarsWidgets, sidebars()...)
}

func sidebars() []wp.Components[string] {
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
			component.Name = widgetComponents.name
			component.Cached = widgetComponents.cached
		}
		component.Order = 10
		return component, true
	})
}
