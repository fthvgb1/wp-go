package theme

import (
	"github.com/fthvgb1/wp-go/internal/theme/common"
	"github.com/fthvgb1/wp-go/internal/theme/twentyfifteen"
)

var themeMap = map[string]func(handle common.Handle){}

func addThemeHookFunc(name string, fn func(handle common.Handle)) {
	if _, ok := themeMap[name]; ok {
		panic("exists same name theme")
	}
	themeMap[name] = fn
}

func Hook(themeName string, handle common.Handle) {
	fn, ok := themeMap[themeName]
	if ok && fn != nil {
		fn(handle)
		return
	}
	themeMap[twentyfifteen.ThemeName](handle)
}
