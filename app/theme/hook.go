package theme

import (
	"github.com/fthvgb1/wp-go/app/theme/twentyfifteen"
	"github.com/fthvgb1/wp-go/app/theme/wp"
	"github.com/fthvgb1/wp-go/safety"
)

var themeMap = safety.NewMap[string, func(*wp.Handle)]()

func AddThemeHookFunc(name string, fn func(handle *wp.Handle)) {
	if _, ok := themeMap.Load(name); ok {
		panic("exists same name theme")
	}
	themeMap.Store(name, fn)
}

func Hook(themeName string, h *wp.Handle) {
	fn, ok := themeMap.Load(themeName)
	if ok && fn != nil {
		fn(h)
		return
	}
	twentyfifteen.Hook(h)
}
