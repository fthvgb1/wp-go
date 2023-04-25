package wphandle

import (
	"github.com/fthvgb1/wp-go/internal/pkg/config"
	"github.com/fthvgb1/wp-go/internal/plugins/wphandle/enlightjs"
	"github.com/fthvgb1/wp-go/internal/plugins/wphandle/hiddenlogin"
	"github.com/fthvgb1/wp-go/internal/theme/wp"
)

var plugins = wp.HandlePlugins{
	"enlightjs":   enlightjs.EnlighterJS,
	"hiddenLogin": hiddenlogin.HiddenLogin,
}

func RegisterPlugins(m wp.HandlePlugins) {
	for k, v := range m {
		if _, ok := plugins[k]; !ok {
			plugins[k] = v
		}
	}
}

func UsePlugins(h *wp.Handle, calls ...string) {
	calls = append(calls, config.GetConfig().Plugins...)
	for _, call := range calls {
		if fn, ok := plugins[call]; ok {
			fn(h)
		}
	}
}
