package wphandle

import (
	"github.com/fthvgb1/wp-go/helper/maps"
	"github.com/fthvgb1/wp-go/internal/plugins/wphandle/enlightjs"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
)

var plugins = wp.HandlePlugins{
	"enlightjs": enlightjs.EnlighterJS,
}

func Plugins() wp.HandlePlugins {
	return maps.Copy(plugins)
}

func RegisterPlugins(h *wp.Handle, calls ...string) {
	for _, call := range calls {
		if fn, ok := plugins[call]; ok {
			fn(h)
		}
	}
}
