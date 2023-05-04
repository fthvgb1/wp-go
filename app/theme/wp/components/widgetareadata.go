package components

import (
	"github.com/fthvgb1/wp-go/app/pkg/constraints"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/app/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/app/wpconfig"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"strings"
)

var widgetFn = map[string]widgetComponent{
	"search":          {fn: widget.Search, name: "widget.search"},
	"recent-posts":    {fn: widget.RecentPosts, name: "widget.recent-posts"},
	"recent-comments": {fn: widget.RecentComments, name: "widget.recent-comments"},
	"archives":        {fn: widget.Archive, name: "widget.archives"},
	"categories":      {fn: widget.Category, name: "widget.categories"},
	"meta":            {fn: widget.Meta, name: "widget.meta", cached: true},
}

type widgetComponent struct {
	fn     func(h *wp.Handle, id string) string
	cached bool
	name   string
}

func WidgetArea(h *wp.Handle) {
	h.PushComponents(constraints.AllScene, constraints.SidebarsWidgets, sidebars()...)
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
			fn, fnName := Block(id)
			if fn == nil {
				return component, false
			}
			component.Fn = fn
			component.Name = str.Join("block.", fnName)
		} else {
			component.Fn = widget.Fn(id, widgetComponents.fn)
			component.Name = widgetComponents.name
			component.Cached = widgetComponents.cached
		}
		component.Order = 10
		return component, true
	})
}
