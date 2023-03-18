package components

import (
	"fmt"
	"github.com/fthvgb1/wp-go/helper/slice"
	"github.com/fthvgb1/wp-go/internal/cmd/reload"
	"github.com/fthvgb1/wp-go/internal/pkg/cache"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints"
	"github.com/fthvgb1/wp-go/internal/pkg/constraints/widgets"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
	"github.com/fthvgb1/wp-go/internal/theme/wp/components/widget"
	"github.com/fthvgb1/wp-go/internal/wpconfig"
	"strings"
)

var widgetFn = map[string]wp.Components[string]{
	"search":          {Fn: widget.Search, CacheKey: "widgetSearch"},
	"recent-posts":    {Fn: widget.RecentPosts},
	"recent-comments": {Fn: widget.RecentComments},
	"archives":        {Fn: widget.Archive},
	"categories":      {Fn: widget.Category},
	"meta":            {Fn: widget.Meta, CacheKey: "widgetMeta"},
}

type Widget struct {
	Fn       func(*wp.Handle) string
	CacheKey string
}

func WidgetArea(h *wp.Handle) {
	sidebar := reload.GetAnyValBys("sidebarWidgets", h, sidebars)
	h.PushComponents(constraints.SidebarsWidgets, sidebar...)
	h.SetData("categories", cache.CategoriesTags(h.C, constraints.Category))
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
			for k, val := range args {
				wp.SetComponentsArgsForMap(h, name, k, val)
			}
			components.Order = 10
			return components, true
		}
		return components, false
	})
}
