package theme

import (
	"github.com/fthvgb1/wp-go/app/theme/twentyfifteen"
	"github.com/fthvgb1/wp-go/app/theme/wp"
)

var themeMap = map[string]func(*wp.Handle){}

func addThemeHookFunc(name string, fn func(handle *wp.Handle)) {
	if _, ok := themeMap[name]; ok {
		panic("exists same name theme")
	}
	themeMap[name] = fn
}

func Hook(themeName string, handle *wp.Handle) {
	fn, ok := themeMap[themeName]
	if ok && fn != nil {
		fn(handle)
		return
	}
	themeMap[twentyfifteen.ThemeName](handle)
}
