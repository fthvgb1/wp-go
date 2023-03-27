package components

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	str "github.com/fthvgb1/wp-go/helper/strings"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var widgetFn = map[string]wp.Components[string]{
	"search":          {Fn: widget.Search},
	"recent-posts":    {Fn: widget.RecentPosts},
	"recent-comments": {Fn: widget.RecentComments},
	"archives":        {Fn: widget.Archive},
	"categories":      {Fn: widget.Category},
	"meta":            {Fn: widget.Meta, CacheKey: "widgetMeta"},
}

func WidgetArea(h *wp.Handle) {
	sidebar := reload.GetAnyValBys("sidebarWidgets", h, sidebars)
	h.PushComponents(constraints.SidebarsWidgets, sidebar...)
}

func sidebars(h *wp.Handle) []wp.Components[string] {
	args := wp.GetComponentsArgs(h, widgets.Widget, map[string]string{})
	beforeWidget, ok := args["{$before_widget}"]
	if !ok {
		beforeWidget = ""
	} else {
		delete(args, "{$before_widget}")
	}
	v := wpconfig.GetPHPArrayVal("sidebars_widgets", []any{}, "sidebar-1")
	return slice.FilterAndMap(v, func(t any) (wp.Components[string], bool) {
		vv := t.(string)
		ss := strings.Split(vv, "-")
		id := ss[len(ss)-1]
		name := strings.Join(ss[0:len(ss)-1], "-")
		components, ok := widgetFn[name]
		if name != "block" && !ok {
			return components, false
		}
		if id != "2" {
			wp.SetComponentsArgsForMap(h, name, "{$id}", id)
		}
		names := str.Join("widget-", name)
		if beforeWidget != "" {
			n := strings.ReplaceAll(name, "-", "_")
			if name == "recent-posts" {
				n = "recent_entries"
			}
			wp.SetComponentsArgsForMap(h, names, "{$before_widget}", fmt.Sprintf(beforeWidget, vv, n))
		}
		for k, val := range args {
			wp.SetComponentsArgsForMap(h, names, k, val)
		}
		if name == "block" {
			fn := Block(id)
			if fn == nil {
				return wp.Components[string]{}, false
			}
			components = wp.Components[string]{Fn: fn, Order: 10}
		}
		components.Order = 10
		return components, true
	})
}
